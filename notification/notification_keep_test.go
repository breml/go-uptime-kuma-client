package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationKeep_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Keep
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Keep Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Keep Alert\",\"webhookURL\":\"https://keep.example.com/webhook\",\"webhookAPIKey\":\"api-key-123\",\"type\":\"Keep\"}"}`,
			),

			want: notification.Keep{
				Base: notification.Base{
					ID:            1,
					Name:          "My Keep Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				KeepDetails: notification.KeepDetails{
					WebhookURL: "https://keep.example.com/webhook",
					APIKey:     "api-key-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Keep Alert","type":"Keep","userId":1,"webhookAPIKey":"api-key-123","webhookURL":"https://keep.example.com/webhook"}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Keep","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Keep\",\"webhookURL\":\"https://api.example.com\",\"webhookAPIKey\":\"key-abc\",\"type\":\"Keep\"}"}`,
			),

			want: notification.Keep{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Keep",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				KeepDetails: notification.KeepDetails{
					WebhookURL: "https://api.example.com",
					APIKey:     "key-abc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Keep","type":"Keep","userId":1,"webhookAPIKey":"key-abc","webhookURL":"https://api.example.com"}`,
		},
		{
			name: "with trailing slash in URL",
			data: []byte(
				`{"id":3,"name":"Keep Trailing","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Keep Trailing\",\"webhookURL\":\"https://alerts.example.com/keep/\",\"webhookAPIKey\":\"key-xyz\",\"type\":\"Keep\"}"}`,
			),

			want: notification.Keep{
				Base: notification.Base{
					ID:            3,
					Name:          "Keep Trailing",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				KeepDetails: notification.KeepDetails{
					WebhookURL: "https://alerts.example.com/keep/",
					APIKey:     "key-xyz",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Keep Trailing","type":"Keep","userId":1,"webhookAPIKey":"key-xyz","webhookURL":"https://alerts.example.com/keep/"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			keep := notification.Keep{}

			err := json.Unmarshal(tc.data, &keep)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, keep)

			data, err := json.Marshal(keep)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
