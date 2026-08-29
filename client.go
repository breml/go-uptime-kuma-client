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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"
	"golang.org/x/sync/errgroup"

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

// pendingListEvents returns the update events that carry the lists the local
// state cache is built from, as a set to strike them off as they arrive. The
// server sends all of them after a login, which is what New waits for and what
// Resync repeats.
func pendingListEvents() map[string]struct{} {
	return map[string]struct{}{
		"monitorList":      empty,
		"maintenanceList":  empty,
		"notificationList": empty,
		"statusPageList":   empty,
		"proxyList":        empty,
		"dockerHostList":   empty,
		"apiKeyList":       empty,
	}
}

// Client represents a connection to an Uptime Kuma server.
type Client struct {
	socketioClient               *socketio.Client
	socketioClientConnectTimeout time.Duration
	socketioLogger               socketio.Logger
	autosetup                    bool

	mu      *sync.Mutex
	updates signals.Signal[string]
	state   state

	// sessionToken is the JWT the server handed out at login. It is what
	// Resync logs in with, see there.
	sessionToken string
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
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},

		mu:      &sync.Mutex{},
		updates: signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
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

	updateSeenMu := sync.Mutex{}
	updateSeenMu.Lock()
	updateSeen := pendingListEvents()
	updateSeenMu.Unlock()

	ready := make(chan struct{})
	closeReady := sync.OnceFunc(func() {
		close(ready)
	})
	defer closeReady()

	c.updates.AddListener(func(_ context.Context, s string) {
		updateSeenMu.Lock()
		defer updateSeenMu.Unlock()

		delete(updateSeen, s)

		if len(updateSeen) == 0 {
			closeReady()
		}
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

	client.OnAny(func(s string, _ []any) {
		if s != "notificationList" && s != "monitorList" && s != "statusPageList" && s != "maintenanceList" &&
			s != "proxyList" &&
			s != "dockerHostList" {
			c.updates.Emit(context.Background(), s)
		}
	})

	errgrp := errgroup.Group{}
	errgrp.Go(func() error {
		return client.Connect(ctx)
	})

	select {
	case <-connect:
	case <-ctx.Done():
		return nil, fmt.Errorf("connect to server: %w", ctx.Err())

	case <-ctxWithConnectTimeout.Done():
		return nil, fmt.Errorf("connect to server: %w", ctxWithConnectTimeout.Err())
	}

	err = errgrp.Wait()
	if err != nil {
		return nil, fmt.Errorf("connect to server: %w", err)
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

	if username != "" && password != "" {
		var loginResponse ackResponse

		loginResponse, err = c.syncEmit(
			ctxWithConnectTimeout,
			"login",
			map[string]any{"username": username, "password": password, "token": ""},
		)
		if err == nil {
			c.setSessionToken(loginResponse.Token)
		}

		if err != nil {
			// Ensure we had the time to receive a potential setup event.
			time.Sleep(10 * time.Millisecond)

			wantSetup := false
			select {
			case <-setupRequired:
				wantSetup = true

			default:
			}

			if (!strings.Contains(err.Error(), "Incorrect username or password") && !strings.Contains(err.Error(), "authIncorrectCreds")) ||
				!wantSetup {
				return nil, fmt.Errorf("login: %w", err)
			}
		}
	}

	for {
		select {
		case <-ready:
			closeOnErr = false
			return c, nil

		case <-setupRequired:
			setupRequired = nil

			if !c.autosetup {
				return nil, errors.New("server does require setup, but autosetup is disabled")
			}

			_, err := c.syncEmit(ctxWithConnectTimeout, "setup", username, password)
			if err != nil {
				return nil, fmt.Errorf("setup: %w", err)
			}

			loginResponse, err := c.syncEmit(
				ctxWithConnectTimeout,
				"login",
				map[string]any{"username": username, "password": password, "token": ""},
			)
			if err != nil {
				return nil, fmt.Errorf("login: %w", err)
			}

			c.setSessionToken(loginResponse.Token)

		case <-ctx.Done():
			return nil, fmt.Errorf("wait for ready: %w", ctx.Err())

		case <-connectTimeoutDone:
			// ctxWithConnectTimeout is derived from ctx, so its Done channel
			// closes whenever the parent ctx is cancelled too. Prefer the
			// parent's error in that case to avoid a misleading
			// "missing events" message on an ordinary cancellation.
			if ctx.Err() != nil {
				return nil, fmt.Errorf("wait for ready: %w", ctx.Err())
			}

			// If all ready events arrived at the exact same instant as the
			// timeout, prefer the success path over the error path.
			select {
			case <-ready:
				closeOnErr = false
				return c, nil

			default:
			}

			updateSeenMu.Lock()
			missing := make([]string, 0, len(updateSeen))
			for event := range updateSeen {
				missing = append(missing, event)
			}
			updateSeenMu.Unlock()

			sort.Strings(missing)

			return nil, fmt.Errorf(
				"wait for ready: %w (missing events: %s)",
				ctxWithConnectTimeout.Err(),
				strings.Join(missing, ", "),
			)
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
// server has resent every list the cache is built from.
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
func (c *Client) Resync(ctx context.Context) error {
	c.mu.Lock()
	token := c.sessionToken
	c.mu.Unlock()

	if token == "" {
		return errors.New("resync: no session token, the client was created without credentials")
	}

	pendingMu := sync.Mutex{}
	pending := pendingListEvents()

	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Registered before the command is emitted, so that no list can arrive
	// unnoticed between the two.
	listenerID := uuid.New()
	c.updates.AddListener(func(_ context.Context, update string) {
		pendingMu.Lock()
		defer pendingMu.Unlock()

		delete(pending, update)

		if len(pending) == 0 {
			closeDone()
		}
	}, listenerID.String())

	defer c.updates.RemoveListener(listenerID.String())

	_, err := c.syncEmit(ctx, "loginByToken", token)
	if err != nil {
		return fmt.Errorf("resync: %w", err)
	}

	select {
	case <-done:
		return nil

	case <-ctx.Done():
		pendingMu.Lock()
		missing := make([]string, 0, len(pending))

		for event := range pending {
			missing = append(missing, event)
		}
		pendingMu.Unlock()

		slices.Sort(missing)

		return fmt.Errorf("resync: %w (missing events: %s)", ctx.Err(), strings.Join(missing, ", "))
	}
}

// setSessionToken records the JWT a login answered with, see Resync.
func (c *Client) setSessionToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessionToken = token
}

type ackResponse struct {
	Msg             string         `json:"msg"`
	OK              bool           `json:"ok"`
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

func (c *Client) syncEmit(ctx context.Context, command string, args ...any) (ackResponse, error) {
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
		if !response.OK {
			return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)
		}

		return response, nil

	case <-ctx.Done():
		return ackResponse{}, fmt.Errorf("%s: %w", command, ctx.Err())
	}
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
