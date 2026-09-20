package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationMilky_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Milky
		wantJSON string
	}{
		{
			name: "group message",
			data: []byte(
				`{"id":1,"name":"Test Milky Group","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"http://localhost:3000\",\"accessToken\":\"test-token\",\"msgType\":\"group\",\"recieverId\":\"123456789\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        1,
					Name:      "Test Milky Group",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "http://localhost:3000",
					AccessToken: "test-token",
					MsgType:     notification.MilkyMessageTypeGroup,
					ReceiverID:  "123456789",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"httpAddr":"http://localhost:3000","id":1,"isDefault":false,"msgType":"group","name":"Test Milky Group","recieverId":"123456789","type":"Milky","userId":42}`,
		},
		{
			name: "private message",
			data: []byte(
				`{"id":2,"name":"Test Milky Private","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"Milky\",\"httpAddr\":\"bot.example.com\",\"accessToken\":\"secret-token\",\"msgType\":\"private\",\"recieverId\":\"987654321\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        2,
					Name:      "Test Milky Private",
					IsActive:  true,
					UserID:    42,
					IsDefault: true,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "bot.example.com",
					AccessToken: "secret-token",
					MsgType:     notification.MilkyMessageTypePrivate,
					ReceiverID:  "987654321",
				},
			},
			wantJSON: `{"accessToken":"secret-token","active":true,"applyExisting":false,"httpAddr":"bot.example.com","id":2,"isDefault":true,"msgType":"private","name":"Test Milky Private","recieverId":"987654321","type":"Milky","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test Milky Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"localhost:3000\",\"accessToken\":\"\",\"recieverId\":\"\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        3,
					Name:      "Test Milky Minimal",
					IsActive:  false,
					UserID:    10,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "localhost:3000",
					AccessToken: "",
					ReceiverID:  "",
				},
			},
			wantJSON: `{"accessToken":"","active":false,"applyExisting":false,"httpAddr":"localhost:3000","id":3,"isDefault":false,"name":"Test Milky Minimal","recieverId":"","type":"Milky","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			milky := notification.Milky{}

			err := json.Unmarshal(tc.data, &milky)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, milky)

			data, err := json.Marshal(milky)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
