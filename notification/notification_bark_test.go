package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationBark_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Bark
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My Bark Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Bark Alert\",\"barkEndpoint\":\"https://bark.example.com\",\"barkGroup\":\"Monitoring\",\"barkSound\":\"alarm\",\"apiVersion\":\"v1\",\"type\":\"bark\"}"}`),

			want: notification.Bark{
				Base: notification.Base{
					ID:            1,
					Name:          "My Bark Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				BarkDetails: notification.BarkDetails{
					Endpoint:   "https://bark.example.com",
					Group:      "Monitoring",
					Sound:      "alarm",
					APIVersion: "v1",
				},
			},
			wantJSON: `{"active":true,"apiVersion":"v1","applyExisting":true,"barkEndpoint":"https://bark.example.com","barkGroup":"Monitoring","barkSound":"alarm","id":1,"isDefault":true,"name":"My Bark Alert","type":"bark","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple Bark","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Bark\",\"barkEndpoint\":\"https://bark.example.com\",\"type\":\"bark\"}"}`),

			want: notification.Bark{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Bark",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BarkDetails: notification.BarkDetails{
					Endpoint: "https://bark.example.com",
				},
			},
			wantJSON: `{"active":true,"apiVersion":"","applyExisting":false,"barkEndpoint":"https://bark.example.com","barkGroup":"","barkSound":"","id":2,"isDefault":false,"name":"Simple Bark","type":"bark","userId":1}`,
		},
		{
			name: "with v2 API version",
			data: []byte(`{"id":3,"name":"Bark v2","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Bark v2\",\"barkEndpoint\":\"https://bark.example.com\",\"barkGroup\":\"UptimeKuma\",\"barkSound\":\"telegraph\",\"apiVersion\":\"v2\",\"type\":\"bark\"}"}`),

			want: notification.Bark{
				Base: notification.Base{
					ID:            3,
					Name:          "Bark v2",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BarkDetails: notification.BarkDetails{
					Endpoint:   "https://bark.example.com",
					Group:      "UptimeKuma",
					Sound:      "telegraph",
					APIVersion: "v2",
				},
			},
			wantJSON: `{"active":true,"apiVersion":"v2","applyExisting":false,"barkEndpoint":"https://bark.example.com","barkGroup":"UptimeKuma","barkSound":"telegraph","id":3,"isDefault":false,"name":"Bark v2","type":"bark","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bark := notification.Bark{}

			err := json.Unmarshal(tc.data, &bark)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, bark)

			data, err := json.Marshal(bark)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
