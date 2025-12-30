package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationZohoCliq_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.ZohoCliq
		wantJSON string
	}{
		{
			name: "success with webhook URL",
			data: []byte(
				`{"id":1,"name":"My ZohoCliq Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My ZohoCliq Alert\",\"webhookUrl\":\"https://zoho-cliq.example.com/webhook/abc123\",\"type\":\"ZohoCliq\"}"}`,
			),

			want: notification.ZohoCliq{
				Base: notification.Base{
					ID:            1,
					Name:          "My ZohoCliq Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ZohoCliqDetails: notification.ZohoCliqDetails{
					WebhookURL: "https://zoho-cliq.example.com/webhook/abc123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My ZohoCliq Alert","type":"ZohoCliq","userId":1,"webhookUrl":"https://zoho-cliq.example.com/webhook/abc123"}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple ZohoCliq","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple ZohoCliq\",\"webhookUrl\":\"https://zoho.example.com/webhook\",\"type\":\"ZohoCliq\"}"}`,
			),

			want: notification.ZohoCliq{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple ZohoCliq",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ZohoCliqDetails: notification.ZohoCliqDetails{
					WebhookURL: "https://zoho.example.com/webhook",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple ZohoCliq","type":"ZohoCliq","userId":1,"webhookUrl":"https://zoho.example.com/webhook"}`,
		},
		{
			name: "inactive notification",
			data: []byte(
				`{"id":3,"name":"Inactive ZohoCliq","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Inactive ZohoCliq\",\"webhookUrl\":\"https://example.com/cliq/webhook/xyz789\",\"type\":\"ZohoCliq\"}"}`,
			),

			want: notification.ZohoCliq{
				Base: notification.Base{
					ID:            3,
					Name:          "Inactive ZohoCliq",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ZohoCliqDetails: notification.ZohoCliqDetails{
					WebhookURL: "https://example.com/cliq/webhook/xyz789",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Inactive ZohoCliq","type":"ZohoCliq","userId":1,"webhookUrl":"https://example.com/cliq/webhook/xyz789"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zohocliq := notification.ZohoCliq{}

			err := json.Unmarshal(tc.data, &zohocliq)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, zohocliq)

			data, err := json.Marshal(zohocliq)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
