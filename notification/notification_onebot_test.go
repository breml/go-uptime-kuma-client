package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOneBot_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.OneBot
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"Test OneBot Group","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"OneBot\",\"httpAddr\":\"http://localhost:5700\",\"accessToken\":\"test-token\",\"msgType\":\"group\",\"recieverId\":\"123456789\"}"}`,
			),

			want: notification.OneBot{
				Base: notification.Base{
					ID:        1,
					Name:      "Test OneBot Group",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				OneBotDetails: notification.OneBotDetails{
					HTTPAddr:    "http://localhost:5700",
					AccessToken: "test-token",
					MsgType:     "group",
					ReceiverID:  "123456789",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"httpAddr":"http://localhost:5700","id":1,"isDefault":false,"msgType":"group","name":"Test OneBot Group","recieverId":"123456789","type":"OneBot","userId":42}`,
		},
		{
			name: "private message",
			data: []byte(
				`{"id":2,"name":"Test OneBot Private","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"OneBot\",\"httpAddr\":\"http://bot.example.com\",\"accessToken\":\"secret-token\",\"msgType\":\"private\",\"recieverId\":\"987654321\"}"}`,
			),

			want: notification.OneBot{
				Base: notification.Base{
					ID:            2,
					Name:          "Test OneBot Private",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				OneBotDetails: notification.OneBotDetails{
					HTTPAddr:    "http://bot.example.com",
					AccessToken: "secret-token",
					MsgType:     "private",
					ReceiverID:  "987654321",
				},
			},
			wantJSON: `{"accessToken":"secret-token","active":true,"applyExisting":false,"httpAddr":"http://bot.example.com","id":2,"isDefault":true,"msgType":"private","name":"Test OneBot Private","recieverId":"987654321","type":"OneBot","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test OneBot Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"OneBot\",\"httpAddr\":\"localhost:5700\",\"accessToken\":\"\",\"msgType\":\"group\",\"recieverId\":\"\"}"}`,
			),

			want: notification.OneBot{
				Base: notification.Base{
					ID:            3,
					Name:          "Test OneBot Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OneBotDetails: notification.OneBotDetails{
					HTTPAddr:    "localhost:5700",
					AccessToken: "",
					MsgType:     "group",
					ReceiverID:  "",
				},
			},
			wantJSON: `{"accessToken":"","active":false,"applyExisting":false,"httpAddr":"localhost:5700","id":3,"isDefault":false,"msgType":"group","name":"Test OneBot Minimal","recieverId":"","type":"OneBot","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			onebot := notification.OneBot{}

			err := json.Unmarshal(tc.data, &onebot)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, onebot)

			data, err := json.Marshal(onebot)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
