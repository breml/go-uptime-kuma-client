package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// Outbound notification proxy environment variables.
//
// Since Uptime Kuma v2 (upstream PR louislam/uptime-kuma#7088), all built-in
// notification providers route their outbound HTTP(S) requests through an
// optional proxy. The proxy is configured globally on the Uptime Kuma server
// process via the [NotificationProxyEnvVar] or [NotificationProxyEnvVarUpper]
// environment variable (the lowercase variant takes precedence).
//
// The value must be a URL. The supported URL schemes are:
//
//   - http://      - HTTP proxy
//   - https://     - HTTPS proxy
//   - socks://     - SOCKS proxy (alias for socks5://)
//   - socks4://    - SOCKS4 proxy
//   - socks5://    - SOCKS5 proxy
//   - socks5h://   - SOCKS5 proxy with remote DNS resolution
//
// Because the proxy is configured on the Uptime Kuma server itself (and not
// per notification or via the settings store exposed by the API), there is
// no corresponding field on the notification types in this client. To enable
// the proxy, set the environment variable on the Uptime Kuma server before
// starting it, for example:
//
//	NOTIFICATION_PROXY=socks5h://proxy.example.com:1080 npm run start-server
const (
	// NotificationProxyEnvVar is the lowercase environment variable name
	// read by the Uptime Kuma server to configure the proxy used for
	// outbound notification requests. It takes precedence over
	// [NotificationProxyEnvVarUpper].
	NotificationProxyEnvVar = "notification_proxy"

	// NotificationProxyEnvVarUpper is the uppercase environment variable
	// name read by the Uptime Kuma server to configure the proxy used for
	// outbound notification requests. It is only consulted if
	// [NotificationProxyEnvVar] is not set.
	NotificationProxyEnvVarUpper = "NOTIFICATION_PROXY"
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
