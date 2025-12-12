package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNtfyNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

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

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
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
	})

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
}

func TestSlackNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

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

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
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
	})

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
}

func TestTeamsNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

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

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
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
	})

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
}

func TestWebhookNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Webhook{
		Base: notification.Base{
			ApplyExisting: false,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Webhook Created",
		},
		WebhookDetails: notification.WebhookDetails{
			WebhookURL:         "https://example.com/webhook",
			WebhookContentType: "json",
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "webhook", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Webhook{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedWebhook := createNotification
		expectedWebhook.ID = id
		expectedWebhook.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedWebhook, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Webhook{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Webhook Updated"
		current.WebhookContentType = "custom"
		current.WebhookCustomBody = `{"title": "Alert - {{ monitorJSON['name'] }}", "message": "{{ msg }}"}`
		current.WebhookAdditionalHeaders = notification.WebhookAdditionalHeaders{
			"Authorization": "Bearer test-token",
			"X-Custom":      "test-value",
		}

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Webhook{}
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
}

func TestSMTPNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.SMTP{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test SMTP Created",
		},
		SMTPDetails: notification.SMTPDetails{
			Host:          "smtp.gmail.com",
			Port:          587,
			Secure:        false,
			IgnoreTLSError: false,
			From:          "noreply@example.com",
			To:            "alerts@example.com",
			CustomSubject: "Alert: {{ monitorJSON['name'] }}",
			CustomBody:    "Status: {{ msg }}",
			HTMLBody:      true,
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "smtp", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.SMTP{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedSMTP := createNotification
		expectedSMTP.ID = id
		expectedSMTP.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedSMTP, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.SMTP{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test SMTP Updated"
		current.Host = "smtp.office365.com"
		current.Port = 25
		current.Username = "user@example.com"
		current.Password = "secretpassword"
		current.CC = "cc@example.com"
		current.BCC = "bcc@example.com"
		current.Secure = true

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.SMTP{}
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
}

func TestTelegramNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Telegram{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Telegram Created",
		},
		TelegramDetails: notification.TelegramDetails{
			BotToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			ChatID:   "@mychannel",
		},
	}

	t.Run("initial_state", func(t *testing.T) {
		notifications := client.GetNotifications(ctx)
		t.Logf("Initial notifications count: %d", len(notifications))
	})

	var id int64
	t.Run("create", func(t *testing.T) {
		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err = client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "telegram", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Telegram{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedTelegram := createNotification
		expectedTelegram.ID = id
		expectedTelegram.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedTelegram, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Telegram{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Telegram Updated"
		current.ChatID = "123456789"
		current.SendSilently = true

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Telegram{}
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
}

func TestWebhookNotificationVariants(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	t.Run("form-data content type", func(t *testing.T) {
		createNotification := notification.Webhook{
			Base: notification.Base{
				ApplyExisting: false,
				IsDefault:     false,
				IsActive:      true,
				Name:          "Test Webhook Form-Data",
			},
			WebhookDetails: notification.WebhookDetails{
				WebhookURL:         "https://example.com/form-webhook",
				WebhookContentType: "form-data",
			},
		}

		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err := client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "webhook", createdNotification.Type())

		specificNotification := notification.Webhook{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		require.Equal(t, "form-data", specificNotification.WebhookContentType)

		// Cleanup
		err = client.DeleteNotification(ctx, id)
		require.NoError(t, err)
	})

	t.Run("with additional headers", func(t *testing.T) {
		createNotification := notification.Webhook{
			Base: notification.Base{
				ApplyExisting: false,
				IsDefault:     false,
				IsActive:      true,
				Name:          "Test Webhook Headers",
			},
			WebhookDetails: notification.WebhookDetails{
				WebhookURL:         "https://api.example.com/webhook",
				WebhookContentType: "json",
				WebhookAdditionalHeaders: notification.WebhookAdditionalHeaders{
					"Authorization": "Bearer secret-token",
					"X-App-ID":      "uptime-kuma",
				},
			},
		}

		initialNotifications := client.GetNotifications(ctx)
		initialCount := len(initialNotifications)

		id, err := client.CreateNotification(ctx, createNotification)
		require.NoError(t, err)
		require.Greater(t, id, int64(0))

		notifications := client.GetNotifications(ctx)
		require.Len(t, notifications, initialCount+1)

		createdNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		specificNotification := notification.Webhook{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		require.Equal(t, 2, len(specificNotification.WebhookAdditionalHeaders))
		require.Equal(t, "Bearer secret-token", specificNotification.WebhookAdditionalHeaders["Authorization"])
		require.Equal(t, "uptime-kuma", specificNotification.WebhookAdditionalHeaders["X-App-ID"])

		// Cleanup
		err = client.DeleteNotification(ctx, id)
		require.NoError(t, err)
	})
}
