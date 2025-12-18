package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPushy_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Pushy
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Pushy Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Pushy Alert\",\"pushyAPIKey\":\"api_key_xxxxx\",\"pushyToken\":\"device_token_xxxxx\",\"type\":\"pushy\"}"}`),

			want: notification.Pushy{
				Base: notification.Base{
					ID:            1,
					Name:          "My Pushy Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PushyDetails: notification.PushyDetails{
					APIKey: "api_key_xxxxx",
					Token:  "device_token_xxxxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Pushy Alert","pushyAPIKey":"api_key_xxxxx","pushyToken":"device_token_xxxxx","type":"pushy","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Pushy","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Pushy\",\"pushyAPIKey\":\"api_key_yyyyy\",\"pushyToken\":\"device_token_yyyyy\",\"type\":\"pushy\"}"}`),

			want: notification.Pushy{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Pushy",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushyDetails: notification.PushyDetails{
					APIKey: "api_key_yyyyy",
					Token:  "device_token_yyyyy",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Pushy","pushyAPIKey":"api_key_yyyyy","pushyToken":"device_token_yyyyy","type":"pushy","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pushy := notification.Pushy{}

			err := json.Unmarshal(tc.data, &pushy)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pushy)

			data, err := json.Marshal(pushy)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
