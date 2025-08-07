package kuma_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/model"
)

func TestNotification(t *testing.T) {
	ctx := t.Context()

	notificationConfig := model.NotificationConfig{
		"applyExisting":            true,
		"isDefault":                true,
		"name":                     "Created",
		"ntfyAuthenticationMethod": "none",
		"ntfyPriority":             5.0,
		"ntfyserverurl":            "https://ntfy.sh",
		"ntfytopic":                "created",
		"type":                     "ntfy",
	}

	// Empty list
	notifications := client.GetNotifications(ctx)
	require.Empty(t, notifications)

	// Create

	id, err := client.CreateNotification(ctx, model.Notification{
		Name:      "Created",
		IsDefault: true,
		Config:    notificationConfig,
	})
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Len(t, notifications, 1)

	createdNotification, err := client.GetNotification(ctx, id)
	require.NoError(t, err)

	require.Equal(t, notificationConfig, createdNotification.Config)

	// Update

	createdNotification.Name = "Updated"
	createdNotification.Config["name"] = "Updated"

	err = client.UpdateNotification(ctx, createdNotification)
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Len(t, notifications, 1)

	updatedNotification, err := client.GetNotification(ctx, id)
	require.NoError(t, err)

	require.Equal(t, createdNotification, updatedNotification)

	// Delete

	err = client.DeleteNotification(ctx, id)
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Empty(t, notifications)
}
