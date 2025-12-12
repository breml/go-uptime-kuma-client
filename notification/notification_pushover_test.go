package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPushover_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Pushover
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Pushover Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Pushover Alert\",\"pushoveruserkey\":\"userkey123\",\"pushoverapptoken\":\"apptoken456\",\"pushoversounds\":\"echo\",\"pushoversounds_up\":\"cashregister\",\"pushoverpriority\":\"1\",\"pushovertitle\":\"Uptime Kuma Alert\",\"pushoverdevice\":\"iphone\",\"pushoverttl\":\"3600\",\"type\":\"pushover\"}"}`),

			want: notification.Pushover{
				Base: notification.Base{
					ID:            1,
					Name:          "My Pushover Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PushoverDetails: notification.PushoverDetails{
					UserKey:  "userkey123",
					AppToken: "apptoken456",
					Sounds:   "echo",
					SoundsUp: "cashregister",
					Priority: "1",
					Title:    "Uptime Kuma Alert",
					Device:   "iphone",
					TTL:      "3600",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Pushover Alert","pushoverapptoken":"apptoken456","pushoverpriority":"1","pushoverdevice":"iphone","pushoversounds":"echo","pushoversounds_up":"cashregister","pushovertitle":"Uptime Kuma Alert","pushoverttl":"3600","pushoveruserkey":"userkey123","type":"pushover","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Pushover","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Pushover\",\"pushoveruserkey\":\"user456\",\"pushoverapptoken\":\"app789\",\"type\":\"pushover\"}"}`),

			want: notification.Pushover{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Pushover",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushoverDetails: notification.PushoverDetails{
					UserKey:  "user456",
					AppToken: "app789",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Pushover","pushoverapptoken":"app789","pushoverpriority":"","pushoverdevice":"","pushoversounds":"","pushoversounds_up":"","pushovertitle":"","pushoverttl":"","pushoveruserkey":"user456","type":"pushover","userId":1}`,
		},
		{
			name: "with optional fields",
			data: []byte(`{"id":3,"name":"Pushover Full","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Pushover Full\",\"pushoveruserkey\":\"userkey789\",\"pushoverapptoken\":\"apptoken012\",\"pushoversounds\":\"siren\",\"pushoversounds_up\":\"bike\",\"pushoverpriority\":\"2\",\"pushovertitle\":\"Server Alert\",\"pushoverdevice\":\"android\",\"pushoverttl\":\"7200\",\"type\":\"pushover\"}"}`),

			want: notification.Pushover{
				Base: notification.Base{
					ID:            3,
					Name:          "Pushover Full",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PushoverDetails: notification.PushoverDetails{
					UserKey:  "userkey789",
					AppToken: "apptoken012",
					Sounds:   "siren",
					SoundsUp: "bike",
					Priority: "2",
					Title:    "Server Alert",
					Device:   "android",
					TTL:      "7200",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Pushover Full","pushoverapptoken":"apptoken012","pushoverpriority":"2","pushoverdevice":"android","pushoversounds":"siren","pushoversounds_up":"bike","pushovertitle":"Server Alert","pushoverttl":"7200","pushoveruserkey":"userkey789","type":"pushover","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pushover := notification.Pushover{}

			err := json.Unmarshal(tc.data, &pushover)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pushover)

			data, err := json.Marshal(pushover)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
