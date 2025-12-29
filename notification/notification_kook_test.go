package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationKook_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Kook
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My KOOK Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My KOOK Alert\",\"kookBotToken\":\"bot-token-123\",\"kookGuildID\":\"guild-456\",\"type\":\"Kook\"}"}`,
			),

			want: notification.Kook{
				Base: notification.Base{
					ID:            1,
					Name:          "My KOOK Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				KookDetails: notification.KookDetails{
					BotToken: "bot-token-123",
					GuildID:  "guild-456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"kookBotToken":"bot-token-123","kookGuildID":"guild-456","name":"My KOOK Alert","type":"Kook","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple KOOK","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple KOOK\",\"kookBotToken\":\"token-abc\",\"kookGuildID\":\"guild-xyz\",\"type\":\"Kook\"}"}`,
			),

			want: notification.Kook{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple KOOK",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				KookDetails: notification.KookDetails{
					BotToken: "token-abc",
					GuildID:  "guild-xyz",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"kookBotToken":"token-abc","kookGuildID":"guild-xyz","name":"Simple KOOK","type":"Kook","userId":1}`,
		},
		{
			name: "with different guild",
			data: []byte(
				`{"id":3,"name":"KOOK Production","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"KOOK Production\",\"kookBotToken\":\"prod-token-xyz\",\"kookGuildID\":\"prod-guild-789\",\"type\":\"Kook\"}"}`,
			),

			want: notification.Kook{
				Base: notification.Base{
					ID:            3,
					Name:          "KOOK Production",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				KookDetails: notification.KookDetails{
					BotToken: "prod-token-xyz",
					GuildID:  "prod-guild-789",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"kookBotToken":"prod-token-xyz","kookGuildID":"prod-guild-789","name":"KOOK Production","type":"Kook","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kook := notification.Kook{}

			err := json.Unmarshal(tc.data, &kook)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, kook)

			data, err := json.Marshal(kook)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
