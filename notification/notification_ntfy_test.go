package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNtfy_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Ntfy
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Ntfy Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Ntfy Alert\",\"ntfyAuthenticationMethod\":\"usernamePassword\",\"ntfyIcon\":\"http://symbol.url\",\"ntfyPriority\":5,\"ntfypassword\":\"password\",\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"user\",\"type\":\"ntfy\"}"}`,
			),

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
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Ntfy Alert","ntfyAuthenticationMethod":"usernamePassword","ntfyIcon":"http://symbol.url","ntfyPriority":5,"ntfyaccesstoken":"","ntfypassword":"password","ntfyserverurl":"https://ntfy.sh","ntfytopic":"topic","ntfyusername":"user","type":"ntfy","userId":1}`,
		},
		{
			name: "with template, call and priority down",
			data: []byte(
				`{"id":2,"name":"My Ntfy Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Ntfy Alert\",\"ntfyAuthenticationMethod\":\"usernamePassword\",\"ntfyCall\":\"+12223334444\",\"ntfyCustomMessage\":\"custom message\",\"ntfyCustomTitle\":\"custom title\",\"ntfyIcon\":\"http://symbol.url\",\"ntfyPriority\":3,\"ntfyPriorityDown\":5,\"ntfyUseTemplate\":true,\"ntfypassword\":\"password\",\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"user\",\"type\":\"ntfy\"}"}`,
			),

			want: notification.Ntfy{
				Base: notification.Base{
					ID:            2,
					Name:          "My Ntfy Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NtfyDetails: notification.NtfyDetails{
					AuthenticationMethod: "usernamePassword",
					Call:                 ptr.To("+12223334444"),
					CustomMessage:        ptr.To("custom message"),
					CustomTitle:          ptr.To("custom title"),
					Icon:                 "http://symbol.url",
					Priority:             3,
					PriorityDown:         5,
					UseTemplate:          ptr.To(true),
					Password:             "password",
					ServerURL:            "https://ntfy.sh",
					Topic:                "topic",
					Username:             "user",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":2,"isDefault":true,"name":"My Ntfy Alert","ntfyAuthenticationMethod":"usernamePassword","ntfyCall":"+12223334444","ntfyCustomMessage":"custom message","ntfyCustomTitle":"custom title","ntfyIcon":"http://symbol.url","ntfyPriority":3,"ntfyPriorityDown":5,"ntfyUseTemplate":true,"ntfyaccesstoken":"","ntfypassword":"password","ntfyserverurl":"https://ntfy.sh","ntfytopic":"topic","ntfyusername":"user","type":"ntfy","userId":1}`,
		},
		{
			name: "explicit empty template fields and disabled toggle",
			data: []byte(
				`{"id":3,"name":"My Ntfy Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Ntfy Alert\",\"ntfyAuthenticationMethod\":\"none\",\"ntfyCall\":\"\",\"ntfyCustomMessage\":\"\",\"ntfyCustomTitle\":\"\",\"ntfyIcon\":\"\",\"ntfyPriority\":5,\"ntfyUseTemplate\":false,\"ntfypassword\":\"\",\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"\",\"type\":\"ntfy\"}"}`,
			),

			want: notification.Ntfy{
				Base: notification.Base{
					ID:            3,
					Name:          "My Ntfy Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NtfyDetails: notification.NtfyDetails{
					AuthenticationMethod: "none",
					Call:                 ptr.To(""),
					CustomMessage:        ptr.To(""),
					CustomTitle:          ptr.To(""),
					Priority:             5,
					ServerURL:            "https://ntfy.sh",
					Topic:                "topic",
					UseTemplate:          ptr.To(false),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":3,"isDefault":true,"name":"My Ntfy Alert","ntfyAuthenticationMethod":"none","ntfyCall":"","ntfyCustomMessage":"","ntfyCustomTitle":"","ntfyIcon":"","ntfyPriority":5,"ntfyUseTemplate":false,"ntfyaccesstoken":"","ntfypassword":"","ntfyserverurl":"https://ntfy.sh","ntfytopic":"topic","ntfyusername":"","type":"ntfy","userId":1}`,
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

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
