package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSC_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SMSC
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My SMSC Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMSC Alert\",\"smscLogin\":\"testuser\",\"smscPassword\":\"testpass\",\"smscToNumber\":\"77123456789\",\"smscSenderName\":\"Uptime\",\"smscTranslit\":\"1\",\"type\":\"smsc\"}"}`),

			want: notification.SMSC{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMSC Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSCDetails: notification.SMSCDetails{
					Login:      "testuser",
					Password:   "testpass",
					ToNumber:   "77123456789",
					SenderName: "Uptime",
					Translit:   "1",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMSC Alert","smscLogin":"testuser","smscPassword":"testpass","smscSenderName":"Uptime","smscToNumber":"77123456789","smscTranslit":"1","type":"smsc","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple SMSC","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMSC\",\"smscLogin\":\"user2\",\"smscPassword\":\"pass2\",\"smscToNumber\":\"77987654321\",\"smscSenderName\":\"\",\"smscTranslit\":\"0\",\"type\":\"smsc\"}"}`),

			want: notification.SMSC{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMSC",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSCDetails: notification.SMSCDetails{
					Login:      "user2",
					Password:   "pass2",
					ToNumber:   "77987654321",
					SenderName: "",
					Translit:   "0",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SMSC","smscLogin":"user2","smscPassword":"pass2","smscSenderName":"","smscToNumber":"77987654321","smscTranslit":"0","type":"smsc","userId":1}`,
		},
		{
			name: "with translit enabled",
			data: []byte(`{"id":3,"name":"SMSC Translit","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSC Translit\",\"smscLogin\":\"user3\",\"smscPassword\":\"pass3\",\"smscToNumber\":\"77111222333\",\"smscSenderName\":\"Alert\",\"smscTranslit\":\"1\",\"type\":\"smsc\"}"}`),

			want: notification.SMSC{
				Base: notification.Base{
					ID:            3,
					Name:          "SMSC Translit",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSCDetails: notification.SMSCDetails{
					Login:      "user3",
					Password:   "pass3",
					ToNumber:   "77111222333",
					SenderName: "Alert",
					Translit:   "1",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SMSC Translit","smscLogin":"user3","smscPassword":"pass3","smscSenderName":"Alert","smscToNumber":"77111222333","smscTranslit":"1","type":"smsc","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smsc := notification.SMSC{}

			err := json.Unmarshal(tc.data, &smsc)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smsc)

			data, err := json.Marshal(smsc)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
