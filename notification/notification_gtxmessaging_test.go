package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGTXMessaging_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.GTXMessaging
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My GTX Messaging Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My GTX Messaging Alert\",\"gtxMessagingApiKey\":\"api-key-123\",\"gtxMessagingFrom\":\"Uptime\",\"gtxMessagingTo\":\"+46701234567\",\"type\":\"gtxmessaging\"}"}`,
			),

			want: notification.GTXMessaging{
				Base: notification.Base{
					ID:            1,
					Name:          "My GTX Messaging Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GTXMessagingDetails: notification.GTXMessagingDetails{
					ApiKey: "api-key-123",
					From:   "Uptime",
					To:     "+46701234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"gtxMessagingApiKey":"api-key-123","gtxMessagingFrom":"Uptime","gtxMessagingTo":"+46701234567","id":1,"isDefault":true,"name":"My GTX Messaging Alert","type":"gtxmessaging","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple GTX","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple GTX\",\"gtxMessagingApiKey\":\"key-abc\",\"gtxMessagingFrom\":\"Alert\",\"gtxMessagingTo\":\"+46700000000\",\"type\":\"gtxmessaging\"}"}`,
			),

			want: notification.GTXMessaging{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple GTX",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GTXMessagingDetails: notification.GTXMessagingDetails{
					ApiKey: "key-abc",
					From:   "Alert",
					To:     "+46700000000",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"gtxMessagingApiKey":"key-abc","gtxMessagingFrom":"Alert","gtxMessagingTo":"+46700000000","id":2,"isDefault":false,"name":"Simple GTX","type":"gtxmessaging","userId":1}`,
		},
		{
			name: "with international number",
			data: []byte(
				`{"id":3,"name":"GTX International","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"GTX International\",\"gtxMessagingApiKey\":\"intl-key-xyz\",\"gtxMessagingFrom\":\"Monitor\",\"gtxMessagingTo\":\"+14155551234\",\"type\":\"gtxmessaging\"}"}`,
			),

			want: notification.GTXMessaging{
				Base: notification.Base{
					ID:            3,
					Name:          "GTX International",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GTXMessagingDetails: notification.GTXMessagingDetails{
					ApiKey: "intl-key-xyz",
					From:   "Monitor",
					To:     "+14155551234",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"gtxMessagingApiKey":"intl-key-xyz","gtxMessagingFrom":"Monitor","gtxMessagingTo":"+14155551234","id":3,"isDefault":false,"name":"GTX International","type":"gtxmessaging","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gtxmessaging := notification.GTXMessaging{}

			err := json.Unmarshal(tc.data, &gtxmessaging)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, gtxmessaging)

			data, err := json.Marshal(gtxmessaging)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
