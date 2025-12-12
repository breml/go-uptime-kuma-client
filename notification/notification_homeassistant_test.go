package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationHomeAssistant_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.HomeAssistant
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Home Assistant Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Home Assistant Alert\",\"homeAssistantUrl\":\"http://192.168.1.100:8123\",\"longLivedAccessToken\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\",\"notificationService\":\"notify.mobile_app_iphone\",\"type\":\"HomeAssistant\"}"}`),

			want: notification.HomeAssistant{
				Base: notification.Base{
					ID:            1,
					Name:          "My Home Assistant Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				HomeAssistantDetails: notification.HomeAssistantDetails{
					HomeAssistantURL:     "http://192.168.1.100:8123",
					LongLivedAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
					NotificationService:  "notify.mobile_app_iphone",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"homeAssistantUrl":"http://192.168.1.100:8123","id":1,"isDefault":true,"longLivedAccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","name":"My Home Assistant Alert","notificationService":"notify.mobile_app_iphone","type":"HomeAssistant","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Home Assistant","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Home Assistant\",\"homeAssistantUrl\":\"http://localhost:8123\",\"longLivedAccessToken\":\"token123\",\"type\":\"HomeAssistant\"}"}`),

			want: notification.HomeAssistant{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Home Assistant",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				HomeAssistantDetails: notification.HomeAssistantDetails{
					HomeAssistantURL:     "http://localhost:8123",
					LongLivedAccessToken: "token123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"homeAssistantUrl":"http://localhost:8123","id":2,"isDefault":false,"longLivedAccessToken":"token123","name":"Simple Home Assistant","notificationService":"","type":"HomeAssistant","userId":1}`,
		},
		{
			name: "with custom notification service",
			data: []byte(`{"id":3,"name":"HA with Custom Service","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"HA with Custom Service\",\"homeAssistantUrl\":\"http://ha.example.com:8123\",\"longLivedAccessToken\":\"custom-token-456\",\"notificationService\":\"notify.persistent_notification\",\"type\":\"HomeAssistant\"}"}`),

			want: notification.HomeAssistant{
				Base: notification.Base{
					ID:            3,
					Name:          "HA with Custom Service",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				HomeAssistantDetails: notification.HomeAssistantDetails{
					HomeAssistantURL:     "http://ha.example.com:8123",
					LongLivedAccessToken: "custom-token-456",
					NotificationService:  "notify.persistent_notification",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"homeAssistantUrl":"http://ha.example.com:8123","id":3,"isDefault":false,"longLivedAccessToken":"custom-token-456","name":"HA with Custom Service","notificationService":"notify.persistent_notification","type":"HomeAssistant","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			homeassistant := notification.HomeAssistant{}

			err := json.Unmarshal(tc.data, &homeassistant)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, homeassistant)

			data, err := json.Marshal(homeassistant)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
