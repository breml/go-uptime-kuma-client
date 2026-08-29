package kuma

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// GetNotifications returns all notifications for the authenticated user.
//
// They are served from the local state cache, which the server keeps up to
// date. Resync rebuilds it if an update event was missed, see
// ErrUpdateEventTimeout.
func (c *Client) GetNotifications(_ context.Context) []notification.Base {
	c.mu.Lock()
	defer c.mu.Unlock()

	notifications := make([]notification.Base, len(c.state.notifications))
	copy(notifications, c.state.notifications)

	return notifications
}

// GetNotification returns a specific notification by ID. Like GetNotifications
// it serves from the local state cache, so a notification whose update event
// was missed is reported as ErrNotFound until Resync.
func (c *Client) GetNotification(_ context.Context, id int64) (notification.Base, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, notif := range c.state.notifications {
		if notif.GetID() == id {
			return notif, nil
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
//
// An error wrapping ErrUpdateEventTimeout means the notification was created and
// only the update event is missing; the returned ID identifies it, as does the
// ID an UpdateEventTimeoutError carries, so retrying the call would create a
// duplicate.
func (c *Client) CreateNotification(ctx context.Context, notif notification.Notification) (int64, error) {
	response, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notif, nil)
	if err != nil {
		if errors.Is(err, ErrUpdateEventTimeout) {
			return response.ID, withCreatedID(err, response.ID)
		}

		return 0, err
	}

	return response.ID, nil
}

// UpdateNotification updates an existing notification.
//
// An error wrapping ErrUpdateEventTimeout means the notification was updated and
// only the update event is missing.
func (c *Client) UpdateNotification(ctx context.Context, notif notification.Notification) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "addNotification", "notificationList", notif, notif.GetID())
	return err
}

// DeleteNotification deletes a notification by ID.
//
// An error wrapping ErrUpdateEventTimeout means the notification was deleted and
// only the update event is missing.
func (c *Client) DeleteNotification(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteNotification", "notificationList", id)
	return err
}

// ErrNotificationTypeNotSupported is returned by TestNotification if the server
// has no notification provider registered for the given notification type.
var ErrNotificationTypeNotSupported = errors.New("notification type is not supported")

// serverMsgNotificationTypeNotSupported is the message the server reports when
// the provider lookup for a notification type fails. It is matched here, in the
// single place that knows about it, so callers can use
// ErrNotificationTypeNotSupported instead of matching on the message themselves.
const serverMsgNotificationTypeNotSupported = "Notification type is not supported"

// TestNotification asks the server to dispatch a test message with the given
// notification and returns the message the notification provider reported.
//
// If the server has no provider for the type of notif, the returned error wraps
// ErrNotificationTypeNotSupported. A provider that fails while dispatching (e.g.
// invalid credentials or an unreachable endpoint) is reported as an error too,
// but only if it raises one: some providers report a delivery failure in their
// message and let the server report success regardless. A nil error therefore
// does not prove the notification was delivered, and callers that need to know
// have to inspect the returned message.
func (c *Client) TestNotification(ctx context.Context, notif notification.Notification) (string, error) {
	response, err := c.syncEmit(ctx, "testNotification", notif)
	if err != nil {
		if strings.Contains(err.Error(), serverMsgNotificationTypeNotSupported) {
			return "", fmt.Errorf("testNotification: %w", ErrNotificationTypeNotSupported)
		}

		return "", err
	}

	return response.Msg, nil
}
