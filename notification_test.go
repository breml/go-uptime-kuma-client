package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// notificationTestCase defines a single notification type's CRUD test scenario.
type notificationTestCase struct {
	name              string                                                                                             // Test name (e.g., "Ntfy", "Slack")
	expectedType      string                                                                                             // Expected type string from API
	create            notification.Notification                                                                          // Notification to create
	updateFunc        func(notification.Notification)                                                                    // Function to modify notification for update test
	verifyCreatedFunc func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) // Function to verify created notification
	createTypedFunc   func(t *testing.T, base notification.Notification) notification.Notification                       // Function to create typed notification
	verifyUpdatedFunc func(t *testing.T, actual notification.Notification, expected notification.Notification)           // Function to verify updated notification
}

func TestNotificationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	testCases := []notificationTestCase{
		{
			name:         "Ntfy",
			expectedType: "ntfy",
			create: notification.Ntfy{
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
			},
			updateFunc: func(n notification.Notification) {
				ntfy := n.(*notification.Ntfy)
				ntfy.Name = "Test Ntfy Updated"
				ntfy.AuthenticationMethod = "usernamePassword"
				ntfy.Username = "testuser"
				ntfy.Password = "testpass"
				ntfy.Priority = 3
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Ntfy)
				require.True(t, ok)
				var ntfy notification.Ntfy
				err := actual.As(&ntfy)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = ntfy.UserID
				require.EqualExportedValues(t, exp, ntfy)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var ntfy notification.Ntfy
				err := base.As(&ntfy)
				require.NoError(t, err)
				return &ntfy
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Ntfy)
				require.True(t, ok)
				var ntfy notification.Ntfy
				err := actual.As(&ntfy)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, ntfy)
			},
		},
		{
			name:         "Slack",
			expectedType: "slack",
			create: notification.Slack{
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
			},
			updateFunc: func(n notification.Notification) {
				slack := n.(*notification.Slack)
				slack.Name = "Test Slack Updated"
				slack.Username = "uptime-bot"
				slack.IconEmoji = ":warning:"
				slack.Channel = "#monitoring"
				slack.ChannelNotify = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Slack)
				require.True(t, ok)
				var slack notification.Slack
				err := actual.As(&slack)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = slack.UserID
				require.EqualExportedValues(t, exp, slack)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var slack notification.Slack
				err := base.As(&slack)
				require.NoError(t, err)
				return &slack
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Slack)
				require.True(t, ok)
				var slack notification.Slack
				err := actual.As(&slack)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, slack)
			},
		},
		{
			name:         "Teams",
			expectedType: "teams",
			create: notification.Teams{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Teams Created",
				},
				TeamsDetails: notification.TeamsDetails{
					WebhookURL: "https://outlook.office.com/webhook/xxx-xxx-xxx/IncomingWebhook/yyy-yyy-yyy",
				},
			},
			updateFunc: func(n notification.Notification) {
				teams := n.(*notification.Teams)
				teams.Name = "Test Teams Updated"
				teams.WebhookURL = "https://outlook.office.com/webhook/updated-xxx-xxx/IncomingWebhook/updated-yyy-yyy"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Teams)
				require.True(t, ok)
				var teams notification.Teams
				err := actual.As(&teams)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = teams.UserID
				require.EqualExportedValues(t, exp, teams)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var teams notification.Teams
				err := base.As(&teams)
				require.NoError(t, err)
				return &teams
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Teams)
				require.True(t, ok)
				var teams notification.Teams
				err := actual.As(&teams)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, teams)
			},
		},
		{
			name:         "Webhook",
			expectedType: "webhook",
			create: notification.Webhook{
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
			},
			updateFunc: func(n notification.Notification) {
				webhook := n.(*notification.Webhook)
				webhook.Name = "Test Webhook Updated"
				webhook.WebhookContentType = "custom"
				webhook.WebhookCustomBody = `{"title": "Alert - {{ monitorJSON['name'] }}", "message": "{{ msg }}"}`
				webhook.WebhookAdditionalHeaders = notification.WebhookAdditionalHeaders{
					"Authorization": "Bearer test-token",
					"X-Custom":      "test-value",
				}
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Webhook)
				require.True(t, ok)
				var webhook notification.Webhook
				err := actual.As(&webhook)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = webhook.UserID
				require.EqualExportedValues(t, exp, webhook)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var webhook notification.Webhook
				err := base.As(&webhook)
				require.NoError(t, err)
				return &webhook
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Webhook)
				require.True(t, ok)
				var webhook notification.Webhook
				err := actual.As(&webhook)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, webhook)
			},
		},
		{
			name:         "SMTP",
			expectedType: "smtp",
			create: notification.SMTP{
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
			},
			updateFunc: func(n notification.Notification) {
				smtp := n.(*notification.SMTP)
				smtp.Name = "Test SMTP Updated"
				smtp.Host = "smtp.office365.com"
				smtp.Port = 25
				smtp.Username = "user@example.com"
				smtp.Password = "secretpassword"
				smtp.CC = "cc@example.com"
				smtp.BCC = "bcc@example.com"
				smtp.Secure = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMTP)
				require.True(t, ok)
				var smtp notification.SMTP
				err := actual.As(&smtp)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smtp.UserID
				require.EqualExportedValues(t, exp, smtp)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smtp notification.SMTP
				err := base.As(&smtp)
				require.NoError(t, err)
				return &smtp
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMTP)
				require.True(t, ok)
				var smtp notification.SMTP
				err := actual.As(&smtp)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smtp)
			},
		},
		{
			name:         "Telegram",
			expectedType: "telegram",
			create: notification.Telegram{
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
			},
			updateFunc: func(n notification.Notification) {
				telegram := n.(*notification.Telegram)
				telegram.Name = "Test Telegram Updated"
				telegram.ChatID = "123456789"
				telegram.SendSilently = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Telegram)
				require.True(t, ok)
				var telegram notification.Telegram
				err := actual.As(&telegram)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = telegram.UserID
				require.EqualExportedValues(t, exp, telegram)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var telegram notification.Telegram
				err := base.As(&telegram)
				require.NoError(t, err)
				return &telegram
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Telegram)
				require.True(t, ok)
				var telegram notification.Telegram
				err := actual.As(&telegram)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, telegram)
			},
		},
		{
			name:         "PagerDuty",
			expectedType: "PagerDuty",
			create: notification.PagerDuty{
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
			},
			updateFunc: func(n notification.Notification) {
				pagerduty := n.(*notification.PagerDuty)
				pagerduty.Name = "Test PagerDuty Updated"
				pagerduty.Priority = "critical"
				pagerduty.AutoResolve = "null"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.PagerDuty)
				require.True(t, ok)
				var pagerduty notification.PagerDuty
				err := actual.As(&pagerduty)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pagerduty.UserID
				require.EqualExportedValues(t, exp, pagerduty)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pagerduty notification.PagerDuty
				err := base.As(&pagerduty)
				require.NoError(t, err)
				return &pagerduty
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.PagerDuty)
				require.True(t, ok)
				var pagerduty notification.PagerDuty
				err := actual.As(&pagerduty)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pagerduty)
			},
		},
		{
			name:         "Signal",
			expectedType: "signal",
			create: notification.Signal{
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
			},
			updateFunc: func(n notification.Notification) {
				signal := n.(*notification.Signal)
				signal.Name = "Test Signal Updated"
				signal.URL = "http://signal-api:9998"
				signal.Recipients = "+1111111111,+2222222222"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Signal)
				require.True(t, ok)
				var signal notification.Signal
				err := actual.As(&signal)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = signal.UserID
				require.EqualExportedValues(t, exp, signal)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var signal notification.Signal
				err := base.As(&signal)
				require.NoError(t, err)
				return &signal
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Signal)
				require.True(t, ok)
				var signal notification.Signal
				err := actual.As(&signal)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, signal)
			},
		},
		{
			name:         "Opsgenie",
			expectedType: "Opsgenie",
			create: notification.Opsgenie{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Opsgenie Created",
				},
				OpsgenieDetails: notification.OpsgenieDetails{
					APIKey:   "test-api-key-123",
					Region:   "us",
					Priority: 3,
				},
			},
			updateFunc: func(n notification.Notification) {
				opsgenie := n.(*notification.Opsgenie)
				opsgenie.Name = "Test Opsgenie Updated"
				opsgenie.Region = "eu"
				opsgenie.Priority = 5
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Opsgenie)
				require.True(t, ok)
				var opsgenie notification.Opsgenie
				err := actual.As(&opsgenie)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = opsgenie.UserID
				require.EqualExportedValues(t, exp, opsgenie)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var opsgenie notification.Opsgenie
				err := base.As(&opsgenie)
				require.NoError(t, err)
				return &opsgenie
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Opsgenie)
				require.True(t, ok)
				var opsgenie notification.Opsgenie
				err := actual.As(&opsgenie)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, opsgenie)
			},
		},
		{
			name:         "HomeAssistant",
			expectedType: "HomeAssistant",
			create: notification.HomeAssistant{
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
			},
			updateFunc: func(n notification.Notification) {
				ha := n.(*notification.HomeAssistant)
				ha.Name = "Test Home Assistant Updated"
				ha.HomeAssistantURL = "http://ha.example.com:8123"
				ha.NotificationService = "notify.persistent_notification"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.HomeAssistant)
				require.True(t, ok)
				var ha notification.HomeAssistant
				err := actual.As(&ha)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = ha.UserID
				require.EqualExportedValues(t, exp, ha)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var ha notification.HomeAssistant
				err := base.As(&ha)
				require.NoError(t, err)
				return &ha
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.HomeAssistant)
				require.True(t, ok)
				var ha notification.HomeAssistant
				err := actual.As(&ha)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, ha)
			},
		},
		{
			name:         "Discord",
			expectedType: "discord",
			create: notification.Discord{
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
			},
			updateFunc: func(n notification.Notification) {
				discord := n.(*notification.Discord)
				discord.Name = "Test Discord Updated"
				discord.Username = "Updated Monitor"
				discord.ChannelType = "createNewForumPost"
				discord.PostName = "System Alert"
				discord.ThreadID = ""
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Discord)
				require.True(t, ok)
				var discord notification.Discord
				err := actual.As(&discord)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = discord.UserID
				require.EqualExportedValues(t, exp, discord)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var discord notification.Discord
				err := base.As(&discord)
				require.NoError(t, err)
				return &discord
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Discord)
				require.True(t, ok)
				var discord notification.Discord
				err := actual.As(&discord)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, discord)
			},
		},
		{
			name:         "Pushbullet",
			expectedType: "pushbullet",
			create: notification.Pushbullet{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Pushbullet Created",
				},
				PushbulletDetails: notification.PushbulletDetails{
					AccessToken: "o.example_access_token",
				},
			},
			updateFunc: func(n notification.Notification) {
				pushbullet := n.(*notification.Pushbullet)
				pushbullet.Name = "Test Pushbullet Updated"
				pushbullet.AccessToken = "o.updated_access_token"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Pushbullet)
				require.True(t, ok)
				var pushbullet notification.Pushbullet
				err := actual.As(&pushbullet)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pushbullet.UserID
				require.EqualExportedValues(t, exp, pushbullet)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pushbullet notification.Pushbullet
				err := base.As(&pushbullet)
				require.NoError(t, err)
				return &pushbullet
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Pushbullet)
				require.True(t, ok)
				var pushbullet notification.Pushbullet
				err := actual.As(&pushbullet)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pushbullet)
			},
		},
		{
			name:         "Pushover",
			expectedType: "pushover",
			create: notification.Pushover{
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
			},
			updateFunc: func(n notification.Notification) {
				pushover := n.(*notification.Pushover)
				pushover.Name = "Test Pushover Updated"
				pushover.Title = "Updated Alert"
				pushover.Priority = "2"
				pushover.Device = "android"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Pushover)
				require.True(t, ok)
				var pushover notification.Pushover
				err := actual.As(&pushover)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pushover.UserID
				require.EqualExportedValues(t, exp, pushover)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pushover notification.Pushover
				err := base.As(&pushover)
				require.NoError(t, err)
				return &pushover
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Pushover)
				require.True(t, ok)
				var pushover notification.Pushover
				err := actual.As(&pushover)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pushover)
			},
		},
		{
			name:         "Gotify",
			expectedType: "gotify",
			create: notification.Gotify{
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
			},
			updateFunc: func(n notification.Notification) {
				gotify := n.(*notification.Gotify)
				gotify.Name = "Test Gotify Updated"
				gotify.Priority = 5
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Gotify)
				require.True(t, ok)
				var gotify notification.Gotify
				err := actual.As(&gotify)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = gotify.UserID
				require.EqualExportedValues(t, exp, gotify)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var gotify notification.Gotify
				err := base.As(&gotify)
				require.NoError(t, err)
				return &gotify
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Gotify)
				require.True(t, ok)
				var gotify notification.Gotify
				err := actual.As(&gotify)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, gotify)
			},
		},
		{
			name:         "GrafanaOncall",
			expectedType: "GrafanaOncall",
			create: notification.GrafanaOncall{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Grafana OnCall Created",
				},
				GrafanaOncallDetails: notification.GrafanaOncallDetails{
					GrafanaOncallURL: "https://alerts.grafana.com/api/v1/incidents/create",
				},
			},
			updateFunc: func(n notification.Notification) {
				grafana := n.(*notification.GrafanaOncall)
				grafana.Name = "Test Grafana OnCall Updated"
				grafana.GrafanaOncallURL = "https://oncall.example.com/api/v1/incidents/create"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.GrafanaOncall)
				require.True(t, ok)
				var grafana notification.GrafanaOncall
				err := actual.As(&grafana)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = grafana.UserID
				require.EqualExportedValues(t, exp, grafana)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var grafana notification.GrafanaOncall
				err := base.As(&grafana)
				require.NoError(t, err)
				return &grafana
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.GrafanaOncall)
				require.True(t, ok)
				var grafana notification.GrafanaOncall
				err := actual.As(&grafana)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, grafana)
			},
		},
		{
			name:         "Twilio",
			expectedType: "twilio",
			create: notification.Twilio{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Twilio Created",
				},
				TwilioDetails: notification.TwilioDetails{
					AccountSID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					APIKey:     "",
					AuthToken:  "test_auth_token",
					ToNumber:   "+15551234567",
					FromNumber: "+15559876543",
				},
			},
			updateFunc: func(n notification.Notification) {
				twilio := n.(*notification.Twilio)
				twilio.Name = "Test Twilio Updated"
				twilio.ToNumber = "+15559999999"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Twilio)
				require.True(t, ok)
				var twilio notification.Twilio
				err := actual.As(&twilio)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = twilio.UserID
				require.EqualExportedValues(t, exp, twilio)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var twilio notification.Twilio
				err := base.As(&twilio)
				require.NoError(t, err)
				return &twilio
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Twilio)
				require.True(t, ok)
				var twilio notification.Twilio
				err := actual.As(&twilio)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, twilio)
			},
		},
		{
			name:         "Mattermost",
			expectedType: "mattermost",
			create: notification.Mattermost{
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
			},
			updateFunc: func(n notification.Notification) {
				mattermost := n.(*notification.Mattermost)
				mattermost.Name = "Test Mattermost Updated"
				mattermost.Username = "Updated Bot"
				mattermost.Channel = "#monitoring"
				mattermost.IconEmoji = ":warning:"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Mattermost)
				require.True(t, ok)
				var mattermost notification.Mattermost
				err := actual.As(&mattermost)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = mattermost.UserID
				require.EqualExportedValues(t, exp, mattermost)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var mattermost notification.Mattermost
				err := base.As(&mattermost)
				require.NoError(t, err)
				return &mattermost
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Mattermost)
				require.True(t, ok)
				var mattermost notification.Mattermost
				err := actual.As(&mattermost)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, mattermost)
			},
		},
		{
			name:         "Matrix",
			expectedType: "matrix",
			create: notification.Matrix{
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
			},
			updateFunc: func(n notification.Notification) {
				matrix := n.(*notification.Matrix)
				matrix.Name = "Test Matrix Updated"
				matrix.InternalRoomID = "!newroomid:example.com"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Matrix)
				require.True(t, ok)
				var matrix notification.Matrix
				err := actual.As(&matrix)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = matrix.UserID
				require.EqualExportedValues(t, exp, matrix)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var matrix notification.Matrix
				err := base.As(&matrix)
				require.NoError(t, err)
				return &matrix
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Matrix)
				require.True(t, ok)
				var matrix notification.Matrix
				err := actual.As(&matrix)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, matrix)
			},
		},
		{
			name:         "RocketChat",
			expectedType: "rocket.chat",
			create: notification.RocketChat{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Rocket.Chat Created",
				},
				RocketChatDetails: notification.RocketChatDetails{
					WebhookURL: "https://rocket.example.com/hooks/xxx",
					Channel:    "#alerts",
					Username:   "Monitor Bot",
					IconEmoji:  ":smiley:",
					Button:     "",
				},
			},
			updateFunc: func(n notification.Notification) {
				rocketchat := n.(*notification.RocketChat)
				rocketchat.Name = "Test Rocket.Chat Updated"
				rocketchat.Channel = "#monitoring"
				rocketchat.Username = "Updated Bot"
				rocketchat.IconEmoji = ":warning:"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.RocketChat)
				require.True(t, ok)
				var rocketchat notification.RocketChat
				err := actual.As(&rocketchat)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = rocketchat.UserID
				require.EqualExportedValues(t, exp, rocketchat)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var rocketchat notification.RocketChat
				err := base.As(&rocketchat)
				require.NoError(t, err)
				return &rocketchat
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.RocketChat)
				require.True(t, ok)
				var rocketchat notification.RocketChat
				err := actual.As(&rocketchat)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, rocketchat)
			},
		},
		{
			name:         "WeCom",
			expectedType: "WeCom",
			create: notification.WeCom{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test WeCom Created",
				},
				WeComDetails: notification.WeComDetails{
					BotKey: "xxxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				wecom := n.(*notification.WeCom)
				wecom.Name = "Test WeCom Updated"
				wecom.BotKey = "yyyy"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.WeCom)
				require.True(t, ok)
				var wecom notification.WeCom
				err := actual.As(&wecom)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = wecom.UserID
				require.EqualExportedValues(t, exp, wecom)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var wecom notification.WeCom
				err := base.As(&wecom)
				require.NoError(t, err)
				return &wecom
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.WeCom)
				require.True(t, ok)
				var wecom notification.WeCom
				err := actual.As(&wecom)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, wecom)
			},
		},
		{
			name:         "Feishu",
			expectedType: "Feishu",
			create: notification.Feishu{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Feishu Created",
				},
				FeishuDetails: notification.FeishuDetails{
					WebHookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				feishu := n.(*notification.Feishu)
				feishu.Name = "Test Feishu Updated"
				feishu.WebHookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/yyyy"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Feishu)
				require.True(t, ok)
				var feishu notification.Feishu
				err := actual.As(&feishu)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = feishu.UserID
				require.EqualExportedValues(t, exp, feishu)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var feishu notification.Feishu
				err := base.As(&feishu)
				require.NoError(t, err)
				return &feishu
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Feishu)
				require.True(t, ok)
				var feishu notification.Feishu
				err := actual.As(&feishu)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, feishu)
			},
		},
		{
			name:         "DingDing",
			expectedType: "DingDing",
			create: notification.DingDing{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test DingDing Created",
				},
				DingDingDetails: notification.DingDingDetails{
					WebHookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxxx",
					SecretKey:  "secret123",
					Mentioning: "@all",
				},
			},
			updateFunc: func(n notification.Notification) {
				dingding := n.(*notification.DingDing)
				dingding.Name = "Test DingDing Updated"
				dingding.Mentioning = ""
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.DingDing)
				require.True(t, ok)
				var dingding notification.DingDing
				err := actual.As(&dingding)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = dingding.UserID
				require.EqualExportedValues(t, exp, dingding)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var dingding notification.DingDing
				err := base.As(&dingding)
				require.NoError(t, err)
				return &dingding
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.DingDing)
				require.True(t, ok)
				var dingding notification.DingDing
				err := actual.As(&dingding)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, dingding)
			},
		},
		{
			name:         "FortySixElks",
			expectedType: "46elks",
			create: notification.FortySixElks{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test 46elks Created",
				},
				FortySixElksDetails: notification.FortySixElksDetails{
					Username:   "user@example.com",
					AuthToken:  "test_token",
					FromNumber: "1234",
					ToNumber:   "0701234567",
				},
			},
			updateFunc: func(n notification.Notification) {
				elks := n.(*notification.FortySixElks)
				elks.Name = "Test 46elks Updated"
				elks.ToNumber = "0709999999"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.FortySixElks)
				require.True(t, ok)
				var elks notification.FortySixElks
				err := actual.As(&elks)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = elks.UserID
				require.EqualExportedValues(t, exp, elks)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var elks notification.FortySixElks
				err := base.As(&elks)
				require.NoError(t, err)
				return &elks
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.FortySixElks)
				require.True(t, ok)
				var elks notification.FortySixElks
				err := actual.As(&elks)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, elks)
			},
		},
		{
			name:         "Alerta",
			expectedType: "alerta",
			create: notification.Alerta{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Alerta Created",
				},
				AlertaDetails: notification.AlertaDetails{
					APIEndpoint:  "https://alerta.example.com/api/alerts",
					APIKey:       "test_api_key",
					Environment:  "Production",
					AlertState:   "critical",
					RecoverState: "cleared",
				},
			},
			updateFunc: func(n notification.Notification) {
				alerta := n.(*notification.Alerta)
				alerta.Name = "Test Alerta Updated"
				alerta.Environment = "Staging"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Alerta)
				require.True(t, ok)
				var alerta notification.Alerta
				err := actual.As(&alerta)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = alerta.UserID
				require.EqualExportedValues(t, exp, alerta)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var alerta notification.Alerta
				err := base.As(&alerta)
				require.NoError(t, err)
				return &alerta
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Alerta)
				require.True(t, ok)
				var alerta notification.Alerta
				err := actual.As(&alerta)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, alerta)
			},
		},
		{
			name:         "AlertNow",
			expectedType: "AlertNow",
			create: notification.AlertNow{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test AlertNow Created",
				},
				AlertNowDetails: notification.AlertNowDetails{
					WebhookURL: "https://alertnow.example.com/api/webhook",
				},
			},
			updateFunc: func(n notification.Notification) {
				alertnow := n.(*notification.AlertNow)
				alertnow.Name = "Test AlertNow Updated"
				alertnow.WebhookURL = "https://alertnow.example.com/api/webhook/updated"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.AlertNow)
				require.True(t, ok)
				var alertnow notification.AlertNow
				err := actual.As(&alertnow)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = alertnow.UserID
				require.EqualExportedValues(t, exp, alertnow)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var alertnow notification.AlertNow
				err := base.As(&alertnow)
				require.NoError(t, err)
				return &alertnow
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.AlertNow)
				require.True(t, ok)
				var alertnow notification.AlertNow
				err := actual.As(&alertnow)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, alertnow)
			},
		},
		{
			name:         "AliyunSMS",
			expectedType: "AliyunSMS",
			create: notification.AliyunSMS{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test AliyunSMS Created",
				},
				AliyunSMSDetails: notification.AliyunSMSDetails{
					AccessKeyID:     "AKIA123",
					SecretAccessKey: "secret123",
					PhoneNumber:     "8613800000001",
					SignName:        "Uptime Kuma",
					TemplateCode:    "SMS_1234567890",
				},
			},
			updateFunc: func(n notification.Notification) {
				aliyunsms := n.(*notification.AliyunSMS)
				aliyunsms.Name = "Test AliyunSMS Updated"
				aliyunsms.PhoneNumber = "8613800000002"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.AliyunSMS)
				require.True(t, ok)
				var aliyunsms notification.AliyunSMS
				err := actual.As(&aliyunsms)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = aliyunsms.UserID
				require.EqualExportedValues(t, exp, aliyunsms)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var aliyunsms notification.AliyunSMS
				err := base.As(&aliyunsms)
				require.NoError(t, err)
				return &aliyunsms
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.AliyunSMS)
				require.True(t, ok)
				var aliyunsms notification.AliyunSMS
				err := actual.As(&aliyunsms)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, aliyunsms)
			},
		},
		{
			name:         "Apprise",
			expectedType: "apprise",
			create: notification.Apprise{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Apprise Created",
				},
				AppriseDetails: notification.AppriseDetails{
					AppriseURL: "json://localhost:8080",
					Title:      "Uptime Kuma Alert",
				},
			},
			updateFunc: func(n notification.Notification) {
				apprise := n.(*notification.Apprise)
				apprise.Name = "Test Apprise Updated"
				apprise.Title = "Updated Alert"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Apprise)
				require.True(t, ok)
				var apprise notification.Apprise
				err := actual.As(&apprise)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = apprise.UserID
				require.EqualExportedValues(t, exp, apprise)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var apprise notification.Apprise
				err := base.As(&apprise)
				require.NoError(t, err)
				return &apprise
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Apprise)
				require.True(t, ok)
				var apprise notification.Apprise
				err := actual.As(&apprise)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, apprise)
			},
		},
		{
			name:         "GoogleChat",
			expectedType: "GoogleChat",
			create: notification.GoogleChat{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Google Chat Created",
				},
				GoogleChatDetails: notification.GoogleChatDetails{
					WebhookURL:  "https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test",
					UseTemplate: false,
					Template:    "",
				},
			},
			updateFunc: func(n notification.Notification) {
				googlechat := n.(*notification.GoogleChat)
				googlechat.Name = "Test Google Chat Updated"
				googlechat.UseTemplate = true
				googlechat.Template = "Updated Template"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.GoogleChat)
				require.True(t, ok)
				var googlechat notification.GoogleChat
				err := actual.As(&googlechat)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = googlechat.UserID
				require.EqualExportedValues(t, exp, googlechat)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var googlechat notification.GoogleChat
				err := base.As(&googlechat)
				require.NoError(t, err)
				return &googlechat
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.GoogleChat)
				require.True(t, ok)
				var googlechat notification.GoogleChat
				err := actual.As(&googlechat)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, googlechat)
			},
		},
		{
			name:         "Bark",
			expectedType: "bark",
			create: notification.Bark{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Bark Created",
				},
				BarkDetails: notification.BarkDetails{
					Endpoint:   "https://bark.example.com",
					Group:      "Monitoring",
					Sound:      "alarm",
					APIVersion: "v1",
				},
			},
			updateFunc: func(n notification.Notification) {
				bark := n.(*notification.Bark)
				bark.Name = "Test Bark Updated"
				bark.Sound = "telegraph"
				bark.APIVersion = "v2"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Bark)
				require.True(t, ok)
				var bark notification.Bark
				err := actual.As(&bark)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = bark.UserID
				require.EqualExportedValues(t, exp, bark)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var bark notification.Bark
				err := base.As(&bark)
				require.NoError(t, err)
				return &bark
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Bark)
				require.True(t, ok)
				var bark notification.Bark
				err := actual.As(&bark)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, bark)
			},
		},
		{
			name:         "Bitrix24",
			expectedType: "Bitrix24",
			create: notification.Bitrix24{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Bitrix24 Created",
				},
				Bitrix24Details: notification.Bitrix24Details{
					WebhookURL:         "https://bitrix24.example.com/rest/1/webhook/",
					NotificationUserID: "user123",
				},
			},
			updateFunc: func(n notification.Notification) {
				bitrix24 := n.(*notification.Bitrix24)
				bitrix24.Name = "Test Bitrix24 Updated"
				bitrix24.NotificationUserID = "admin"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Bitrix24)
				require.True(t, ok)
				var bitrix24 notification.Bitrix24
				err := actual.As(&bitrix24)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = bitrix24.UserID
				require.EqualExportedValues(t, exp, bitrix24)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var bitrix24 notification.Bitrix24
				err := base.As(&bitrix24)
				require.NoError(t, err)
				return &bitrix24
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Bitrix24)
				require.True(t, ok)
				var bitrix24 notification.Bitrix24
				err := actual.As(&bitrix24)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, bitrix24)
			},
		},
		{
			name:         "Brevo",
			expectedType: "brevo",
			create: notification.Brevo{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Brevo Created",
				},
				BrevoDetails: notification.BrevoDetails{
					APIKey:    "test-api-key",
					ToEmail:   "recipient@example.com",
					FromEmail: "sender@example.com",
					FromName:  "Uptime Kuma",
					Subject:   "Alert Notification",
				},
			},
			updateFunc: func(n notification.Notification) {
				brevo := n.(*notification.Brevo)
				brevo.Name = "Test Brevo Updated"
				brevo.ToEmail = "updated@example.com"
				brevo.FromName = "Updated System"
				brevo.Subject = "Updated Alert"
				brevo.CCEmail = "cc@example.com"
				brevo.BCCEmail = "bcc@example.com"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Brevo)
				require.True(t, ok)
				var brevo notification.Brevo
				err := actual.As(&brevo)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = brevo.UserID
				require.EqualExportedValues(t, exp, brevo)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var brevo notification.Brevo
				err := base.As(&brevo)
				require.NoError(t, err)
				return &brevo
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Brevo)
				require.True(t, ok)
				var brevo notification.Brevo
				err := actual.As(&brevo)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, brevo)
			},
		},
		{
			name:         "CallMeBot",
			expectedType: "CallMeBot",
			create: notification.CallMeBot{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test CallMeBot Created",
				},
				CallMeBotDetails: notification.CallMeBotDetails{
					Endpoint: "https://api.callmebot.com/start",
				},
			},
			updateFunc: func(n notification.Notification) {
				callmebot := n.(*notification.CallMeBot)
				callmebot.Name = "Test CallMeBot Updated"
				callmebot.Endpoint = "https://custom.callmebot.endpoint.com/start"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.CallMeBot)
				require.True(t, ok)
				var callmebot notification.CallMeBot
				err := actual.As(&callmebot)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = callmebot.UserID
				require.EqualExportedValues(t, exp, callmebot)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var callmebot notification.CallMeBot
				err := base.As(&callmebot)
				require.NoError(t, err)
				return &callmebot
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.CallMeBot)
				require.True(t, ok)
				var callmebot notification.CallMeBot
				err := actual.As(&callmebot)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, callmebot)
			},
		},
		{
			name:         "Cellsynt",
			expectedType: "Cellsynt",
			create: notification.Cellsynt{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Cellsynt Created",
				},
				CellsyntDetails: notification.CellsyntDetails{
					Login:          "testuser",
					Password:       "testpass",
					Destination:    "46701234567",
					Originator:     "Uptime",
					OriginatorType: "Numeric",
					AllowLongSMS:   false,
				},
			},
			updateFunc: func(n notification.Notification) {
				cellsynt := n.(*notification.Cellsynt)
				cellsynt.Name = "Test Cellsynt Updated"
				cellsynt.Destination = "46709876543"
				cellsynt.Originator = "Updated"
				cellsynt.OriginatorType = "Alphanumeric"
				cellsynt.AllowLongSMS = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Cellsynt)
				require.True(t, ok)
				var cellsynt notification.Cellsynt
				err := actual.As(&cellsynt)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = cellsynt.UserID
				require.EqualExportedValues(t, exp, cellsynt)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var cellsynt notification.Cellsynt
				err := base.As(&cellsynt)
				require.NoError(t, err)
				return &cellsynt
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Cellsynt)
				require.True(t, ok)
				var cellsynt notification.Cellsynt
				err := actual.As(&cellsynt)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, cellsynt)
			},
		},
		{
			name:         "ClickSendSMS",
			expectedType: "clicksendsms",
			create: notification.ClickSendSMS{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test ClickSend Created",
				},
				ClickSendSMSDetails: notification.ClickSendSMSDetails{
					Login:      "testuser",
					Password:   "apikey123",
					ToNumber:   "61412345678",
					SenderName: "Uptime",
				},
			},
			updateFunc: func(n notification.Notification) {
				clicksendsms := n.(*notification.ClickSendSMS)
				clicksendsms.Name = "Test ClickSend Updated"
				clicksendsms.ToNumber = "61487654321"
				clicksendsms.SenderName = "Updated Monitor"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.ClickSendSMS)
				require.True(t, ok)
				var clicksendsms notification.ClickSendSMS
				err := actual.As(&clicksendsms)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = clicksendsms.UserID
				require.EqualExportedValues(t, exp, clicksendsms)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var clicksendsms notification.ClickSendSMS
				err := base.As(&clicksendsms)
				require.NoError(t, err)
				return &clicksendsms
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.ClickSendSMS)
				require.True(t, ok)
				var clicksendsms notification.ClickSendSMS
				err := actual.As(&clicksendsms)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, clicksendsms)
			},
		},
		{
			name:         "Evolution",
			expectedType: "EvolutionApi",
			create: notification.Evolution{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Evolution Created",
				},
				EvolutionDetails: notification.EvolutionDetails{
					APIURL:       "https://evolapicloud.com",
					InstanceName: "myinstance",
					AuthToken:    "token123",
					Recipient:    "5511999999999",
				},
			},
			updateFunc: func(n notification.Notification) {
				evolution := n.(*notification.Evolution)
				evolution.Name = "Test Evolution Updated"
				evolution.APIURL = "https://custom.api.com"
				evolution.InstanceName = "newinstance"
				evolution.Recipient = "5521987654321"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Evolution)
				require.True(t, ok)
				var evolution notification.Evolution
				err := actual.As(&evolution)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = evolution.UserID
				require.EqualExportedValues(t, exp, evolution)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var evolution notification.Evolution
				err := base.As(&evolution)
				require.NoError(t, err)
				return &evolution
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Evolution)
				require.True(t, ok)
				var evolution notification.Evolution
				err := actual.As(&evolution)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, evolution)
			},
		},
		{
			name:         "FlashDuty",
			expectedType: "FlashDuty",
			create: notification.FlashDuty{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test FlashDuty Created",
				},
				FlashDutyDetails: notification.FlashDutyDetails{
					IntegrationKey: "test_key_123",
					Severity:       "Critical",
				},
			},
			updateFunc: func(n notification.Notification) {
				flashduty := n.(*notification.FlashDuty)
				flashduty.Name = "Test FlashDuty Updated"
				flashduty.Severity = "Warning"
				flashduty.IntegrationKey = "updated_key_456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.FlashDuty)
				require.True(t, ok)
				var flashduty notification.FlashDuty
				err := actual.As(&flashduty)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = flashduty.UserID
				require.EqualExportedValues(t, exp, flashduty)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var flashduty notification.FlashDuty
				err := base.As(&flashduty)
				require.NoError(t, err)
				return &flashduty
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.FlashDuty)
				require.True(t, ok)
				var flashduty notification.FlashDuty
				err := actual.As(&flashduty)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, flashduty)
			},
		},
		{
			name:         "GoAlert",
			expectedType: "GoAlert",
			create: notification.GoAlert{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test GoAlert Created",
				},
				GoAlertDetails: notification.GoAlertDetails{
					BaseURL: "https://goalert.example.com",
					Token:   "test-token-123",
				},
			},
			updateFunc: func(n notification.Notification) {
				goalert := n.(*notification.GoAlert)
				goalert.Name = "Test GoAlert Updated"
				goalert.Token = "updated-token-456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.GoAlert)
				require.True(t, ok)
				var goalert notification.GoAlert
				err := actual.As(&goalert)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = goalert.UserID
				require.EqualExportedValues(t, exp, goalert)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var goalert notification.GoAlert
				err := base.As(&goalert)
				require.NoError(t, err)
				return &goalert
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.GoAlert)
				require.True(t, ok)
				var goalert notification.GoAlert
				err := actual.As(&goalert)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, goalert)
			},
		},
		{
			name:         "Gorush",
			expectedType: "gorush",
			create: notification.Gorush{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Gorush Created",
				},
				GorushDetails: notification.GorushDetails{
					ServerURL:   "https://gorush.example.com",
					DeviceToken: "test-device-token",
					Platform:    "ios",
					Title:       "Uptime Alert",
					Priority:    "high",
					Retry:       3,
					Topic:       "com.example.app",
				},
			},
			updateFunc: func(n notification.Notification) {
				gorush := n.(*notification.Gorush)
				gorush.Name = "Test Gorush Updated"
				gorush.Priority = "critical"
				gorush.Retry = 5
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Gorush)
				require.True(t, ok)
				var gorush notification.Gorush
				err := actual.As(&gorush)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = gorush.UserID
				require.EqualExportedValues(t, exp, gorush)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var gorush notification.Gorush
				err := base.As(&gorush)
				require.NoError(t, err)
				return &gorush
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Gorush)
				require.True(t, ok)
				var gorush notification.Gorush
				err := actual.As(&gorush)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, gorush)
			},
		},
		{
			name:         "GTXMessaging",
			expectedType: "gtxmessaging",
			create: notification.GTXMessaging{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test GTX Messaging Created",
				},
				GTXMessagingDetails: notification.GTXMessagingDetails{
					APIKey: "test-api-key",
					From:   "Uptime",
					To:     "+46701234567",
				},
			},
			updateFunc: func(n notification.Notification) {
				gtx := n.(*notification.GTXMessaging)
				gtx.Name = "Test GTX Messaging Updated"
				gtx.From = "Monitor"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.GTXMessaging)
				require.True(t, ok)
				var gtx notification.GTXMessaging
				err := actual.As(&gtx)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = gtx.UserID
				require.EqualExportedValues(t, exp, gtx)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var gtx notification.GTXMessaging
				err := base.As(&gtx)
				require.NoError(t, err)
				return &gtx
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.GTXMessaging)
				require.True(t, ok)
				var gtx notification.GTXMessaging
				err := actual.As(&gtx)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, gtx)
			},
		},
		{
			name:         "HeiiOnCall",
			expectedType: "HeiiOnCall",
			create: notification.HeiiOnCall{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Heii On-Call Created",
				},
				HeiiOnCallDetails: notification.HeiiOnCallDetails{
					APIKey:    "test-api-key",
					TriggerID: "test-trigger-id",
				},
			},
			updateFunc: func(n notification.Notification) {
				heii := n.(*notification.HeiiOnCall)
				heii.Name = "Test Heii On-Call Updated"
				heii.APIKey = "updated-api-key"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.HeiiOnCall)
				require.True(t, ok)
				var heii notification.HeiiOnCall
				err := actual.As(&heii)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = heii.UserID
				require.EqualExportedValues(t, exp, heii)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var heii notification.HeiiOnCall
				err := base.As(&heii)
				require.NoError(t, err)
				return &heii
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.HeiiOnCall)
				require.True(t, ok)
				var heii notification.HeiiOnCall
				err := actual.As(&heii)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, heii)
			},
		},
		{
			name:         "Keep",
			expectedType: "Keep",
			create: notification.Keep{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Keep Created",
				},
				KeepDetails: notification.KeepDetails{
					WebhookURL: "https://keep.example.com/webhook",
					APIKey:     "test-api-key",
				},
			},
			updateFunc: func(n notification.Notification) {
				keep := n.(*notification.Keep)
				keep.Name = "Test Keep Updated"
				keep.APIKey = "updated-api-key"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Keep)
				require.True(t, ok)
				var keep notification.Keep
				err := actual.As(&keep)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = keep.UserID
				require.EqualExportedValues(t, exp, keep)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var keep notification.Keep
				err := base.As(&keep)
				require.NoError(t, err)
				return &keep
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Keep)
				require.True(t, ok)
				var keep notification.Keep
				err := actual.As(&keep)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, keep)
			},
		},
		{
			name:         "Kook",
			expectedType: "Kook",
			create: notification.Kook{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Kook Created",
				},
				KookDetails: notification.KookDetails{
					BotToken: "test-bot-token",
					GuildID:  "test-guild-id",
				},
			},
			updateFunc: func(n notification.Notification) {
				kook := n.(*notification.Kook)
				kook.Name = "Test Kook Updated"
				kook.BotToken = "updated-bot-token"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Kook)
				require.True(t, ok)
				var kook notification.Kook
				err := actual.As(&kook)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = kook.UserID
				require.EqualExportedValues(t, exp, kook)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var kook notification.Kook
				err := base.As(&kook)
				require.NoError(t, err)
				return &kook
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Kook)
				require.True(t, ok)
				var kook notification.Kook
				err := actual.As(&kook)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, kook)
			},
		},
		{
			name:         "FreeMobile",
			expectedType: "FreeMobile",
			create: notification.FreeMobile{
				Base: notification.Base{
					ApplyExisting: false,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Free Mobile Created",
				},
				FreeMobileDetails: notification.FreeMobileDetails{
					User: "12345678",
					Pass: "abcdef123456",
				},
			},
			updateFunc: func(n notification.Notification) {
				freemobile := n.(*notification.FreeMobile)
				freemobile.Name = "Test Free Mobile Updated"
				freemobile.Pass = "updated123456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.FreeMobile)
				require.True(t, ok)
				var freemobile notification.FreeMobile
				err := actual.As(&freemobile)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = freemobile.UserID
				require.EqualExportedValues(t, exp, freemobile)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var freemobile notification.FreeMobile
				err := base.As(&freemobile)
				require.NoError(t, err)
				return &freemobile
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.FreeMobile)
				require.True(t, ok)
				var freemobile notification.FreeMobile
				err := actual.As(&freemobile)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, freemobile)
			},
		},
		{
			name:         "Line",
			expectedType: "line",
			create: notification.Line{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test LINE Created",
				},
				LineDetails: notification.LineDetails{
					ChannelAccessToken: "channel-access-token-test",
					UserID:             "U1234567890abcdef1234567890abcdef",
				},
			},
			updateFunc: func(n notification.Notification) {
				line := n.(*notification.Line)
				line.Name = "Test LINE Updated"
				line.ChannelAccessToken = "updated-token-123"
				line.LineDetails.UserID = "U9876543210fedcba9876543210fedcba"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Line)
				require.True(t, ok)
				var line notification.Line
				err := actual.As(&line)
				require.NoError(t, err)
				exp.ID = id
				exp.Base.UserID = line.Base.UserID
				require.EqualExportedValues(t, exp, line)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var line notification.Line
				err := base.As(&line)
				require.NoError(t, err)
				return &line
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Line)
				require.True(t, ok)
				var line notification.Line
				err := actual.As(&line)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, line)
			},
		},
		{
			name:         "LineNotify",
			expectedType: "LineNotify",
			create: notification.LineNotify{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test LINE Notify Created",
				},
				LineNotifyDetails: notification.LineNotifyDetails{
					AccessToken: "access-token-test",
				},
			},
			updateFunc: func(n notification.Notification) {
				linenotify := n.(*notification.LineNotify)
				linenotify.Name = "Test LINE Notify Updated"
				linenotify.AccessToken = "updated-token-123"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.LineNotify)
				require.True(t, ok)
				var linenotify notification.LineNotify
				err := actual.As(&linenotify)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = linenotify.UserID
				require.EqualExportedValues(t, exp, linenotify)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var linenotify notification.LineNotify
				err := base.As(&linenotify)
				require.NoError(t, err)
				return &linenotify
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.LineNotify)
				require.True(t, ok)
				var linenotify notification.LineNotify
				err := actual.As(&linenotify)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, linenotify)
			},
		},
		{
			name:         "LunaSea",
			expectedType: "lunasea",
			create: notification.LunaSea{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test LunaSea Created",
				},
				LunaSeaDetails: notification.LunaSeaDetails{
					Target:        "user",
					LunaSeaUserID: "user-test-123",
					Device:        "",
				},
			},
			updateFunc: func(n notification.Notification) {
				lunasea := n.(*notification.LunaSea)
				lunasea.Name = "Test LunaSea Updated"
				lunasea.Target = "device"
				lunasea.LunaSeaUserID = ""
				lunasea.Device = "device-456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.LunaSea)
				require.True(t, ok)
				var lunasea notification.LunaSea
				err := actual.As(&lunasea)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = lunasea.UserID
				require.EqualExportedValues(t, exp, lunasea)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var lunasea notification.LunaSea
				err := base.As(&lunasea)
				require.NoError(t, err)
				return &lunasea
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.LunaSea)
				require.True(t, ok)
				var lunasea notification.LunaSea
				err := actual.As(&lunasea)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, lunasea)
			},
		},
		{
			name:         "NextcloudTalk",
			expectedType: "NextcloudTalk",
			create: notification.NextcloudTalk{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Nextcloud Talk Created",
				},
				NextcloudTalkDetails: notification.NextcloudTalkDetails{
					Host:              "https://nextcloud.example.com",
					ConversationToken: "token-test",
					BotSecret:         "secret-test",
					SendSilentUp:      true,
					SendSilentDown:    false,
				},
			},
			updateFunc: func(n notification.Notification) {
				nextcloudtalk := n.(*notification.NextcloudTalk)
				nextcloudtalk.Name = "Test Nextcloud Talk Updated"
				nextcloudtalk.Host = "https://updated.example.com"
				nextcloudtalk.ConversationToken = "token-updated"
				nextcloudtalk.SendSilentUp = false
				nextcloudtalk.SendSilentDown = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.NextcloudTalk)
				require.True(t, ok)
				var nextcloudtalk notification.NextcloudTalk
				err := actual.As(&nextcloudtalk)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = nextcloudtalk.UserID
				require.EqualExportedValues(t, exp, nextcloudtalk)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var nextcloudtalk notification.NextcloudTalk
				err := base.As(&nextcloudtalk)
				require.NoError(t, err)
				return &nextcloudtalk
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.NextcloudTalk)
				require.True(t, ok)
				var nextcloudtalk notification.NextcloudTalk
				err := actual.As(&nextcloudtalk)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, nextcloudtalk)
			},
		},
		{
			name:         "Nostr",
			expectedType: "nostr",
			create: notification.Nostr{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Nostr Created",
				},
				NostrDetails: notification.NostrDetails{
					Sender:     "nsec1test-sender",
					Recipients: "npub1recipient1\nnpub1recipient2",
					Relays:     "wss://relay1.example.com\nwss://relay2.example.com",
				},
			},
			updateFunc: func(n notification.Notification) {
				nostr := n.(*notification.Nostr)
				nostr.Name = "Test Nostr Updated"
				nostr.Sender = "nsec1updated-sender"
				nostr.Recipients = "npub1updated-recipient"
				nostr.Relays = "wss://updated-relay.example.com"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Nostr)
				require.True(t, ok)
				var nostr notification.Nostr
				err := actual.As(&nostr)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = nostr.UserID
				require.EqualExportedValues(t, exp, nostr)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var nostr notification.Nostr
				err := base.As(&nostr)
				require.NoError(t, err)
				return &nostr
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Nostr)
				require.True(t, ok)
				var nostr notification.Nostr
				err := actual.As(&nostr)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, nostr)
			},
		},
		{
			name:         "OneBot",
			expectedType: "OneBot",
			create: notification.OneBot{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test OneBot Created",
				},
				OneBotDetails: notification.OneBotDetails{
					HTTPAddr:    "http://localhost:5700",
					AccessToken: "test-token",
					MsgType:     "group",
					ReceiverID:  "123456789",
				},
			},
			updateFunc: func(n notification.Notification) {
				onebot := n.(*notification.OneBot)
				onebot.Name = "Test OneBot Updated"
				onebot.AccessToken = "updated-token"
				onebot.MsgType = "private"
				onebot.ReceiverID = "987654321"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.OneBot)
				require.True(t, ok)
				var onebot notification.OneBot
				err := actual.As(&onebot)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = onebot.UserID
				require.EqualExportedValues(t, exp, onebot)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var onebot notification.OneBot
				err := base.As(&onebot)
				require.NoError(t, err)
				return &onebot
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.OneBot)
				require.True(t, ok)
				var onebot notification.OneBot
				err := actual.As(&onebot)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, onebot)
			},
		},
		{
			name:         "Octopush",
			expectedType: "octopush",
			create: notification.Octopush{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Octopush V2 Created",
				},
				OctopushDetails: notification.OctopushDetails{
					Version:     "2",
					APIKey:      "test-api-key",
					Login:       "testuser",
					PhoneNumber: "+33612345678",
					SMSType:     "sms_premium",
					SenderName:  "AlertBot",
				},
			},
			updateFunc: func(n notification.Notification) {
				octopush := n.(*notification.Octopush)
				octopush.Name = "Test Octopush V2 Updated"
				octopush.APIKey = "updated-api-key"
				octopush.PhoneNumber = "+33698765432"
				octopush.SenderName = "UpdatedBot"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Octopush)
				require.True(t, ok)
				var octopush notification.Octopush
				err := actual.As(&octopush)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = octopush.UserID
				require.EqualExportedValues(t, exp, octopush)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var octopush notification.Octopush
				err := base.As(&octopush)
				require.NoError(t, err)
				return &octopush
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Octopush)
				require.True(t, ok)
				var octopush notification.Octopush
				err := actual.As(&octopush)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, octopush)
			},
		},
		{
			name:         "OneChat",
			expectedType: "OneChat",
			create: notification.OneChat{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test OneChat Created",
				},
				OneChatDetails: notification.OneChatDetails{
					AccessToken: "test-token",
					ReceiverID:  "user123",
					BotID:       "bot456",
				},
			},
			updateFunc: func(n notification.Notification) {
				onechat := n.(*notification.OneChat)
				onechat.Name = "Test OneChat Updated"
				onechat.AccessToken = "updated-token"
				onechat.ReceiverID = "group789"
				onechat.BotID = "botgroup"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.OneChat)
				require.True(t, ok)
				var onechat notification.OneChat
				err := actual.As(&onechat)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = onechat.UserID
				require.EqualExportedValues(t, exp, onechat)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var onechat notification.OneChat
				err := base.As(&onechat)
				require.NoError(t, err)
				return &onechat
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.OneChat)
				require.True(t, ok)
				var onechat notification.OneChat
				err := actual.As(&onechat)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, onechat)
			},
		},
		{
			name:         "Notifery",
			expectedType: "notifery",
			create: notification.Notifery{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Notifery Created",
				},
				NotiferyDetails: notification.NotiferyDetails{
					APIKey: "test-api-key",
					Title:  "Uptime Alert",
					Group:  "monitoring",
				},
			},
			updateFunc: func(n notification.Notification) {
				notifery := n.(*notification.Notifery)
				notifery.Name = "Test Notifery Updated"
				notifery.APIKey = "updated-api-key"
				notifery.Title = "Critical Alert"
				notifery.Group = "critical"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Notifery)
				require.True(t, ok)
				var notifery notification.Notifery
				err := actual.As(&notifery)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = notifery.UserID
				require.EqualExportedValues(t, exp, notifery)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var notifery notification.Notifery
				err := base.As(&notifery)
				require.NoError(t, err)
				return &notifery
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Notifery)
				require.True(t, ok)
				var notifery notification.Notifery
				err := actual.As(&notifery)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, notifery)
			},
		},
		{
			name:         "OneSender",
			expectedType: "onesender",
			create: notification.OneSender{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test OneSender Created",
				},
				OneSenderDetails: notification.OneSenderDetails{
					URL:          "https://api.onesender.com/send",
					Token:        "test-token",
					Receiver:     "5511999999999",
					TypeReceiver: "private",
				},
			},
			updateFunc: func(n notification.Notification) {
				onesender := n.(*notification.OneSender)
				onesender.Name = "Test OneSender Updated"
				onesender.Token = "updated-token"
				onesender.Receiver = "120363123456789-1234567890"
				onesender.TypeReceiver = "group"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.OneSender)
				require.True(t, ok)
				var onesender notification.OneSender
				err := actual.As(&onesender)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = onesender.UserID
				require.EqualExportedValues(t, exp, onesender)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var onesender notification.OneSender
				err := base.As(&onesender)
				require.NoError(t, err)
				return &onesender
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.OneSender)
				require.True(t, ok)
				var onesender notification.OneSender
				err := actual.As(&onesender)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, onesender)
			},
		},
		{
			name:         "PagerTree",
			expectedType: "PagerTree",
			create: notification.PagerTree{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test PagerTree Created",
				},
				PagerTreeDetails: notification.PagerTreeDetails{
					IntegrationURL: "https://api.pagertree.com/api/v2/events",
					Urgency:        "high",
					AutoResolve:    "resolve",
				},
			},
			updateFunc: func(n notification.Notification) {
				pagertree := n.(*notification.PagerTree)
				pagertree.Name = "Test PagerTree Updated"
				pagertree.Urgency = "medium"
				pagertree.AutoResolve = ""
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.PagerTree)
				require.True(t, ok)
				var pagertree notification.PagerTree
				err := actual.As(&pagertree)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pagertree.UserID
				require.EqualExportedValues(t, exp, pagertree)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pagertree notification.PagerTree
				err := base.As(&pagertree)
				require.NoError(t, err)
				return &pagertree
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.PagerTree)
				require.True(t, ok)
				var pagertree notification.PagerTree
				err := actual.As(&pagertree)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pagertree)
			},
		},
		{
			name:         "PromoSMS",
			expectedType: "promosms",
			create: notification.PromoSMS{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test PromoSMS Created",
				},
				PromoSMSDetails: notification.PromoSMSDetails{
					Login:        "user@example.com",
					Password:     "password123",
					PhoneNumber:  "+48123456789",
					SenderName:   "UptimeKuma",
					SMSType:      "1",
					AllowLongSMS: true,
				},
			},
			updateFunc: func(n notification.Notification) {
				promosms := n.(*notification.PromoSMS)
				promosms.Name = "Test PromoSMS Updated"
				promosms.Password = "newpassword"
				promosms.PhoneNumber = "+48987654321"
				promosms.AllowLongSMS = false
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.PromoSMS)
				require.True(t, ok)
				var promosms notification.PromoSMS
				err := actual.As(&promosms)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = promosms.UserID
				require.EqualExportedValues(t, exp, promosms)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var promosms notification.PromoSMS
				err := base.As(&promosms)
				require.NoError(t, err)
				return &promosms
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.PromoSMS)
				require.True(t, ok)
				var promosms notification.PromoSMS
				err := actual.As(&promosms)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, promosms)
			},
		},
		{
			name:         "Pumble",
			expectedType: "Pumble",
			create: notification.Pumble{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Pumble Created",
				},
				PumbleDetails: notification.PumbleDetails{
					WebhookURL: "https://pumble.com/webhook/test",
				},
			},
			updateFunc: func(n notification.Notification) {
				pumble := n.(*notification.Pumble)
				pumble.Name = "Test Pumble Updated"
				pumble.WebhookURL = "https://pumble.com/webhook/updated"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Pumble)
				require.True(t, ok)
				var pumble notification.Pumble
				err := actual.As(&pumble)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pumble.UserID
				require.EqualExportedValues(t, exp, pumble)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pumble notification.Pumble
				err := base.As(&pumble)
				require.NoError(t, err)
				return &pumble
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Pumble)
				require.True(t, ok)
				var pumble notification.Pumble
				err := actual.As(&pumble)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pumble)
			},
		},
		{
			name:         "PushDeer",
			expectedType: "PushDeer",
			create: notification.PushDeer{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test PushDeer Created",
				},
				PushDeerDetails: notification.PushDeerDetails{
					Key:    "PDxxxxxxxxxxxxxx",
					Server: "https://api2.pushdeer.com",
				},
			},
			updateFunc: func(n notification.Notification) {
				pushdeer := n.(*notification.PushDeer)
				pushdeer.Name = "Test PushDeer Updated"
				pushdeer.Key = "PDyyyyyyyyyyyyyyyy"
				pushdeer.Server = "https://custom.pushdeer.com"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.PushDeer)
				require.True(t, ok)
				var pushdeer notification.PushDeer
				err := actual.As(&pushdeer)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pushdeer.UserID
				require.EqualExportedValues(t, exp, pushdeer)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pushdeer notification.PushDeer
				err := base.As(&pushdeer)
				require.NoError(t, err)
				return &pushdeer
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.PushDeer)
				require.True(t, ok)
				var pushdeer notification.PushDeer
				err := actual.As(&pushdeer)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pushdeer)
			},
		},
		{
			name:         "PushPlus",
			expectedType: "PushPlus",
			create: notification.PushPlus{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test PushPlus Created",
				},
				PushPlusDetails: notification.PushPlusDetails{
					SendKey: "test_send_key_xxxxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				pushplus := n.(*notification.PushPlus)
				pushplus.Name = "Test PushPlus Updated"
				pushplus.SendKey = "updated_send_key_yyyyy"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.PushPlus)
				require.True(t, ok)
				var pushplus notification.PushPlus
				err := actual.As(&pushplus)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pushplus.UserID
				require.EqualExportedValues(t, exp, pushplus)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pushplus notification.PushPlus
				err := base.As(&pushplus)
				require.NoError(t, err)
				return &pushplus
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.PushPlus)
				require.True(t, ok)
				var pushplus notification.PushPlus
				err := actual.As(&pushplus)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pushplus)
			},
		},
		{
			name:         "Pushy",
			expectedType: "pushy",
			create: notification.Pushy{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Pushy Created",
				},
				PushyDetails: notification.PushyDetails{
					APIKey: "test_api_key_xxxxx",
					Token:  "test_device_token_xxxxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				pushy := n.(*notification.Pushy)
				pushy.Name = "Test Pushy Updated"
				pushy.APIKey = "updated_api_key_yyyyy"
				pushy.Token = "updated_device_token_yyyyy"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Pushy)
				require.True(t, ok)
				var pushy notification.Pushy
				err := actual.As(&pushy)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = pushy.UserID
				require.EqualExportedValues(t, exp, pushy)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var pushy notification.Pushy
				err := base.As(&pushy)
				require.NoError(t, err)
				return &pushy
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Pushy)
				require.True(t, ok)
				var pushy notification.Pushy
				err := actual.As(&pushy)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, pushy)
			},
		},
		{
			name:         "SendGrid",
			expectedType: "SendGrid",
			create: notification.SendGrid{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SendGrid Created",
				},
				SendGridDetails: notification.SendGridDetails{
					APIKey:    "SG.test_api_key_xxxxx",
					ToEmail:   "test@example.com",
					FromEmail: "sender@example.com",
					Subject:   "Test Subject",
					CcEmail:   "cc@example.com",
					BccEmail:  "bcc@example.com",
				},
			},
			updateFunc: func(n notification.Notification) {
				sendgrid := n.(*notification.SendGrid)
				sendgrid.Name = "Test SendGrid Updated"
				sendgrid.APIKey = "SG.updated_api_key_yyyyy"
				sendgrid.ToEmail = "updated@example.com"
				sendgrid.Subject = "Updated Subject"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SendGrid)
				require.True(t, ok)
				var sendgrid notification.SendGrid
				err := actual.As(&sendgrid)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = sendgrid.UserID
				require.EqualExportedValues(t, exp, sendgrid)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var sendgrid notification.SendGrid
				err := base.As(&sendgrid)
				require.NoError(t, err)
				return &sendgrid
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SendGrid)
				require.True(t, ok)
				var sendgrid notification.SendGrid
				err := actual.As(&sendgrid)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, sendgrid)
			},
		},
		{
			name:         "ServerChan",
			expectedType: "ServerChan",
			create: notification.ServerChan{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test ServerChan Created",
				},
				ServerChanDetails: notification.ServerChanDetails{
					SendKey: "SCT123456789abcdefghijklmnopqrst",
				},
			},
			updateFunc: func(n notification.Notification) {
				serverchan := n.(*notification.ServerChan)
				serverchan.Name = "Test ServerChan Updated"
				serverchan.SendKey = "SCT000000000000000000000000000000"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.ServerChan)
				require.True(t, ok)
				var serverchan notification.ServerChan
				err := actual.As(&serverchan)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = serverchan.UserID
				require.EqualExportedValues(t, exp, serverchan)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var serverchan notification.ServerChan
				err := base.As(&serverchan)
				require.NoError(t, err)
				return &serverchan
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.ServerChan)
				require.True(t, ok)
				var serverchan notification.ServerChan
				err := actual.As(&serverchan)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, serverchan)
			},
		},
		{
			name:         "SerwerSMS",
			expectedType: "serwersms",
			create: notification.SerwerSMS{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SerwerSMS Created",
				},
				SerwerSMSDetails: notification.SerwerSMSDetails{
					Username:    "testuser",
					Password:    "testpass",
					PhoneNumber: "48123456789",
					SenderName:  "Uptime",
				},
			},
			updateFunc: func(n notification.Notification) {
				serwersms := n.(*notification.SerwerSMS)
				serwersms.Name = "Test SerwerSMS Updated"
				serwersms.PhoneNumber = "48987654321"
				serwersms.SenderName = "UpdatedAlert"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SerwerSMS)
				require.True(t, ok)
				var serwersms notification.SerwerSMS
				err := actual.As(&serwersms)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = serwersms.UserID
				require.EqualExportedValues(t, exp, serwersms)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var serwersms notification.SerwerSMS
				err := base.As(&serwersms)
				require.NoError(t, err)
				return &serwersms
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SerwerSMS)
				require.True(t, ok)
				var serwersms notification.SerwerSMS
				err := actual.As(&serwersms)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, serwersms)
			},
		},
		{
			name:         "SevenIO",
			expectedType: "sevenio",
			create: notification.SevenIO{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SevenIO Created",
				},
				SevenIODetails: notification.SevenIODetails{
					APIKey: "test-api-key",
					Sender: "UptimeKuma",
					To:     "49123456789",
				},
			},
			updateFunc: func(n notification.Notification) {
				sevenio := n.(*notification.SevenIO)
				sevenio.Name = "Test SevenIO Updated"
				sevenio.To = "49987654321"
				sevenio.Sender = "UpdatedAlert"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SevenIO)
				require.True(t, ok)
				var sevenio notification.SevenIO
				err := actual.As(&sevenio)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = sevenio.UserID
				require.EqualExportedValues(t, exp, sevenio)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var sevenio notification.SevenIO
				err := base.As(&sevenio)
				require.NoError(t, err)
				return &sevenio
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SevenIO)
				require.True(t, ok)
				var sevenio notification.SevenIO
				err := actual.As(&sevenio)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, sevenio)
			},
		},
		{
			name:         "SIGNL4",
			expectedType: "SIGNL4",
			create: notification.SIGNL4{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SIGNL4 Created",
				},
				SIGNL4Details: notification.SIGNL4Details{
					WebhookURL: "https://connect.signl4.com/webhook/test-webhook",
				},
			},
			updateFunc: func(n notification.Notification) {
				signl4 := n.(*notification.SIGNL4)
				signl4.Name = "Test SIGNL4 Updated"
				signl4.WebhookURL = "https://connect.signl4.com/webhook/updated-webhook"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SIGNL4)
				require.True(t, ok)
				var signl4 notification.SIGNL4
				err := actual.As(&signl4)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = signl4.UserID
				require.EqualExportedValues(t, exp, signl4)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var signl4 notification.SIGNL4
				err := base.As(&signl4)
				require.NoError(t, err)
				return &signl4
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SIGNL4)
				require.True(t, ok)
				var signl4 notification.SIGNL4
				err := actual.As(&signl4)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, signl4)
			},
		},
		{
			name:         "SMSC",
			expectedType: "smsc",
			create: notification.SMSC{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SMSC Created",
				},
				SMSCDetails: notification.SMSCDetails{
					Login:      "testuser",
					Password:   "testpass",
					ToNumber:   "77123456789",
					SenderName: "Uptime",
					Translit:   "1",
				},
			},
			updateFunc: func(n notification.Notification) {
				smsc := n.(*notification.SMSC)
				smsc.Name = "Test SMSC Updated"
				smsc.SenderName = "Updated"
				smsc.Translit = "0"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMSC)
				require.True(t, ok)
				var smsc notification.SMSC
				err := actual.As(&smsc)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smsc.UserID
				require.EqualExportedValues(t, exp, smsc)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smsc notification.SMSC
				err := base.As(&smsc)
				require.NoError(t, err)
				return &smsc
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMSC)
				require.True(t, ok)
				var smsc notification.SMSC
				err := actual.As(&smsc)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smsc)
			},
		},
		{
			name:         "SMSEagle",
			expectedType: "SMSEagle",
			create: notification.SMSEagle{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SMSEagle Created",
				},
				SMSEagleDetails: notification.SMSEagleDetails{
					URL:           "https://smseagle.example.com",
					Token:         "test-token-123",
					RecipientType: "smseagle-to",
					Recipient:     "1234567890",
					MsgType:       "smseagle-sms",
					Priority:      1,
					Encoding:      false,
					APIType:       "smseagle-apiv1",
				},
			},
			updateFunc: func(n notification.Notification) {
				smseagle := n.(*notification.SMSEagle)
				smseagle.Name = "Test SMSEagle Updated"
				smseagle.Priority = 2
				smseagle.Encoding = true
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMSEagle)
				require.True(t, ok)
				var smseagle notification.SMSEagle
				err := actual.As(&smseagle)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smseagle.UserID
				require.EqualExportedValues(t, exp, smseagle)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smseagle notification.SMSEagle
				err := base.As(&smseagle)
				require.NoError(t, err)
				return &smseagle
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMSEagle)
				require.True(t, ok)
				var smseagle notification.SMSEagle
				err := actual.As(&smseagle)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smseagle)
			},
		},
		{
			name:         "SMSManager",
			expectedType: "SMSManager",
			create: notification.SMSManager{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SMSManager Created",
				},
				SMSManagerDetails: notification.SMSManagerDetails{
					APIKey:      "test-api-key",
					Numbers:     "420777123456",
					MessageType: "1",
				},
			},
			updateFunc: func(n notification.Notification) {
				smsmanager := n.(*notification.SMSManager)
				smsmanager.Name = "Test SMSManager Updated"
				smsmanager.Numbers = "420999888777"
				smsmanager.MessageType = "2"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMSManager)
				require.True(t, ok)
				var smsmanager notification.SMSManager
				err := actual.As(&smsmanager)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smsmanager.UserID
				require.EqualExportedValues(t, exp, smsmanager)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smsmanager notification.SMSManager
				err := base.As(&smsmanager)
				require.NoError(t, err)
				return &smsmanager
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMSManager)
				require.True(t, ok)
				var smsmanager notification.SMSManager
				err := actual.As(&smsmanager)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smsmanager)
			},
		},
		{
			name:         "SMSPartner",
			expectedType: "SMSPartner",
			create: notification.SMSPartner{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SMSPartner Created",
				},
				SMSPartnerDetails: notification.SMSPartnerDetails{
					APIKey:      "test-api-key",
					PhoneNumber: "33612345678",
					SenderName:  "Uptime",
				},
			},
			updateFunc: func(n notification.Notification) {
				smspartner := n.(*notification.SMSPartner)
				smspartner.Name = "Test SMSPartner Updated"
				smspartner.PhoneNumber = "33687654321"
				smspartner.SenderName = "Monitor"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMSPartner)
				require.True(t, ok)
				var smspartner notification.SMSPartner
				err := actual.As(&smspartner)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smspartner.UserID
				require.EqualExportedValues(t, exp, smspartner)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smspartner notification.SMSPartner
				err := base.As(&smspartner)
				require.NoError(t, err)
				return &smspartner
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMSPartner)
				require.True(t, ok)
				var smspartner notification.SMSPartner
				err := actual.As(&smspartner)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smspartner)
			},
		},
		{
			name:         "SMSPlanet",
			expectedType: "SMSPlanet",
			create: notification.SMSPlanet{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SMS Planet Created",
				},
				SMSPlanetDetails: notification.SMSPlanetDetails{
					APIToken:     "test-token-123",
					PhoneNumbers: "48123456789",
					SenderName:   "Uptime Kuma",
				},
			},
			updateFunc: func(n notification.Notification) {
				smsplanet := n.(*notification.SMSPlanet)
				smsplanet.Name = "Test SMS Planet Updated"
				smsplanet.PhoneNumbers = "48987654321"
				smsplanet.SenderName = "Monitor"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SMSPlanet)
				require.True(t, ok)
				var smsplanet notification.SMSPlanet
				err := actual.As(&smsplanet)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = smsplanet.UserID
				require.EqualExportedValues(t, exp, smsplanet)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var smsplanet notification.SMSPlanet
				err := base.As(&smsplanet)
				require.NoError(t, err)
				return &smsplanet
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SMSPlanet)
				require.True(t, ok)
				var smsplanet notification.SMSPlanet
				err := actual.As(&smsplanet)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, smsplanet)
			},
		},
		{
			name:         "Splunk",
			expectedType: "Splunk",
			create: notification.Splunk{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Splunk Created",
				},
				SplunkDetails: notification.SplunkDetails{
					RestURL:        "https://api.victorops.com/api/v1/incidents",
					Severity:       "CRITICAL",
					AutoResolve:    "RECOVERY",
					IntegrationKey: "test-routing-key",
				},
			},
			updateFunc: func(n notification.Notification) {
				splunk := n.(*notification.Splunk)
				splunk.Name = "Test Splunk Updated"
				splunk.Severity = "HIGH"
				splunk.AutoResolve = "ACKNOWLEDGED"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Splunk)
				require.True(t, ok)
				var splunk notification.Splunk
				err := actual.As(&splunk)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = splunk.UserID
				require.EqualExportedValues(t, exp, splunk)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var splunk notification.Splunk
				err := base.As(&splunk)
				require.NoError(t, err)
				return &splunk
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Splunk)
				require.True(t, ok)
				var splunk notification.Splunk
				err := actual.As(&splunk)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, splunk)
			},
		},
		{
			name:         "SpugPush",
			expectedType: "SpugPush",
			create: notification.SpugPush{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test SpugPush Created",
				},
				SpugPushDetails: notification.SpugPushDetails{
					TemplateKey: "test-template-key",
				},
			},
			updateFunc: func(n notification.Notification) {
				spugpush := n.(*notification.SpugPush)
				spugpush.Name = "Test SpugPush Updated"
				spugpush.TemplateKey = "updated-template-key"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.SpugPush)
				require.True(t, ok)
				var spugpush notification.SpugPush
				err := actual.As(&spugpush)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = spugpush.UserID
				require.EqualExportedValues(t, exp, spugpush)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var spugpush notification.SpugPush
				err := base.As(&spugpush)
				require.NoError(t, err)
				return &spugpush
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.SpugPush)
				require.True(t, ok)
				var spugpush notification.SpugPush
				err := actual.As(&spugpush)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, spugpush)
			},
		},
		{
			name:         "Squadcast",
			expectedType: "squadcast",
			create: notification.Squadcast{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Squadcast Created",
				},
				SquadcastDetails: notification.SquadcastDetails{
					WebhookURL: "https://api.squadcast.com/api/v3/incidents/webhook",
				},
			},
			updateFunc: func(n notification.Notification) {
				squadcast := n.(*notification.Squadcast)
				squadcast.Name = "Test Squadcast Updated"
				squadcast.WebhookURL = "https://updated.squadcast.com/webhook"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Squadcast)
				require.True(t, ok)
				var squadcast notification.Squadcast
				err := actual.As(&squadcast)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = squadcast.UserID
				require.EqualExportedValues(t, exp, squadcast)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var squadcast notification.Squadcast
				err := base.As(&squadcast)
				require.NoError(t, err)
				return &squadcast
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Squadcast)
				require.True(t, ok)
				var squadcast notification.Squadcast
				err := actual.As(&squadcast)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, squadcast)
			},
		},
		{
			name:         "Stackfield",
			expectedType: "stackfield",
			create: notification.Stackfield{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Stackfield Created",
				},
				StackfieldDetails: notification.StackfieldDetails{
					WebhookURL: "https://app.stackfield.com/webhook/v1/xxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				stackfield := n.(*notification.Stackfield)
				stackfield.Name = "Test Stackfield Updated"
				stackfield.WebhookURL = "https://updated.stackfield.com/webhook"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Stackfield)
				require.True(t, ok)
				var stackfield notification.Stackfield
				err := actual.As(&stackfield)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = stackfield.UserID
				require.EqualExportedValues(t, exp, stackfield)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var stackfield notification.Stackfield
				err := base.As(&stackfield)
				require.NoError(t, err)
				return &stackfield
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Stackfield)
				require.True(t, ok)
				var stackfield notification.Stackfield
				err := actual.As(&stackfield)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, stackfield)
			},
		},
		{
			name:         "TechulusPush",
			expectedType: "PushByTechulus",
			create: notification.TechulusPush{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test TechulusPush Created",
				},
				TechulusPushDetails: notification.TechulusPushDetails{
					APIKey:        "test-api-key",
					Title:         "Alert Title",
					Sound:         "default",
					Channel:       "alerts",
					TimeSensitive: true,
				},
			},
			updateFunc: func(n notification.Notification) {
				techuluspush := n.(*notification.TechulusPush)
				techuluspush.Name = "Test TechulusPush Updated"
				techuluspush.Title = "Updated Title"
				techuluspush.Sound = "bell"
				techuluspush.Channel = "monitoring"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.TechulusPush)
				require.True(t, ok)
				var techuluspush notification.TechulusPush
				err := actual.As(&techuluspush)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = techuluspush.UserID
				require.EqualExportedValues(t, exp, techuluspush)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var techuluspush notification.TechulusPush
				err := base.As(&techuluspush)
				require.NoError(t, err)
				return &techuluspush
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.TechulusPush)
				require.True(t, ok)
				var techuluspush notification.TechulusPush
				err := actual.As(&techuluspush)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, techuluspush)
			},
		},
		{
			name:         "Threema",
			expectedType: "threema",
			create: notification.Threema{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Threema Created",
				},
				ThreemaDetails: notification.ThreemaDetails{
					SenderIdentity: "GATEWAY1",
					Secret:         "test-secret",
					Recipient:      "USERID123",
					RecipientType:  "identity",
				},
			},
			updateFunc: func(n notification.Notification) {
				threema := n.(*notification.Threema)
				threema.Name = "Test Threema Updated"
				threema.Recipient = "+41791234567"
				threema.RecipientType = "phone"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Threema)
				require.True(t, ok)
				var threema notification.Threema
				err := actual.As(&threema)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = threema.UserID
				require.EqualExportedValues(t, exp, threema)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var threema notification.Threema
				err := base.As(&threema)
				require.NoError(t, err)
				return &threema
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Threema)
				require.True(t, ok)
				var threema notification.Threema
				err := actual.As(&threema)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, threema)
			},
		},
		{
			name:         "WAHA",
			expectedType: "waha",
			create: notification.WAHA{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test WAHA Created",
				},
				WAHADetails: notification.WAHADetails{
					APIURL:  "https://waha.example.com",
					Session: "default",
					ChatID:  "5511999999999",
					APIKey:  "test-api-key",
				},
			},
			updateFunc: func(n notification.Notification) {
				waha := n.(*notification.WAHA)
				waha.Name = "Test WAHA Updated"
				waha.Session = "alerts"
				waha.ChatID = "+5511987654321"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.WAHA)
				require.True(t, ok)
				var waha notification.WAHA
				err := actual.As(&waha)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = waha.UserID
				require.EqualExportedValues(t, exp, waha)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var waha notification.WAHA
				err := base.As(&waha)
				require.NoError(t, err)
				return &waha
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.WAHA)
				require.True(t, ok)
				var waha notification.WAHA
				err := actual.As(&waha)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, waha)
			},
		},
		{
			name:         "Whapi",
			expectedType: "whapi",
			create: notification.Whapi{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test Whapi Created",
				},
				WhapiDetails: notification.WhapiDetails{
					APIURL:    "https://gate.whapi.cloud",
					AuthToken: "test-auth-token",
					Recipient: "5511999999999",
				},
			},
			updateFunc: func(n notification.Notification) {
				whapi := n.(*notification.Whapi)
				whapi.Name = "Test Whapi Updated"
				whapi.Recipient = "+5511987654321"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.Whapi)
				require.True(t, ok)
				var whapi notification.Whapi
				err := actual.As(&whapi)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = whapi.UserID
				require.EqualExportedValues(t, exp, whapi)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var whapi notification.Whapi
				err := base.As(&whapi)
				require.NoError(t, err)
				return &whapi
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.Whapi)
				require.True(t, ok)
				var whapi notification.Whapi
				err := actual.As(&whapi)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, whapi)
			},
		},
		{
			name:         "WPush",
			expectedType: "WPush",
			create: notification.WPush{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test WPush Created",
				},
				WPushDetails: notification.WPushDetails{
					APIKey:  "test-api-key-123",
					Channel: "channel-alerts",
				},
			},
			updateFunc: func(n notification.Notification) {
				wpush := n.(*notification.WPush)
				wpush.Name = "Test WPush Updated"
				wpush.Channel = "channel-monitoring"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.WPush)
				require.True(t, ok)
				var wpush notification.WPush
				err := actual.As(&wpush)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = wpush.UserID
				require.EqualExportedValues(t, exp, wpush)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var wpush notification.WPush
				err := base.As(&wpush)
				require.NoError(t, err)
				return &wpush
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.WPush)
				require.True(t, ok)
				var wpush notification.WPush
				err := actual.As(&wpush)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, wpush)
			},
		},
		{
			name:         "YZJ",
			expectedType: "YZJ",
			create: notification.YZJ{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test YZJ Created",
				},
				YZJDetails: notification.YZJDetails{
					WebHookURL: "https://api.yzj.cn/webhook",
					Token:      "test-token-123",
				},
			},
			updateFunc: func(n notification.Notification) {
				yzj := n.(*notification.YZJ)
				yzj.Name = "Test YZJ Updated"
				yzj.Token = "updated-token-456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.YZJ)
				require.True(t, ok)
				var yzj notification.YZJ
				err := actual.As(&yzj)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = yzj.UserID
				require.EqualExportedValues(t, exp, yzj)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var yzj notification.YZJ
				err := base.As(&yzj)
				require.NoError(t, err)
				return &yzj
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.YZJ)
				require.True(t, ok)
				var yzj notification.YZJ
				err := actual.As(&yzj)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, yzj)
			},
		},
		{
			name:         "ZohoCliq",
			expectedType: "ZohoCliq",
			create: notification.ZohoCliq{
				Base: notification.Base{
					ApplyExisting: true,
					IsDefault:     false,
					IsActive:      true,
					Name:          "Test ZohoCliq Created",
				},
				ZohoCliqDetails: notification.ZohoCliqDetails{
					WebhookURL: "https://zoho-cliq.example.com/webhook/test123",
				},
			},
			updateFunc: func(n notification.Notification) {
				zohocliq := n.(*notification.ZohoCliq)
				zohocliq.Name = "Test ZohoCliq Updated"
				zohocliq.WebhookURL = "https://zoho-cliq.example.com/webhook/updated456"
			},
			verifyCreatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
				t.Helper()
				exp, ok := expected.(notification.ZohoCliq)
				require.True(t, ok)
				var zohocliq notification.ZohoCliq
				err := actual.As(&zohocliq)
				require.NoError(t, err)
				exp.ID = id
				exp.UserID = zohocliq.UserID
				require.EqualExportedValues(t, exp, zohocliq)
			},
			createTypedFunc: func(t *testing.T, base notification.Notification) notification.Notification {
				t.Helper()
				var zohocliq notification.ZohoCliq
				err := base.As(&zohocliq)
				require.NoError(t, err)
				return &zohocliq
			},
			verifyUpdatedFunc: func(t *testing.T, actual notification.Notification, expected notification.Notification) {
				t.Helper()
				exp, ok := expected.(*notification.ZohoCliq)
				require.True(t, ok)
				var zohocliq notification.ZohoCliq
				err := actual.As(&zohocliq)
				require.NoError(t, err)
				require.EqualExportedValues(t, *exp, zohocliq)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()

			var err error

			t.Run("initial_state", func(t *testing.T) {
				notifications := client.GetNotifications(ctx)
				t.Logf("Initial notifications count: %d", len(notifications))
			})

			var id int64
			t.Run("create", func(t *testing.T) {
				initialNotifications := client.GetNotifications(ctx)
				initialCount := len(initialNotifications)

				id, err = client.CreateNotification(ctx, tc.create)
				require.NoError(t, err)
				require.Greater(t, id, int64(0))

				notifications := client.GetNotifications(ctx)
				require.Len(t, notifications, initialCount+1)

				createdNotification, err := client.GetNotification(ctx, id)
				require.NoError(t, err)
				require.Equal(t, tc.expectedType, createdNotification.Type())
				require.Equal(t, id, createdNotification.GetID())

				tc.verifyCreatedFunc(t, createdNotification, tc.create, id)
			})

			t.Run("update", func(t *testing.T) {
				require.NotZero(t, id, "create test failed, unable to test update")

				currentNotification, err := client.GetNotification(ctx, id)
				require.NoError(t, err)

				updated := tc.createTypedFunc(t, currentNotification)
				tc.updateFunc(updated)

				err = client.UpdateNotification(ctx, updated)
				require.NoError(t, err)

				retrievedNotification, err := client.GetNotification(ctx, id)
				require.NoError(t, err)

				tc.verifyUpdatedFunc(t, retrievedNotification, updated)
			})

			t.Run("delete", func(t *testing.T) {
				require.NotZero(t, id, "create test failed, unable to test delete")

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
