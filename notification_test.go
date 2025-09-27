package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNtfyNotificationCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	createNotification := notification.Ntfy{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     true,
			IsActive:      true,
			Name:          "Test Ntfy Created",
		},
		NtfyDetails: notification.NtfyDetails{
			AuthenticationMethod: "none",
			Priority:             5,
			ServerURL:            "https://ntfy.sh",
			Topic:                "test-topic",
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err := client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "ntfy", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Ntfy{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedNtfy := createNotification
		expectedNtfy.ID = id
		expectedNtfy.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedNtfy, specificNotification)

		t.Run("update", func(t *testing.T) {
			currentNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			current := notification.Ntfy{}
			err = currentNotification.As(&current)
			require.NoError(t, err)

			current.Name = "Test Ntfy Updated"
			current.AuthenticationMethod = "usernamePassword"
			current.Username = "testuser"
			current.Password = "testpass"
			current.Priority = 3

			err = client.UpdateNotification(ctx, current)
			require.NoError(t, err)

			retrievedNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			retrieved := notification.Ntfy{}
			err = retrievedNotification.As(&retrieved)
			require.NoError(t, err)
			require.EqualExportedValues(t, current, retrieved)
		})

		t.Run("delete", func(t *testing.T) {
			preDeleteNotifications := client.GetNotifications(ctx)
			preDeleteCount := len(preDeleteNotifications)

			err := client.DeleteNotification(ctx, id)
			require.NoError(t, err)

			notifications := client.GetNotifications(ctx)
			require.Len(t, notifications, preDeleteCount-1)

			_, err = client.GetNotification(ctx, id)
			require.Error(t, err)
		})
	})
}

func TestSlackNotificationCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	createNotification := notification.Slack{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Slack Created",
		},
		SlackDetails: notification.SlackDetails{
			WebhookURL:    "https://hooks.slack.com/services/test/webhook/url",
			Username:      "uptime-kuma",
			IconEmoji:     ":robot_face:",
			Channel:       "#alerts",
			RichMessage:   true,
			ChannelNotify: false,
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err := client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "slack", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Slack{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedSlack := createNotification
		expectedSlack.ID = id
		expectedSlack.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedSlack, specificNotification)

		t.Run("update", func(t *testing.T) {
			currentNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			current := notification.Slack{}
			err = currentNotification.As(&current)
			require.NoError(t, err)

			current.Name = "Test Slack Updated"
			current.Username = "uptime-bot"
			current.IconEmoji = ":warning:"
			current.Channel = "#monitoring"
			current.ChannelNotify = true

			err = client.UpdateNotification(ctx, current)
			require.NoError(t, err)

			retrievedNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			retrieved := notification.Slack{}
			err = retrievedNotification.As(&retrieved)
			require.NoError(t, err)
			require.EqualExportedValues(t, current, retrieved)
		})

		t.Run("delete", func(t *testing.T) {
			preDeleteNotifications := client.GetNotifications(ctx)
			preDeleteCount := len(preDeleteNotifications)

			err := client.DeleteNotification(ctx, id)
			require.NoError(t, err)

			notifications := client.GetNotifications(ctx)
			require.Len(t, notifications, preDeleteCount-1)

			_, err = client.GetNotification(ctx, id)
			require.Error(t, err)
		})
	})
}

func TestTeamsNotificationCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	createNotification := notification.Teams{
		Base: notification.Base{
			ApplyExisting: false,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Teams Created",
		},
		TeamsDetails: notification.TeamsDetails{
			WebhookURL: "https://outlook.office.com/webhook/xxx-xxx-xxx/IncomingWebhook/yyy-yyy-yyy",
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err := client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "teams", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Teams{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedTeams := createNotification
		expectedTeams.ID = id
		expectedTeams.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedTeams, specificNotification)

		t.Run("update", func(t *testing.T) {
			currentNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			current := notification.Teams{}
			err = currentNotification.As(&current)
			require.NoError(t, err)

			current.Name = "Test Teams Updated"
			current.WebhookURL = "https://outlook.office.com/webhook/updated-xxx-xxx/IncomingWebhook/updated-yyy-yyy"

			err = client.UpdateNotification(ctx, current)
			require.NoError(t, err)

			retrievedNotification, err := client.GetNotification(ctx, id)
			require.NoError(t, err)

			retrieved := notification.Teams{}
			err = retrievedNotification.As(&retrieved)
			require.NoError(t, err)
			require.EqualExportedValues(t, current, retrieved)
		})

		t.Run("delete", func(t *testing.T) {
			preDeleteNotifications := client.GetNotifications(ctx)
			preDeleteCount := len(preDeleteNotifications)

			err := client.DeleteNotification(ctx, id)
			require.NoError(t, err)

			notifications := client.GetNotifications(ctx)
			require.Len(t, notifications, preDeleteCount-1)

			_, err = client.GetNotification(ctx, id)
			require.Error(t, err)
		})
	})
}
