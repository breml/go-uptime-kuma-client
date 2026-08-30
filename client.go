package kuma

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"

	"github.com/breml/go-uptime-kuma-client/dockerhost"
	"github.com/breml/go-uptime-kuma-client/maintenance"
	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
	"github.com/breml/go-uptime-kuma-client/proxy"
	"github.com/breml/go-uptime-kuma-client/statuspage"
)

// ErrNotFound is returned when a requested resource is not found.
var ErrNotFound = errors.New("not found")

// ErrUpdateEventTimeout is returned when the server acknowledged a mutating
// command successfully, but the update event that carries the change did not
// arrive before the context was done. The write landed on the server; what is
// missing is only the broadcast.
//
// An error wrapping ErrUpdateEventTimeout also wraps the context error, so
// errors.Is(err, context.DeadlineExceeded) still reports an expired deadline
// and errors.Is(err, context.Canceled) a cancelled context. Use
// errors.Is(err, ErrUpdateEventTimeout) to tell the two cases apart: a plain
// context error means the server never confirmed the command and it may or may
// not have been applied, while ErrUpdateEventTimeout means it was applied.
//
// The Create methods return the ID the server assigned alongside such an error,
// so a caller can adopt the created resource instead of retrying and creating a
// duplicate. That ID is zero only if the ack carried none, which the server does
// not do for a command it reports as successful. Because the idiomatic way to
// propagate an error drops the values next to it, the ID is also carried by the
// error itself, see UpdateEventTimeoutError.
//
// The local state cache is refreshed from the update event, so on this path it
// is not up to date yet. The cache is maintained by the socket.io event
// handlers and not by the caller's context, so it catches up on its own once
// the event arrives; if the event is lost for good, the cache stays stale until
// the next broadcast for that resource. Until then a getter that serves from
// the cache does not report the resource; Resync makes the server resend the
// lists instead of waiting for that broadcast.
var ErrUpdateEventTimeout = errors.New("update event not received")

// UpdateEventTimeoutError is the error the client returns for a command that
// the server acknowledged without the update event arriving, see
// ErrUpdateEventTimeout. It wraps both ErrUpdateEventTimeout and the context
// error, so errors.Is keeps reporting either.
//
// It exists so that the ID of a created resource survives the way errors are
// usually propagated: a caller that writes the idiomatic
//
//	id, err := client.CreateNotification(ctx, notif)
//	if err != nil {
//		return 0, err
//	}
//
// discards the returned ID, and with it the only handle on a notification the
// server did create. Recovering the error keeps that handle:
//
//	var timeoutErr *kuma.UpdateEventTimeoutError
//	if errors.As(err, &timeoutErr) && timeoutErr.ID != 0 {
//		// The resource exists, adopt it instead of creating it again.
//	}
type UpdateEventTimeoutError struct {
	// Command is the socket.io command the server acknowledged.
	Command string

	// ID is the ID the server assigned to the created resource. It is zero for
	// a command that creates nothing, and for an ack that carried no ID.
	ID int64

	// Err is the error of the context that was done before the update event
	// arrived.
	Err error
}

func (e *UpdateEventTimeoutError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Command, ErrUpdateEventTimeout, e.Err)
}

func (e *UpdateEventTimeoutError) Unwrap() []error {
	return []error{ErrUpdateEventTimeout, e.Err}
}

// withCreatedID records id in the UpdateEventTimeoutError err is wrapping, if it
// is one, and returns err unchanged otherwise. The command layer knows which
// field of the ack carries the ID, the layer that builds the error does not.
func withCreatedID(err error, id int64) error {
	var timeoutErr *UpdateEventTimeoutError
	if errors.As(err, &timeoutErr) {
		timeoutErr.ID = id
	}

	return err
}

// Log level constants for configuring socket.io client logging verbosity.
const (
	LogLevelDebug = utils.DEBUG
	LogLevelInfo  = utils.INFO
	LogLevelWarn  = utils.WARN
	LogLevelError = utils.ERROR
	LogLevelNone  = utils.NONE
)

// LogLevel converts a string log level to its corresponding integer constant.
func LogLevel(level string) int {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LogLevelDebug

	case "INFO":
		return LogLevelInfo

	case "WARN":
		return LogLevelWarn

	case "ERROR":
		return LogLevelError

	default:
		return LogLevelNone
	}
}

//nolint:gochecknoglobals // empty is used as a placeholder value in maps to represent set membership.
var empty = struct{}{}

// knownReadyEvents lists all initial list events an Uptime Kuma server may emit
// after login. The subset a given server actually emits depends on its version
// and on any reverse proxy in front of it (older versions and some proxies never
// emit certain events, e.g. apiKeyList).
//
//nolint:gochecknoglobals // package-level constant list of the known ready events.
var knownReadyEvents = []string{
	"monitorList",
	"notificationList",
	"statusPageList",
	"maintenanceList",
	"proxyList",
	"dockerHostList",
	"apiKeyList",
}

// defaultReadyEvents is the set of events New waits for before returning. These
// are the primary entities this library manages and are emitted by all
// supported Uptime Kuma versions. The remaining knownReadyEvents are treated as
// best-effort: they are cached when they arrive, but never block New.
//
//nolint:gochecknoglobals // package-level default for the required ready events.
var defaultReadyEvents = []string{
	"monitorList",
	"notificationList",
	"statusPageList",
}

// defaultReadyGracePeriod bounds how long New waits for the best-effort
// (optional) ready events once all required events have been received. On a
// healthy server all events arrive within milliseconds, so this only elapses
// when an optional event is genuinely never sent.
const defaultReadyGracePeriod = 500 * time.Millisecond

type entryPageResponse struct {
	Type string `json:"type"`
}

type dbConfig struct {
	Type string `json:"type"`
}

type setupDatabaseRequest struct {
	DBConfig dbConfig `json:"dbConfig"`
}

type setupDatabaseResponse struct {
	OK bool `json:"ok"`
}

type state struct {
	notifications []notification.Base
	monitors      []monitor.Base
	statusPages   map[int64]statuspage.StatusPage
	maintenances  []maintenance.Maintenance
	proxies       []proxy.Proxy
	dockerHosts   []dockerhost.DockerHost
}

// Client represents a connection to an Uptime Kuma server.
type Client struct {
	socketioClient               *socketio.Client
	socketioClientConnectTimeout time.Duration
	socketioLogger               socketio.Logger
	autosetup                    bool
	readyEvents                  []string
	readyGracePeriod             time.Duration

	mu      *sync.Mutex
	updates signals.Signal[string]
	state   state

	// sessionToken is the JWT the server handed out at login. It is what
	// Resync logs in with, see there.
	sessionToken string

	// autoLoggedIn records that the server has authentication disabled and
	// logged the client in itself, which leaves it without a session token,
	// see Resync.
	autoLoggedIn bool

	// totpCode produces the one-time code for an account with two-factor
	// authentication enabled, see WithTOTPSecret and WithTOTPCode.
	// totpCodeSet tells a configured callback from none, totpSources counts
	// how many of the two options set one, and totpErr carries a secret New
	// has to reject.
	totpCode    func(ctx context.Context) (string, error)
	totpCodeSet bool
	totpSources int
	totpErr     error

	// readyEventsMissing are the best-effort ready events the server did not
	// emit while New was connecting, see MissingReadyEvents.
	readyEventsMissing []string
}

// Option is a functional option for configuring a Client.
type Option func(c *Client)

// WithAutosetup enables automatic server setup during client connection.
func WithAutosetup() Option {
	return func(c *Client) {
		c.autosetup = true
	}
}

// WithLogLevel sets the socket.io client logging level.
func WithLogLevel(level int) Option {
	return func(c *Client) {
		if level >= utils.DEBUG && level <= utils.NONE {
			c.socketioLogger = &utils.DefaultLogger{Level: level}
		}
	}
}

// WithConnectTimeout sets the socket.io client connection timeout.
// This timeout defines the overall duration, with is allowed for establishing
// the connection to Uptime Kuma.
// In the case of autosetup with an uninitialized Uptime Kuma this timeout
// also includes the time required for the initial setup.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.socketioClientConnectTimeout = timeout
	}
}

// WithReadyEvents overrides the set of initial list events New waits for before
// returning. By default New waits for monitorList, notificationList and
// statusPageList; the remaining known events (e.g. maintenanceList, apiKeyList)
// are best-effort. Use this to widen or narrow the required set for servers or
// reverse proxies that emit a different subset of events.
//
// An optional event that never arrives (server never sends it, or a proxy
// drops it) leaves the corresponding client state (e.g. GetMaintenances,
// GetProxyList, GetDockerHostList) populated as empty, indistinguishable from
// a genuinely empty list on the server. MissingReadyEvents names those events,
// and a warning goes to the logger WithLogLevel configures, which is silent by
// default. Widen readyEvents to require an event if its state must be
// trustworthy.
//
// New rejects an event that is not one of the known ready events, rather than
// waiting out the connect timeout for something the server never emits.
func WithReadyEvents(events ...string) Option {
	return func(c *Client) {
		c.readyEvents = slices.Clone(events)
	}
}

// WithReadyGracePeriod sets how long New waits for the best-effort (optional)
// ready events after all required events have been received. A value of zero
// makes New return as soon as the required events arrive, without waiting for
// the optional ones.
func WithReadyGracePeriod(d time.Duration) Option {
	return func(c *Client) {
		c.readyGracePeriod = d
	}
}

// setupDatabase handles the database setup phase for Uptime Kuma v2.
// It checks if database setup is needed and configures SQLite if required.
// The function will wait for the server to restart after database configuration.
//
//nolint:revive // Complexity is necessary for complete database setup logic
func setupDatabase(ctx context.Context, baseURL string) error {
	// Convert socket.io URL to HTTP URL
	httpURL := strings.Replace(baseURL, "ws://", "http://", 1)
	httpURL = strings.Replace(httpURL, "wss://", "https://", 1)

	// Check if database setup is needed
	entryPageURL := httpURL + "/api/entry-page"

	var entryPage entryPageResponse

	// Check if parent context is already cancelled
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())

	default:
	}

	// Check entry-page without retry - let the caller (pool.Retry) handle retries
	// Use a longer timeout for the HTTP request itself, independent of parent context
	httpCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodGet, entryPageURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("create entry-page request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Return connection errors as-is so caller can retry
		return fmt.Errorf("entry-page request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("entry-page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read entry-page response: %w", err)
	}

	err = json.Unmarshal(body, &entryPage)
	if err != nil {
		return fmt.Errorf("parse entry-page response: %w", err)
	}

	// If database setup is not needed, return early
	if entryPage.Type != "setup-database" {
		return nil
	}

	// Configure database with SQLite
	setupDBURL := httpURL + "/setup-database"
	setupReq := setupDatabaseRequest{
		DBConfig: dbConfig{
			Type: "sqlite",
		},
	}

	reqBody, err := json.Marshal(setupReq)
	if err != nil {
		return fmt.Errorf("marshal setup request: %w", err)
	}

	httpCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err = http.NewRequestWithContext(httpCtx, http.MethodPost, setupDBURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create setup-database request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("setup database: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("setup-database returned status %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read setup-database response: %w", err)
	}

	var setupResp setupDatabaseResponse
	err = json.Unmarshal(body, &setupResp)
	if err != nil {
		return fmt.Errorf("parse setup-database response: %w", err)
	}

	if !setupResp.OK {
		return errors.New("setup-database failed")
	}

	// Wait for server to restart by polling entry-page until it changes
	// The server should transition from "setup-database" to "setup" (user setup)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for server restart: %w", ctx.Err())

		case <-timeout:
			return errors.New("timeout waiting for server restart")

		case <-ticker.C:
			// Use a short timeout for each poll attempt
			pollCtx, pollCancel := context.WithTimeout(context.Background(), 2*time.Second)
			pollReq, err := http.NewRequestWithContext(pollCtx, http.MethodGet, entryPageURL, http.NoBody)
			if err != nil {
				pollCancel()
				continue
			}

			pollResp, err := http.DefaultClient.Do(pollReq)
			if err != nil {
				pollCancel()
				continue
			}

			pollBody, err := io.ReadAll(pollResp.Body)
			_ = pollResp.Body.Close()
			pollCancel()
			if err != nil {
				continue
			}

			var checkEntryPage entryPageResponse
			err = json.Unmarshal(pollBody, &checkEntryPage)
			if err != nil {
				continue
			}

			// If entry page type changed from "setup-database", server has restarted
			if checkEntryPage.Type != "setup-database" {
				return nil
			}
		}
	}
}

// New creates a new Client connected to an Uptime Kuma server.
//
//nolint:revive // Complexity is necessary for complete client initialization and event setup
func New(ctx context.Context, baseURL string, username string, password string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger:   &utils.DefaultLogger{Level: utils.NONE},
		readyEvents:      defaultReadyEvents,
		readyGracePeriod: defaultReadyGracePeriod,

		mu:      &sync.Mutex{},
		updates: signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
	}

	unknown := unknownReadyEvents(c.readyEvents)
	if len(unknown) > 0 {
		return nil, fmt.Errorf("ready events: unknown: %s", strings.Join(unknown, ", "))
	}

	if (username == "") != (password == "") {
		return nil, errors.New("credentials: username and password have to be set together")
	}

	// The count comes first: a caller who configured two sources is told that
	// before being sent to fix whichever of them also failed to decode.
	if c.totpSources > 1 {
		return nil, errors.New("totp: at most one of WithTOTPSecret and WithTOTPCode may be used")
	}

	if c.totpErr != nil {
		return nil, c.totpErr
	}

	ctxWithConnectTimeout := ctx

	// connectTimeoutDone is non-nil only when WithConnectTimeout is
	// configured. Using nil keeps the ready-wait select deterministic when
	// ctxWithConnectTimeout and ctx are the same context (no timeout set).
	var connectTimeoutDone <-chan struct{}

	if c.socketioClientConnectTimeout != 0 {
		var cancel func()
		ctxWithConnectTimeout, cancel = context.WithTimeout(ctx, c.socketioClientConnectTimeout)
		defer cancel()

		connectTimeoutDone = ctxWithConnectTimeout.Done()
	}

	// Handle database setup for Uptime Kuma v2 if autosetup is enabled
	if c.autosetup {
		err := setupDatabase(ctxWithConnectTimeout, baseURL)
		if err != nil {
			return nil, fmt.Errorf("database setup: %w", err)
		}
	}

	client, err := socketio.NewClient(
		socketio.WithRawURL(baseURL),
		socketio.WithLogger(c.socketioLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %w", err)
	}

	c.socketioClient = client

	gate := newReadyGate(c.splitReadyEvents())

	c.updates.AddListener(func(_ context.Context, event string) {
		gate.observe(event)
	}, "connect-ready")
	defer c.updates.RemoveListener("connect-ready")

	client.On("notificationList", func(notificationList []notification.Base) {
		c.mu.Lock()
		c.state.notifications = notificationList
		defer c.mu.Unlock()

		c.updates.Emit(context.Background(), "notificationList")
	})

	client.On("monitorList", func(monitorMap map[string]monitor.Base) {
		c.mu.Lock()
		defer c.mu.Unlock()

		// Convert map to slice
		monitors := make([]monitor.Base, 0, len(monitorMap))
		for _, monitor := range monitorMap {
			monitors = append(monitors, monitor)
		}

		c.state.monitors = monitors

		c.updates.Emit(context.Background(), "monitorList")
	})

	// Uptime Kuma v2 sends updateMonitorIntoList for individual monitor updates (add/edit/pause/resume)
	client.On("updateMonitorIntoList", func(monitorMap map[string]monitor.Base) {
		c.mu.Lock()
		defer c.mu.Unlock()

		// Update or add the monitors in the map to our state
		for _, updatedMonitor := range monitorMap {
			found := false
			for i, existingMonitor := range c.state.monitors {
				if existingMonitor.ID == updatedMonitor.ID {
					c.state.monitors[i] = updatedMonitor
					found = true
					break
				}
			}

			if !found {
				c.state.monitors = append(c.state.monitors, updatedMonitor)
			}
		}

		c.updates.Emit(context.Background(), "updateMonitorIntoList")
	})

	// Uptime Kuma v2 sends deleteMonitorFromList when a monitor is deleted
	client.On("deleteMonitorFromList", func(monitorID int64) {
		c.mu.Lock()
		defer c.mu.Unlock()

		// Remove the monitor from our state
		for i, existingMonitor := range c.state.monitors {
			if existingMonitor.ID == monitorID {
				c.state.monitors = append(c.state.monitors[:i], c.state.monitors[i+1:]...)
				break
			}
		}

		c.updates.Emit(context.Background(), "deleteMonitorFromList")
	})

	client.On("statusPageList", func(statusPageMap map[int64]statuspage.StatusPage) {
		c.mu.Lock()
		c.state.statusPages = statusPageMap
		defer c.mu.Unlock()

		c.updates.Emit(context.Background(), "statusPageList")
	})

	client.On("maintenanceList", func(maintenanceMap map[string]maintenance.Maintenance) {
		c.mu.Lock()
		defer c.mu.Unlock()

		// Convert map to slice
		maintenances := make([]maintenance.Maintenance, 0, len(maintenanceMap))
		for _, m := range maintenanceMap {
			maintenances = append(maintenances, m)
		}

		c.state.maintenances = maintenances

		c.updates.Emit(context.Background(), "maintenanceList")
	})

	client.On("proxyList", func(proxyList []proxy.Proxy) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.state.proxies = proxyList

		c.updates.Emit(context.Background(), "proxyList")
	})

	client.On("dockerHostList", func(dockerHostList []dockerhost.DockerHost) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.state.dockerHosts = dockerHostList

		c.updates.Emit(context.Background(), "dockerHostList")
	})

	connect := make(chan struct{})
	closeConnect := sync.OnceFunc(func() {
		close(connect)
	})
	defer closeConnect()

	client.On("connect", func() {
		closeConnect()
	})

	setupRequired := make(chan struct{})
	closeSetupRequired := sync.OnceFunc(func() {
		close(setupRequired)
	})
	defer closeSetupRequired()

	if c.autosetup {
		client.On("setup", func() {
			closeSetupRequired()
		})
	}

	// The server states how it wants the connection to be authenticated as the
	// last thing it does for it, so these are registered before the connect and
	// unconditionally, see authBarrier.
	loginRequired := make(chan struct{})
	closeLoginRequired := sync.OnceFunc(func() {
		close(loginRequired)
	})
	defer closeLoginRequired()

	client.On("loginRequired", func() {
		closeLoginRequired()
	})

	autoLogin := make(chan struct{})
	closeAutoLogin := sync.OnceFunc(func() {
		close(autoLogin)
	})
	defer closeAutoLogin()

	client.On("autoLogin", func() {
		closeAutoLogin()
	})

	barrier := authBarrier{
		loginRequired: loginRequired,
		autoLogin:     autoLogin,
		setupRequired: setupRequired,
	}

	client.OnAny(func(s string, _ []any) {
		if s != "notificationList" && s != "monitorList" && s != "statusPageList" && s != "maintenanceList" &&
			s != "proxyList" &&
			s != "dockerHostList" {
			c.updates.Emit(context.Background(), s)
		}
	})

	// client.Connect returns as soon as the transport is running and the
	// handshake has been requested; a non-nil error is the real transport
	// failure, which is surfaced instead of waiting out the timeout. It gets
	// ctx and not ctxWithConnectTimeout, because the socket.io client keeps
	// reading from the connection until that context is done and
	// ctxWithConnectTimeout is canceled as soon as New returns.
	connectErr := make(chan error, 1)
	go func() {
		connectErr <- client.Connect(ctx)
	}()

connectLoop:
	for {
		select {
		case <-connect:
			break connectLoop

		case err := <-connectErr:
			if err != nil {
				return nil, fmt.Errorf("connect to server: %w", err)
			}

			// The dial returned; keep waiting for the connect event. Setting
			// connectErr to nil makes this case block forever, and marks the
			// channel as spent by a transport that is up.
			connectErr = nil

		case <-ctx.Done():
			return nil, fmt.Errorf("connect to server: %w", ctx.Err())

		case <-ctxWithConnectTimeout.Done():
			if connectErr == nil {
				c.disconnectOrphan()
			} else {
				dialErr := c.abandonDial(connectErr)
				if dialErr != nil {
					return nil, fmt.Errorf("connect to server: %w", dialErr)
				}
			}

			return nil, fmt.Errorf("connect to server: %w", ctxWithConnectTimeout.Err())
		}
	}

	// The socket.io client is now connected. On any subsequent error path
	// the caller receives no *Client handle and therefore cannot call
	// Disconnect themselves. Trigger a best-effort async close so that
	// goroutines and connections are eventually cleaned up.
	// Cleared to false on the success path so the caller takes ownership.
	closeOnErr := true
	defer func() {
		if closeOnErr {
			go func() {
				disconnectErr := c.Disconnect()
				if disconnectErr != nil {
					c.socketioLogger.Errorf("disconnect after New() error: %s", disconnectErr)
				}
			}()
		}
	}()

	err = c.authenticate(
		ctxWithConnectTimeout,
		credentials{username: username, password: password},
		barrier,
	)
	if err != nil {
		// A server that is not set up yet rejects the credentials it does not
		// have a user for; running the setup is what makes them work, so that
		// is not a failure when the server asked for one.
		if !errors.Is(err, ErrInvalidCredentials) || !setupPending(setupRequired) {
			return nil, err
		}
	}

	// Both success paths hand the client over through returnReady, so neither
	// can skip the grace window the best-effort events get. connectTimeoutDone
	// keeps that wait inside the caller's WithConnectTimeout budget.
	returnReady := func() *Client {
		c.awaitOptionalReadyEvents(ctx, "New", gate, connectTimeoutDone)
		c.setMissingReadyEvents(gate.missingOptional())

		closeOnErr = false

		return c
	}

	for {
		// A server that needs setup says so right after the connect, while a
		// required set that is empty (see WithReadyEvents) opens the ready gate
		// before that. Checking setup first keeps it from losing that race and
		// handing back a client for a server that is not set up.
		select {
		case <-setupRequired:
			// Handled below the select.

		default:
			select {
			case <-gate.requiredDone:
				return returnReady(), nil

			case <-setupRequired:
				// Handled below the select.

			case <-ctx.Done():
				return nil, fmt.Errorf("wait for ready: %w", ctx.Err())

			case <-connectTimeoutDone:
				// ctxWithConnectTimeout is derived from ctx, so its Done
				// channel closes whenever the parent ctx is cancelled too.
				// Prefer the parent's error in that case to avoid a misleading
				// "missing events" message on an ordinary cancellation.
				if ctx.Err() != nil {
					return nil, fmt.Errorf("wait for ready: %w", ctx.Err())
				}

				// If all ready events arrived at the exact same instant as the
				// timeout, prefer the success path over the error path.
				select {
				case <-gate.requiredDone:
					return returnReady(), nil

				default:
				}

				return nil, fmt.Errorf(
					"wait for ready: %w (missing events: %s)",
					ctxWithConnectTimeout.Err(),
					strings.Join(gate.missingRequired(), ", "),
				)
			}
		}

		setupRequired = nil

		if !c.autosetup {
			return nil, errors.New("server does require setup, but autosetup is disabled")
		}

		err = c.runSetup(ctxWithConnectTimeout, username, password)
		if err != nil {
			return nil, err
		}
	}
}

// Disconnect closes the connection to the Uptime Kuma server.
func (c *Client) Disconnect() error {
	err := c.socketioClient.Close()
	if err != nil {
		return fmt.Errorf("close socket.io client: %w", err)
	}

	return nil
}

// Resync rebuilds the local state cache from the server and returns once the
// server has resent the lists the cache is built from.
//
// It is the way out of the stale cache an operation that failed with
// ErrUpdateEventTimeout leaves behind: the getters that serve from the cache do
// not report a resource until its list is broadcast again, and for
// notifications, proxies and Docker hosts nothing else triggers that broadcast.
//
// The resync is all or nothing, because the server has no command to re-request
// a single list. What it does have is a login, which answers with all of them,
// so Resync logs in again with the token from the login New performed. A client
// created without credentials therefore cannot resync.
//
// Which lists Resync waits for follows New: the events configured with
// WithReadyEvents are required and the remaining known events are best-effort,
// awaited for the ready grace period. The best-effort events New never saw are
// skipped altogether, because a server (or a proxy in front of it) that does
// not emit one while connecting does not emit it here either.
func (c *Client) Resync(ctx context.Context) error {
	c.mu.Lock()
	token := c.sessionToken
	autoLoggedIn := c.autoLoggedIn
	c.mu.Unlock()

	if token == "" {
		if autoLoggedIn {
			// The server hands a session token only to a login it performed
			// for a client that asked for one, and the commands it does have
			// resend some of the lists but not all of them.
			return errors.New(
				"resync: the server logged the client in itself (authentication disabled) " +
					"and offers no command to resend the lists",
			)
		}

		return errors.New("resync: no session token, the client was created without credentials")
	}

	gate := newReadyGate(c.resyncReadyEvents())

	// Registered before the command is emitted, so that no list can arrive
	// unnoticed between the two.
	listenerID := uuid.New()
	c.updates.AddListener(func(_ context.Context, update string) {
		gate.observe(update)
	}, listenerID.String())

	defer c.updates.RemoveListener(listenerID.String())

	response, err := c.emitAck(ctx, "loginByToken", token)
	if err != nil {
		return fmt.Errorf("resync: %w", err)
	}

	if !response.OK {
		// A token the server no longer accepts is the one failure a caller can
		// act on: it means the password was changed or the user was removed,
		// and a new login is the only way back.
		return fmt.Errorf("resync: %w", loginError("loginByToken", response))
	}

	select {
	case <-gate.requiredDone:

	case <-ctx.Done():
		return fmt.Errorf(
			"resync: %w (missing events: %s)",
			ctx.Err(),
			strings.Join(gate.missingRequired(), ", "),
		)
	}

	// The required lists are in. Give the best-effort ones the same grace
	// window New grants them, so the cache they feed is refreshed too.
	c.awaitOptionalReadyEvents(ctx, "Resync", gate, nil)

	return nil
}

// MissingReadyEvents returns the best-effort ready events the server did not
// emit while New was connecting, sorted.
//
// The state such an event carries (maintenances, proxies, Docker hosts, API
// keys) is empty in the local cache for a reason the getters cannot tell apart
// from a genuinely empty server, so a caller that depends on one of them can
// check here — or require the event with WithReadyEvents and have New fail
// instead.
func (c *Client) MissingReadyEvents() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return slices.Clone(c.readyEventsMissing)
}

// abandonDial cleans up the dial New walks away from when the connect timeout
// fires, and returns the transport error if the dial already failed with one.
// New returns no *Client on this path, so a connection that was, or still gets,
// established has no owner: it is disconnected instead of left running until
// ctx is done.
func (c *Client) abandonDial(connectErr <-chan error) error {
	select {
	case err := <-connectErr:
		if err != nil {
			return err
		}

		c.disconnectOrphan()

	default:
		// The dial is still in flight (it uses ctx, so this timeout does not
		// abort it). Drain connectErr in the background so a late failure is
		// not silently discarded and a late success is not leaked.
		go func() {
			err := <-connectErr
			if err != nil {
				c.socketioLogger.Warnf("New: dial failed after connect timeout: %s", err)

				return
			}

			c.disconnectOrphan()
		}()
	}

	return nil
}

// disconnectOrphan closes a connection no caller holds a handle to. The close
// is best-effort and asynchronous, because Disconnect waits for the transport
// and the message loop to wind down, which a server that leaves a long poll
// hanging can stall for as long as it likes — and the caller it would block is
// on its way out.
func (c *Client) disconnectOrphan() {
	go func() {
		err := c.Disconnect()
		if err != nil {
			c.socketioLogger.Warnf("disconnect orphaned connection: %s", err)
		}
	}()
}

// runSetup sets the server up and logs in again afterwards, which is what a
// server that answers the connect with a setup event asks for.
func (c *Client) runSetup(ctx context.Context, username string, password string) error {
	if username == "" || password == "" {
		return errors.New("setup: the server requires setup, which needs a username and password")
	}

	_, err := c.syncEmit(ctx, "setup", username, password)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	// The user the setup just created has no two-factor authentication, so
	// this login is answered without the server ever asking for a code.
	err = c.loginWithPassword(ctx, username, password)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}

// setMissingReadyEvents records what New waited for in vain, see
// MissingReadyEvents and resyncReadyEvents.
func (c *Client) setMissingReadyEvents(missing []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.readyEventsMissing = missing
}

// setSessionToken records the JWT a login answered with, see Resync.
func (c *Client) setSessionToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessionToken = token
}

// splitReadyEvents returns the update events that carry the lists the local
// state cache is built from, as two sets to strike them off as they arrive:
// the events the client waits for (its configured readyEvents) and the
// best-effort rest, which are cached when they arrive but never block. The
// server sends them after a login, which is what New waits for and what Resync
// repeats.
func (c *Client) splitReadyEvents() (required map[string]struct{}, optional map[string]struct{}) {
	required = make(map[string]struct{}, len(c.readyEvents))
	for _, event := range c.readyEvents {
		required[event] = empty
	}

	optional = make(map[string]struct{}, len(knownReadyEvents))

	for _, event := range knownReadyEvents {
		if _, isRequired := required[event]; !isRequired {
			optional[event] = empty
		}
	}

	return required, optional
}

// resyncReadyEvents is splitReadyEvents for a resync: the best-effort events
// the server did not emit while New was connecting are dropped, because it does
// not emit them here either and waiting them out would cost every resync the
// full ready grace period for nothing.
func (c *Client) resyncReadyEvents() (required map[string]struct{}, optional map[string]struct{}) {
	required, optional = c.splitReadyEvents()

	for _, event := range c.MissingReadyEvents() {
		delete(optional, event)
	}

	return required, optional
}

// unknownReadyEvents returns the events that are not among knownReadyEvents,
// sorted.
func unknownReadyEvents(events []string) []string {
	unknown := make([]string, 0, len(events))

	for _, event := range events {
		if !slices.Contains(knownReadyEvents, event) {
			unknown = append(unknown, event)
		}
	}

	slices.Sort(unknown)

	return unknown
}

// awaitOptionalReadyEvents gives the best-effort ready events the ready grace
// period to arrive, so the state they populate is current when the caller named
// by caller returns, and warns about the ones that never came. A grace period
// of zero skips the wait, and with it the warning: the caller asked to return
// as soon as the required events are in, which is too early to tell an event
// that is missing from one that is merely late.
//
// The wait cannot fail the caller: the required events are already in, so every
// way out of it is a success. Besides the optional events arriving and the
// grace period elapsing, those are a done context and abort, which is how New
// keeps the wait within the caller's WithConnectTimeout budget; a nil abort
// never fires.
func (c *Client) awaitOptionalReadyEvents(
	ctx context.Context,
	caller string,
	gate *readyGate,
	abort <-chan struct{},
) {
	if c.readyGracePeriod <= 0 {
		return
	}

	graceTimer := time.NewTimer(c.readyGracePeriod)
	defer graceTimer.Stop()

	select {
	case <-gate.optionalDone:
	case <-graceTimer.C:
	case <-abort:
	case <-ctx.Done():
	}

	missing := gate.missingOptional()
	if len(missing) > 0 {
		c.socketioLogger.Warnf(
			"%s: optional ready events did not arrive within the grace period: %s",
			caller,
			strings.Join(missing, ", "),
		)
	}
}

// readyGate tracks the initial list events a login answers with, so New and
// Resync can wait for the required ones and give the best-effort rest a grace
// period. Events are struck off both sets as they arrive, and a set that runs
// empty opens its gate.
type readyGate struct {
	mu       sync.Mutex
	required map[string]struct{}
	optional map[string]struct{}

	requiredDone  chan struct{}
	optionalDone  chan struct{}
	closeRequired func()
	closeOptional func()
}

func newReadyGate(required map[string]struct{}, optional map[string]struct{}) *readyGate {
	gate := &readyGate{
		required:     required,
		optional:     optional,
		requiredDone: make(chan struct{}),
		optionalDone: make(chan struct{}),
	}

	gate.closeRequired = sync.OnceFunc(func() { close(gate.requiredDone) })
	gate.closeOptional = sync.OnceFunc(func() { close(gate.optionalDone) })

	// A set that starts out empty has nothing to wait for, e.g. the required
	// set of a WithReadyEvents() with no events at all. No listener exists yet,
	// so this needs no lock.
	gate.openEmptyGates()

	return gate
}

// observe strikes event off both sets and opens the gates that ran empty.
func (g *readyGate) observe(event string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.required, event)
	delete(g.optional, event)

	g.openEmptyGates()
}

// openEmptyGates opens the gate of every set that has nothing left to wait for.
// Every caller but newReadyGate holds g.mu.
func (g *readyGate) openEmptyGates() {
	if len(g.required) == 0 {
		g.closeRequired()
	}

	if len(g.optional) == 0 {
		g.closeOptional()
	}
}

// missingRequired returns the required events that have not arrived, sorted.
func (g *readyGate) missingRequired() []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	return sortedEvents(g.required)
}

// missingOptional returns the best-effort events that have not arrived, sorted.
func (g *readyGate) missingOptional() []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	return sortedEvents(g.optional)
}

func sortedEvents(events map[string]struct{}) []string {
	sorted := make([]string, 0, len(events))
	for event := range events {
		sorted = append(sorted, event)
	}

	slices.Sort(sorted)

	return sorted
}

type ackResponse struct {
	Msg string `json:"msg"`
	OK  bool   `json:"ok"`

	// TokenRequired is how a login for an account with two-factor
	// authentication enabled is answered. Such an ack carries neither an ok
	// nor a message, so it has to be recognized by this field alone.
	TokenRequired bool `json:"tokenRequired"`

	// MsgI18n reports that Msg is a translation key rather than a sentence.
	MsgI18n bool `json:"msgi18n"`

	// URI is the otpauth:// URI a prepare2FA answers with.
	URI string `json:"uri"`

	// Status is whether two-factor authentication is on, from a twoFAStatus.
	Status bool `json:"status"`

	Token           string         `json:"token"`
	ID              int64          `json:"id"`
	MonitorID       int64          `json:"monitorID"`
	MaintenanceID   int64          `json:"maintenanceID"`
	Maintenance     map[string]any `json:"maintenance"`
	Monitors        []any          `json:"monitors"`
	StatusPages     []any          `json:"statusPages"`
	Monitor         map[string]any `json:"monitor"`
	Data            map[string]any `json:"data"`
	Tags            []any          `json:"tags"`
	Tag             map[string]any `json:"tag"`
	Config          map[string]any `json:"config"`
	PublicGroupList []any          `json:"publicGroupList"`
	Incident        map[string]any `json:"incident"`
}

// emitAck emits command and returns the ack the server answered with, whether
// it reports success or not.
//
// Unlike syncEmit it leaves the interpretation of the ack to the caller, which
// is what the login path needs: a login that wants a one-time code is answered
// with neither an ok nor a message, so there is nothing for syncEmit to report.
// The returned error is reserved for an ack that never arrived.
func (c *Client) emitAck(ctx context.Context, command string, args ...any) (ackResponse, error) {
	// Buffered and never closed, so a late ack (e.g. after the context
	// expired) neither leaks the ack goroutine forever nor panics sending on
	// a closed channel.
	res := make(chan ackResponse, 1)

	args = append(args, emit.WithAck(func(response ackResponse) {
		res <- response
	}))

	err := c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	select {
	case response := <-res:
		return response, nil

	case <-ctx.Done():
		return ackResponse{}, fmt.Errorf("%s: %w", command, ctx.Err())
	}
}

func (c *Client) syncEmit(ctx context.Context, command string, args ...any) (ackResponse, error) {
	response, err := c.emitAck(ctx, command, args...)
	if err != nil {
		return ackResponse{}, err
	}

	if !response.OK {
		return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)
	}

	return response, nil
}

// syncEmitWithUpdateEvent emits command and returns once the server has
// acknowledged it and has broadcast an updateEvent, which is what refreshes the
// local state cache. The listener matches the event by name, so the client
// cannot tell which write caused the broadcast it observes.
//
// If the ack reports success but no update event arrives before ctx is done,
// the response is returned together with an error wrapping
// ErrUpdateEventTimeout and the context error: the command was applied, only
// the broadcast is missing. Discarding the response here would tell the caller
// the write failed for a write the server performed.
func (c *Client) syncEmitWithUpdateEvent(
	ctx context.Context,
	command string,
	updateEvent string,
	args ...any,
) (ackResponse, error) {
	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Register listener for notifications updates.
	// Signal done, if update is received and remove listener.
	listenerID := uuid.New()
	c.updates.AddListener(func(_ context.Context, update string) {
		if update == updateEvent {
			closeDone()
		}
	}, listenerID.String())
	defer c.updates.RemoveListener(listenerID.String())

	// Buffered and never closed, for the same reason as in syncEmit: the ack
	// runs on a goroutine of its own, so a send with no receiver left leaks it,
	// and closing the channel underneath it panics the whole process. Checking
	// ctx.Err() in the callback does not help, because the context can expire
	// while the send is already blocked.
	res := make(chan ackResponse, 1)

	args = append(args, emit.WithAck(func(response ackResponse) {
		res <- response
	}))
	err := c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	return awaitAckAndUpdateEvent(ctx, command, done, res)
}

// awaitAckAndUpdateEvent waits for the ack delivered on res and for the update
// event that closes done, and reports which of the two arrived before ctx was
// done.
//
// An ack that reports a failure is returned as the server's error even when ctx
// is done as well, because the rejection is the more useful answer and it is
// what the caller would have received had the ack arrived a moment earlier. An
// update event without an ack stays a plain context error: without the ack the
// client knows neither the outcome of the command nor the ID the server
// assigned.
func awaitAckAndUpdateEvent(
	ctx context.Context,
	command string,
	done <-chan struct{},
	res <-chan ackResponse,
) (ackResponse, error) {
	var (
		response ackResponse
		acked    bool
	)

	// Ensure, we have received both signals: done and ack
	// Setting channel to nil blocks forever, thisway we ensure, that
	// we also receive the second signal.
	for done != nil || res != nil {
		select {
		case <-done:
			done = nil

		case response = <-res:
			if !response.OK {
				return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)
			}

			acked = true
			res = nil

		case <-ctx.Done():
			// A signal that has already been delivered and the context being
			// done can become ready in the same select, which picks between
			// ready cases at random. Both signals are therefore collected
			// explicitly, instead of letting that coin flip decide whether
			// what they report is seen. Receiving from a nil channel is never
			// ready, so the default case covers the signals already taken.
			if !acked {
				select {
				case response = <-res:
					acked = true

				default:
				}
			}

			select {
			case <-done:
				done = nil

			default:
			}

			return resultOnContextDone(ctx, command, response, acked, done == nil)
		}
	}

	return response, nil
}

// resultOnContextDone reports the outcome of a command whose context is done,
// from the signals that are in hand: acked tells whether the ack in response
// arrived, updated whether the update event did.
func resultOnContextDone(
	ctx context.Context,
	command string,
	response ackResponse,
	acked bool,
	updated bool,
) (ackResponse, error) {
	switch {
	case acked && !response.OK:
		return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)

	case acked && updated:
		// Both signals are in hand and the context merely expired while they
		// were collected. Nothing is missing, so this is the same success the
		// caller's loop returns.
		return response, nil

	case acked:
		// The server applied the command, only its broadcast is missing.
		// Returning the response lets the caller keep what the ack carried,
		// e.g. the ID of a created resource.
		return response, &UpdateEventTimeoutError{Command: command, Err: ctx.Err()}

	default:
		return ackResponse{}, fmt.Errorf("%s: %w", command, ctx.Err())
	}
}
