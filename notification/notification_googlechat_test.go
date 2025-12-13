package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGoogleChat_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.GoogleChat
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Google Chat Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Google Chat Alert\",\"googleChatWebhookURL\":\"https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test\",\"googleChatUseTemplate\":true,\"googleChatTemplate\":\"Test Template\",\"type\":\"GoogleChat\"}"}`),

			want: notification.GoogleChat{
				Base: notification.Base{
					ID:            1,
					Name:          "My Google Chat Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GoogleChatDetails: notification.GoogleChatDetails{
					WebhookURL:  "https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test",
					UseTemplate: true,
					Template:    "Test Template",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"googleChatTemplate":"Test Template","googleChatUseTemplate":true,"googleChatWebhookURL":"https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test","id":1,"isDefault":true,"name":"My Google Chat Alert","type":"GoogleChat","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Google Chat","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Google Chat\",\"googleChatWebhookURL\":\"https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test\",\"googleChatUseTemplate\":false,\"googleChatTemplate\":\"\",\"type\":\"GoogleChat\"}"}`),

			want: notification.GoogleChat{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Google Chat",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GoogleChatDetails: notification.GoogleChatDetails{
					WebhookURL:  "https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test",
					UseTemplate: false,
					Template:    "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"googleChatTemplate":"","googleChatUseTemplate":false,"googleChatWebhookURL":"https://chat.googleapis.com/v1/spaces/AAAAAA/messages?key=test","id":2,"isDefault":false,"name":"Simple Google Chat","type":"GoogleChat","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			googlechat := notification.GoogleChat{}

			err := json.Unmarshal(tc.data, &googlechat)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, googlechat)

			data, err := json.Marshal(googlechat)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
