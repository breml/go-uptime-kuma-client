package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSPlanet_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.SMSPlanet
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My SMS Planet Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMS Planet Alert\",\"smsplanetApiToken\":\"test-token-123\",\"smsplanetPhoneNumbers\":\"48123456789\",\"smsplanetSenderName\":\"Uptime Kuma\",\"type\":\"SMSPlanet\"}"}`),

			want: notification.SMSPlanet{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMS Planet Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSPlanetDetails: notification.SMSPlanetDetails{
					APIToken:     "test-token-123",
					PhoneNumbers: "48123456789",
					SenderName:   "Uptime Kuma",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMS Planet Alert","smsplanetApiToken":"test-token-123","smsplanetPhoneNumbers":"48123456789","smsplanetSenderName":"Uptime Kuma","type":"SMSPlanet","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple SMS Planet","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMS Planet\",\"smsplanetApiToken\":\"minimal-token\",\"smsplanetPhoneNumbers\":\"48987654321\",\"smsplanetSenderName\":\"Alert\",\"type\":\"SMSPlanet\"}"}`),

			want: notification.SMSPlanet{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMS Planet",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSPlanetDetails: notification.SMSPlanetDetails{
					APIToken:     "minimal-token",
					PhoneNumbers: "48987654321",
					SenderName:   "Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SMS Planet","smsplanetApiToken":"minimal-token","smsplanetPhoneNumbers":"48987654321","smsplanetSenderName":"Alert","type":"SMSPlanet","userId":1}`,
		},
		{
			name: "with multiple phone numbers",
			data: []byte(`{"id":3,"name":"SMS Planet Multi","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMS Planet Multi\",\"smsplanetApiToken\":\"multi-token\",\"smsplanetPhoneNumbers\":\"48111111111,48222222222\",\"smsplanetSenderName\":\"Monitor\",\"type\":\"SMSPlanet\"}"}`),

			want: notification.SMSPlanet{
				Base: notification.Base{
					ID:            3,
					Name:          "SMS Planet Multi",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSPlanetDetails: notification.SMSPlanetDetails{
					APIToken:     "multi-token",
					PhoneNumbers: "48111111111,48222222222",
					SenderName:   "Monitor",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SMS Planet Multi","smsplanetApiToken":"multi-token","smsplanetPhoneNumbers":"48111111111,48222222222","smsplanetSenderName":"Monitor","type":"SMSPlanet","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smsplanet := notification.SMSPlanet{}

			err := json.Unmarshal(tc.data, &smsplanet)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smsplanet)

			data, err := json.Marshal(smsplanet)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
