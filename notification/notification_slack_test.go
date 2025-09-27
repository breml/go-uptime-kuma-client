package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSlack_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Slack
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Slack Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Slack Alert\",\"slackwebhookURL\":\"https://hooks.slack.com/services/xxx/yyy/zzz\",\"slackusername\":\"uptime-kuma\",\"slackiconemo\":\":ghost:\",\"slackchannel\":\"#alerts\",\"slackrichmessage\":true,\"slackchannelnotify\":false,\"type\":\"slack\"}"}`),

			want: notification.Slack{
				Base: notification.Base{
					ID:            1,
					Name:          "My Slack Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SlackDetails: notification.SlackDetails{
					WebhookURL:    "https://hooks.slack.com/services/xxx/yyy/zzz",
					Username:      "uptime-kuma",
					IconEmoji:     ":ghost:",
					Channel:       "#alerts",
					RichMessage:   true,
					ChannelNotify: false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Slack Alert","slackchannel":"#alerts","slackchannelnotify":false,"slackiconemo":":ghost:","slackrichmessage":true,"slackusername":"uptime-kuma","slackwebhookURL":"https://hooks.slack.com/services/xxx/yyy/zzz","type":"slack","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slack := notification.Slack{}

			err := json.Unmarshal(tc.data, &slack)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, slack)

			data, err := json.Marshal(slack)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
