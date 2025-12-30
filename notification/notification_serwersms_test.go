package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSerwerSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SerwerSMS
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My SerwerSMS Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SerwerSMS Alert\",\"serwersmsUsername\":\"testuser\",\"serwersmsPassword\":\"testpass\",\"serwersmsPhoneNumber\":\"48123456789\",\"serwersmsSenderName\":\"Uptime\",\"type\":\"serwersms\"}"}`,
			),

			want: notification.SerwerSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My SerwerSMS Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SerwerSMSDetails: notification.SerwerSMSDetails{
					Username:    "testuser",
					Password:    "testpass",
					PhoneNumber: "48123456789",
					SenderName:  "Uptime",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SerwerSMS Alert","serwersmsPassword":"testpass","serwersmsPhoneNumber":"48123456789","serwersmsSenderName":"Uptime","serwersmsUsername":"testuser","type":"serwersms","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple SerwerSMS","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SerwerSMS\",\"serwersmsUsername\":\"user\",\"serwersmsPassword\":\"pass\",\"serwersmsPhoneNumber\":\"48999999999\",\"serwersmsSenderName\":\"Alert\",\"type\":\"serwersms\"}"}`,
			),

			want: notification.SerwerSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SerwerSMS",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SerwerSMSDetails: notification.SerwerSMSDetails{
					Username:    "user",
					Password:    "pass",
					PhoneNumber: "48999999999",
					SenderName:  "Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SerwerSMS","serwersmsPassword":"pass","serwersmsPhoneNumber":"48999999999","serwersmsSenderName":"Alert","serwersmsUsername":"user","type":"serwersms","userId":1}`,
		},
		{
			name: "with different phone number",
			data: []byte(
				`{"id":3,"name":"SerwerSMS International","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SerwerSMS International\",\"serwersmsUsername\":\"intluser\",\"serwersmsPassword\":\"intlpass\",\"serwersmsPhoneNumber\":\"358501234567\",\"serwersmsSenderName\":\"GlobalAlert\",\"type\":\"serwersms\"}"}`,
			),

			want: notification.SerwerSMS{
				Base: notification.Base{
					ID:            3,
					Name:          "SerwerSMS International",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SerwerSMSDetails: notification.SerwerSMSDetails{
					Username:    "intluser",
					Password:    "intlpass",
					PhoneNumber: "358501234567",
					SenderName:  "GlobalAlert",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SerwerSMS International","serwersmsPassword":"intlpass","serwersmsPhoneNumber":"358501234567","serwersmsSenderName":"GlobalAlert","serwersmsUsername":"intluser","type":"serwersms","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serwersms := notification.SerwerSMS{}

			err := json.Unmarshal(tc.data, &serwersms)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, serwersms)

			data, err := json.Marshal(serwersms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
