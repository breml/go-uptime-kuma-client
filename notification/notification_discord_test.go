package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationDiscord_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Discord
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Discord Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Discord Alert\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz\",\"discordUsername\":\"Uptime Monitor\",\"discordChannelType\":\"postToThread\",\"threadId\":\"987654321\",\"discordPrefixMessage\":\"Alert:\",\"disableUrl\":true,\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            1,
					Name:          "My Discord Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
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
			wantJSON: `{"active":true,"applyExisting":true,"disableUrl":true,"discordChannelType":"postToThread","discordPrefixMessage":"Alert:","discordUsername":"Uptime Monitor","discordWebhookUrl":"https://discordapp.com/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz","id":1,"isDefault":true,"name":"My Discord Alert","postName":"","threadId":"987654321","type":"discord","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Discord","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Discord\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/xyz/abc\",\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Discord",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL: "https://discordapp.com/api/webhooks/xyz/abc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"","discordPrefixMessage":"","discordUsername":"","discordWebhookUrl":"https://discordapp.com/api/webhooks/xyz/abc","id":2,"isDefault":false,"name":"Simple Discord","postName":"","threadId":"","type":"discord","userId":1}`,
		},
		{
			name: "with forum post",
			data: []byte(
				`{"id":3,"name":"Discord Forum","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Discord Forum\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/forum/webhook\",\"discordUsername\":\"Monitor Bot\",\"discordChannelType\":\"createNewForumPost\",\"postName\":\"System Alert\",\"disableUrl\":false,\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            3,
					Name:          "Discord Forum",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL:  "https://discordapp.com/api/webhooks/forum/webhook",
					Username:    "Monitor Bot",
					ChannelType: "createNewForumPost",
					PostName:    "System Alert",
					DisableURL:  false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"createNewForumPost","discordPrefixMessage":"","discordUsername":"Monitor Bot","discordWebhookUrl":"https://discordapp.com/api/webhooks/forum/webhook","id":3,"isDefault":false,"name":"Discord Forum","postName":"System Alert","threadId":"","type":"discord","userId":1}`,
		},
		{
			name: "with suppress notifications",
			data: []byte(
				`{"id":4,"name":"Silent Discord","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Silent Discord\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/silent/webhook\",\"discordSuppressNotifications\":true,\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            4,
					Name:          "Silent Discord",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL:            "https://discordapp.com/api/webhooks/silent/webhook",
					SuppressNotifications: true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"","discordPrefixMessage":"","discordSuppressNotifications":true,"discordUsername":"","discordWebhookUrl":"https://discordapp.com/api/webhooks/silent/webhook","id":4,"isDefault":false,"name":"Silent Discord","postName":"","threadId":"","type":"discord","userId":1}`,
		},
		{
			name: "with minimalist message format",
			data: []byte(
				`{"id":5,"name":"Minimalist Discord","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Minimalist Discord\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/min/webhook\",\"discordMessageFormat\":\"minimalist\",\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            5,
					Name:          "Minimalist Discord",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL:    "https://discordapp.com/api/webhooks/min/webhook",
					MessageFormat: "minimalist",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"","discordMessageFormat":"minimalist","discordPrefixMessage":"","discordUsername":"","discordWebhookUrl":"https://discordapp.com/api/webhooks/min/webhook","id":5,"isDefault":false,"name":"Minimalist Discord","postName":"","threadId":"","type":"discord","userId":1}`,
		},
		{
			name: "with custom message format and template",
			data: []byte(
				`{"id":6,"name":"Custom Discord","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Custom Discord\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/custom/webhook\",\"discordMessageFormat\":\"custom\",\"discordMessageTemplate\":\"Service {{name}} is {{status}}\",\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            6,
					Name:          "Custom Discord",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL:      "https://discordapp.com/api/webhooks/custom/webhook",
					MessageFormat:   "custom",
					MessageTemplate: "Service {{name}} is {{status}}",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"","discordMessageFormat":"custom","discordMessageTemplate":"Service {{name}} is {{status}}","discordPrefixMessage":"","discordUsername":"","discordWebhookUrl":"https://discordapp.com/api/webhooks/custom/webhook","id":6,"isDefault":false,"name":"Custom Discord","postName":"","threadId":"","type":"discord","userId":1}`,
		},
		{
			name: "with legacy use message template",
			data: []byte(
				`{"id":7,"name":"Legacy Template Discord","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Legacy Template Discord\",\"discordWebhookUrl\":\"https://discordapp.com/api/webhooks/legacy/webhook\",\"discordUseMessageTemplate\":true,\"discordMessageTemplate\":\"Alert: {{name}}\",\"type\":\"discord\"}"}`,
			),

			want: notification.Discord{
				Base: notification.Base{
					ID:            7,
					Name:          "Legacy Template Discord",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DiscordDetails: notification.DiscordDetails{
					WebhookURL:         "https://discordapp.com/api/webhooks/legacy/webhook",
					UseMessageTemplate: true,
					MessageTemplate:    "Alert: {{name}}",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"discordChannelType":"","discordMessageTemplate":"Alert: {{name}}","discordPrefixMessage":"","discordUseMessageTemplate":true,"discordUsername":"","discordWebhookUrl":"https://discordapp.com/api/webhooks/legacy/webhook","id":7,"isDefault":false,"name":"Legacy Template Discord","postName":"","threadId":"","type":"discord","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			discord := notification.Discord{}

			err := json.Unmarshal(tc.data, &discord)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, discord)

			data, err := json.Marshal(discord)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
