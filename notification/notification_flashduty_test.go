package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFlashDuty_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.FlashDuty
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My FlashDuty Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My FlashDuty Alert\",\"flashdutyIntegrationKey\":\"key123\",\"flashdutySeverity\":\"Critical\",\"type\":\"FlashDuty\"}"}`),

			want: notification.FlashDuty{
				Base: notification.Base{
					ID:            1,
					Name:          "My FlashDuty Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FlashDutyDetails: notification.FlashDutyDetails{
					IntegrationKey: "key123",
					Severity:       "Critical",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"flashdutyIntegrationKey":"key123","flashdutySeverity":"Critical","id":1,"isDefault":true,"name":"My FlashDuty Alert","type":"FlashDuty","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple FlashDuty","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple FlashDuty\",\"flashdutyIntegrationKey\":\"key456\",\"flashdutySeverity\":\"Warning\",\"type\":\"FlashDuty\"}"}`),

			want: notification.FlashDuty{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple FlashDuty",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FlashDutyDetails: notification.FlashDutyDetails{
					IntegrationKey: "key456",
					Severity:       "Warning",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"flashdutyIntegrationKey":"key456","flashdutySeverity":"Warning","id":2,"isDefault":false,"name":"Simple FlashDuty","type":"FlashDuty","userId":1}`,
		},
		{
			name: "with URL integration key",
			data: []byte(`{"id":3,"name":"FlashDuty URL","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"FlashDuty URL\",\"flashdutyIntegrationKey\":\"https://api.flashcat.cloud/event/push/alert/standard?integration_key=key789\",\"flashdutySeverity\":\"Info\",\"type\":\"FlashDuty\"}"}`),

			want: notification.FlashDuty{
				Base: notification.Base{
					ID:            3,
					Name:          "FlashDuty URL",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FlashDutyDetails: notification.FlashDutyDetails{
					IntegrationKey: "https://api.flashcat.cloud/event/push/alert/standard?integration_key=key789",
					Severity:       "Info",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"flashdutyIntegrationKey":"https://api.flashcat.cloud/event/push/alert/standard?integration_key=key789","flashdutySeverity":"Info","id":3,"isDefault":false,"name":"FlashDuty URL","type":"FlashDuty","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flashduty := notification.FlashDuty{}

			err := json.Unmarshal(tc.data, &flashduty)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, flashduty)

			data, err := json.Marshal(flashduty)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
