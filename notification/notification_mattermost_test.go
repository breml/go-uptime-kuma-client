package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationMattermost_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Mattermost
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Mattermost Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Mattermost Alert\",\"mattermostWebhookUrl\":\"https://mattermost.example.com/hooks/xxx\",\"mattermostusername\":\"Monitor Bot\",\"mattermostchannel\":\"#alerts\",\"mattermosticonemo\":\":smiley: :frowning:\",\"mattermosticonurl\":\"https://example.com/icon.png\",\"type\":\"mattermost\"}"}`,
			),

			want: notification.Mattermost{
				Base: notification.Base{
					ID:            1,
					Name:          "My Mattermost Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				MattermostDetails: notification.MattermostDetails{
					WebhookURL: "https://mattermost.example.com/hooks/xxx",
					Username:   "Monitor Bot",
					Channel:    "#alerts",
					IconEmoji:  ":smiley: :frowning:",
					IconURL:    "https://example.com/icon.png",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Mattermost Alert","mattermostWebhookUrl":"https://mattermost.example.com/hooks/xxx","mattermostusername":"Monitor Bot","mattermostchannel":"#alerts","mattermosticonemo":":smiley: :frowning:","mattermosticonurl":"https://example.com/icon.png","type":"mattermost","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Mattermost","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Mattermost\",\"mattermostWebhookUrl\":\"https://mattermost.org/hooks/xxx\",\"type\":\"mattermost\"}"}`,
			),

			want: notification.Mattermost{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Mattermost",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				MattermostDetails: notification.MattermostDetails{
					WebhookURL: "https://mattermost.org/hooks/xxx",
					Username:   "",
					Channel:    "",
					IconEmoji:  "",
					IconURL:    "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Mattermost","mattermostWebhookUrl":"https://mattermost.org/hooks/xxx","mattermostusername":"","mattermostchannel":"","mattermosticonemo":"","mattermosticonurl":"","type":"mattermost","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mattermost := notification.Mattermost{}

			err := json.Unmarshal(tc.data, &mattermost)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, mattermost)

			data, err := json.Marshal(mattermost)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
