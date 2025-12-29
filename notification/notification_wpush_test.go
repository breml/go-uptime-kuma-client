package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWPush_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.WPush
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My WPush Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My WPush Alert\",\"wpushAPIkey\":\"test-api-key-123\",\"wpushChannel\":\"channel-alerts\",\"type\":\"WPush\"}"}`,
			),

			want: notification.WPush{
				Base: notification.Base{
					ID:            1,
					Name:          "My WPush Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WPushDetails: notification.WPushDetails{
					APIKey:  "test-api-key-123",
					Channel: "channel-alerts",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My WPush Alert","type":"WPush","userId":1,"wpushAPIkey":"test-api-key-123","wpushChannel":"channel-alerts"}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple WPush","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple WPush\",\"wpushAPIkey\":\"minimal-key\",\"wpushChannel\":\"default\",\"type\":\"WPush\"}"}`,
			),

			want: notification.WPush{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple WPush",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WPushDetails: notification.WPushDetails{
					APIKey:  "minimal-key",
					Channel: "default",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple WPush","type":"WPush","userId":1,"wpushAPIkey":"minimal-key","wpushChannel":"default"}`,
		},
		{
			name: "with different channel",
			data: []byte(
				`{"id":3,"name":"WPush Multi Channel","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"WPush Multi Channel\",\"wpushAPIkey\":\"production-api-key\",\"wpushChannel\":\"monitoring-channel\",\"type\":\"WPush\"}"}`,
			),

			want: notification.WPush{
				Base: notification.Base{
					ID:            3,
					Name:          "WPush Multi Channel",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WPushDetails: notification.WPushDetails{
					APIKey:  "production-api-key",
					Channel: "monitoring-channel",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"WPush Multi Channel","type":"WPush","userId":1,"wpushAPIkey":"production-api-key","wpushChannel":"monitoring-channel"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wpush := notification.WPush{}

			err := json.Unmarshal(tc.data, &wpush)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, wpush)

			data, err := json.Marshal(wpush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
