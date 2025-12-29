package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationClickSendSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.ClickSendSMS
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My ClickSend Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My ClickSend Alert\",\"clicksendsmsLogin\":\"testuser\",\"clicksendsmsPassword\":\"apikey123\",\"clicksendsmsToNumber\":\"61412345678\",\"clicksendsmsSenderName\":\"Uptime\",\"type\":\"clicksendsms\"}"}`,
			),

			want: notification.ClickSendSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My ClickSend Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ClickSendSMSDetails: notification.ClickSendSMSDetails{
					Login:      "testuser",
					Password:   "apikey123",
					ToNumber:   "61412345678",
					SenderName: "Uptime",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"clicksendsmsLogin":"testuser","clicksendsmsPassword":"apikey123","clicksendsmsSenderName":"Uptime","clicksendsmsToNumber":"61412345678","id":1,"isDefault":true,"name":"My ClickSend Alert","type":"clicksendsms","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple ClickSend","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple ClickSend\",\"clicksendsmsLogin\":\"user\",\"clicksendsmsPassword\":\"key\",\"clicksendsmsToNumber\":\"61400000000\",\"clicksendsmsSenderName\":\"Alert\",\"type\":\"clicksendsms\"}"}`,
			),

			want: notification.ClickSendSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple ClickSend",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ClickSendSMSDetails: notification.ClickSendSMSDetails{
					Login:      "user",
					Password:   "key",
					ToNumber:   "61400000000",
					SenderName: "Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"clicksendsmsLogin":"user","clicksendsmsPassword":"key","clicksendsmsSenderName":"Alert","clicksendsmsToNumber":"61400000000","id":2,"isDefault":false,"name":"Simple ClickSend","type":"clicksendsms","userId":1}`,
		},
		{
			name: "with different recipient",
			data: []byte(
				`{"id":3,"name":"ClickSend US","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"ClickSend US\",\"clicksendsmsLogin\":\"ususer\",\"clicksendsmsPassword\":\"usakey\",\"clicksendsmsToNumber\":\"12125551234\",\"clicksendsmsSenderName\":\"Monitor\",\"type\":\"clicksendsms\"}"}`,
			),

			want: notification.ClickSendSMS{
				Base: notification.Base{
					ID:            3,
					Name:          "ClickSend US",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ClickSendSMSDetails: notification.ClickSendSMSDetails{
					Login:      "ususer",
					Password:   "usakey",
					ToNumber:   "12125551234",
					SenderName: "Monitor",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"clicksendsmsLogin":"ususer","clicksendsmsPassword":"usakey","clicksendsmsSenderName":"Monitor","clicksendsmsToNumber":"12125551234","id":3,"isDefault":false,"name":"ClickSend US","type":"clicksendsms","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clicksendsms := notification.ClickSendSMS{}

			err := json.Unmarshal(tc.data, &clicksendsms)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, clicksendsms)

			data, err := json.Marshal(clicksendsms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
