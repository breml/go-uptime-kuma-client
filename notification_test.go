package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// notificationTestCase defines a single notification type's CRUD test scenario
type notificationTestCase struct {
	name         string                          // Test name (e.g., "Ntfy", "Slack")
	expectedType string                          // Expected type string from API
	create       notification.Notification       // Notification to create
	updateFunc   func(notification.Notification) // Function to modify notification for update test
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
					ApiKey:   "test-api-key-123",
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
					ApiKey:     "",
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
					BotKey: "abc123def456",
				},
			},
			updateFunc: func(n notification.Notification) {
				wecom := n.(*notification.WeCom)
				wecom.Name = "Test WeCom Updated"
				wecom.BotKey = "xyz789abc123"
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
					WebHookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
				},
			},
			updateFunc: func(n notification.Notification) {
				feishu := n.(*notification.Feishu)
				feishu.Name = "Test Feishu Updated"
				feishu.WebHookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/yyy"
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
					WebHookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
					SecretKey:  "secret123",
					Mentioning: "everyone",
				},
			},
			updateFunc: func(n notification.Notification) {
				dingding := n.(*notification.DingDing)
				dingding.Name = "Test DingDing Updated"
				dingding.Mentioning = ""
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

				verifyCreatedNotification(t, createdNotification, tc.create, id)
			})

			t.Run("update", func(t *testing.T) {
				currentNotification, err := client.GetNotification(ctx, id)
				require.NoError(t, err)

				updated := createTypedNotification(t, currentNotification, tc.create)
				tc.updateFunc(updated)

				err = client.UpdateNotification(ctx, updated)
				require.NoError(t, err)

				retrievedNotification, err := client.GetNotification(ctx, id)
				require.NoError(t, err)

				verifyUpdatedNotification(t, retrievedNotification, updated)
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
}

// verifyCreatedNotification checks that the created notification matches expected values
func verifyCreatedNotification(t *testing.T, actual notification.Notification, expected notification.Notification, id int64) {
	t.Helper()

	switch exp := expected.(type) {
	case notification.Ntfy:
		var ntfy notification.Ntfy
		err := actual.As(&ntfy)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = ntfy.UserID
		require.EqualExportedValues(t, exp, ntfy)
	case notification.Slack:
		var slack notification.Slack
		err := actual.As(&slack)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = slack.UserID
		require.EqualExportedValues(t, exp, slack)
	case notification.Teams:
		var teams notification.Teams
		err := actual.As(&teams)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = teams.UserID
		require.EqualExportedValues(t, exp, teams)
	case notification.Webhook:
		var webhook notification.Webhook
		err := actual.As(&webhook)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = webhook.UserID
		require.EqualExportedValues(t, exp, webhook)
	case notification.SMTP:
		var smtp notification.SMTP
		err := actual.As(&smtp)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = smtp.UserID
		require.EqualExportedValues(t, exp, smtp)
	case notification.Telegram:
		var telegram notification.Telegram
		err := actual.As(&telegram)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = telegram.UserID
		require.EqualExportedValues(t, exp, telegram)
	case notification.PagerDuty:
		var pagerduty notification.PagerDuty
		err := actual.As(&pagerduty)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = pagerduty.UserID
		require.EqualExportedValues(t, exp, pagerduty)
	case notification.Signal:
		var signal notification.Signal
		err := actual.As(&signal)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = signal.UserID
		require.EqualExportedValues(t, exp, signal)
	case notification.Opsgenie:
		var opsgenie notification.Opsgenie
		err := actual.As(&opsgenie)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = opsgenie.UserID
		require.EqualExportedValues(t, exp, opsgenie)
	case notification.HomeAssistant:
		var ha notification.HomeAssistant
		err := actual.As(&ha)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = ha.UserID
		require.EqualExportedValues(t, exp, ha)
	case notification.Discord:
		var discord notification.Discord
		err := actual.As(&discord)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = discord.UserID
		require.EqualExportedValues(t, exp, discord)
	case notification.Pushbullet:
		var pushbullet notification.Pushbullet
		err := actual.As(&pushbullet)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = pushbullet.UserID
		require.EqualExportedValues(t, exp, pushbullet)
	case notification.Pushover:
		var pushover notification.Pushover
		err := actual.As(&pushover)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = pushover.UserID
		require.EqualExportedValues(t, exp, pushover)
	case notification.Gotify:
		var gotify notification.Gotify
		err := actual.As(&gotify)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = gotify.UserID
		require.EqualExportedValues(t, exp, gotify)
	case notification.GrafanaOncall:
		var grafana notification.GrafanaOncall
		err := actual.As(&grafana)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = grafana.UserID
		require.EqualExportedValues(t, exp, grafana)
	case notification.Twilio:
		var twilio notification.Twilio
		err := actual.As(&twilio)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = twilio.UserID
		require.EqualExportedValues(t, exp, twilio)
	case notification.Mattermost:
		var mattermost notification.Mattermost
		err := actual.As(&mattermost)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = mattermost.UserID
		require.EqualExportedValues(t, exp, mattermost)
	case notification.Matrix:
		var matrix notification.Matrix
		err := actual.As(&matrix)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = matrix.UserID
		require.EqualExportedValues(t, exp, matrix)
	case notification.RocketChat:
		var rocketchat notification.RocketChat
		err := actual.As(&rocketchat)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = rocketchat.UserID
		require.EqualExportedValues(t, exp, rocketchat)
	case notification.WeCom:
		var wecom notification.WeCom
		err := actual.As(&wecom)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = wecom.UserID
		require.EqualExportedValues(t, exp, wecom)
	case notification.Feishu:
		var feishu notification.Feishu
		err := actual.As(&feishu)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = feishu.UserID
		require.EqualExportedValues(t, exp, feishu)
	case notification.DingDing:
		var dingding notification.DingDing
		err := actual.As(&dingding)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = dingding.UserID
		require.EqualExportedValues(t, exp, dingding)
	case notification.Apprise:
		var apprise notification.Apprise
		err := actual.As(&apprise)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = apprise.UserID
		require.EqualExportedValues(t, exp, apprise)
	case notification.GoogleChat:
		var googlechat notification.GoogleChat
		err := actual.As(&googlechat)
		require.NoError(t, err)
		exp.ID = id
		exp.UserID = googlechat.UserID
		require.EqualExportedValues(t, exp, googlechat)
	default:
		t.Fatalf("unknown notification type: %T", expected)
	}
}

// createTypedNotification creates a properly typed notification from base notification
func createTypedNotification(t *testing.T, base notification.Notification, template notification.Notification) notification.Notification {
	t.Helper()

	switch template.(type) {
	case notification.Ntfy:
		var ntfy notification.Ntfy
		err := base.As(&ntfy)
		require.NoError(t, err)
		return &ntfy
	case notification.Slack:
		var slack notification.Slack
		err := base.As(&slack)
		require.NoError(t, err)
		return &slack
	case notification.Teams:
		var teams notification.Teams
		err := base.As(&teams)
		require.NoError(t, err)
		return &teams
	case notification.Webhook:
		var webhook notification.Webhook
		err := base.As(&webhook)
		require.NoError(t, err)
		return &webhook
	case notification.SMTP:
		var smtp notification.SMTP
		err := base.As(&smtp)
		require.NoError(t, err)
		return &smtp
	case notification.Telegram:
		var telegram notification.Telegram
		err := base.As(&telegram)
		require.NoError(t, err)
		return &telegram
	case notification.PagerDuty:
		var pagerduty notification.PagerDuty
		err := base.As(&pagerduty)
		require.NoError(t, err)
		return &pagerduty
	case notification.Signal:
		var signal notification.Signal
		err := base.As(&signal)
		require.NoError(t, err)
		return &signal
	case notification.Opsgenie:
		var opsgenie notification.Opsgenie
		err := base.As(&opsgenie)
		require.NoError(t, err)
		return &opsgenie
	case notification.HomeAssistant:
		var ha notification.HomeAssistant
		err := base.As(&ha)
		require.NoError(t, err)
		return &ha
	case notification.Discord:
		var discord notification.Discord
		err := base.As(&discord)
		require.NoError(t, err)
		return &discord
	case notification.Pushbullet:
		var pushbullet notification.Pushbullet
		err := base.As(&pushbullet)
		require.NoError(t, err)
		return &pushbullet
	case notification.Pushover:
		var pushover notification.Pushover
		err := base.As(&pushover)
		require.NoError(t, err)
		return &pushover
	case notification.Gotify:
		var gotify notification.Gotify
		err := base.As(&gotify)
		require.NoError(t, err)
		return &gotify
	case notification.GrafanaOncall:
		var grafana notification.GrafanaOncall
		err := base.As(&grafana)
		require.NoError(t, err)
		return &grafana
	case notification.Twilio:
		var twilio notification.Twilio
		err := base.As(&twilio)
		require.NoError(t, err)
		return &twilio
	case notification.Mattermost:
		var mattermost notification.Mattermost
		err := base.As(&mattermost)
		require.NoError(t, err)
		return &mattermost
	case notification.Matrix:
		var matrix notification.Matrix
		err := base.As(&matrix)
		require.NoError(t, err)
		return &matrix
	case notification.RocketChat:
		var rocketchat notification.RocketChat
		err := base.As(&rocketchat)
		require.NoError(t, err)
		return &rocketchat
	case notification.WeCom:
		var wecom notification.WeCom
		err := base.As(&wecom)
		require.NoError(t, err)
		return &wecom
	case notification.Feishu:
		var feishu notification.Feishu
		err := base.As(&feishu)
		require.NoError(t, err)
		return &feishu
	case notification.DingDing:
		var dingding notification.DingDing
		err := base.As(&dingding)
		require.NoError(t, err)
		return &dingding
	case notification.Apprise:
		var apprise notification.Apprise
		err := base.As(&apprise)
		require.NoError(t, err)
		return &apprise
	case notification.GoogleChat:
		var googlechat notification.GoogleChat
		err := base.As(&googlechat)
		require.NoError(t, err)
		return &googlechat
	default:
		t.Fatalf("unknown notification type: %T", template)
		return nil
	}
}

// verifyUpdatedNotification checks that the updated notification matches expected values
func verifyUpdatedNotification(t *testing.T, actual notification.Notification, expected notification.Notification) {
	t.Helper()

	switch exp := expected.(type) {
	case *notification.Ntfy:
		var ntfy notification.Ntfy
		err := actual.As(&ntfy)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, ntfy)
	case *notification.Slack:
		var slack notification.Slack
		err := actual.As(&slack)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, slack)
	case *notification.Teams:
		var teams notification.Teams
		err := actual.As(&teams)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, teams)
	case *notification.Webhook:
		var webhook notification.Webhook
		err := actual.As(&webhook)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, webhook)
	case *notification.SMTP:
		var smtp notification.SMTP
		err := actual.As(&smtp)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, smtp)
	case *notification.Telegram:
		var telegram notification.Telegram
		err := actual.As(&telegram)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, telegram)
	case *notification.PagerDuty:
		var pagerduty notification.PagerDuty
		err := actual.As(&pagerduty)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, pagerduty)
	case *notification.Signal:
		var signal notification.Signal
		err := actual.As(&signal)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, signal)
	case *notification.Opsgenie:
		var opsgenie notification.Opsgenie
		err := actual.As(&opsgenie)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, opsgenie)
	case *notification.HomeAssistant:
		var ha notification.HomeAssistant
		err := actual.As(&ha)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, ha)
	case *notification.Discord:
		var discord notification.Discord
		err := actual.As(&discord)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, discord)
	case *notification.Pushbullet:
		var pushbullet notification.Pushbullet
		err := actual.As(&pushbullet)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, pushbullet)
	case *notification.Pushover:
		var pushover notification.Pushover
		err := actual.As(&pushover)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, pushover)
	case *notification.Gotify:
		var gotify notification.Gotify
		err := actual.As(&gotify)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, gotify)
	case *notification.GrafanaOncall:
		var grafana notification.GrafanaOncall
		err := actual.As(&grafana)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, grafana)
	case *notification.Twilio:
		var twilio notification.Twilio
		err := actual.As(&twilio)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, twilio)
	case *notification.Mattermost:
		var mattermost notification.Mattermost
		err := actual.As(&mattermost)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, mattermost)
	case *notification.Matrix:
		var matrix notification.Matrix
		err := actual.As(&matrix)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, matrix)
	case *notification.RocketChat:
		var rocketchat notification.RocketChat
		err := actual.As(&rocketchat)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, rocketchat)
	case *notification.WeCom:
		var wecom notification.WeCom
		err := actual.As(&wecom)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, wecom)
	case *notification.Feishu:
		var feishu notification.Feishu
		err := actual.As(&feishu)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, feishu)
	case *notification.DingDing:
		var dingding notification.DingDing
		err := actual.As(&dingding)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, dingding)
	case *notification.Apprise:
		var apprise notification.Apprise
		err := actual.As(&apprise)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, apprise)
	case *notification.GoogleChat:
		var googlechat notification.GoogleChat
		err := actual.As(&googlechat)
		require.NoError(t, err)
		require.EqualExportedValues(t, *exp, googlechat)
	default:
		t.Fatalf("unknown notification type: %T", expected)
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
