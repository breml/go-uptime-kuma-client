package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationStackfield_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Stackfield
		wantJSON string
	}{
		{
			name: "success with webhook URL",
			data: []byte(`{"id":1,"name":"My Stackfield Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Stackfield Alert\",\"stackfieldwebhookURL\":\"https://app.stackfield.com/webhook/v1/xxx\",\"type\":\"stackfield\"}"}`),

			want: notification.Stackfield{
				Base: notification.Base{
					ID:            1,
					Name:          "My Stackfield Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				StackfieldDetails: notification.StackfieldDetails{
					WebhookURL: "https://app.stackfield.com/webhook/v1/xxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Stackfield Alert","stackfieldwebhookURL":"https://app.stackfield.com/webhook/v1/xxx","type":"stackfield","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple Stackfield","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Stackfield\",\"stackfieldwebhookURL\":\"https://webhook.example.com\",\"type\":\"stackfield\"}"}`),

			want: notification.Stackfield{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Stackfield",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				StackfieldDetails: notification.StackfieldDetails{
					WebhookURL: "https://webhook.example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Stackfield","stackfieldwebhookURL":"https://webhook.example.com","type":"stackfield","userId":1}`,
		},
		{
			name: "with different webhook endpoint",
			data: []byte(`{"id":3,"name":"Stackfield Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Stackfield Custom\",\"stackfieldwebhookURL\":\"https://custom.stackfield.instance/webhook/abc123\",\"type\":\"stackfield\"}"}`),

			want: notification.Stackfield{
				Base: notification.Base{
					ID:            3,
					Name:          "Stackfield Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				StackfieldDetails: notification.StackfieldDetails{
					WebhookURL: "https://custom.stackfield.instance/webhook/abc123",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Stackfield Custom","stackfieldwebhookURL":"https://custom.stackfield.instance/webhook/abc123","type":"stackfield","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stackfield := notification.Stackfield{}

			err := json.Unmarshal(tc.data, &stackfield)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, stackfield)

			data, err := json.Marshal(stackfield)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
