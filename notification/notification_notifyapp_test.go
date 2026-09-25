package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNotifyApp_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.NotifyApp
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Notify Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Notify Alert\",\"notifyAppDeviceId\":\"ABC12345\",\"notifyAppToken\":\"token-123\",\"notifyAppIconUrl\":\"https://example.com/icon.png\",\"type\":\"notifyapp\"}"}`,
			),

			want: notification.NotifyApp{
				Base: notification.Base{
					ID:            1,
					Name:          "My Notify Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NotifyAppDetails: notification.NotifyAppDetails{
					DeviceID: "ABC12345",
					Token:    "token-123",
					IconURL:  new("https://example.com/icon.png"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Notify Alert","notifyAppDeviceId":"ABC12345","notifyAppIconUrl":"https://example.com/icon.png","notifyAppToken":"token-123","type":"notifyapp","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Notify","active":false,"userId":1,"isDefault":false,"config":"{\"name\":\"Simple Notify\",\"notifyAppDeviceId\":\"XYZ98765\",\"notifyAppToken\":\"token-abc\",\"type\":\"notifyapp\"}"}`,
			),

			want: notification.NotifyApp{
				Base: notification.Base{
					ID:       2,
					Name:     "Simple Notify",
					IsActive: false,
					UserID:   1,
				},
				NotifyAppDetails: notification.NotifyAppDetails{
					DeviceID: "XYZ98765",
					Token:    "token-abc",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Notify","notifyAppDeviceId":"XYZ98765","notifyAppToken":"token-abc","type":"notifyapp","userId":1}`,
		},
		{
			// The device id and token are required in the Uptime Kuma form, so
			// neither carries omitempty and an empty value survives the round
			// trip as an empty key.
			name: "empty required fields are kept",
			data: []byte(
				`{"id":3,"name":"Notify Empty","active":false,"userId":1,"isDefault":false,"config":"{\"notifyAppDeviceId\":\"\",\"notifyAppToken\":\"\",\"type\":\"notifyapp\"}"}`,
			),

			want: notification.NotifyApp{
				Base: notification.Base{
					ID:     3,
					Name:   "Notify Empty",
					UserID: 1,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Notify Empty","notifyAppDeviceId":"","notifyAppToken":"","type":"notifyapp","userId":1}`,
		},
		{
			name:    "missing config field",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false}`),
			wantErr: true,
		},
		{
			name:    "invalid config json",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"not-json"}`),
			wantErr: true,
		},
		{
			name: "invalid config detail type",
			data: []byte(
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"notifyAppDeviceId\":123,\"type\":\"notifyapp\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			notifyApp := notification.NotifyApp{}

			err := json.Unmarshal(tc.data, &notifyApp)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, notifyApp)

			data, err := json.Marshal(notifyApp)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
