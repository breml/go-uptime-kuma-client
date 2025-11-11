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
	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"

	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
)

var ErrNotFound = fmt.Errorf("not found")

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
}

type Client struct {
	socketioClient *socketio.Client
	socketioLogger socketio.Logger
	autosetup      bool

	mu      *sync.Mutex
	updates signals.Signal[string]
	state   state
}

type Option func(c *Client)

func WithAutosetup() Option {
	return func(c *Client) {
		c.autosetup = true
	}
}

func WithLogLevel(level int) Option {
	return func(c *Client) {
		if level >= utils.DEBUG && level <= utils.NONE {
			c.socketioLogger = &utils.DefaultLogger{Level: level}
		}
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
			resp.Body.Close()
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

func New(ctx context.Context, baseURL string, username string, password string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},

		mu:      &sync.Mutex{},
		updates: signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
	}

	// Handle database setup for Uptime Kuma v2 if autosetup is enabled
	if c.autosetup {
		if err := setupDatabase(ctx, baseURL); err != nil {
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

	updateSeenMu := sync.Mutex{}
	updateSeenMu.Lock()
	updateSeen := map[string]struct{}{
		"monitorList":      empty,
		"maintenanceList":  empty,
		"notificationList": empty,
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

	client.On("notificationList", func(notificationList []notification.Base) {
		c.mu.Lock()
		c.state.notifications = notificationList
		defer c.mu.Unlock()

		c.updates.Emit(ctx, "notificationList")
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

		c.updates.Emit(ctx, "monitorList")
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

	client.OnAny(func(s string, i []any) {
		if s != "notificationList" && s != "monitorList" {
			c.updates.Emit(ctx, s)
		}
	})

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %v", err)
	}

	select {
	case <-connect:
	case <-ctx.Done():
		return nil, fmt.Errorf("connect to server: %v", ctx.Err())
	}

	if username != "" && password != "" {
		_, err = c.syncEmit(ctx, "login", map[string]any{"username": username, "password": password, "token": ""})
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
	}

	for {
		select {
		case <-ready:
			return c, nil

		case <-setupRequired:
			setupRequired = nil

			if !c.autosetup {
				return nil, fmt.Errorf("server does require setup, but autosetup is disabled")
			}

			_, err := c.syncEmit(ctx, "setup", username, password)
			if err != nil {
				return nil, fmt.Errorf("setup: %v", err)
			}

			_, err = c.syncEmit(ctx, "login", map[string]any{"username": username, "password": password, "token": ""})
			if err != nil {
				return nil, fmt.Errorf("login: %v", err)
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("wait for ready: %v", ctx.Err())
		}
	}
}

func (c *Client) Disconnect() error {
	return c.socketioClient.Close()
}

type ackResponse struct {
	Msg       string         `json:"msg"`
	OK        bool           `json:"ok"`
	ID        int64          `json:"id"`
	MonitorID int64          `json:"monitorID"`
	Monitor   map[string]any `json:"monitor"`
	Data      map[string]any `json:"data"`
	Tags      []any          `json:"tags"`
	Tag       map[string]any `json:"tag"`
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
