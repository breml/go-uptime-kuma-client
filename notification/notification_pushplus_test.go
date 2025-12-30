package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPushPlus_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.PushPlus
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My PushPlus Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My PushPlus Alert\",\"pushPlusSendKey\":\"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"type\":\"PushPlus\"}"}`,
			),

			want: notification.PushPlus{
				Base: notification.Base{
					ID:            1,
					Name:          "My PushPlus Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PushPlusDetails: notification.PushPlusDetails{
					SendKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My PushPlus Alert","pushPlusSendKey":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","type":"PushPlus","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple PushPlus","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple PushPlus\",\"pushPlusSendKey\":\"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy\",\"type\":\"PushPlus\"}"}`,
			),

			want: notification.PushPlus{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple PushPlus",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushPlusDetails: notification.PushPlusDetails{
					SendKey: "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple PushPlus","pushPlusSendKey":"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy","type":"PushPlus","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pushplus := notification.PushPlus{}

			err := json.Unmarshal(tc.data, &pushplus)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pushplus)

			data, err := json.Marshal(pushplus)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
