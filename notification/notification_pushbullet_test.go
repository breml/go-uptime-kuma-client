package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPushbullet_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Pushbullet
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Pushbullet Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Pushbullet Alert\",\"pushbulletAccessToken\":\"o.example_token_12345\",\"type\":\"pushbullet\"}"}`,
			),

			want: notification.Pushbullet{
				Base: notification.Base{
					ID:            1,
					Name:          "My Pushbullet Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PushbulletDetails: notification.PushbulletDetails{
					AccessToken: "o.example_token_12345",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Pushbullet Alert","pushbulletAccessToken":"o.example_token_12345","type":"pushbullet","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Pushbullet","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Pushbullet\",\"pushbulletAccessToken\":\"o.token123\",\"type\":\"pushbullet\"}"}`,
			),

			want: notification.Pushbullet{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Pushbullet",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushbulletDetails: notification.PushbulletDetails{
					AccessToken: "o.token123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Pushbullet","pushbulletAccessToken":"o.token123","type":"pushbullet","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pushbullet := notification.Pushbullet{}

			err := json.Unmarshal(tc.data, &pushbullet)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pushbullet)

			data, err := json.Marshal(pushbullet)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
