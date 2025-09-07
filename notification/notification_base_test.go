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

		want notification.Base
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"Created","active":true,"userId":1,"isDefault":true,"config":"{\"id\":0,\"name\":\"Created\",\"active\":true,\"userId\":0,\"isDefault\":true,\"ntfyaccesstoken\":\"\",\"ntfyAuthenticationMethod\":\"none\",\"ntfyIcon\":\"\",\"ntfypassword\":\"\",\"ntfyPriority\":5,\"ntfyserverurl\":\"https://ntfy.sh\",\"ntfytopic\":\"topic\",\"ntfyusername\":\"\",\"type\":\"ntfy\",\"applyExisting\":true}"}`),

			want: notification.Base{
				ID:            1,
				Name:          "Created",
				IsActive:      true,
				UserID:        1,
				IsDefault:     true,
				ApplyExisting: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ntfy := notification.Base{}

			err := json.Unmarshal(tc.data, &ntfy)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, ntfy)
		})
	}
}
