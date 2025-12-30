package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSplunk_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Splunk
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Splunk Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Splunk Alert\",\"splunkRestURL\":\"https://api.victorops.com/api/v1/incidents\",\"splunkSeverity\":\"CRITICAL\",\"splunkAutoResolve\":\"RECOVERY\",\"pagerdutyIntegrationKey\":\"test-routing-key\",\"type\":\"Splunk\"}"}`,
			),

			want: notification.Splunk{
				Base: notification.Base{
					ID:            1,
					Name:          "My Splunk Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SplunkDetails: notification.SplunkDetails{
					RestURL:        "https://api.victorops.com/api/v1/incidents",
					Severity:       "CRITICAL",
					AutoResolve:    "RECOVERY",
					IntegrationKey: "test-routing-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Splunk Alert","pagerdutyIntegrationKey":"test-routing-key","splunkAutoResolve":"RECOVERY","splunkRestURL":"https://api.victorops.com/api/v1/incidents","splunkSeverity":"CRITICAL","type":"Splunk","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Splunk","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Splunk\",\"splunkRestURL\":\"https://api.victorops.com/api/v1/incidents\",\"splunkSeverity\":\"HIGH\",\"splunkAutoResolve\":\"0\",\"pagerdutyIntegrationKey\":\"routing-key\",\"type\":\"Splunk\"}"}`,
			),

			want: notification.Splunk{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Splunk",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SplunkDetails: notification.SplunkDetails{
					RestURL:        "https://api.victorops.com/api/v1/incidents",
					Severity:       "HIGH",
					AutoResolve:    "0",
					IntegrationKey: "routing-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Splunk","pagerdutyIntegrationKey":"routing-key","splunkAutoResolve":"0","splunkRestURL":"https://api.victorops.com/api/v1/incidents","splunkSeverity":"HIGH","type":"Splunk","userId":1}`,
		},
		{
			name: "with warning severity",
			data: []byte(
				`{"id":3,"name":"Splunk Warning","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Splunk Warning\",\"splunkRestURL\":\"https://custom.victorops.com/api/v1/incidents\",\"splunkSeverity\":\"WARNING\",\"splunkAutoResolve\":\"ACKNOWLEDGED\",\"pagerdutyIntegrationKey\":\"custom-key\",\"type\":\"Splunk\"}"}`,
			),

			want: notification.Splunk{
				Base: notification.Base{
					ID:            3,
					Name:          "Splunk Warning",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SplunkDetails: notification.SplunkDetails{
					RestURL:        "https://custom.victorops.com/api/v1/incidents",
					Severity:       "WARNING",
					AutoResolve:    "ACKNOWLEDGED",
					IntegrationKey: "custom-key",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Splunk Warning","pagerdutyIntegrationKey":"custom-key","splunkAutoResolve":"ACKNOWLEDGED","splunkRestURL":"https://custom.victorops.com/api/v1/incidents","splunkSeverity":"WARNING","type":"Splunk","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			splunk := notification.Splunk{}

			err := json.Unmarshal(tc.data, &splunk)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, splunk)

			data, err := json.Marshal(splunk)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
