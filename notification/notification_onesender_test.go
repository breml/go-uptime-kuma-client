package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOneSender_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.OneSender
		wantJSON string
	}{
		{
			name: "success private",
			data: []byte(
				`{"id":1,"name":"Test OneSender Private","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"onesender\",\"onesenderURL\":\"https://api.onesender.com/send\",\"onesenderToken\":\"test-token\",\"onesenderReceiver\":\"5511999999999\",\"onesenderTypeReceiver\":\"private\"}"}`,
			),

			want: notification.OneSender{
				Base: notification.Base{
					ID:        1,
					Name:      "Test OneSender Private",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				OneSenderDetails: notification.OneSenderDetails{
					URL:          "https://api.onesender.com/send",
					Token:        "test-token",
					Receiver:     "5511999999999",
					TypeReceiver: "private",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"Test OneSender Private","onesenderReceiver":"5511999999999","onesenderToken":"test-token","onesenderTypeReceiver":"private","onesenderURL":"https://api.onesender.com/send","type":"onesender","userId":42}`,
		},
		{
			name: "success group",
			data: []byte(
				`{"id":2,"name":"Test OneSender Group","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"onesender\",\"onesenderURL\":\"https://api.onesender.com/send\",\"onesenderToken\":\"secret-token\",\"onesenderReceiver\":\"120363123456789-1234567890\",\"onesenderTypeReceiver\":\"group\"}"}`,
			),

			want: notification.OneSender{
				Base: notification.Base{
					ID:            2,
					Name:          "Test OneSender Group",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				OneSenderDetails: notification.OneSenderDetails{
					URL:          "https://api.onesender.com/send",
					Token:        "secret-token",
					Receiver:     "120363123456789-1234567890",
					TypeReceiver: "group",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":true,"name":"Test OneSender Group","onesenderReceiver":"120363123456789-1234567890","onesenderToken":"secret-token","onesenderTypeReceiver":"group","onesenderURL":"https://api.onesender.com/send","type":"onesender","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test OneSender Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"onesender\",\"onesenderURL\":\"\",\"onesenderToken\":\"\",\"onesenderReceiver\":\"\",\"onesenderTypeReceiver\":\"\"}"}`,
			),

			want: notification.OneSender{
				Base: notification.Base{
					ID:            3,
					Name:          "Test OneSender Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OneSenderDetails: notification.OneSenderDetails{
					URL:          "",
					Token:        "",
					Receiver:     "",
					TypeReceiver: "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Test OneSender Minimal","onesenderReceiver":"","onesenderToken":"","onesenderTypeReceiver":"","onesenderURL":"","type":"onesender","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			onesender := notification.OneSender{}

			err := json.Unmarshal(tc.data, &onesender)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, onesender)

			data, err := json.Marshal(onesender)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
