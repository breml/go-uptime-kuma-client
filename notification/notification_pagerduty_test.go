package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPagerDuty_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.PagerDuty
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My PagerDuty Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My PagerDuty Alert\",\"pagerdutyIntegrationUrl\":\"https://events.pagerduty.com/v2/enqueue\",\"pagerdutyIntegrationKey\":\"test-integration-key-123\",\"pagerdutyPriority\":\"warning\",\"pagerdutyAutoResolve\":\"resolve\",\"type\":\"PagerDuty\"}"}`,
			),

			want: notification.PagerDuty{
				Base: notification.Base{
					ID:            1,
					Name:          "My PagerDuty Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PagerDutyDetails: notification.PagerDutyDetails{
					IntegrationURL: "https://events.pagerduty.com/v2/enqueue",
					IntegrationKey: "test-integration-key-123",
					Priority:       "warning",
					AutoResolve:    "resolve",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My PagerDuty Alert","pagerdutyIntegrationUrl":"https://events.pagerduty.com/v2/enqueue","pagerdutyIntegrationKey":"test-integration-key-123","pagerdutyPriority":"warning","pagerdutyAutoResolve":"resolve","type":"PagerDuty","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple PagerDuty","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple PagerDuty\",\"pagerdutyIntegrationKey\":\"abc-123\",\"type\":\"PagerDuty\"}"}`,
			),

			want: notification.PagerDuty{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple PagerDuty",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PagerDutyDetails: notification.PagerDutyDetails{
					IntegrationKey: "abc-123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple PagerDuty","pagerdutyIntegrationUrl":"","pagerdutyIntegrationKey":"abc-123","pagerdutyPriority":"","pagerdutyAutoResolve":"","type":"PagerDuty","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pagerduty := notification.PagerDuty{}

			err := json.Unmarshal(tc.data, &pagerduty)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pagerduty)

			data, err := json.Marshal(pagerduty)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
