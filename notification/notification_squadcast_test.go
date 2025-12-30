package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSquadcast_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Squadcast
		wantJSON string
	}{
		{
			name: "success with webhook URL",
			data: []byte(
				`{"id":1,"name":"My Squadcast Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Squadcast Alert\",\"squadcastWebhookURL\":\"https://api.squadcast.com/api/v3/incidents/webhook\",\"type\":\"squadcast\"}"}`,
			),

			want: notification.Squadcast{
				Base: notification.Base{
					ID:            1,
					Name:          "My Squadcast Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SquadcastDetails: notification.SquadcastDetails{
					WebhookURL: "https://api.squadcast.com/api/v3/incidents/webhook",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Squadcast Alert","squadcastWebhookURL":"https://api.squadcast.com/api/v3/incidents/webhook","type":"squadcast","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Squadcast","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Squadcast\",\"squadcastWebhookURL\":\"https://webhook.example.com\",\"type\":\"squadcast\"}"}`,
			),

			want: notification.Squadcast{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Squadcast",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SquadcastDetails: notification.SquadcastDetails{
					WebhookURL: "https://webhook.example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Squadcast","squadcastWebhookURL":"https://webhook.example.com","type":"squadcast","userId":1}`,
		},
		{
			name: "with different webhook endpoint",
			data: []byte(
				`{"id":3,"name":"Squadcast Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Squadcast Custom\",\"squadcastWebhookURL\":\"https://custom.squadcast.instance.com/webhook/xyz123\",\"type\":\"squadcast\"}"}`,
			),

			want: notification.Squadcast{
				Base: notification.Base{
					ID:            3,
					Name:          "Squadcast Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SquadcastDetails: notification.SquadcastDetails{
					WebhookURL: "https://custom.squadcast.instance.com/webhook/xyz123",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Squadcast Custom","squadcastWebhookURL":"https://custom.squadcast.instance.com/webhook/xyz123","type":"squadcast","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			squadcast := notification.Squadcast{}

			err := json.Unmarshal(tc.data, &squadcast)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, squadcast)

			data, err := json.Marshal(squadcast)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
