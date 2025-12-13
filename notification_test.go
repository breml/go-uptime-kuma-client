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
			Host:           "smtp.gmail.com",
			Port:           587,
			Secure:         false,
			IgnoreTLSError: false,
			From:           "noreply@example.com",
			To:             "alerts@example.com",
			CustomSubject:  "Alert: {{ monitorJSON['name'] }}",
			CustomBody:     "Status: {{ msg }}",
			HTMLBody:       true,
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

func TestPagerDutyNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.PagerDuty{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test PagerDuty Created",
		},
		PagerDutyDetails: notification.PagerDutyDetails{
			IntegrationURL: "https://events.pagerduty.com/v2/enqueue",
			IntegrationKey: "test-integration-key-123",
			Priority:       "warning",
			AutoResolve:    "resolve",
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
		require.Equal(t, "PagerDuty", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.PagerDuty{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedPagerDuty := createNotification
		expectedPagerDuty.ID = id
		expectedPagerDuty.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedPagerDuty, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.PagerDuty{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test PagerDuty Updated"
		current.Priority = "critical"
		current.AutoResolve = "null"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.PagerDuty{}
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

func TestSignalNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Signal{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Signal Created",
		},
		SignalDetails: notification.SignalDetails{
			URL:        "http://localhost:9998",
			Number:     "+1234567890",
			Recipients: "+9876543210,+1112223333",
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
		require.Equal(t, "signal", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Signal{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedSignal := createNotification
		expectedSignal.ID = id
		expectedSignal.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedSignal, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Signal{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Signal Updated"
		current.URL = "http://signal-api:9998"
		current.Recipients = "+1111111111,+2222222222"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Signal{}
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

func TestOpsgenieNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Opsgenie{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Opsgenie Created",
		},
		OpsgenieDetails: notification.OpsgenieDetails{
			ApiKey:   "test-api-key-123",
			Region:   "us",
			Priority: 3,
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
		require.Equal(t, "Opsgenie", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Opsgenie{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedOpsgenie := createNotification
		expectedOpsgenie.ID = id
		expectedOpsgenie.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedOpsgenie, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Opsgenie{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Opsgenie Updated"
		current.Region = "eu"
		current.Priority = 5

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Opsgenie{}
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

func TestHomeAssistantNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.HomeAssistant{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Home Assistant Created",
		},
		HomeAssistantDetails: notification.HomeAssistantDetails{
			HomeAssistantURL:     "http://192.168.1.100:8123",
			LongLivedAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			NotificationService:  "notify.mobile_app_iphone",
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
		require.Equal(t, "HomeAssistant", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.HomeAssistant{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedHomeAssistant := createNotification
		expectedHomeAssistant.ID = id
		expectedHomeAssistant.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedHomeAssistant, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.HomeAssistant{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Home Assistant Updated"
		current.HomeAssistantURL = "http://ha.example.com:8123"
		current.NotificationService = "notify.persistent_notification"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.HomeAssistant{}
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

func TestDiscordNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Discord{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Discord Created",
		},
		DiscordDetails: notification.DiscordDetails{
			WebhookURL:    "https://discordapp.com/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz",
			Username:      "Uptime Monitor",
			ChannelType:   "postToThread",
			ThreadID:      "987654321",
			PrefixMessage: "Alert:",
			DisableURL:    true,
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
		require.Equal(t, "discord", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Discord{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedDiscord := createNotification
		expectedDiscord.ID = id
		expectedDiscord.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedDiscord, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Discord{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Discord Updated"
		current.Username = "Updated Monitor"
		current.ChannelType = "createNewForumPost"
		current.PostName = "System Alert"
		current.ThreadID = ""

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Discord{}
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

func TestPushbulletNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Pushbullet{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Pushbullet Created",
		},
		PushbulletDetails: notification.PushbulletDetails{
			AccessToken: "o.example_access_token",
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
		require.Equal(t, "pushbullet", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Pushbullet{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedPushbullet := createNotification
		expectedPushbullet.ID = id
		expectedPushbullet.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedPushbullet, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Pushbullet{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Pushbullet Updated"
		current.AccessToken = "o.updated_access_token"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Pushbullet{}
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

func TestPushoverNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Pushover{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Pushover Created",
		},
		PushoverDetails: notification.PushoverDetails{
			UserKey:  "userkey123",
			AppToken: "apptoken456",
			Sounds:   "echo",
			SoundsUp: "cashregister",
			Priority: "1",
			Title:    "Uptime Kuma Alert",
			Device:   "iphone",
			TTL:      "3600",
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
		require.Equal(t, "pushover", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Pushover{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedPushover := createNotification
		expectedPushover.ID = id
		expectedPushover.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedPushover, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Pushover{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Pushover Updated"
		current.Title = "Updated Alert"
		current.Priority = "2"
		current.Device = "android"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Pushover{}
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

func TestGotifyNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Gotify{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Gotify Created",
		},
		GotifyDetails: notification.GotifyDetails{
			ServerURL:        "https://gotify.example.com",
			ApplicationToken: "test-token",
			Priority:         8,
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
		require.Equal(t, "gotify", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Gotify{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedGotify := createNotification
		expectedGotify.ID = id
		expectedGotify.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedGotify, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Gotify{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Gotify Updated"
		current.Priority = 5

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Gotify{}
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

func TestGrafanaOncallNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.GrafanaOncall{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Grafana OnCall Created",
		},
		GrafanaOncallDetails: notification.GrafanaOncallDetails{
			GrafanaOncallURL: "https://alerts.grafana.com/api/v1/incidents/create",
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
		require.Equal(t, "GrafanaOncall", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.GrafanaOncall{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedGrafanaOncall := createNotification
		expectedGrafanaOncall.ID = id
		expectedGrafanaOncall.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedGrafanaOncall, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.GrafanaOncall{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Grafana OnCall Updated"
		current.GrafanaOncallURL = "https://oncall.example.com/api/v1/incidents/create"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.GrafanaOncall{}
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

func TestTwilioNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Twilio{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Twilio Created",
		},
		TwilioDetails: notification.TwilioDetails{
			AccountSID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ApiKey:     "",
			AuthToken:  "test_auth_token",
			ToNumber:   "+15551234567",
			FromNumber: "+15559876543",
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
		require.Equal(t, "twilio", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Twilio{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedTwilio := createNotification
		expectedTwilio.ID = id
		expectedTwilio.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedTwilio, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Twilio{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Twilio Updated"
		current.ToNumber = "+15559999999"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Twilio{}
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

func TestMattermostNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Mattermost{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Mattermost Created",
		},
		MattermostDetails: notification.MattermostDetails{
			WebhookURL: "https://mattermost.example.com/hooks/xxx",
			Username:   "Monitor Bot",
			Channel:    "#alerts",
			IconEmoji:  ":smiley:",
			IconURL:    "https://example.com/icon.png",
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
		require.Equal(t, "mattermost", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Mattermost{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedMattermost := createNotification
		expectedMattermost.ID = id
		expectedMattermost.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedMattermost, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Mattermost{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Mattermost Updated"
		current.Username = "Updated Bot"
		current.Channel = "#monitoring"
		current.IconEmoji = ":warning:"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Mattermost{}
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

func TestMatrixNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	createNotification := notification.Matrix{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     false,
			IsActive:      true,
			Name:          "Test Matrix Created",
		},
		MatrixDetails: notification.MatrixDetails{
			HomeserverURL:  "https://matrix.example.com",
			InternalRoomID: "!roomid:example.com",
			AccessToken:    "test_access_token",
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
		require.Equal(t, "matrix", createdNotification.Type())
		require.Equal(t, id, createdNotification.GetID())

		specificNotification := notification.Matrix{}
		err = createdNotification.As(&specificNotification)
		require.NoError(t, err)

		expectedMatrix := createNotification
		expectedMatrix.ID = id
		expectedMatrix.UserID = specificNotification.UserID
		require.EqualExportedValues(t, expectedMatrix, specificNotification)
	})

	t.Run("update", func(t *testing.T) {
		currentNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		current := notification.Matrix{}
		err = currentNotification.As(&current)
		require.NoError(t, err)

		current.Name = "Test Matrix Updated"
		current.InternalRoomID = "!newroomid:example.com"

		err = client.UpdateNotification(ctx, current)
		require.NoError(t, err)

		retrievedNotification, err := client.GetNotification(ctx, id)
		require.NoError(t, err)

		retrieved := notification.Matrix{}
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
