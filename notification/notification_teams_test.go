package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTeams_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Teams
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Teams Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Teams Alert\",\"webhookUrl\":\"https://outlook.office.com/webhook/xxx\",\"type\":\"teams\"}"}`),

			want: notification.Teams{
				Base: notification.Base{
					ID:            1,
					Name:          "My Teams Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TeamsDetails: notification.TeamsDetails{
					WebhookURL: "https://outlook.office.com/webhook/xxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Teams Alert","type":"teams","userId":1,"webhookUrl":"https://outlook.office.com/webhook/xxx"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teams := notification.Teams{}

			err := json.Unmarshal(tc.data, &teams)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, teams)

			data, err := json.Marshal(teams)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
