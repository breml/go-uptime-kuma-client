package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTelnyx_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Telnyx
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Telnyx Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Telnyx Alert\",\"telnyxApiKey\":\"KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"telnyxMessagingProfileId\":\"4001763e-7f7d-4c87-a8b1-1c5a0e5a3f48\",\"telnyxPhoneNumber\":\"+15559876543\",\"telnyxToNumber\":\"+15551234567\",\"type\":\"telnyx\"}"}`,
			),

			want: notification.Telnyx{
				Base: notification.Base{
					ID:            1,
					Name:          "My Telnyx Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TelnyxDetails: notification.TelnyxDetails{
					APIKey:             "KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					MessagingProfileID: "4001763e-7f7d-4c87-a8b1-1c5a0e5a3f48",
					PhoneNumber:        "+15559876543",
					ToNumber:           "+15551234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Telnyx Alert","telnyxApiKey":"KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","telnyxMessagingProfileId":"4001763e-7f7d-4c87-a8b1-1c5a0e5a3f48","telnyxPhoneNumber":"+15559876543","telnyxToNumber":"+15551234567","type":"telnyx","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Telnyx","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Telnyx\",\"telnyxApiKey\":\"KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"telnyxPhoneNumber\":\"+15559876543\",\"telnyxToNumber\":\"+15551234567\",\"type\":\"telnyx\"}"}`,
			),

			want: notification.Telnyx{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Telnyx",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TelnyxDetails: notification.TelnyxDetails{
					APIKey:      "KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					PhoneNumber: "+15559876543",
					ToNumber:    "+15551234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Telnyx","telnyxApiKey":"KEYxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","telnyxMessagingProfileId":"","telnyxPhoneNumber":"+15559876543","telnyxToNumber":"+15551234567","type":"telnyx","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			telnyx := notification.Telnyx{}

			err := json.Unmarshal(tc.data, &telnyx)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, telnyx)

			data, err := json.Marshal(telnyx)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
