package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGorush_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Gorush
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Gorush Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Gorush Alert\",\"gorushServerURL\":\"https://gorush.example.com\",\"gorushDeviceToken\":\"device-token-123\",\"gorushPlatform\":\"ios\",\"gorushTitle\":\"Alert\",\"gorushPriority\":\"high\",\"gorushRetry\":3,\"gorushTopic\":\"com.example.app\",\"type\":\"gorush\"}"}`,
			),

			want: notification.Gorush{
				Base: notification.Base{
					ID:            1,
					Name:          "My Gorush Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GorushDetails: notification.GorushDetails{
					ServerURL:   "https://gorush.example.com",
					DeviceToken: "device-token-123",
					Platform:    "ios",
					Title:       "Alert",
					Priority:    "high",
					Retry:       3,
					Topic:       "com.example.app",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"gorushDeviceToken":"device-token-123","gorushPlatform":"ios","gorushPriority":"high","gorushRetry":3,"gorushServerURL":"https://gorush.example.com","gorushTitle":"Alert","gorushTopic":"com.example.app","id":1,"isDefault":true,"name":"My Gorush Alert","type":"gorush","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Gorush","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Gorush\",\"gorushServerURL\":\"https://push.example.com\",\"gorushDeviceToken\":\"token-abc\",\"gorushPlatform\":\"android\",\"gorushTitle\":\"\",\"gorushPriority\":\"\",\"gorushRetry\":0,\"gorushTopic\":\"\",\"type\":\"gorush\"}"}`,
			),

			want: notification.Gorush{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Gorush",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GorushDetails: notification.GorushDetails{
					ServerURL:   "https://push.example.com",
					DeviceToken: "token-abc",
					Platform:    "android",
					Title:       "",
					Priority:    "",
					Retry:       0,
					Topic:       "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"gorushDeviceToken":"token-abc","gorushPlatform":"android","gorushPriority":"","gorushRetry":0,"gorushServerURL":"https://push.example.com","gorushTitle":"","gorushTopic":"","id":2,"isDefault":false,"name":"Simple Gorush","type":"gorush","userId":1}`,
		},
		{
			name: "with huawei platform",
			data: []byte(
				`{"id":3,"name":"Gorush Huawei","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Gorush Huawei\",\"gorushServerURL\":\"https://notification.example.com\",\"gorushDeviceToken\":\"huawei-token-xyz\",\"gorushPlatform\":\"huawei\",\"gorushTitle\":\"System Alert\",\"gorushPriority\":\"critical\",\"gorushRetry\":5,\"gorushTopic\":\"\",\"type\":\"gorush\"}"}`,
			),

			want: notification.Gorush{
				Base: notification.Base{
					ID:            3,
					Name:          "Gorush Huawei",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GorushDetails: notification.GorushDetails{
					ServerURL:   "https://notification.example.com",
					DeviceToken: "huawei-token-xyz",
					Platform:    "huawei",
					Title:       "System Alert",
					Priority:    "critical",
					Retry:       5,
					Topic:       "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"gorushDeviceToken":"huawei-token-xyz","gorushPlatform":"huawei","gorushPriority":"critical","gorushRetry":5,"gorushServerURL":"https://notification.example.com","gorushTitle":"System Alert","gorushTopic":"","id":3,"isDefault":false,"name":"Gorush Huawei","type":"gorush","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gorush := notification.Gorush{}

			err := json.Unmarshal(tc.data, &gorush)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, gorush)

			data, err := json.Marshal(gorush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
