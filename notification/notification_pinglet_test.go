package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPinglet_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Pinglet
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Pinglet Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Pinglet Alert\",\"pingletPublishUrl\":\"https://app.pinglet.co.uk/my-namespace/alerts\",\"pingletApiKey\":\"test-api-key\",\"type\":\"pinglet\"}"}`,
			),

			want: notification.Pinglet{
				Base: notification.Base{
					ID:            1,
					Name:          "My Pinglet Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PingletDetails: notification.PingletDetails{
					PublishURL: "https://app.pinglet.co.uk/my-namespace/alerts",
					APIKey:     "test-api-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Pinglet Alert","pingletPublishUrl":"https://app.pinglet.co.uk/my-namespace/alerts","pingletApiKey":"test-api-key","type":"pinglet","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Pinglet","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Pinglet\",\"pingletPublishUrl\":\"https://app.pinglet.co.uk/simple/alerts\",\"pingletApiKey\":\"simple-key\",\"type\":\"pinglet\"}"}`,
			),

			want: notification.Pinglet{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Pinglet",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PingletDetails: notification.PingletDetails{
					PublishURL: "https://app.pinglet.co.uk/simple/alerts",
					APIKey:     "simple-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Pinglet","pingletPublishUrl":"https://app.pinglet.co.uk/simple/alerts","pingletApiKey":"simple-key","type":"pinglet","userId":1}`,
		},
		{
			name: "publish url with trailing slash",
			data: []byte(
				`{"id":3,"name":"Pinglet Trailing Slash","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Pinglet Trailing Slash\",\"pingletPublishUrl\":\"https://pinglet.example.com/team/alerts/\",\"pingletApiKey\":\"trailing-key-123\",\"type\":\"pinglet\"}"}`,
			),

			want: notification.Pinglet{
				Base: notification.Base{
					ID:            3,
					Name:          "Pinglet Trailing Slash",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PingletDetails: notification.PingletDetails{
					PublishURL: "https://pinglet.example.com/team/alerts/",
					APIKey:     "trailing-key-123",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Pinglet Trailing Slash","pingletPublishUrl":"https://pinglet.example.com/team/alerts/","pingletApiKey":"trailing-key-123","type":"pinglet","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pinglet := notification.Pinglet{}

			err := json.Unmarshal(tc.data, &pinglet)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pinglet)

			data, err := json.Marshal(pinglet)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
