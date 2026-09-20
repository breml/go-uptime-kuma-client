package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSGateway_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.SMSGateway
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My SMS Gateway Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMS Gateway Alert\",\"smsgatewayUrl\":\"http://localhost:8080\",\"smsgatewayApiKey\":\"test-api-key\",\"smsgatewayTo\":\"+15551234567, +15559876543\",\"type\":\"SMSGateway\"}"}`,
			),

			want: notification.SMSGateway{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMS Gateway Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSGatewayDetails: notification.SMSGatewayDetails{
					URL:        "http://localhost:8080",
					APIKey:     "test-api-key",
					Recipients: "+15551234567, +15559876543",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMS Gateway Alert","smsgatewayUrl":"http://localhost:8080","smsgatewayApiKey":"test-api-key","smsgatewayTo":"+15551234567, +15559876543","type":"SMSGateway","userId":1}`,
		},
		{
			name: "minimal configuration with single recipient",
			data: []byte(
				`{"id":2,"name":"Simple SMS Gateway","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMS Gateway\",\"smsgatewayUrl\":\"https://gateway.example.com\",\"smsgatewayApiKey\":\"simple-key\",\"smsgatewayTo\":\"+15551234567\",\"type\":\"SMSGateway\"}"}`,
			),

			want: notification.SMSGateway{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMS Gateway",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSGatewayDetails: notification.SMSGatewayDetails{
					URL:        "https://gateway.example.com",
					APIKey:     "simple-key",
					Recipients: "+15551234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SMS Gateway","smsgatewayUrl":"https://gateway.example.com","smsgatewayApiKey":"simple-key","smsgatewayTo":"+15551234567","type":"SMSGateway","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smsgateway := notification.SMSGateway{}

			err := json.Unmarshal(tc.data, &smsgateway)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smsgateway)

			data, err := json.Marshal(smsgateway)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
