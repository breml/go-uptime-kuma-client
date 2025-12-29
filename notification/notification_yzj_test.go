package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationYZJ_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.YZJ
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My YZJ Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My YZJ Alert\",\"yzjWebHookUrl\":\"https://api.yzj.cn/webhook\",\"yzjToken\":\"test-token-123\",\"type\":\"YZJ\"}"}`,
			),

			want: notification.YZJ{
				Base: notification.Base{
					ID:            1,
					Name:          "My YZJ Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				YZJDetails: notification.YZJDetails{
					WebHookURL: "https://api.yzj.cn/webhook",
					Token:      "test-token-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My YZJ Alert","type":"YZJ","userId":1,"yzjToken":"test-token-123","yzjWebHookUrl":"https://api.yzj.cn/webhook"}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple YZJ","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple YZJ\",\"yzjWebHookUrl\":\"https://yunzhijia.com/webhook\",\"yzjToken\":\"simple-token\",\"type\":\"YZJ\"}"}`,
			),

			want: notification.YZJ{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple YZJ",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				YZJDetails: notification.YZJDetails{
					WebHookURL: "https://yunzhijia.com/webhook",
					Token:      "simple-token",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple YZJ","type":"YZJ","userId":1,"yzjToken":"simple-token","yzjWebHookUrl":"https://yunzhijia.com/webhook"}`,
		},
		{
			name: "with different webhook URL",
			data: []byte(
				`{"id":3,"name":"YZJ Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"YZJ Custom\",\"yzjWebHookUrl\":\"https://custom.webhook.yzj.io/notify\",\"yzjToken\":\"custom-token-456\",\"type\":\"YZJ\"}"}`,
			),

			want: notification.YZJ{
				Base: notification.Base{
					ID:            3,
					Name:          "YZJ Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				YZJDetails: notification.YZJDetails{
					WebHookURL: "https://custom.webhook.yzj.io/notify",
					Token:      "custom-token-456",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"YZJ Custom","type":"YZJ","userId":1,"yzjToken":"custom-token-456","yzjWebHookUrl":"https://custom.webhook.yzj.io/notify"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			yzj := notification.YZJ{}

			err := json.Unmarshal(tc.data, &yzj)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, yzj)

			data, err := json.Marshal(yzj)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
