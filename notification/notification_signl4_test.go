package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSIGNL4_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SIGNL4
		wantJSON string
	}{
		{
			name: "success with webhook URL",
			data: []byte(
				`{"id":1,"name":"My SIGNL4 Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SIGNL4 Alert\",\"webhookURL\":\"https://connect.signl4.com/webhook/webhook-uuid-here\",\"type\":\"SIGNL4\"}"}`,
			),

			want: notification.SIGNL4{
				Base: notification.Base{
					ID:            1,
					Name:          "My SIGNL4 Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SIGNL4Details: notification.SIGNL4Details{
					WebhookURL: "https://connect.signl4.com/webhook/webhook-uuid-here",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SIGNL4 Alert","type":"SIGNL4","userId":1,"webhookURL":"https://connect.signl4.com/webhook/webhook-uuid-here"}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple SIGNL4","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SIGNL4\",\"webhookURL\":\"https://connect.signl4.com/webhook/simple-webhook\",\"type\":\"SIGNL4\"}"}`,
			),

			want: notification.SIGNL4{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SIGNL4",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SIGNL4Details: notification.SIGNL4Details{
					WebhookURL: "https://connect.signl4.com/webhook/simple-webhook",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SIGNL4","type":"SIGNL4","userId":1,"webhookURL":"https://connect.signl4.com/webhook/simple-webhook"}`,
		},
		{
			name: "with different webhook endpoint",
			data: []byte(
				`{"id":3,"name":"SIGNL4 Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SIGNL4 Custom\",\"webhookURL\":\"https://custom.signl4.endpoint.com/webhook/custom-id\",\"type\":\"SIGNL4\"}"}`,
			),

			want: notification.SIGNL4{
				Base: notification.Base{
					ID:            3,
					Name:          "SIGNL4 Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SIGNL4Details: notification.SIGNL4Details{
					WebhookURL: "https://custom.signl4.endpoint.com/webhook/custom-id",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SIGNL4 Custom","type":"SIGNL4","userId":1,"webhookURL":"https://custom.signl4.endpoint.com/webhook/custom-id"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signl4 := notification.SIGNL4{}

			err := json.Unmarshal(tc.data, &signl4)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, signl4)

			data, err := json.Marshal(signl4)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
