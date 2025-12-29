package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestBase_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Base
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"Created","active":true,"userId":1,"isDefault":true,"config":"{\"name\":\"Created\",\"active\":true,\"isDefault\":true,\"ntfyaccesstoken\":\"\",\"ntfyAuthenticationMethod\":\"none\",\"ntfyIcon\":\"\",\"ntfypassword\":\"\",\"ntfyPriority\":5,\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"\",\"type\":\"ntfy\",\"applyExisting\":true}"}`,
			),

			want: notification.Base{
				ID:            1,
				Name:          "Created",
				IsActive:      true,
				UserID:        1,
				IsDefault:     true,
				ApplyExisting: true,
			},
			wantJSON: `{"active":true, "applyExisting":true, "id":1, "isDefault":true, "name":"Created", "ntfyAuthenticationMethod":"none", "ntfyIcon":"", "ntfyPriority":5, "ntfyaccesstoken":"", "ntfypassword":"", "ntfyserverurl":"https://ntfy.sh", "ntfytopic":"topic", "ntfyusername":"", "type":"ntfy", "userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			someNotification := notification.Base{}

			err := json.Unmarshal(tc.data, &someNotification)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, someNotification)

			data, err := json.Marshal(someNotification)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
