package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFortySixElks_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.FortySixElks
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My 46elks Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My 46elks Alert\",\"elksUsername\":\"username@example.com\",\"elksAuthToken\":\"auth_token_secret\",\"elksFromNumber\":\"1234\",\"elksToNumber\":\"0701234567\",\"type\":\"46elks\"}"}`,
			),

			want: notification.FortySixElks{
				Base: notification.Base{
					ID:            1,
					Name:          "My 46elks Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FortySixElksDetails: notification.FortySixElksDetails{
					Username:   "username@example.com",
					AuthToken:  "auth_token_secret",
					FromNumber: "1234",
					ToNumber:   "0701234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"elksAuthToken":"auth_token_secret","elksFromNumber":"1234","elksToNumber":"0701234567","elksUsername":"username@example.com","id":1,"isDefault":true,"name":"My 46elks Alert","type":"46elks","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple 46elks","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple 46elks\",\"elksUsername\":\"user\",\"elksAuthToken\":\"token\",\"elksFromNumber\":\"1234\",\"elksToNumber\":\"0701234567\",\"type\":\"46elks\"}"}`,
			),

			want: notification.FortySixElks{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple 46elks",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FortySixElksDetails: notification.FortySixElksDetails{
					Username:   "user",
					AuthToken:  "token",
					FromNumber: "1234",
					ToNumber:   "0701234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"elksAuthToken":"token","elksFromNumber":"1234","elksToNumber":"0701234567","elksUsername":"user","id":2,"isDefault":false,"name":"Simple 46elks","type":"46elks","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			elks := notification.FortySixElks{}

			err := json.Unmarshal(tc.data, &elks)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, elks)

			data, err := json.Marshal(elks)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
