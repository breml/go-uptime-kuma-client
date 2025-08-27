package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotification(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	// Empty list
	notifications := client.GetNotifications(ctx)
	require.Empty(t, notifications)

	// Create
	noti := notification.Ntfy{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     true,
			IsActive:      true,
			Name:          "Created",
		},
		NtfyDetails: notification.NtfyDetails{
			AuthenticationMethod: "none",
			Priority:             5,
			ServerURL:            "https://ntfy.sh",
			Topic:                "topic",
		},
	}

	id, err := client.CreateNotification(ctx, noti)
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Len(t, notifications, 1)

	for _, n := range notifications {
		t.Log(n)
	}

	createdNotification, err := client.GetNotification(ctx, id)
	require.NoError(t, err)

	t.Log(createdNotification)

	createdNotificationNtfy := notification.Ntfy{}
	err = createdNotification.As(&createdNotificationNtfy)
	require.NoError(t, err)

	t.Log(createdNotificationNtfy)

	// Update original notification
	noti.ID = id
	noti.UserID = createdNotificationNtfy.UserID

	require.EqualExportedValues(t, noti, createdNotificationNtfy)

	// Update
	createdNotificationNtfy.Name = "Updated"
	createdNotificationNtfy.AuthenticationMethod = "usernamePassword"
	createdNotificationNtfy.Username = "user"
	createdNotificationNtfy.Password = "password"

	err = client.UpdateNotification(ctx, createdNotificationNtfy)
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Len(t, notifications, 1)

	updatedNotification, err := client.GetNotification(ctx, id)
	require.NoError(t, err)

	updatedNotificationNtfy := notification.Ntfy{}
	err = updatedNotification.As(&updatedNotificationNtfy)
	require.NoError(t, err)

	require.EqualExportedValues(t, createdNotificationNtfy, updatedNotificationNtfy)

	// Delete
	err = client.DeleteNotification(ctx, id)
	require.NoError(t, err)

	notifications = client.GetNotifications(ctx)
	require.Empty(t, notifications)
}
