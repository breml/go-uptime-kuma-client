package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationMax_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Max
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Max Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Max Alert\",\"maxApiUrl\":\"https://platform-api.max.ru\",\"maxBotToken\":\"bot-token-123\",\"maxChatID\":\"-12345\",\"maxUseTemplate\":true,\"maxTemplate\":\"Monitor: {monitor.name}\\nStatus: {monitor.status}\",\"maxTemplateFormat\":\"markdown\",\"type\":\"max\"}"}`,
			),

			want: notification.Max{
				Base: notification.Base{
					ID:            1,
					Name:          "My Max Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				MaxDetails: notification.MaxDetails{
					APIURL:         "https://platform-api.max.ru",
					BotToken:       "bot-token-123",
					ChatID:         "-12345",
					UseTemplate:    true,
					Template:       "Monitor: {monitor.name}\nStatus: {monitor.status}",
					TemplateFormat: "markdown",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Max Alert","maxApiUrl":"https://platform-api.max.ru","maxBotToken":"bot-token-123","maxChatID":"-12345","maxUseTemplate":true,"maxTemplate":"Monitor: {monitor.name}\nStatus: {monitor.status}","maxTemplateFormat":"markdown","type":"max","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Max","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Max\",\"maxApiUrl\":\"https://platform-api.max.ru\",\"maxBotToken\":\"bot-token-123\",\"maxChatID\":\"-12345\",\"type\":\"max\"}"}`,
			),

			want: notification.Max{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Max",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				MaxDetails: notification.MaxDetails{
					APIURL:   "https://platform-api.max.ru",
					BotToken: "bot-token-123",
					ChatID:   "-12345",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Max","maxApiUrl":"https://platform-api.max.ru","maxBotToken":"bot-token-123","maxChatID":"-12345","maxUseTemplate":false,"maxTemplate":"","maxTemplateFormat":"","type":"max","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			maxn := notification.Max{}

			err := json.Unmarshal(tc.data, &maxn)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, maxn)

			data, err := json.Marshal(maxn)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
