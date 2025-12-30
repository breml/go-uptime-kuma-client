package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFeishu_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Feishu
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Feishu Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Feishu Alert\",\"feishuWebHookUrl\":\"https://open.feishu.cn/open-apis/bot/v2/hook/xxx\",\"type\":\"Feishu\"}"}`,
			),

			want: notification.Feishu{
				Base: notification.Base{
					ID:            1,
					Name:          "My Feishu Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FeishuDetails: notification.FeishuDetails{
					WebHookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"feishuWebHookUrl":"https://open.feishu.cn/open-apis/bot/v2/hook/xxx","id":1,"isDefault":true,"name":"My Feishu Alert","type":"Feishu","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Feishu","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Feishu\",\"feishuWebHookUrl\":\"https://open.feishu.cn/open-apis/bot/v2/hook/yyy\",\"type\":\"Feishu\"}"}`,
			),

			want: notification.Feishu{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Feishu",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FeishuDetails: notification.FeishuDetails{
					WebHookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/yyy",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"feishuWebHookUrl":"https://open.feishu.cn/open-apis/bot/v2/hook/yyy","id":2,"isDefault":false,"name":"Simple Feishu","type":"Feishu","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			feishu := notification.Feishu{}

			err := json.Unmarshal(tc.data, &feishu)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, feishu)

			data, err := json.Marshal(feishu)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
