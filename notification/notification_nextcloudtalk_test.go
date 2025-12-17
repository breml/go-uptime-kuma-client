package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNextcloudTalk_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.NextcloudTalk
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My Nextcloud Talk Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Nextcloud Talk Alert\",\"host\":\"https://nextcloud.example.com\",\"conversationToken\":\"token123\",\"botSecret\":\"secret-key-123\",\"sendSilentUp\":true,\"sendSilentDown\":false,\"type\":\"NextcloudTalk\"}"}`),

			want: notification.NextcloudTalk{
				Base: notification.Base{
					ID:            1,
					Name:          "My Nextcloud Talk Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NextcloudTalkDetails: notification.NextcloudTalkDetails{
					Host:              "https://nextcloud.example.com",
					ConversationToken: "token123",
					BotSecret:         "secret-key-123",
					SendSilentUp:      true,
					SendSilentDown:    false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"botSecret":"secret-key-123","conversationToken":"token123","host":"https://nextcloud.example.com","id":1,"isDefault":true,"name":"My Nextcloud Talk Alert","sendSilentDown":false,"sendSilentUp":true,"type":"NextcloudTalk","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple Nextcloud Talk","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Nextcloud Talk\",\"host\":\"https://nc.example.com\",\"conversationToken\":\"abc123\",\"botSecret\":\"secret-abc\",\"sendSilentUp\":false,\"sendSilentDown\":false,\"type\":\"NextcloudTalk\"}"}`),

			want: notification.NextcloudTalk{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Nextcloud Talk",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				NextcloudTalkDetails: notification.NextcloudTalkDetails{
					Host:              "https://nc.example.com",
					ConversationToken: "abc123",
					BotSecret:         "secret-abc",
					SendSilentUp:      false,
					SendSilentDown:    false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"botSecret":"secret-abc","conversationToken":"abc123","host":"https://nc.example.com","id":2,"isDefault":false,"name":"Simple Nextcloud Talk","sendSilentDown":false,"sendSilentUp":false,"type":"NextcloudTalk","userId":1}`,
		},
		{
			name: "with silent down enabled",
			data: []byte(`{"id":3,"name":"Nextcloud Talk Silent","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Nextcloud Talk Silent\",\"host\":\"https://cloud.example.com\",\"conversationToken\":\"xyz789\",\"botSecret\":\"prod-secret\",\"sendSilentUp\":false,\"sendSilentDown\":true,\"type\":\"NextcloudTalk\"}"}`),

			want: notification.NextcloudTalk{
				Base: notification.Base{
					ID:            3,
					Name:          "Nextcloud Talk Silent",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				NextcloudTalkDetails: notification.NextcloudTalkDetails{
					Host:              "https://cloud.example.com",
					ConversationToken: "xyz789",
					BotSecret:         "prod-secret",
					SendSilentUp:      false,
					SendSilentDown:    true,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"botSecret":"prod-secret","conversationToken":"xyz789","host":"https://cloud.example.com","id":3,"isDefault":false,"name":"Nextcloud Talk Silent","sendSilentDown":true,"sendSilentUp":false,"type":"NextcloudTalk","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextcloudtalk := notification.NextcloudTalk{}

			err := json.Unmarshal(tc.data, &nextcloudtalk)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, nextcloudtalk)

			data, err := json.Marshal(nextcloudtalk)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
