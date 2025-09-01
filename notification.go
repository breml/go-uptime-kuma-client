package kuma

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func (c *Client) GetNotifications(_ context.Context) []notification.Base {
	c.mu.Lock()
	defer c.mu.Unlock()

	notifications := make([]notification.Base, len(c.state.notifications))
	copy(notifications, c.state.notifications)

	return notifications
}

func (c *Client) GetNotification(_ context.Context, id int) (notification.Base, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, notification := range c.state.notifications {
		if notification.GetID() == id {
			return notification, nil
		}
	}

	return notification.Base{}, fmt.Errorf("get notification: %w", ErrNotFound)
}

func (c *Client) GetNotificationAs(ctx context.Context, id int, target any) error {
	notification, err := c.GetNotification(ctx, id)
	if err != nil {
		return err
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return json.Unmarshal(notificationJSON, target)
}

func (c *Client) CreateNotification(ctx context.Context, notification notification.Notification) (int, error) {
	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	res := make(chan ackResponse)
	defer close(res)

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
		notification,
		nil, // no ID, create new entry.
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	)
	if err != nil {
		return 0, fmt.Errorf("create notification: %v", err)
	}

	var id int

	// Ensure, we have received both signals: done and ack
	// Setting channel to nil blocks forever, thisway we ensure, that
	// we also receive the second signal.
	for done != nil || res != nil {
		select {
		case <-done:
			done = nil
		case response := <-res:
			if !response.OK {
				return 0, fmt.Errorf("create notification: %s", response.Msg)
			}

			id = response.ID

			res = nil
		case <-ctx.Done():
			return 0, fmt.Errorf("create notification: %v", ctx.Err())
		}
	}

	return id, nil
}

func (c *Client) UpdateNotification(ctx context.Context, notification notification.Notification) error {
	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	res := make(chan ackResponse)
	defer close(res)

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
		notification,
		notification.GetID(),
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	)
	if err != nil {
		return fmt.Errorf("update notification: %v", err)
	}

	// Ensure, we have received both signals: done and ack
	// Setting channel to nil blocks forever, thisway we ensure, that
	// we also receive the second signal.
	for done != nil || res != nil {
		select {
		case <-done:
			done = nil
		case response := <-res:
			if !response.OK {
				return fmt.Errorf("update notification: %s", response.Msg)
			}

			res = nil
		case <-ctx.Done():
			return fmt.Errorf("update notification: %v", ctx.Err())
		}
	}

	return nil
}

func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() {
		close(done)
	})
	defer closeDone()

	res := make(chan ackResponse)
	defer close(res)

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
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	)
	if err != nil {
		return fmt.Errorf("delete notification: %v", err)
	}

	// Ensure, we have received both signals: done and ack
	// Setting channel to nil blocks forever, thisway we ensure, that
	// we also receive the second signal.
	for done != nil || res != nil {
		select {
		case <-done:
			done = nil
		case response := <-res:
			if !response.OK {
				return fmt.Errorf("delete notification: %s", response.Msg)
			}

			res = nil
		case <-ctx.Done():
			return fmt.Errorf("delete notification: %v", ctx.Err())
		}
	}

	return nil
}
