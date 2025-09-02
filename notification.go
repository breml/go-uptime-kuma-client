package kuma

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func (c *Client) GetNotifications(_ context.Context) []notification.Base {
	c.mu.Lock()
	defer c.mu.Unlock()

	notifications := make([]notification.Base, len(c.state.notifications))
	copy(notifications, c.state.notifications)

	return notifications
}

func (c *Client) GetNotification(_ context.Context, id int64) (notification.Base, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, notification := range c.state.notifications {
		if notification.GetID() == id {
			return notification, nil
		}
	}

	return notification.Base{}, fmt.Errorf("get notification: %w", ErrNotFound)
}

func (c *Client) GetNotificationAs(ctx context.Context, id int64, target any) error {
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

func (c *Client) CreateNotification(ctx context.Context, notification notification.Notification) (int64, error) {
	response, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notification, nil)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

func (c *Client) UpdateNotification(ctx context.Context, notification notification.Notification) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notification, notification.GetID())
	return err
}

func (c *Client) DeleteNotification(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteNotification", "notificationList", id)
	return err
}
