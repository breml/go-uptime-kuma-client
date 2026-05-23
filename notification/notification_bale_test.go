package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationBale_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Bale
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Bale Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Bale Alert\",\"baleBotToken\":\"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11\",\"baleChatID\":\"@mychannel\",\"type\":\"bale\"}"}`,
			),

			want: notification.Bale{
				Base: notification.Base{
					ID:            1,
					Name:          "My Bale Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				BaleDetails: notification.BaleDetails{
					BotToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
					ChatID:   "@mychannel",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"baleBotToken":"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11","baleChatID":"@mychannel","id":1,"isDefault":true,"name":"My Bale Alert","type":"bale","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Bale","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Bale\",\"baleBotToken\":\"123456:ABC\",\"baleChatID\":\"123456789\",\"type\":\"bale\"}"}`,
			),

			want: notification.Bale{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Bale",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BaleDetails: notification.BaleDetails{
					BotToken: "123456:ABC",
					ChatID:   "123456789",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"baleBotToken":"123456:ABC","baleChatID":"123456789","id":2,"isDefault":false,"name":"Simple Bale","type":"bale","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bale := notification.Bale{}

			err := json.Unmarshal(tc.data, &bale)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, bale)

			data, err := json.Marshal(bale)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
