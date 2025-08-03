package kuma

import (
	"context"
	"fmt"
	"sync"

	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"

	"github.com/breml/go-uptime-kuma-client/model"
)

var ErrNotFound = fmt.Errorf("not found")

var empty = struct{}{}

type state struct {
	notifications []model.Notification
}

type Client struct {
	socketioClient *socketio.Client

	mu      *sync.Mutex
	updates signals.Signal[string]
	state   state
}

func New(ctx context.Context, baseURL string, username string, password string) (*Client, error) {
	client, err := socketio.NewClient(
		socketio.WithRawURL(baseURL),
		socketio.WithLogger(&utils.DefaultLogger{Level: utils.ERROR}),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %v", err)
	}

	c := Client{
		socketioClient: client,

		mu:      &sync.Mutex{},
		updates: signals.New[string](),
	}

	updateSeenMu := sync.Mutex{}
	updateSeenMu.Lock()
	updateSeen := map[string]struct{}{
		"notificationList": empty,
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

	client.On("notificationList", func(notificationList []model.Notification) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.state.notifications = notificationList
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

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %v", err)
	}

	select {
	case <-connect:
	case <-ctx.Done():
		return nil, fmt.Errorf("connect to server: %v", ctx.Err())
	}

	err = client.Emit("login", map[string]any{"username": username, "password": password, "token": ""},
		emit.WithAck(func(loginResponse map[string]any) {}),
	)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	select {
	case <-ready:
		return &c, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for ready: %v", ctx.Err())
	}
}

func (c *Client) Disconnect() error {
	return c.socketioClient.Close()
}
