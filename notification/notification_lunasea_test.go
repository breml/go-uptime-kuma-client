package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationLunaSea_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.LunaSea
		wantJSON string
	}{
		{
			name: "success with user target",
			data: []byte(`{"id":1,"name":"My LunaSea Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My LunaSea Alert\",\"lunaseaTarget\":\"user\",\"lunaseaUserID\":\"user123\",\"lunaseaDevice\":\"\",\"type\":\"lunasea\"}"}`),

			want: notification.LunaSea{
				Base: notification.Base{
					ID:            1,
					Name:          "My LunaSea Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				LunaSeaDetails: notification.LunaSeaDetails{
					Target: "user",
					UserID: "user123",
					Device: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"lunaseaDevice":"","lunaseaTarget":"user","lunaseaUserID":"user123","name":"My LunaSea Alert","type":"lunasea","userId":1}`,
		},
		{
			name: "with device target",
			data: []byte(`{"id":2,"name":"Simple LunaSea","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple LunaSea\",\"lunaseaTarget\":\"device\",\"lunaseaUserID\":\"\",\"lunaseaDevice\":\"device456\",\"type\":\"lunasea\"}"}`),

			want: notification.LunaSea{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple LunaSea",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LunaSeaDetails: notification.LunaSeaDetails{
					Target: "device",
					UserID: "",
					Device: "device456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"lunaseaDevice":"device456","lunaseaTarget":"device","lunaseaUserID":"","name":"Simple LunaSea","type":"lunasea","userId":1}`,
		},
		{
			name: "with different user",
			data: []byte(`{"id":3,"name":"LunaSea Production","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"LunaSea Production\",\"lunaseaTarget\":\"user\",\"lunaseaUserID\":\"prod-user-789\",\"lunaseaDevice\":\"\",\"type\":\"lunasea\"}"}`),

			want: notification.LunaSea{
				Base: notification.Base{
					ID:            3,
					Name:          "LunaSea Production",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LunaSeaDetails: notification.LunaSeaDetails{
					Target: "user",
					UserID: "prod-user-789",
					Device: "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"lunaseaDevice":"","lunaseaTarget":"user","lunaseaUserID":"prod-user-789","name":"LunaSea Production","type":"lunasea","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lunasea := notification.LunaSea{}

			err := json.Unmarshal(tc.data, &lunasea)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, lunasea)

			data, err := json.Marshal(lunasea)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
