package kuma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
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

var ErrNotFound = fmt.Errorf("not found")
var ErrTwoFactorRequired = fmt.Errorf("two-factor authentication required but no TOTP secret provided")

const (
	LogLevelDebug = utils.DEBUG
	LogLevelInfo  = utils.INFO
	LogLevelWarn  = utils.WARN
	LogLevelError = utils.ERROR
	LogLevelNone  = utils.NONE
)

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

type Client struct {
	socketioClient               *socketio.Client
	socketioClientConnectTimeout time.Duration
	socketioLogger               socketio.Logger
	autosetup                    bool
	autosetupEnable2FA           bool

	mu      *sync.Mutex
	updates signals.Signal[string]
	state   state

	// Authentication fields
	authToken  string // Current session token
	has2FA     bool   // Whether user has 2FA enabled
	totpSecret string // TOTP secret for automatic token generation
}

type Option func(c *Client)

func WithAutosetup(enable2FA bool) Option {
	return func(c *Client) {
		c.autosetup = true
		c.autosetupEnable2FA = enable2FA
	}
}

func WithTOTPSecret(secret string) Option {
	return func(c *Client) {
		c.totpSecret = secret
	}
}

func WithLogLevel(level int) Option {
	return func(c *Client) {
		if level >= utils.DEBUG && level <= utils.NONE {
			c.socketioLogger = &utils.DefaultLogger{Level: level}
		}
	}
}

func WithConnectTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.socketioClientConnectTimeout = timeout
	}
}

// setupDatabase handles the database setup phase for Uptime Kuma v2.
// It checks if database setup is needed and configures SQLite if required.
// The function will wait for the server to restart after database configuration.
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
		return ctx.Err()
	default:
	}

	// Check entry-page without retry - let the caller (pool.Retry) handle retries
	// Use a longer timeout for the HTTP request itself, independent of parent context
	httpCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, "GET", entryPageURL, nil)
	if err != nil {
		return fmt.Errorf("create entry-page request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Return connection errors as-is so caller can retry
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("entry-page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read entry-page response: %w", err)
	}

	if err := json.Unmarshal(body, &entryPage); err != nil {
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

	req, err = http.NewRequestWithContext(httpCtx, "POST", setupDBURL, bytes.NewReader(reqBody))
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
	if err := json.Unmarshal(body, &setupResp); err != nil {
		return fmt.Errorf("parse setup-database response: %w", err)
	}

	if !setupResp.OK {
		return fmt.Errorf("setup-database failed")
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
			return fmt.Errorf("timeout waiting for server restart")
		case <-ticker.C:
			// Use a short timeout for each poll attempt
			pollCtx, pollCancel := context.WithTimeout(context.Background(), 2*time.Second)
			req, err := http.NewRequestWithContext(pollCtx, "GET", entryPageURL, nil)
			if err != nil {
				pollCancel()
				continue
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				pollCancel()
				continue
			}

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			pollCancel()
			if err != nil {
				continue
			}

			var checkEntryPage entryPageResponse
			if err := json.Unmarshal(body, &checkEntryPage); err != nil {
				continue
			}

			// If entry page type changed from "setup-database", server has restarted
			if checkEntryPage.Type != "setup-database" {
				return nil
			}
		}
	}
}

// setupEventHandlers configures all Socket.IO event handlers for the client.
func (c *Client) setupEventHandlers(ctx context.Context) error {
	c.socketioClient.On("notificationList", func(notificationList []notification.Base) {
		c.mu.Lock()
		c.state.notifications = notificationList
		defer c.mu.Unlock()

		c.updates.Emit(context.Background(), "notificationList")
	})

	c.socketioClient.On("monitorList", func(monitorMap map[string]monitor.Base) {
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
	c.socketioClient.On("updateMonitorIntoList", func(monitorMap map[string]monitor.Base) {
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
	c.socketioClient.On("deleteMonitorFromList", func(monitorID int64) {
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

	c.socketioClient.On("statusPageList", func(statusPageMap map[int64]statuspage.StatusPage) {
		c.mu.Lock()
		c.state.statusPages = statusPageMap
		defer c.mu.Unlock()

		c.updates.Emit(context.Background(), "statusPageList")
	})

	c.socketioClient.On("maintenanceList", func(maintenanceMap map[string]maintenance.Maintenance) {
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

	c.socketioClient.On("proxyList", func(proxyList []proxy.Proxy) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.state.proxies = proxyList

		c.updates.Emit(context.Background(), "proxyList")
	})

	c.socketioClient.On("dockerHostList", func(dockerHostList []dockerhost.DockerHost) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.state.dockerHosts = dockerHostList

		c.updates.Emit(context.Background(), "dockerHostList")
	})

	c.socketioClient.OnAny(func(s string, i []any) {
		if s != "notificationList" && s != "monitorList" && s != "statusPageList" && s != "maintenanceList" && s != "proxyList" && s != "dockerHostList" {
			c.updates.Emit(context.Background(), s)
		}
	})

	return nil
}

// setupConnectHandler sets up the connect event handler and returns a channel that will be closed when connected.
func (c *Client) setupConnectHandler() chan struct{} {
	connect := make(chan struct{})
	closeConnect := sync.OnceFunc(func() {
		close(connect)
	})

	c.socketioClient.On("connect", func() {
		closeConnect()
	})

	return connect
}

// waitForConnect waits for the Socket.IO connection to be established.
func (c *Client) waitForConnect(ctx context.Context, connect chan struct{}) error {
	select {
	case <-connect:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("connect to server: %v", ctx.Err())
	}
}

// waitForReady waits for all initial data to be loaded.
func (c *Client) waitForReady(ctx context.Context) error {
	updateSeenMu := sync.Mutex{}
	updateSeenMu.Lock()
	updateSeen := map[string]struct{}{
		"monitorList":      empty,
		"maintenanceList":  empty,
		"notificationList": empty,
		"statusPageList":   empty,
		"proxyList":        empty,
		"dockerHostList":   empty,
		"apiKeyList":       empty,
	}
	updateSeenMu.Unlock()

	ready := make(chan struct{})
	closeReady := sync.OnceFunc(func() {
		close(ready)
	})
	defer closeReady()

	c.updates.AddListener(func(ctx context.Context, s string) {
		updateSeenMu.Lock()
		defer updateSeenMu.Unlock()

		delete(updateSeen, s)

		if len(updateSeen) == 0 {
			closeReady()
		}
	}, "connect-ready")
	defer c.updates.RemoveListener("connect-ready")

	select {
	case <-ready:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for ready: %v", ctx.Err())
	}
}

// generateTOTPToken generates a TOTP token from the given secret.
// The secret should be base32-encoded (as provided by Uptime Kuma).
func generateTOTPToken(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}

func New(ctx context.Context, baseURL string, username string, password string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},
		mu:             &sync.Mutex{},
		updates:        signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
	}

	ctxWithConnectTimeout := ctx
	if c.socketioClientConnectTimeout != 0 {
		var cancel func()
		ctxWithConnectTimeout, cancel = context.WithTimeout(ctx, c.socketioClientConnectTimeout)
		defer cancel()
	}

	// Handle database setup for Uptime Kuma v2 if autosetup is enabled
	if c.autosetup {
		if err := setupDatabase(ctxWithConnectTimeout, baseURL); err != nil {
			return nil, fmt.Errorf("database setup: %w", err)
		}
	}

	client, err := socketio.NewClient(
		socketio.WithRawURL(baseURL),
		socketio.WithLogger(c.socketioLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %v", err)
	}

	c.socketioClient = client

	if err := c.setupEventHandlers(ctx); err != nil {
		return nil, err
	}

	connect := c.setupConnectHandler()

	setupRequired := make(chan struct{})
	closeSetupRequired := sync.OnceFunc(func() {
		close(setupRequired)
	})
	defer closeSetupRequired()

	if c.autosetup {
		c.socketioClient.On("setup", func() {
			closeSetupRequired()
		})
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %v", err)
	}

	if err := c.waitForConnect(ctx, connect); err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		// Generate TOTP token if secret is provided
		token := ""
		if c.totpSecret != "" {
			token, err = generateTOTPToken(c.totpSecret)
			if err != nil {
				return nil, fmt.Errorf("generate TOTP token: %w", err)
			}
		}

		resp, err := c.syncEmit(ctxWithConnectTimeout, "login", map[string]any{"username": username, "password": password, "token": token})
		if err != nil {
			// Ensure we had the time to receive a potential setup event.
			time.Sleep(10 * time.Millisecond)

			wantSetup := false
			select {
			case <-setupRequired:
				wantSetup = true
			default:
			}

			if (!strings.Contains(err.Error(), "Incorrect username or password") && !strings.Contains(err.Error(), "authIncorrectCreds")) || !wantSetup {
				return nil, fmt.Errorf("login: %v", err)
			}
		}

		// Check if 2FA is required but no TOTP secret was provided
		if resp.TokenRequired && c.totpSecret == "" {
			return nil, ErrTwoFactorRequired
		}
	}

	// Handle setup if required
	select {
	case <-setupRequired:
		if !c.autosetup {
			return nil, fmt.Errorf("server does require setup, but autosetup is disabled")
		}

		_, err := c.syncEmit(ctxWithConnectTimeout, "setup", username, password)
		if err != nil {
			return nil, fmt.Errorf("setup: %v", err)
		}

		// If 2FA should be enabled during autosetup, set it up now
		if c.autosetupEnable2FA {
			if err := c.setup2FA(ctxWithConnectTimeout, password); err != nil {
				return nil, fmt.Errorf("setup 2FA: %w", err)
			}
		}

		// Generate TOTP token for login if secret is available (set by setup2FA)
		token := ""
		if c.totpSecret != "" {
			token, err = generateTOTPToken(c.totpSecret)
			if err != nil {
				return nil, fmt.Errorf("generate TOTP token: %w", err)
			}
		}

		_, err = c.syncEmit(ctxWithConnectTimeout, "login", map[string]any{"username": username, "password": password, "token": token})
		if err != nil {
			return nil, fmt.Errorf("login: %v", err)
		}
	default:
	}

	if err := c.waitForReady(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Disconnect() error {
	return c.socketioClient.Close()
}

// prepare2FA initiates 2FA setup and returns the secret and URI.
// Internal method used during autosetup.
func (c *Client) prepare2FA(ctx context.Context, password string) (secret string, uri string, err error) {
	resp, err := c.syncEmit(ctx, "prepare2FA", password)
	if err != nil {
		return "", "", err
	}

	uriVal, ok := resp.Data["uri"].(string)
	if !ok {
		return "", "", fmt.Errorf("prepare2FA: invalid response, missing uri")
	}

	return "", uriVal, nil
}

// save2FA activates 2FA for the current user.
// Internal method used during autosetup.
func (c *Client) save2FA(ctx context.Context, password string) error {
	_, err := c.syncEmit(ctx, "save2FA", password)
	return err
}

// twoFAStatus checks if 2FA is enabled for the current user.
// Internal method used during autosetup.
func (c *Client) twoFAStatus(ctx context.Context) (bool, error) {
	resp, err := c.syncEmit(ctx, "twoFAStatus")
	if err != nil {
		return false, err
	}

	status, ok := resp.Data["status"].(bool)
	if !ok {
		return false, fmt.Errorf("twoFAStatus: invalid response")
	}

	return status, nil
}

// setup2FA handles the complete 2FA setup flow during autosetup.
// It calls prepare2FA, extracts the secret, calls save2FA, and stores the secret.
func (c *Client) setup2FA(ctx context.Context, password string) error {
	// Prepare 2FA - this generates the secret
	_, uri, err := c.prepare2FA(ctx, password)
	if err != nil {
		return fmt.Errorf("prepare 2FA: %w", err)
	}

	// Extract secret from URI
	// Format: otpauth://totp/Uptime%20Kuma:username?secret=BASE32SECRET
	secret, err := extractSecretFromURI(uri)
	if err != nil {
		return fmt.Errorf("extract secret from URI: %w", err)
	}

	// Save 2FA - this activates it
	if err := c.save2FA(ctx, password); err != nil {
		return fmt.Errorf("save 2FA: %w", err)
	}

	// Store the secret for future use
	c.totpSecret = secret

	// Verify 2FA is now enabled
	status, err := c.twoFAStatus(ctx)
	if err != nil {
		return fmt.Errorf("verify 2FA status: %w", err)
	}
	if !status {
		return fmt.Errorf("2FA setup completed but status check failed")
	}

	return nil
}

// extractSecretFromURI extracts the base32-encoded secret from an otpauth:// URI.
func extractSecretFromURI(uri string) (string, error) {
	// URI format: otpauth://totp/Uptime%20Kuma:username?secret=BASE32SECRET
	if !strings.HasPrefix(uri, "otpauth://") {
		return "", fmt.Errorf("invalid otpauth URI")
	}

	// Find the secret parameter
	parts := strings.Split(uri, "?")
	if len(parts) < 2 {
		return "", fmt.Errorf("no query parameters in URI")
	}

	// Parse query parameters
	params := strings.Split(parts[1], "&")
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 && kv[0] == "secret" {
			return kv[1], nil
		}
	}

	return "", fmt.Errorf("secret parameter not found in URI")
}

type ackResponse struct {
	Msg             string         `json:"msg"`
	OK              bool           `json:"ok"`
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
	TokenRequired   bool           `json:"tokenRequired"`
}

func (c *Client) syncEmit(ctx context.Context, command string, args ...any) (ackResponse, error) {
	res := make(chan ackResponse)
	defer close(res)

	args = append(args, emit.WithAck(func(response ackResponse) {
		res <- response
	}))

	err := c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %v", command, err)
	}

	select {
	case response := <-res:
		if !response.OK {
			return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)
		}

		return response, nil
	case <-ctx.Done():
		return ackResponse{}, fmt.Errorf("%s: %v", command, ctx.Err())
	}
}

func (c *Client) syncEmitWithUpdateEvent(ctx context.Context, command string, updateEvent string, args ...any) (ackResponse, error) {
	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Register listener for notifications updates.
	// Signal done, if update is received and remove listener.
	listenerID := uuid.New()
	c.updates.AddListener(func(ctx context.Context, update string) {
		if update == updateEvent {
			c.updates.RemoveListener(listenerID.String())
			closeDone()
		}
	}, listenerID.String())

	res := make(chan ackResponse)
	defer close(res)

	args = append(args, emit.WithAck(func(response ackResponse) {
		res <- response
	}))
	err := c.socketioClient.Emit(command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %v", command, err)
	}

	var response ackResponse
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

			res = nil
		case <-ctx.Done():
			return ackResponse{}, fmt.Errorf("%s: %v", command, ctx.Err())
		}
	}

	return response, nil
}
