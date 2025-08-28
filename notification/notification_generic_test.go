package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGeneric_String(t *testing.T) {
	tests := []struct {
		name string

		wantString string
		wantJSON   string
	}{
		{
			name: "success",

			wantString: `id: 1, name: "My Ntfy Alert", active: true, isDefault: true, applyExisting: true, userId: 1, type: "ntfy", ntfyAuthenticationMethod: "usernamePassword", ntfyIcon: "http://symbol.url", ntfyPriority: 5, ntfypassword: "password", ntfyserverurl: "https://ntfy.sh", ntfytopic: "topic", ntfyusername: "user"`,
			wantJSON:   `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Ntfy Alert","ntfyAuthenticationMethod":"usernamePassword","ntfyIcon":"http://symbol.url","ntfyPriority":5,"ntfypassword":"password","ntfyserverurl":"https://ntfy.sh","ntfytopic":"topic","ntfyusername":"user","type":"ntfy","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			noti := notification.Generic{
				Base: notification.Base{
					ID:            1,
					Name:          "My Ntfy Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TypeName: "ntfy",
				GenericDetails: notification.GenericDetails{
					"ntfyAuthenticationMethod": "usernamePassword",
					"ntfyIcon":                 "http://symbol.url",
					"ntfyPriority":             5,
					"ntfypassword":             "password",
					"ntfyserverurl":            "https://ntfy.sh",
					"ntfytopic":                "topic",
					"ntfyusername":             "user",
				},
			}

			require.Equal(t, tc.wantString, noti.String())

			data, err := json.Marshal(noti)
			require.NoError(t, err)

			require.Equal(t, tc.wantJSON, string(data))
		})
	}
}
