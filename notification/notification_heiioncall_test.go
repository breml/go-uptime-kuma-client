package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationHeiiOnCall_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.HeiiOnCall
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Heii On-Call Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Heii On-Call Alert\",\"heiiOnCallApiKey\":\"api-key-123\",\"heiiOnCallTriggerId\":\"trigger-456\",\"type\":\"HeiiOnCall\"}"}`,
			),

			want: notification.HeiiOnCall{
				Base: notification.Base{
					ID:            1,
					Name:          "My Heii On-Call Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				HeiiOnCallDetails: notification.HeiiOnCallDetails{
					APIKey:    "api-key-123",
					TriggerID: "trigger-456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"heiiOnCallApiKey":"api-key-123","heiiOnCallTriggerId":"trigger-456","id":1,"isDefault":true,"name":"My Heii On-Call Alert","type":"HeiiOnCall","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Heii","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Heii\",\"heiiOnCallApiKey\":\"key-abc\",\"heiiOnCallTriggerId\":\"trig-123\",\"type\":\"HeiiOnCall\"}"}`,
			),

			want: notification.HeiiOnCall{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Heii",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				HeiiOnCallDetails: notification.HeiiOnCallDetails{
					APIKey:    "key-abc",
					TriggerID: "trig-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"heiiOnCallApiKey":"key-abc","heiiOnCallTriggerId":"trig-123","id":2,"isDefault":false,"name":"Simple Heii","type":"HeiiOnCall","userId":1}`,
		},
		{
			name: "with different credentials",
			data: []byte(
				`{"id":3,"name":"Heii Production","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Heii Production\",\"heiiOnCallApiKey\":\"prod-key-xyz\",\"heiiOnCallTriggerId\":\"prod-trigger-789\",\"type\":\"HeiiOnCall\"}"}`,
			),

			want: notification.HeiiOnCall{
				Base: notification.Base{
					ID:            3,
					Name:          "Heii Production",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				HeiiOnCallDetails: notification.HeiiOnCallDetails{
					APIKey:    "prod-key-xyz",
					TriggerID: "prod-trigger-789",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"heiiOnCallApiKey":"prod-key-xyz","heiiOnCallTriggerId":"prod-trigger-789","id":3,"isDefault":false,"name":"Heii Production","type":"HeiiOnCall","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			heiioncall := notification.HeiiOnCall{}

			err := json.Unmarshal(tc.data, &heiioncall)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, heiioncall)

			data, err := json.Marshal(heiioncall)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
