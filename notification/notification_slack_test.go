package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
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
			data: []byte(
				`{"id":1,"name":"My Slack Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Slack Alert\",\"slackwebhookURL\":\"https://hooks.slack.com/services/xxx/yyy/zzz\",\"slackusername\":\"uptime-kuma\",\"slackiconemo\":\":ghost:\",\"slackchannel\":\"#alerts\",\"slackrichmessage\":true,\"slackchannelnotify\":false,\"type\":\"slack\"}"}`,
			),

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
		{
			name: "with include group name",
			data: []byte(
				`{"id":2,"name":"Slack Group Path","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Slack Group Path\",\"slackwebhookURL\":\"https://hooks.slack.com/services/aaa/bbb/ccc\",\"slackrichmessage\":true,\"slackchannelnotify\":false,\"slackIncludeGroupName\":true,\"type\":\"slack\"}"}`,
			),

			want: notification.Slack{
				Base: notification.Base{
					ID:            2,
					Name:          "Slack Group Path",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SlackDetails: notification.SlackDetails{
					WebhookURL:       "https://hooks.slack.com/services/aaa/bbb/ccc",
					RichMessage:      true,
					ChannelNotify:    false,
					IncludeGroupName: ptr.To(true),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Slack Group Path","slackIncludeGroupName":true,"slackchannel":"","slackchannelnotify":false,"slackiconemo":"","slackrichmessage":true,"slackusername":"","slackwebhookURL":"https://hooks.slack.com/services/aaa/bbb/ccc","type":"slack","userId":1}`,
		},
		{
			name: "with template",
			data: []byte(
				`{"id":3,"name":"Slack Template","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Slack Template\",\"slackwebhookURL\":\"https://hooks.slack.com/services/ddd/eee/fff\",\"slackrichmessage\":false,\"slackchannelnotify\":false,\"slackUseTemplate\":true,\"slackTemplate\":\"Alert: {{ msg }}\",\"type\":\"slack\"}"}`,
			),

			want: notification.Slack{
				Base: notification.Base{
					ID:            3,
					Name:          "Slack Template",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SlackDetails: notification.SlackDetails{
					WebhookURL:    "https://hooks.slack.com/services/ddd/eee/fff",
					RichMessage:   false,
					ChannelNotify: false,
					UseTemplate:   ptr.To(true),
					Template:      ptr.To("Alert: {{ msg }}"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Slack Template","slackTemplate":"Alert: {{ msg }}","slackUseTemplate":true,"slackchannel":"","slackchannelnotify":false,"slackiconemo":"","slackrichmessage":false,"slackusername":"","slackwebhookURL":"https://hooks.slack.com/services/ddd/eee/fff","type":"slack","userId":1}`,
		},
		{
			name: "pointer to false and empty string are serialized",
			data: []byte(
				`{"id":4,"name":"Slack Explicit False","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Slack Explicit False\",\"slackwebhookURL\":\"https://hooks.slack.com/services/ggg/hhh/iii\",\"slackrichmessage\":false,\"slackchannelnotify\":false,\"slackIncludeGroupName\":false,\"slackUseTemplate\":false,\"slackTemplate\":\"\",\"type\":\"slack\"}"}`,
			),

			want: notification.Slack{
				Base: notification.Base{
					ID:            4,
					Name:          "Slack Explicit False",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SlackDetails: notification.SlackDetails{
					WebhookURL:       "https://hooks.slack.com/services/ggg/hhh/iii",
					RichMessage:      false,
					ChannelNotify:    false,
					IncludeGroupName: ptr.To(false),
					UseTemplate:      ptr.To(false),
					Template:         ptr.To(""),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Slack Explicit False","slackIncludeGroupName":false,"slackTemplate":"","slackUseTemplate":false,"slackchannel":"","slackchannelnotify":false,"slackiconemo":"","slackrichmessage":false,"slackusername":"","slackwebhookURL":"https://hooks.slack.com/services/ggg/hhh/iii","type":"slack","userId":1}`,
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
