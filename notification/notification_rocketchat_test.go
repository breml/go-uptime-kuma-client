package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationRocketChat_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.RocketChat
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Rocket.Chat Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Rocket.Chat Alert\",\"rocketwebhookURL\":\"https://rocket.example.com/hooks/xxx\",\"rocketchannel\":\"#alerts\",\"rocketusername\":\"Monitor Bot\",\"rocketiconemo\":\":smiley:\",\"rocketbutton\":\"\",\"type\":\"rocket.chat\"}"}`),

			want: notification.RocketChat{
				Base: notification.Base{
					ID:            1,
					Name:          "My Rocket.Chat Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				RocketChatDetails: notification.RocketChatDetails{
					WebhookURL: "https://rocket.example.com/hooks/xxx",
					Channel:    "#alerts",
					Username:   "Monitor Bot",
					IconEmoji:  ":smiley:",
					Button:     "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Rocket.Chat Alert","rocketbutton":"","rocketchannel":"#alerts","rocketiconemo":":smiley:","rocketusername":"Monitor Bot","rocketwebhookURL":"https://rocket.example.com/hooks/xxx","type":"rocket.chat","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Rocket.Chat","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Rocket.Chat\",\"rocketwebhookURL\":\"https://rocket.org/hooks/xxx\",\"type\":\"rocket.chat\"}"}`),

			want: notification.RocketChat{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Rocket.Chat",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				RocketChatDetails: notification.RocketChatDetails{
					WebhookURL: "https://rocket.org/hooks/xxx",
					Channel:    "",
					Username:   "",
					IconEmoji:  "",
					Button:     "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Rocket.Chat","rocketbutton":"","rocketchannel":"","rocketiconemo":"","rocketusername":"","rocketwebhookURL":"https://rocket.org/hooks/xxx","type":"rocket.chat","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rocketchat := notification.RocketChat{}

			err := json.Unmarshal(tc.data, &rocketchat)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, rocketchat)

			data, err := json.Marshal(rocketchat)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
