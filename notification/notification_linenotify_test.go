package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationLineNotify_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.LineNotify
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My LINE Notify Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My LINE Notify Alert\",\"lineNotifyAccessToken\":\"access-token-123\",\"type\":\"LineNotify\"}"}`,
			),

			want: notification.LineNotify{
				Base: notification.Base{
					ID:            1,
					Name:          "My LINE Notify Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				LineNotifyDetails: notification.LineNotifyDetails{
					AccessToken: "access-token-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"lineNotifyAccessToken":"access-token-123","name":"My LINE Notify Alert","type":"LineNotify","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple LINE Notify","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple LINE Notify\",\"lineNotifyAccessToken\":\"token-abc\",\"type\":\"LineNotify\"}"}`,
			),

			want: notification.LineNotify{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple LINE Notify",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LineNotifyDetails: notification.LineNotifyDetails{
					AccessToken: "token-abc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"lineNotifyAccessToken":"token-abc","name":"Simple LINE Notify","type":"LineNotify","userId":1}`,
		},
		{
			name: "with different token",
			data: []byte(
				`{"id":3,"name":"LINE Notify Production","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"LINE Notify Production\",\"lineNotifyAccessToken\":\"prod-token-xyz\",\"type\":\"LineNotify\"}"}`,
			),

			want: notification.LineNotify{
				Base: notification.Base{
					ID:            3,
					Name:          "LINE Notify Production",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LineNotifyDetails: notification.LineNotifyDetails{
					AccessToken: "prod-token-xyz",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"lineNotifyAccessToken":"prod-token-xyz","name":"LINE Notify Production","type":"LineNotify","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			linenotify := notification.LineNotify{}

			err := json.Unmarshal(tc.data, &linenotify)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, linenotify)

			data, err := json.Marshal(linenotify)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
