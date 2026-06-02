package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSIR_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.SMSIR
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My SMSIR Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMSIR Alert\",\"smsirApiKey\":\"test-api-key\",\"smsirNumber\":\"9123456789\",\"smsirTemplate\":\"12345\",\"type\":\"smsir\"}"}`,
			),

			want: notification.SMSIR{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMSIR Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSIRDetails: notification.SMSIRDetails{
					APIKey:   "test-api-key",
					Number:   "9123456789",
					Template: "12345",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMSIR Alert","smsirApiKey":"test-api-key","smsirNumber":"9123456789","smsirTemplate":"12345","type":"smsir","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple SMSIR","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMSIR\",\"smsirApiKey\":\"minimal-key\",\"smsirNumber\":\"09123456789\",\"smsirTemplate\":\"54321\",\"type\":\"smsir\"}"}`,
			),

			want: notification.SMSIR{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMSIR",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSIRDetails: notification.SMSIRDetails{
					APIKey:   "minimal-key",
					Number:   "09123456789",
					Template: "54321",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SMSIR","smsirApiKey":"minimal-key","smsirNumber":"09123456789","smsirTemplate":"54321","type":"smsir","userId":1}`,
		},
		{
			name: "multiple recipients",
			data: []byte(
				`{"id":3,"name":"SMSIR Multi","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSIR Multi\",\"smsirApiKey\":\"api-key-multi\",\"smsirNumber\":\"9123456789,09987654321\",\"smsirTemplate\":\"99999\",\"type\":\"smsir\"}"}`,
			),

			want: notification.SMSIR{
				Base: notification.Base{
					ID:            3,
					Name:          "SMSIR Multi",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSIRDetails: notification.SMSIRDetails{
					APIKey:   "api-key-multi",
					Number:   "9123456789,09987654321",
					Template: "99999",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SMSIR Multi","smsirApiKey":"api-key-multi","smsirNumber":"9123456789,09987654321","smsirTemplate":"99999","type":"smsir","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smsir := notification.SMSIR{}

			err := json.Unmarshal(tc.data, &smsir)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smsir)

			data, err := json.Marshal(smsir)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
