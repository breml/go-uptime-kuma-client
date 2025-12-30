package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationCallMeBot_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.CallMeBot
		wantJSON string
	}{
		{
			name: "success with endpoint",
			data: []byte(
				`{"id":1,"name":"My CallMeBot Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My CallMeBot Alert\",\"callMeBotEndpoint\":\"https://api.callmebot.com/start\",\"type\":\"CallMeBot\"}"}`,
			),

			want: notification.CallMeBot{
				Base: notification.Base{
					ID:            1,
					Name:          "My CallMeBot Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				CallMeBotDetails: notification.CallMeBotDetails{
					Endpoint: "https://api.callmebot.com/start",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"callMeBotEndpoint":"https://api.callmebot.com/start","id":1,"isDefault":true,"name":"My CallMeBot Alert","type":"CallMeBot","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple CallMeBot","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple CallMeBot\",\"callMeBotEndpoint\":\"https://api.callmebot.com/start\",\"type\":\"CallMeBot\"}"}`,
			),

			want: notification.CallMeBot{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple CallMeBot",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				CallMeBotDetails: notification.CallMeBotDetails{
					Endpoint: "https://api.callmebot.com/start",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"callMeBotEndpoint":"https://api.callmebot.com/start","id":2,"isDefault":false,"name":"Simple CallMeBot","type":"CallMeBot","userId":1}`,
		},
		{
			name: "with different endpoint",
			data: []byte(
				`{"id":3,"name":"Custom Endpoint","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Custom Endpoint\",\"callMeBotEndpoint\":\"https://custom.endpoint.com/webhook\",\"type\":\"CallMeBot\"}"}`,
			),

			want: notification.CallMeBot{
				Base: notification.Base{
					ID:            3,
					Name:          "Custom Endpoint",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				CallMeBotDetails: notification.CallMeBotDetails{
					Endpoint: "https://custom.endpoint.com/webhook",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"callMeBotEndpoint":"https://custom.endpoint.com/webhook","id":3,"isDefault":false,"name":"Custom Endpoint","type":"CallMeBot","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			callmebot := notification.CallMeBot{}

			err := json.Unmarshal(tc.data, &callmebot)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, callmebot)

			data, err := json.Marshal(callmebot)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
