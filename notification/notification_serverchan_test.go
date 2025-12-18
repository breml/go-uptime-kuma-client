package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationServerChan_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.ServerChan
		wantJSON string
	}{
		{
			name: "success with send key",
			data: []byte(`{"id":1,"name":"My ServerChan Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My ServerChan Alert\",\"serverChanSendKey\":\"SCT123456789abcdefghijklmnopqrst\",\"type\":\"ServerChan\"}"}`),

			want: notification.ServerChan{
				Base: notification.Base{
					ID:            1,
					Name:          "My ServerChan Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ServerChanDetails: notification.ServerChanDetails{
					SendKey: "SCT123456789abcdefghijklmnopqrst",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My ServerChan Alert","serverChanSendKey":"SCT123456789abcdefghijklmnopqrst","type":"ServerChan","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple ServerChan","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple ServerChan\",\"serverChanSendKey\":\"SCT000000000000000000000000000000\",\"type\":\"ServerChan\"}"}`),

			want: notification.ServerChan{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple ServerChan",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ServerChanDetails: notification.ServerChanDetails{
					SendKey: "SCT000000000000000000000000000000",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple ServerChan","serverChanSendKey":"SCT000000000000000000000000000000","type":"ServerChan","userId":1}`,
		},
		{
			name: "with serverchan3 format",
			data: []byte(`{"id":3,"name":"ServerChan v3","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"ServerChan v3\",\"serverChanSendKey\":\"sctp123456t789abcdefghijklmnopqrst\",\"type\":\"ServerChan\"}"}`),

			want: notification.ServerChan{
				Base: notification.Base{
					ID:            3,
					Name:          "ServerChan v3",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ServerChanDetails: notification.ServerChanDetails{
					SendKey: "sctp123456t789abcdefghijklmnopqrst",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"ServerChan v3","serverChanSendKey":"sctp123456t789abcdefghijklmnopqrst","type":"ServerChan","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serverchan := notification.ServerChan{}

			err := json.Unmarshal(tc.data, &serverchan)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, serverchan)

			data, err := json.Marshal(serverchan)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
