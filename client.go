package kuma

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"

	"github.com/breml/go-uptime-kuma-client/notification"
)

var (
	ErrNotFound        = fmt.Errorf("not found")
	ErrInvalidResponse = fmt.Errorf("invalid response")
)

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

type state struct {
	notifications []notification.Base
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

func New(ctx context.Context, baseURL string, username string, password string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},

		mu:      &sync.Mutex{},
		updates: signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
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
		// fmt.Printf("%s: %#v", s, i)
		if s != "notificationList" {
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

			if !strings.Contains(err.Error(), "Incorrect username or password") || !wantSetup {
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

func (c *Client) syncEmit(ctx context.Context, event any, args ...any) (map[string]any, error) {
	resp := make(chan map[string]any)

	args = append(args,
		emit.WithAck(func(loginResponse map[string]any) {
			resp <- loginResponse
		}),
	)

	err := c.socketioClient.Emit(event, args...,
	)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	select {
	case response := <-resp:
		okValue, ok := response["ok"]
		if !ok {
			return response, ErrInvalidResponse
		}

		success, ok := okValue.(bool)
		if !ok {
			return response, ErrInvalidResponse
		}

		if !success {
			msg := response["msg"]
			return response, fmt.Errorf("%v failed: %v", event, msg)
		}

		return response, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("syncEmit wait for response to: %v (args: %v)", event, args)
	}
}
