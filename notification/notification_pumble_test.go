package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPumble_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Pumble
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Pumble Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Pumble Alert\",\"webhookURL\":\"https://pumble.com/webhook/xxx\",\"type\":\"Pumble\"}"}`),

			want: notification.Pumble{
				Base: notification.Base{
					ID:            1,
					Name:          "My Pumble Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PumbleDetails: notification.PumbleDetails{
					WebhookURL: "https://pumble.com/webhook/xxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Pumble Alert","type":"Pumble","userId":1,"webhookURL":"https://pumble.com/webhook/xxx"}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Pumble","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Pumble\",\"webhookURL\":\"https://pumble.com/webhook/yyy\",\"type\":\"Pumble\"}"}`),

			want: notification.Pumble{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Pumble",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PumbleDetails: notification.PumbleDetails{
					WebhookURL: "https://pumble.com/webhook/yyy",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Pumble","type":"Pumble","userId":1,"webhookURL":"https://pumble.com/webhook/yyy"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pumble := notification.Pumble{}

			err := json.Unmarshal(tc.data, &pumble)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pumble)

			data, err := json.Marshal(pumble)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
