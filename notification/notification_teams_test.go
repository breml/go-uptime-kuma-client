package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
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
			data: []byte(
				`{"id":1,"name":"My Teams Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Teams Alert\",\"webhookUrl\":\"https://outlook.office.com/webhook/xxx\",\"type\":\"teams\"}"}`,
			),

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
		{
			name: "with tags enabled",
			data: []byte(
				`{"id":2,"name":"Teams With Tags","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Teams With Tags\",\"webhookUrl\":\"https://outlook.office.com/webhook/yyy\",\"teamsEnableTags\":true,\"type\":\"teams\"}"}`,
			),

			want: notification.Teams{
				Base: notification.Base{
					ID:            2,
					Name:          "Teams With Tags",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TeamsDetails: notification.TeamsDetails{
					WebhookURL: "https://outlook.office.com/webhook/yyy",
					EnableTags: ptr.To(true),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Teams With Tags","teamsEnableTags":true,"type":"teams","userId":1,"webhookUrl":"https://outlook.office.com/webhook/yyy"}`,
		},
		{
			name: "with tags explicitly disabled",
			data: []byte(
				`{"id":3,"name":"Teams Without Tags","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Teams Without Tags\",\"webhookUrl\":\"https://outlook.office.com/webhook/zzz\",\"teamsEnableTags\":false,\"type\":\"teams\"}"}`,
			),

			want: notification.Teams{
				Base: notification.Base{
					ID:            3,
					Name:          "Teams Without Tags",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TeamsDetails: notification.TeamsDetails{
					WebhookURL: "https://outlook.office.com/webhook/zzz",
					EnableTags: ptr.To(false),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Teams Without Tags","teamsEnableTags":false,"type":"teams","userId":1,"webhookUrl":"https://outlook.office.com/webhook/zzz"}`,
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
