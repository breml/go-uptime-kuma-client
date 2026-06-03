package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationVK_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.VK
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My VK Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My VK Alert\",\"vkAccessToken\":\"vk1.a.abcdefghijklmnopqrstuvwxyz\",\"vkPeerId\":\"2000000001\",\"vkApiVersion\":\"5.199\",\"vkDontParseLinks\":true,\"type\":\"VK\"}"}`,
			),

			want: notification.VK{
				Base: notification.Base{
					ID:            1,
					Name:          "My VK Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				VKDetails: notification.VKDetails{
					AccessToken:    "vk1.a.abcdefghijklmnopqrstuvwxyz",
					PeerID:         "2000000001",
					APIVersion:     "5.199",
					DontParseLinks: true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My VK Alert","vkAccessToken":"vk1.a.abcdefghijklmnopqrstuvwxyz","vkPeerId":"2000000001","vkApiVersion":"5.199","vkDontParseLinks":true,"type":"VK","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple VK","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple VK\",\"vkAccessToken\":\"vk1.a.token\",\"vkPeerId\":\"12345\",\"vkApiVersion\":\"5.199\",\"type\":\"VK\"}"}`,
			),

			want: notification.VK{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple VK",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				VKDetails: notification.VKDetails{
					AccessToken: "vk1.a.token",
					PeerID:      "12345",
					APIVersion:  "5.199",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple VK","vkAccessToken":"vk1.a.token","vkPeerId":"12345","vkApiVersion":"5.199","vkDontParseLinks":false,"type":"VK","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vk := notification.VK{}

			err := json.Unmarshal(tc.data, &vk)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, vk)

			data, err := json.Marshal(vk)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
