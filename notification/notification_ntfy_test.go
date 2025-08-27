package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNtfy_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want notification.Ntfy
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Ntfy Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Ntfy Alert\",\"ntfyAuthenticationMethod\":\"usernamePassword\",\"ntfyIcon\":\"http://symbol.url\",\"ntfyPriority\":5,\"ntfypassword\":\"password\",\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"user\",\"type\":\"ntfy\"}"}`),

			want: notification.Ntfy{
				Base: notification.Base{
					ID:            1,
					Name:          "My Ntfy Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NtfyDetails: notification.NtfyDetails{
					AuthenticationMethod: "usernamePassword",
					Icon:                 "http://symbol.url",
					Priority:             5,
					Password:             "password",
					ServerURL:            "https://ntfy.sh",
					Topic:                "topic",
					Username:             "user",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ntfy := notification.Ntfy{}

			err := json.Unmarshal(tc.data, &ntfy)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, ntfy)

			data, err := json.Marshal(ntfy)
			require.NoError(t, err)

			t.Log(string(data))
		})
	}
}
