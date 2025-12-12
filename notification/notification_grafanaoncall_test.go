package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGrafanaOncall_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.GrafanaOncall
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Grafana OnCall Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Grafana OnCall Alert\",\"grafanaOncallURL\":\"https://alerts.grafana.com/api/v1/incidents/create\",\"type\":\"GrafanaOncall\"}"}`),

			want: notification.GrafanaOncall{
				Base: notification.Base{
					ID:            1,
					Name:          "My Grafana OnCall Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GrafanaOncallDetails: notification.GrafanaOncallDetails{
					GrafanaOncallURL: "https://alerts.grafana.com/api/v1/incidents/create",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"grafanaOncallURL":"https://alerts.grafana.com/api/v1/incidents/create","id":1,"isDefault":true,"name":"My Grafana OnCall Alert","type":"GrafanaOncall","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Grafana OnCall","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Grafana OnCall\",\"grafanaOncallURL\":\"https://oncall.example.com/api/v1/incidents/create\",\"type\":\"GrafanaOncall\"}"}`),

			want: notification.GrafanaOncall{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Grafana OnCall",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GrafanaOncallDetails: notification.GrafanaOncallDetails{
					GrafanaOncallURL: "https://oncall.example.com/api/v1/incidents/create",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"grafanaOncallURL":"https://oncall.example.com/api/v1/incidents/create","id":2,"isDefault":false,"name":"Simple Grafana OnCall","type":"GrafanaOncall","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grafanaoncall := notification.GrafanaOncall{}

			err := json.Unmarshal(tc.data, &grafanaoncall)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, grafanaoncall)

			data, err := json.Marshal(grafanaoncall)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
