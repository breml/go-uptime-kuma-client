package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// GetNotifications returns all notifications for the authenticated user.
func (c *Client) GetNotifications(_ context.Context) []notification.Base {
	c.mu.Lock()
	defer c.mu.Unlock()

	notifications := make([]notification.Base, len(c.state.notifications))
	copy(notifications, c.state.notifications)

	return notifications
}

// GetNotification returns a specific notification by ID.
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

// GetNotificationAs returns a specific notification by ID and coverts it to the target type.
func (c *Client) GetNotificationAs(ctx context.Context, id int64, target any) error {
	notif, err := c.GetNotification(ctx, id)
	if err != nil {
		return err
	}

	err = notif.As(target)
	if err != nil {
		return fmt.Errorf("get monitor %d as %t: %w", id, target, err)
	}

	return nil
}

// CreateNotification creates a new notification.
func (c *Client) CreateNotification(ctx context.Context, notif notification.Notification) (int64, error) {
	response, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notif, nil)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

// UpdateNotification updates an existing notification.
func (c *Client) UpdateNotification(ctx context.Context, notif notification.Notification) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notif, notif.GetID())
	return err
}

// DeleteNotification deletes a notification by ID.
func (c *Client) DeleteNotification(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteNotification", "notificationList", id)
	return err
}
