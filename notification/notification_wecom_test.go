package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWeCom_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.WeCom
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My WeCom Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My WeCom Alert\",\"weComBotKey\":\"abc123def456\",\"type\":\"WeCom\"}"}`),

			want: notification.WeCom{
				Base: notification.Base{
					ID:            1,
					Name:          "My WeCom Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WeComDetails: notification.WeComDetails{
					BotKey: "abc123def456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My WeCom Alert","type":"WeCom","userId":1,"weComBotKey":"abc123def456"}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple WeCom","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple WeCom\",\"weComBotKey\":\"xyz789abc\",\"type\":\"WeCom\"}"}`),

			want: notification.WeCom{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple WeCom",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WeComDetails: notification.WeComDetails{
					BotKey: "xyz789abc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple WeCom","type":"WeCom","userId":1,"weComBotKey":"xyz789abc"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wecom := notification.WeCom{}

			err := json.Unmarshal(tc.data, &wecom)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, wecom)

			data, err := json.Marshal(wecom)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
