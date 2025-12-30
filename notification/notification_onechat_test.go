package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOneChat_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.OneChat
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"Test OneChat","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"OneChat\",\"accessToken\":\"test-token\",\"recieverId\":\"user123\",\"botId\":\"bot456\"}"}`,
			),

			want: notification.OneChat{
				Base: notification.Base{
					ID:        1,
					Name:      "Test OneChat",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				OneChatDetails: notification.OneChatDetails{
					AccessToken: "test-token",
					ReceiverID:  "user123",
					BotID:       "bot456",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"botId":"bot456","id":1,"isDefault":false,"name":"Test OneChat","recieverId":"user123","type":"OneChat","userId":42}`,
		},
		{
			name: "with group receiver",
			data: []byte(
				`{"id":2,"name":"Test OneChat Group","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"OneChat\",\"accessToken\":\"secret-token\",\"recieverId\":\"group789\",\"botId\":\"botgroup\"}"}`,
			),

			want: notification.OneChat{
				Base: notification.Base{
					ID:            2,
					Name:          "Test OneChat Group",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				OneChatDetails: notification.OneChatDetails{
					AccessToken: "secret-token",
					ReceiverID:  "group789",
					BotID:       "botgroup",
				},
			},
			wantJSON: `{"accessToken":"secret-token","active":true,"applyExisting":false,"botId":"botgroup","id":2,"isDefault":true,"name":"Test OneChat Group","recieverId":"group789","type":"OneChat","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test OneChat Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"OneChat\",\"accessToken\":\"\",\"recieverId\":\"\",\"botId\":\"\"}"}`,
			),

			want: notification.OneChat{
				Base: notification.Base{
					ID:            3,
					Name:          "Test OneChat Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OneChatDetails: notification.OneChatDetails{
					AccessToken: "",
					ReceiverID:  "",
					BotID:       "",
				},
			},
			wantJSON: `{"accessToken":"","active":false,"applyExisting":false,"botId":"","id":3,"isDefault":false,"name":"Test OneChat Minimal","recieverId":"","type":"OneChat","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			onechat := notification.OneChat{}

			err := json.Unmarshal(tc.data, &onechat)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, onechat)

			data, err := json.Marshal(onechat)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
