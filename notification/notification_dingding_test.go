package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationDingDing_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.DingDing
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My DingDing Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My DingDing Alert\",\"webHookUrl\":\"https://oapi.dingtalk.com/robot/send?access_token=xxx\",\"secretKey\":\"secret123\",\"mentioning\":\"everyone\",\"type\":\"DingDing\"}"}`),

			want: notification.DingDing{
				Base: notification.Base{
					ID:            1,
					Name:          "My DingDing Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				DingDingDetails: notification.DingDingDetails{
					WebHookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
					SecretKey:  "secret123",
					Mentioning: "everyone",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"mentioning":"everyone","name":"My DingDing Alert","secretKey":"secret123","type":"DingDing","userId":1,"webHookUrl":"https://oapi.dingtalk.com/robot/send?access_token=xxx"}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple DingDing","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple DingDing\",\"webHookUrl\":\"https://oapi.dingtalk.com/robot/send?access_token=yyy\",\"secretKey\":\"\",\"mentioning\":\"\",\"type\":\"DingDing\"}"}`),

			want: notification.DingDing{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple DingDing",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				DingDingDetails: notification.DingDingDetails{
					WebHookURL: "https://oapi.dingtalk.com/robot/send?access_token=yyy",
					SecretKey:  "",
					Mentioning: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"mentioning":"","name":"Simple DingDing","secretKey":"","type":"DingDing","userId":1,"webHookUrl":"https://oapi.dingtalk.com/robot/send?access_token=yyy"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dingding := notification.DingDing{}

			err := json.Unmarshal(tc.data, &dingding)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, dingding)

			data, err := json.Marshal(dingding)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
