package kuma

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"

	"github.com/breml/go-uptime-kuma-client/model"
)

type notificationResponse struct {
	Msg string `json:"msg"`
	OK  bool   `json:"ok"`
	ID  int    `json:"id"`
}

func (c *Client) GetNotifications(_ context.Context) []model.Notification {
	c.mu.Lock()
	defer c.mu.Unlock()

	notifications := make([]model.Notification, len(c.state.notifications))
	copy(notifications, c.state.notifications)

	return notifications
}

func (c *Client) GetNotification(_ context.Context, id int) (model.Notification, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, notification := range c.state.notifications {
		if notification.ID == id {
			return notification, nil
		}
	}

	return model.Notification{}, fmt.Errorf("get notification: %w", ErrNotFound)
}

func (c *Client) CreateNotification(ctx context.Context, notification model.Notification) (int, error) {
	errChan := make(chan error)
	defer close(errChan)

	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Register listener for notifications updates.
	// Signal done, if update is received and remove listener.
	listenerID := uuid.New()
	c.updates.AddListener(func(ctx context.Context, update string) {
		if update == "notificationList" {
			c.updates.RemoveListener(listenerID.String())
			closeDone()
		}
	}, listenerID.String())

	var id int
	idMu := sync.Mutex{}

	err := c.socketioClient.Emit("addNotification",
		notification.Writeable(),
		nil, // no ID, create new entry.
		emit.WithAck(func(response notificationResponse) {
			if !response.OK {
				errChan <- fmt.Errorf("ack: %s", response.Msg)
			}

			idMu.Lock()
			defer idMu.Unlock()
			id = response.ID
		}),
	)
	if err != nil {
		return 0, fmt.Errorf("create notification: %v", err)
	}

	select {
	case <-done:
		idMu.Lock()
		defer idMu.Unlock()

		return id, nil
	case err := <-errChan:
		return 0, fmt.Errorf("create notification: %v", err)
	case <-ctx.Done():
		return 0, fmt.Errorf("create notification: %v", ctx.Err())
	}
}

func (c *Client) UpdateNotification(ctx context.Context, notification model.Notification) error {
	errChan := make(chan error)
	defer close(errChan)

	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Register listener for notifications updates.
	// Signal done, if update is received and remove listener.
	listenerID := uuid.New()
	c.updates.AddListener(func(ctx context.Context, update string) {
		if update == "notificationList" {
			c.updates.RemoveListener(listenerID.String())
			closeDone()
		}
	}, listenerID.String())

	err := c.socketioClient.Emit("addNotification",
		notification.Writeable(),
		notification.ID,
		emit.WithAck(func(response notificationResponse) {
			if !response.OK {
				errChan <- fmt.Errorf("ack: %s", response.Msg)
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("update notification: %v", err)
	}

	select {
	case <-done:
		return nil
	case err := <-errChan:
		return fmt.Errorf("update notification: %v", err)
	case <-ctx.Done():
		return fmt.Errorf("update notification: %v", ctx.Err())
	}
}

func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	errChan := make(chan error)
	defer close(errChan)

	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	// Register listener for notifications updates.
	// Signal done, if update is received and remove listener.
	listenerID := uuid.New()
	c.updates.AddListener(func(ctx context.Context, update string) {
		if update == "notificationList" {
			c.updates.RemoveListener(listenerID.String())
			closeDone()
		}
	}, listenerID.String())

	err := c.socketioClient.Emit("deleteNotification",
		id,
		emit.WithAck(func(response notificationResponse) {
			if !response.OK {
				errChan <- fmt.Errorf("ack: %s", response.Msg)
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("delete notification: %v", err)
	}

	select {
	case <-done:
		return nil
	case err := <-errChan:
		return fmt.Errorf("delete notification: %v", err)
	case <-ctx.Done():
		return fmt.Errorf("delete notification: %v", ctx.Err())
	}
}
