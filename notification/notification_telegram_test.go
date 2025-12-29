package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTelegram_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Telegram
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Telegram Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Telegram Alert\",\"telegramBotToken\":\"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11\",\"telegramChatID\":\"@mychannel\",\"telegramServerUrl\":\"https://api.telegram.org\",\"telegramSendSilently\":true,\"telegramProtectContent\":false,\"telegramMessageThreadID\":\"42\",\"telegramUseTemplate\":true,\"telegramTemplate\":\"Monitor: {monitor.name}\\nStatus: {monitor.status}\",\"telegramTemplateParseMode\":\"Markdown\",\"type\":\"telegram\"}"}`,
			),

			want: notification.Telegram{
				Base: notification.Base{
					ID:            1,
					Name:          "My Telegram Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TelegramDetails: notification.TelegramDetails{
					BotToken:          "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
					ChatID:            "@mychannel",
					ServerURL:         "https://api.telegram.org",
					SendSilently:      true,
					ProtectContent:    false,
					MessageThreadID:   "42",
					UseTemplate:       true,
					Template:          "Monitor: {monitor.name}\nStatus: {monitor.status}",
					TemplateParseMode: "Markdown",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Telegram Alert","telegramBotToken":"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11","telegramChatID":"@mychannel","telegramServerUrl":"https://api.telegram.org","telegramSendSilently":true,"telegramProtectContent":false,"telegramMessageThreadID":"42","telegramUseTemplate":true,"telegramTemplate":"Monitor: {monitor.name}\nStatus: {monitor.status}","telegramTemplateParseMode":"Markdown","type":"telegram","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Telegram","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Telegram\",\"telegramBotToken\":\"123456:ABC\",\"telegramChatID\":\"123456789\",\"type\":\"telegram\"}"}`,
			),

			want: notification.Telegram{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Telegram",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TelegramDetails: notification.TelegramDetails{
					BotToken: "123456:ABC",
					ChatID:   "123456789",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Telegram","telegramBotToken":"123456:ABC","telegramChatID":"123456789","telegramServerUrl":"","telegramSendSilently":false,"telegramProtectContent":false,"telegramMessageThreadID":"","telegramUseTemplate":false,"telegramTemplate":"","telegramTemplateParseMode":"","type":"telegram","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			telegram := notification.Telegram{}

			err := json.Unmarshal(tc.data, &telegram)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, telegram)

			data, err := json.Marshal(telegram)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
