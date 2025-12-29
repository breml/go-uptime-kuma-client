package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGotify_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Gotify
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Gotify Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Gotify Alert\",\"gotifyserverurl\":\"https://gotify.example.com\",\"gotifyapplicationToken\":\"app-token-123\",\"gotifyPriority\":8,\"type\":\"gotify\"}"}`,
			),

			want: notification.Gotify{
				Base: notification.Base{
					ID:            1,
					Name:          "My Gotify Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GotifyDetails: notification.GotifyDetails{
					ServerURL:        "https://gotify.example.com",
					ApplicationToken: "app-token-123",
					Priority:         8,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Gotify Alert","gotifyserverurl":"https://gotify.example.com","gotifyapplicationToken":"app-token-123","gotifyPriority":8,"type":"gotify","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Gotify","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Gotify\",\"gotifyserverurl\":\"https://gotify.org\",\"gotifyapplicationToken\":\"token\",\"type\":\"gotify\"}"}`,
			),

			want: notification.Gotify{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Gotify",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GotifyDetails: notification.GotifyDetails{
					ServerURL:        "https://gotify.org",
					ApplicationToken: "token",
					Priority:         0,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Gotify","gotifyserverurl":"https://gotify.org","gotifyapplicationToken":"token","gotifyPriority":0,"type":"gotify","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotify := notification.Gotify{}

			err := json.Unmarshal(tc.data, &gotify)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, gotify)

			data, err := json.Marshal(gotify)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
