package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationApprise_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Apprise
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Apprise Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Apprise Alert\",\"appriseURL\":\"json://localhost:8080\",\"title\":\"Uptime Kuma Alert\",\"type\":\"apprise\"}"}`),

			want: notification.Apprise{
				Base: notification.Base{
					ID:            1,
					Name:          "My Apprise Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				AppriseDetails: notification.AppriseDetails{
					AppriseURL: "json://localhost:8080",
					Title:      "Uptime Kuma Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"appriseURL":"json://localhost:8080","id":1,"isDefault":true,"name":"My Apprise Alert","title":"Uptime Kuma Alert","type":"apprise","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Apprise","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Apprise\",\"appriseURL\":\"json://localhost:8080\",\"type\":\"apprise\"}"}`),

			want: notification.Apprise{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Apprise",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				AppriseDetails: notification.AppriseDetails{
					AppriseURL: "json://localhost:8080",
					Title:      "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"appriseURL":"json://localhost:8080","id":2,"isDefault":false,"name":"Simple Apprise","title":"","type":"apprise","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apprise := notification.Apprise{}

			err := json.Unmarshal(tc.data, &apprise)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, apprise)

			data, err := json.Marshal(apprise)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
