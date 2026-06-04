package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationHaloPSA_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.HaloPSA
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My HaloPSA Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My HaloPSA Alert\",\"halowebhookurl\":\"https://example.halopsa.com/api/v1/webhook\",\"haloUsername\":\"admin\",\"haloPassword\":\"secret\",\"type\":\"HaloPSA\"}"}`,
			),

			want: notification.HaloPSA{
				Base: notification.Base{
					ID:            1,
					Name:          "My HaloPSA Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				HaloPSADetails: notification.HaloPSADetails{
					WebhookURL: "https://example.halopsa.com/api/v1/webhook",
					Username:   "admin",
					Password:   "secret",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"haloPassword":"secret","haloUsername":"admin","halowebhookurl":"https://example.halopsa.com/api/v1/webhook","id":1,"isDefault":true,"name":"My HaloPSA Alert","type":"HaloPSA","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple HaloPSA","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple HaloPSA\",\"halowebhookurl\":\"https://halopsa.example.com/webhook\",\"haloUsername\":\"\",\"haloPassword\":\"\",\"type\":\"HaloPSA\"}"}`,
			),

			want: notification.HaloPSA{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple HaloPSA",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				HaloPSADetails: notification.HaloPSADetails{
					WebhookURL: "https://halopsa.example.com/webhook",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"haloPassword":"","haloUsername":"","halowebhookurl":"https://halopsa.example.com/webhook","id":2,"isDefault":false,"name":"Simple HaloPSA","type":"HaloPSA","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			halopsa := notification.HaloPSA{}

			err := json.Unmarshal(tc.data, &halopsa)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, halopsa)

			data, err := json.Marshal(halopsa)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
