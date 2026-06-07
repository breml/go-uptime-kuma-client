package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationVKTeams_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.VKTeams
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My VKTeams Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My VKTeams Alert\",\"vkteamsBotToken\":\"001.1234567890.abcdefgh:1000000001\",\"vkteamsChatId\":\"*****@chat.agent\",\"vkteamsBaseUrl\":\"https://myteam.mail.ru\",\"vkteamsUseTemplate\":true,\"vkteamsTemplate\":\"Monitor: {monitor.name}\\nStatus: {monitor.status}\",\"vkteamsTemplateFormat\":\"MarkdownV2\",\"type\":\"VKTeams\"}"}`,
			),

			want: notification.VKTeams{
				Base: notification.Base{
					ID:            1,
					Name:          "My VKTeams Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				VKTeamsDetails: notification.VKTeamsDetails{
					BotToken:       "001.1234567890.abcdefgh:1000000001",
					ChatID:         "*****@chat.agent",
					BaseURL:        "https://myteam.mail.ru",
					UseTemplate:    true,
					Template:       "Monitor: {monitor.name}\nStatus: {monitor.status}",
					TemplateFormat: "MarkdownV2",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My VKTeams Alert","vkteamsBotToken":"001.1234567890.abcdefgh:1000000001","vkteamsChatId":"*****@chat.agent","vkteamsBaseUrl":"https://myteam.mail.ru","vkteamsUseTemplate":true,"vkteamsTemplate":"Monitor: {monitor.name}\nStatus: {monitor.status}","vkteamsTemplateFormat":"MarkdownV2","type":"VKTeams","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple VKTeams","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple VKTeams\",\"vkteamsBotToken\":\"001.1234567890.abcdefgh:1000000001\",\"vkteamsChatId\":\"*****@chat.agent\",\"type\":\"VKTeams\"}"}`,
			),

			want: notification.VKTeams{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple VKTeams",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				VKTeamsDetails: notification.VKTeamsDetails{
					BotToken: "001.1234567890.abcdefgh:1000000001",
					ChatID:   "*****@chat.agent",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple VKTeams","vkteamsBotToken":"001.1234567890.abcdefgh:1000000001","vkteamsChatId":"*****@chat.agent","vkteamsBaseUrl":"","vkteamsUseTemplate":false,"vkteamsTemplate":"","vkteamsTemplateFormat":"","type":"VKTeams","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vkteams := notification.VKTeams{}

			err := json.Unmarshal(tc.data, &vkteams)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, vkteams)

			data, err := json.Marshal(vkteams)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
