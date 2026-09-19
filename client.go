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
// An error wrapping ErrUpdateEventTimeout also wraps the error that ended the
// wait, so errors.Is(err, context.DeadlineExceeded) still reports an expired
// deadline - the caller's own or the client's, see ErrOperationTimeout - and
// errors.Is(err, context.Canceled) a cancelled context. Use
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

	// Err describes the context that was done before the update event
	// arrived: its cause when the client's own operation timeout expired, see
	// ErrOperationTimeout, and otherwise its Err(). Compare it with errors.Is
	// rather than by equality, which only the latter case satisfies.
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

// ErrOperationTimeout is returned when a command exceeded the client's own
// per-operation budget, see WithOperationTimeout, rather than a deadline the
// caller set. The server did not answer in time: either the ack never arrived,
// or, together with ErrUpdateEventTimeout, the ack arrived but the write never
// became visible in the state cache.
//
// A lost ack says nothing about the server: the command was already on the
// wire when the budget expired, so it may well have been applied. Unlike an
// error wrapping ErrUpdateEventTimeout, this one carries no ID to adopt a
// created resource by, so re-read before retrying a create, and see Resync for
// a state cache that missed a broadcast.
//
// An error wrapping it also wraps context.DeadlineExceeded, so the context
// checks a caller already has keep reporting. Use errors.Is(err,
// ErrOperationTimeout) to tell a client-imposed expiry from the caller's own.
var ErrOperationTimeout = errors.New("operation timeout exceeded")

// operationTimeoutError is the cause recorded on a context derived by
// operationContext, which is its only construction site and always passes a
// positive budget. It names that budget, and wraps context.DeadlineExceeded so
// that replacing the plain context error with it changes what the message
// says, not what errors.Is reports.
//
// It stays unexported because the only thing it carries is the timeout the
// caller configured: ErrOperationTimeout is the whole public contract.
type operationTimeoutError struct {
	timeout time.Duration
}

func (e *operationTimeoutError) Error() string {
	return fmt.Sprintf("%s: %s", ErrOperationTimeout, e.timeout)
}

func (*operationTimeoutError) Unwrap() []error {
	return []error{ErrOperationTimeout, context.DeadlineExceeded}
}

// contextErr describes a done context: the operation timeout recorded by
// operationContext when that is what expired, and otherwise ctx.Err(). Call it
// only on a context that is done, as both call sites do; for a live one it
// returns nil like ctx.Err() does.
//
// It matches the client's own cause rather than any cause, because
// context.Cause reports the cause of the first cancelled ancestor: a caller
// that built its context with context.WithCancelCause would otherwise see its
// own error in place of context.Canceled, and the contract on
// ErrUpdateEventTimeout promises it stays reportable.
//
//nolint:wrapcheck // Both returns are the context's own error, which the callers wrap with the command they belong to.
func contextErr(ctx context.Context) error {
	var operationTimeout *operationTimeoutError
	if errors.As(context.Cause(ctx), &operationTimeout) {
		return operationTimeout
	}

	return ctx.Err()
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

// defaultOperationTimeout is the budget WithOperationTimeout applies when the
// caller configures none. Uptime Kuma answers the commands that only touch its
// own database in milliseconds, so a minute is far beyond any healthy round
// trip and elapses only when an ack or a broadcast is genuinely lost. The
// commands that drive an outbound request instead are exempt rather than
// budgeted, see unboundedCommands. See WithOperationTimeout.
//
// It is a variable so the tests can shorten it and exercise the default that
// New applies, see export_test.go.
//
//nolint:gochecknoglobals // shortened by the tests, constant everywhere else.
var defaultOperationTimeout = 60 * time.Second

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
//
// A Client must be created by New. The zero value is not usable and its methods
// panic.
type Client struct {
	socketioClient               *socketio.Client
	socketioClientConnectTimeout time.Duration
	socketioLogger               socketio.Logger
	autosetup                    bool
	readyEvents                  []string
	readyGracePeriod             time.Duration

	// operationTimeout is set once in New, before any goroutine exists, and so
	// is read without the lock. See WithOperationTimeout.
	operationTimeout time.Duration

	mu *sync.Mutex

	// listWrites serializes the writes to each list-replacing update event,
	// see wholeListUpdateEvents. One semaphore per event name, created in New.
	listWrites map[string]chan struct{}
	updates    signals.Signal[string]
	state      state

	// sessionToken is the JWT the server handed out at login. It is what
	// Resync logs in with, see there.
	sessionToken string

	// sessionTokenPreset is the token WithSessionToken configured, which New
	// tries before the password login.
	sessionTokenPreset string

	// sessionTokenRejected records that the server refused sessionTokenPreset
	// and the password login took over, see Client.SessionTokenRejected.
	sessionTokenRejected bool

	// sessionTokenRejectedCallback is told about that fallback as it happens,
	// see WithSessionTokenRejectedCallback.
	sessionTokenRejectedCallback func(err error)

	// autoLoggedIn records that the server has authentication disabled and
	// logged the client in itself. Credentials given anyway are still used, so
	// this says nothing about whether there is a session token; what it is for
	// is telling a client that has none why, see Resync.
	autoLoggedIn bool

	// totpCode produces the one-time code for an account with two-factor
	// authentication enabled, and is nil for a client configured with neither
	// WithTOTPSecret nor WithTOTPCode. totpSources counts how many of the two
	// options set one, so New can reject a caller that used both, and totpErr
	// carries an option failure for New to report. Two failing options would
	// leave only the later error, but New rejects the two-source case before it
	// looks at totpErr, so at most one option's failure is ever reported.
	totpCode    func(ctx context.Context) (string, error)
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
// It covers the login as well, including the short wait for the server to
// state how it wants to be authenticated. A server that states nothing is
// waited for with no more than half of the remaining budget, so the login
// still has time to run.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.socketioClientConnectTimeout = timeout
	}
}

// WithOperationTimeout bounds each individual command: the wait for the
// server's ack, and the wait for the update event that confirms the command
// reached the local state cache. Resync, which is a command and the lists it
// resends, is bounded as a whole. It defaults to one minute; a value of zero
// or less disables it, leaving every command bounded only by the caller's
// context.
//
// It exists because an ack the server never sends, or an update event lost on
// the way, would otherwise block the caller for as long as its context lives.
// Callers whose context has no deadline - a Terraform provider is handed one
// per RPC and Terraform sets none - would block forever.
//
// The bound is on the round trip, not on the call: a write that has to queue
// behind the other writes to the list it broadcasts, see wholeListUpdateEvents,
// spends no budget while it waits its turn, and is bounded there by the
// caller's context alone. A caller that wants a bound on the whole call sets
// its own deadline, which the queue honours.
//
// TestNotification and TestDockerHost are exempt, see unboundedCommands.
//
// WithConnectTimeout bounds establishing the connection as a whole; the login
// and setup it covers are commands, so this budget bounds them too. Either way
// a caller that passes its own deadline keeps it: the effective bound is
// whichever expires first.
func WithOperationTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.operationTimeout = timeout
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
// username and password have to be set together or both be empty. Passing
// neither is the way to connect to a server with authentication disabled, and
// to one that WithSessionToken authenticates instead; a server that does want a
// login answers a client with nothing to offer with ErrAuthRequired.
//
// The login it performs reports what the server refused through the package's
// sentinel errors, so a caller can tell the cases apart with errors.Is:
// ErrAuthRequired, ErrInvalidCredentials, ErrTwoFactorRequired,
// ErrInvalidTOTPCode, ErrInvalidSessionToken, ErrUserInactive and
// ErrRateLimited. Those are the rejections an Uptime Kuma 2.x server states;
// 1.x is not supported, see the package documentation.
//
//nolint:revive // Complexity is necessary for complete client initialization and event setup
func New(ctx context.Context, baseURL string, username string, password string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger:   &utils.DefaultLogger{Level: utils.NONE},
		readyEvents:      defaultReadyEvents,
		readyGracePeriod: defaultReadyGracePeriod,
		operationTimeout: defaultOperationTimeout,

		mu:         &sync.Mutex{},
		updates:    signals.New[string](),
		listWrites: newListWriteLocks(),
	}

	for _, opt := range opts {
		opt(c)
	}

	unknown := unknownReadyEvents(c.readyEvents)
	if len(unknown) > 0 {
		return nil, fmt.Errorf("ready events: unknown: %s", strings.Join(unknown, ", "))
	}

	creds, err := newCredentials(username, password, c.sessionTokenPreset)
	if err != nil {
		return nil, err
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

	client.On("setup", func() {
		closeSetupRequired()
	})

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

	err = c.authenticate(ctxWithConnectTimeout, creds, barrier)
	if err != nil {
		// A server that is not set up yet rejects the credentials it does not
		// have a user for; running the setup is what makes them work, so that
		// is not a failure when the server asked for one.
		if !errors.Is(err, ErrInvalidCredentials) ||
			!setupPending(ctxWithConnectTimeout, setupRequired) {
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
// so Resync logs in again with the session token the client holds - the one the
// login New performed handed out, or the one WithSessionToken supplied and the
// server accepted. A client New returned without logging it in never obtained
// one and cannot resync: that is a client created without credentials, against
// a server that either has authentication disabled or never stated that it
// wants a login.
//
// A token the server has since stopped accepting is reported as
// ErrInvalidSessionToken, wrapped by ErrUserInactive when the account it names
// was deactivated or deleted. Neither is recoverable here: a new New is what
// obtains a fresh token.
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
					"and offers no command to resend the lists, pass credentials to New " +
					"for a session token",
			)
		}

		return errors.New(
			"resync: no session token, the client was created without credentials",
		)
	}

	// Resync is one command and the lists it resends, so the operation budget
	// covers both, see WithOperationTimeout. emitAck derives one of its own,
	// which bounds the ack alone; a server that answers the token and then
	// never resends a list would otherwise block a caller that set no deadline
	// forever - in the very call ErrUpdateEventTimeout points it at.
	ctx, cancel := c.operationContext(ctx, "resync")
	defer cancel()

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
			contextErr(ctx),
			strings.Join(gate.missingRequired(), ", "),
		)
	}

	// The required lists are in. Give the best-effort ones the same grace
	// window New grants them, so the cache they feed is refreshed too.
	c.awaitOptionalReadyEvents(ctx, "Resync", gate, nil)

	return nil
}

// SessionToken returns the session token the client is authenticated with: the
// one the server handed out for the login New performed, or the one
// WithSessionToken supplied and the server accepted. It is the empty string for
// a client New returned without logging it in, which is a client created
// without credentials, against a server that either has authentication disabled
// or never stated that it wants a login.
//
// It is the credential WithSessionToken takes, so a caller can persist it and
// reconnect later without the password and without a one-time code. It is a
// bearer credential that does not expire, see WithSessionToken.
func (c *Client) SessionToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.sessionToken
}

// SessionTokenRejected reports whether the server refused the token
// WithSessionToken configured and New logged in with the username and password
// instead.
//
// It is what tells a caller that its stored token is dead, because nothing else
// does: the login succeeds either way, and SessionToken then returns the fresh
// token of the password login, which the caller cannot tell from the one it
// presented. A caller that persists tokens checks this and writes the new one
// back, see WithSessionToken.
func (c *Client) SessionTokenRejected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.sessionTokenRejected
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

// setSessionTokenRejected records that the token WithSessionToken configured
// was refused and the password login recovered, see
// Client.SessionTokenRejected.
func (c *Client) setSessionTokenRejected() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessionTokenRejected = true
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

// unboundedCommands names the commands the operation budget does not apply to:
// testNotification, behind TestNotification, and testDockerHost, behind
// TestDockerHost.
//
// Every other command only touches Uptime Kuma's own database and is answered
// in milliseconds. These two instead wait on the server while it makes a
// request of its own - to a notification provider, to a Docker daemon - and an
// unreachable one takes as long as that provider's own timeouts, which can run
// well past a minute. The answer is the point of asking, so the client does not
// cut it short; the caller's context is the bound.
func unboundedCommands() []string {
	return []string{
		"testNotification",
		"testDockerHost",
	}
}

// operationContext derives the context command runs under, see
// WithOperationTimeout. It returns ctx unchanged when the budget is disabled or
// command is exempt from it, and context.WithTimeoutCause already keeps
// whichever of the two deadlines comes first, so a caller's own deadline is
// never extended.
func (c *Client) operationContext(ctx context.Context, command string) (context.Context, context.CancelFunc) {
	if c.operationTimeout <= 0 || slices.Contains(unboundedCommands(), command) {
		return ctx, func() {}
	}

	return context.WithTimeoutCause(ctx, c.operationTimeout, &operationTimeoutError{timeout: c.operationTimeout})
}

// ackOptions returns the emit options that deliver the ack on res, together
// with the expiry that reclaims the registration the socket.io client keeps
// for it. That registration is otherwise dropped only by the ack itself, so a
// command abandoned on a done context - the routine outcome now that an
// operation timeout bounds one, see WithOperationTimeout - would retain it,
// and the closure it holds, for the life of the socket. The expiry needs no
// callback of its own: the waiting select reports what happened.
//
// A context with no deadline gets no expiry, because there is no duration to
// give one; that is the WithOperationTimeout(0) opt-out, where a command is
// abandoned only if the caller abandons it.
func ackOptions(ctx context.Context, res chan<- ackResponse) []any {
	options := []any{
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	}

	if deadline, ok := ctx.Deadline(); ok {
		options = append(options, emit.WithTimeout(time.Until(deadline), nil))
	}

	return options
}

// emitAck emits command and returns the ack the server answered with, whether
// it reports success or not.
//
// Unlike syncEmit it leaves the interpretation of the ack to the caller, which
// is what the login path needs: a login that wants a one-time code is answered
// with neither an ok nor a message, so there is nothing for syncEmit to report.
// The returned error is reserved for an ack that never arrived.
func (c *Client) emitAck(ctx context.Context, command string, args ...any) (ackResponse, error) {
	ctx, cancel := c.operationContext(ctx, command)
	defer cancel()

	// Buffered and never closed, so a late ack (e.g. after the context
	// expired) neither leaks the ack goroutine forever nor panics sending on
	// a closed channel.
	res := make(chan ackResponse, 1)

	args = append(args, ackOptions(ctx, res)...)

	err := c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	select {
	case response := <-res:
		return response, nil

	case <-ctx.Done():
		return ackResponse{}, fmt.Errorf("%s: %w", command, contextErr(ctx))
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

// wholeListUpdateEvents names the update events whose payload is a complete
// list that replaces the cached one.
//
// Such a broadcast is only as fresh as the moment the server read the list, and
// Uptime Kuma reads it after the write it belongs to and emits it afterwards,
// over a connection that preserves that order. Two writes in flight at once
// therefore interleave into a read-without, emit-after order, and the older
// snapshot lands last and drops the newer write from the cache again. Writes to
// these lists are serialized so that no two of them can interleave that way.
//
// updateMonitorIntoList and deleteMonitorFromList are absent on purpose: they
// carry a single monitor and are merged into the cache rather than replacing
// it, so they cannot go stale this way and stay concurrent.
func wholeListUpdateEvents() []string {
	return []string{
		"notificationList",
		"proxyList",
		"dockerHostList",
		"maintenanceList",
		"monitorList",
		"statusPageList",
	}
}

// newListWriteLocks builds the write semaphore per list-replacing update
// event. A buffered channel with room for one token rather than a mutex,
// because a caller waiting its turn has to be able to give up, see
// lockListWrite.
func newListWriteLocks() map[string]chan struct{} {
	locks := make(map[string]chan struct{}, len(wholeListUpdateEvents()))
	for _, event := range wholeListUpdateEvents() {
		locks[event] = make(chan struct{}, 1)
	}

	return locks
}

// lockListWrite takes the write lock of updateEvent and returns the release. It
// returns a nil release for an event that carries no whole list and needs no
// serializing, and an error for a caller whose context was done before its turn
// came.
//
// The wait is on ctx and not on the operation budget: the budget bounds a
// command's own round trip and is taken once the lock is won, see
// syncEmitWithConfirmedUpdateEvent, so what bounds the queue is whatever bounds
// the caller. Giving up here is the only way out of a queue behind writes that
// are each waiting out a server that stopped answering.
func (c *Client) lockListWrite(ctx context.Context, updateEvent string) (func(), error) {
	lock, ok := c.listWrites[updateEvent]
	if !ok {
		return nil, nil //nolint:nilnil // no lock to take is not a failure.
	}

	select {
	case lock <- struct{}{}:
		return func() { <-lock }, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("waiting for the other writes to %s: %w", updateEvent, contextErr(ctx))
	}
}

// cachedByID reports whether list holds an entry with the given ID. It is what
// the create and delete confirmations are built from: a create is settled once
// the ID the server assigned is in the cache, a delete once it is gone.
func cachedByID[T interface{ GetID() int64 }](list []T, id int64) bool {
	for i := range list {
		if list[i].GetID() == id {
			return true
		}
	}

	return false
}

// updateEventConfirm reports whether the local state cache already reflects the
// write the server acknowledged with response.
//
// It exists because the update events are broadcasts with no reference to the
// write that caused them: they carry a whole list, and the listener sees only
// its name. With several writes in flight on one connection - a Terraform apply
// runs its resource operations concurrently over a single client - the first
// broadcast a waiting write observes may well be the one a neighbouring write
// triggered, taken from a list that predates its own change. Accepting it
// leaves the cache one broadcast behind, and the cache-served getters then
// report a resource the server has as missing.
//
// A confirmation closes that gap by looking at what the cache actually holds,
// so waiting continues until the write is visible there.
type updateEventConfirm func(response ackResponse) bool

// syncEmitWithUpdateEvent emits command and returns once the server has
// acknowledged it and has broadcast updateEvent, which is what refreshes the
// local state cache.
//
// The broadcast is matched by name alone, so the caller cannot tell which write
// caused the one it observes. Callers that can recognize their own write pass a
// confirmation to syncEmitWithConfirmedUpdateEvent instead.
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
	return c.syncEmitWithConfirmedUpdateEvent(ctx, command, updateEvent, nil, args...)
}

// syncEmitWithConfirmedUpdateEvent is syncEmitWithUpdateEvent for a caller that
// can recognize its own write in the state cache.
//
// Every broadcast of updateEvent makes it re-evaluate confirm, and it returns
// once the ack has arrived and confirm reports the write is visible. A confirm
// that never does is reported like a broadcast that never arrives, as
// ErrUpdateEventTimeout: the server applied the command either way, and Resync
// is what rebuilds the cache.
//
// A nil confirm keeps the plain behaviour of returning on the first broadcast.
func (c *Client) syncEmitWithConfirmedUpdateEvent(
	ctx context.Context,
	command string,
	updateEvent string,
	confirm updateEventConfirm,
	args ...any,
) (ackResponse, error) {
	// Held for the whole command, so that the server never has two reads of
	// this list in flight and cannot emit them out of order.
	unlock, err := c.lockListWrite(ctx, updateEvent)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	if unlock != nil {
		defer unlock()
	}

	// Derived after the lock, not before it: the timeout bounds this command's
	// own round trip, and queueing behind the other writes to this list is not
	// part of it. Taking it earlier would let a busy list spend the budget
	// before the command is even sent.
	ctx, cancel := c.operationContext(ctx, command)
	defer cancel()

	// Buffered with room for one token and written to without blocking: the
	// listener runs on the goroutine that applied the state and must never be
	// held up, and a token that is already pending says everything a second
	// one would. The waiter re-reads the cache after every token it takes, so
	// a broadcast that is coalesced away cannot hide a state change.
	updated := make(chan struct{}, 1)

	listenerID := uuid.New()
	c.updates.AddListener(func(_ context.Context, update string) {
		if update != updateEvent {
			return
		}

		select {
		case updated <- struct{}{}:
		default:
		}
	}, listenerID.String())
	defer c.updates.RemoveListener(listenerID.String())

	// Buffered and never closed, for the same reason as in syncEmit: the ack
	// runs on a goroutine of its own, so a send with no receiver left leaks it,
	// and closing the channel underneath it panics the whole process. Checking
	// ctx.Err() in the callback does not help, because the context can expire
	// while the send is already blocked.
	res := make(chan ackResponse, 1)

	args = append(args, ackOptions(ctx, res)...)

	err = c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	return awaitAckAndUpdateEvent(ctx, command, updated, res, confirm)
}

// awaitAckAndUpdateEvent waits for the ack delivered on res and for the update
// event that confirm accepts, or, when confirm is nil, for the first update
// event signalled on updated.
//
// An ack that reports a failure is returned as the server's error even when ctx
// is done as well, because the rejection is the more useful answer and it is
// what the caller would have received had the ack arrived a moment earlier. An
// update event without an ack stays a plain context error: without the ack the
// client knows neither the outcome of the command nor the ID the server
// assigned, and a confirmation cannot run without the ID either.
func awaitAckAndUpdateEvent(
	ctx context.Context,
	command string,
	updated <-chan struct{},
	res <-chan ackResponse,
	confirm updateEventConfirm,
) (ackResponse, error) {
	var (
		response  ackResponse
		acked     bool
		sawUpdate bool
	)

	for {
		if acked && confirmUpdate(confirm, response, sawUpdate) {
			return response, nil
		}

		select {
		case <-updated:
			sawUpdate = true

		case response = <-res:
			if !response.OK {
				return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)
			}

			acked = true
			// Nil blocks forever, so the loop keeps waiting on the remaining
			// signals instead of spinning on a channel that is done.
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
			case <-updated:
				sawUpdate = true

			default:
			}

			return resultOnContextDone(
				ctx, command, response, acked, acked && confirmUpdate(confirm, response, sawUpdate),
			)
		}
	}
}

// confirmUpdate reports whether the write acknowledged by response is settled:
// visible in the state cache when the caller supplied a confirmation, and
// otherwise simply broadcast at least once.
func confirmUpdate(confirm updateEventConfirm, response ackResponse, sawUpdate bool) bool {
	if confirm == nil {
		return sawUpdate
	}

	return confirm(response)
}

// resultOnContextDone reports the outcome of a command whose context is done,
// from the signals that are in hand: acked tells whether the ack in response
// arrived, settled whether the write is visible in the state cache - the
// update event for a command that waits for one by name, the confirmation for
// a command that recognizes its own write.
func resultOnContextDone(
	ctx context.Context,
	command string,
	response ackResponse,
	acked bool,
	settled bool,
) (ackResponse, error) {
	switch {
	case acked && !response.OK:
		return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)

	case acked && settled:
		// Both signals are in hand and the context merely expired while they
		// were collected. Nothing is missing, so this is the same success the
		// caller's loop returns.
		return response, nil

	case acked:
		// The server applied the command, only the state cache has not caught
		// up with it. Returning the response lets the caller keep what the ack
		// carried, e.g. the ID of a created resource.
		return response, &UpdateEventTimeoutError{Command: command, Err: contextErr(ctx)}

	default:
		return ackResponse{}, fmt.Errorf("%s: %w", command, contextErr(ctx))
	}
}
