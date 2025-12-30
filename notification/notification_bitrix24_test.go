package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationBitrix24_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Bitrix24
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Bitrix24 Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Bitrix24 Alert\",\"bitrix24WebhookURL\":\"https://bitrix24.example.com/rest/1/webhook/\",\"bitrix24UserID\":\"user123\",\"type\":\"Bitrix24\"}"}`,
			),

			want: notification.Bitrix24{
				Base: notification.Base{
					ID:            1,
					Name:          "My Bitrix24 Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				Bitrix24Details: notification.Bitrix24Details{
					WebhookURL:         "https://bitrix24.example.com/rest/1/webhook/",
					NotificationUserID: "user123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"bitrix24UserID":"user123","bitrix24WebhookURL":"https://bitrix24.example.com/rest/1/webhook/","id":1,"isDefault":true,"name":"My Bitrix24 Alert","type":"Bitrix24","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Bitrix24","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Bitrix24\",\"bitrix24WebhookURL\":\"https://bitrix24.example.com/rest/1/webhook/\",\"type\":\"Bitrix24\"}"}`,
			),

			want: notification.Bitrix24{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Bitrix24",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Bitrix24Details: notification.Bitrix24Details{
					WebhookURL: "https://bitrix24.example.com/rest/1/webhook/",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"bitrix24UserID":"","bitrix24WebhookURL":"https://bitrix24.example.com/rest/1/webhook/","id":2,"isDefault":false,"name":"Simple Bitrix24","type":"Bitrix24","userId":1}`,
		},
		{
			name: "inactive notification",
			data: []byte(
				`{"id":3,"name":"Inactive Bitrix24","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Inactive Bitrix24\",\"bitrix24WebhookURL\":\"https://bitrix24.example.com/rest/1/webhook/\",\"bitrix24UserID\":\"admin\",\"type\":\"Bitrix24\"}"}`,
			),

			want: notification.Bitrix24{
				Base: notification.Base{
					ID:            3,
					Name:          "Inactive Bitrix24",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Bitrix24Details: notification.Bitrix24Details{
					WebhookURL:         "https://bitrix24.example.com/rest/1/webhook/",
					NotificationUserID: "admin",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"bitrix24UserID":"admin","bitrix24WebhookURL":"https://bitrix24.example.com/rest/1/webhook/","id":3,"isDefault":false,"name":"Inactive Bitrix24","type":"Bitrix24","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bitrix24 := notification.Bitrix24{}

			err := json.Unmarshal(tc.data, &bitrix24)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, bitrix24)

			data, err := json.Marshal(bitrix24)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
