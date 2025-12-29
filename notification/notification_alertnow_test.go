package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationAlertNow_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.AlertNow
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My AlertNow Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My AlertNow Alert\",\"alertNowWebhookURL\":\"https://alertnow.example.com/api/webhook\",\"type\":\"AlertNow\"}"}`,
			),

			want: notification.AlertNow{
				Base: notification.Base{
					ID:            1,
					Name:          "My AlertNow Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				AlertNowDetails: notification.AlertNowDetails{
					WebhookURL: "https://alertnow.example.com/api/webhook",
				},
			},
			wantJSON: `{"active":true,"alertNowWebhookURL":"https://alertnow.example.com/api/webhook","applyExisting":true,"id":1,"isDefault":true,"name":"My AlertNow Alert","type":"AlertNow","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple AlertNow","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple AlertNow\",\"alertNowWebhookURL\":\"https://api.alertnow.io/webhook\",\"type\":\"AlertNow\"}"}`,
			),

			want: notification.AlertNow{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple AlertNow",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				AlertNowDetails: notification.AlertNowDetails{
					WebhookURL: "https://api.alertnow.io/webhook",
				},
			},
			wantJSON: `{"active":true,"alertNowWebhookURL":"https://api.alertnow.io/webhook","applyExisting":false,"id":2,"isDefault":false,"name":"Simple AlertNow","type":"AlertNow","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			alertnow := notification.AlertNow{}

			err := json.Unmarshal(tc.data, &alertnow)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, alertnow)

			data, err := json.Marshal(alertnow)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
