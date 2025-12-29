package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTwilio_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Twilio
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Twilio Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Twilio Alert\",\"twilioAccountSID\":\"ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"twilioApiKey\":\"SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"twilioAuthToken\":\"your_auth_token\",\"twilioToNumber\":\"+15551234567\",\"twilioFromNumber\":\"+15559876543\",\"type\":\"twilio\"}"}`,
			),

			want: notification.Twilio{
				Base: notification.Base{
					ID:            1,
					Name:          "My Twilio Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TwilioDetails: notification.TwilioDetails{
					AccountSID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					ApiKey:     "SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					AuthToken:  "your_auth_token",
					ToNumber:   "+15551234567",
					FromNumber: "+15559876543",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Twilio Alert","twilioAccountSID":"ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","twilioApiKey":"SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","twilioAuthToken":"your_auth_token","twilioToNumber":"+15551234567","twilioFromNumber":"+15559876543","type":"twilio","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Twilio","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Twilio\",\"twilioAccountSID\":\"ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"twilioAuthToken\":\"auth_token\",\"twilioToNumber\":\"+15551234567\",\"twilioFromNumber\":\"+15559876543\",\"type\":\"twilio\"}"}`,
			),

			want: notification.Twilio{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Twilio",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TwilioDetails: notification.TwilioDetails{
					AccountSID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					AuthToken:  "auth_token",
					ToNumber:   "+15551234567",
					FromNumber: "+15559876543",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Twilio","twilioAccountSID":"ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","twilioApiKey":"","twilioAuthToken":"auth_token","twilioToNumber":"+15551234567","twilioFromNumber":"+15559876543","type":"twilio","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			twilio := notification.Twilio{}

			err := json.Unmarshal(tc.data, &twilio)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, twilio)

			data, err := json.Marshal(twilio)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
