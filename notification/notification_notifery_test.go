package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNotifery_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Notifery
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"Test Notifery","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"notifery\",\"notiferyApiKey\":\"test-api-key\",\"notiferyTitle\":\"Alert\",\"notiferyGroup\":\"monitoring\"}"}`),

			want: notification.Notifery{
				Base: notification.Base{
					ID:        1,
					Name:      "Test Notifery",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				NotiferyDetails: notification.NotiferyDetails{
					APIKey: "test-api-key",
					Title:  "Alert",
					Group:  "monitoring",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"Test Notifery","notiferyApiKey":"test-api-key","notiferyGroup":"monitoring","notiferyTitle":"Alert","type":"notifery","userId":42}`,
		},
		{
			name: "with default title",
			data: []byte(`{"id":2,"name":"Test Notifery Default","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"notifery\",\"notiferyApiKey\":\"secret-key\",\"notiferyTitle\":\"Uptime Kuma Alert\",\"notiferyGroup\":\"\"}"}`),

			want: notification.Notifery{
				Base: notification.Base{
					ID:            2,
					Name:          "Test Notifery Default",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				NotiferyDetails: notification.NotiferyDetails{
					APIKey: "secret-key",
					Title:  "Uptime Kuma Alert",
					Group:  "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":true,"name":"Test Notifery Default","notiferyApiKey":"secret-key","notiferyGroup":"","notiferyTitle":"Uptime Kuma Alert","type":"notifery","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":3,"name":"Test Notifery Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"notifery\",\"notiferyApiKey\":\"\",\"notiferyTitle\":\"\",\"notiferyGroup\":\"\"}"}`),

			want: notification.Notifery{
				Base: notification.Base{
					ID:            3,
					Name:          "Test Notifery Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				NotiferyDetails: notification.NotiferyDetails{
					APIKey: "",
					Title:  "",
					Group:  "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Test Notifery Minimal","notiferyApiKey":"","notiferyGroup":"","notiferyTitle":"","type":"notifery","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			notifery := notification.Notifery{}

			err := json.Unmarshal(tc.data, &notifery)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, notifery)

			data, err := json.Marshal(notifery)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
