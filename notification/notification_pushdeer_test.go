package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPushDeer_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.PushDeer
		wantJSON string
	}{
		{
			name: "success with custom server",
			data: []byte(
				`{"id":1,"name":"My PushDeer Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My PushDeer Alert\",\"pushdeerKey\":\"PDxxxxxxxxxxxxxx\",\"pushdeerServer\":\"https://custom.pushdeer.com\",\"type\":\"PushDeer\"}"}`,
			),

			want: notification.PushDeer{
				Base: notification.Base{
					ID:            1,
					Name:          "My PushDeer Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PushDeerDetails: notification.PushDeerDetails{
					Key:    "PDxxxxxxxxxxxxxx",
					Server: "https://custom.pushdeer.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My PushDeer Alert","pushdeerKey":"PDxxxxxxxxxxxxxx","pushdeerServer":"https://custom.pushdeer.com","type":"PushDeer","userId":1}`,
		},
		{
			name: "minimal without server",
			data: []byte(
				`{"id":2,"name":"Simple PushDeer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple PushDeer\",\"pushdeerKey\":\"PDyyyyyyyyyyyyyy\",\"pushdeerServer\":\"\",\"type\":\"PushDeer\"}"}`,
			),

			want: notification.PushDeer{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple PushDeer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushDeerDetails: notification.PushDeerDetails{
					Key:    "PDyyyyyyyyyyyyyy",
					Server: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple PushDeer","pushdeerKey":"PDyyyyyyyyyyyyyy","pushdeerServer":"","type":"PushDeer","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pushdeer := notification.PushDeer{}

			err := json.Unmarshal(tc.data, &pushdeer)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pushdeer)

			data, err := json.Marshal(pushdeer)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
