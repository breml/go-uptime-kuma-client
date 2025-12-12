package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOpsgenie_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Opsgenie
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Opsgenie Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Opsgenie Alert\",\"opsgenieApiKey\":\"test-api-key-123\",\"opsgenieRegion\":\"us\",\"opsgeniePriority\":3,\"type\":\"Opsgenie\"}"}`),

			want: notification.Opsgenie{
				Base: notification.Base{
					ID:            1,
					Name:          "My Opsgenie Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				OpsgenieDetails: notification.OpsgenieDetails{
					ApiKey:   "test-api-key-123",
					Region:   "us",
					Priority: 3,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Opsgenie Alert","opsgenieApiKey":"test-api-key-123","opsgenieRegion":"us","opsgeniePriority":3,"type":"Opsgenie","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Opsgenie","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Opsgenie\",\"opsgenieApiKey\":\"abc-123\",\"type\":\"Opsgenie\"}"}`),

			want: notification.Opsgenie{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Opsgenie",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OpsgenieDetails: notification.OpsgenieDetails{
					ApiKey: "abc-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Opsgenie","opsgenieApiKey":"abc-123","opsgenieRegion":"","opsgeniePriority":0,"type":"Opsgenie","userId":1}`,
		},
		{
			name: "with eu region",
			data: []byte(`{"id":3,"name":"EU Opsgenie","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"EU Opsgenie\",\"opsgenieApiKey\":\"eu-api-key\",\"opsgenieRegion\":\"eu\",\"opsgeniePriority\":5,\"type\":\"Opsgenie\"}"}`),

			want: notification.Opsgenie{
				Base: notification.Base{
					ID:            3,
					Name:          "EU Opsgenie",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OpsgenieDetails: notification.OpsgenieDetails{
					ApiKey:   "eu-api-key",
					Region:   "eu",
					Priority: 5,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"EU Opsgenie","opsgenieApiKey":"eu-api-key","opsgenieRegion":"eu","opsgeniePriority":5,"type":"Opsgenie","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opsgenie := notification.Opsgenie{}

			err := json.Unmarshal(tc.data, &opsgenie)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, opsgenie)

			data, err := json.Marshal(opsgenie)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
