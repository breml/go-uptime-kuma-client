package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPinglet_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Pinglet
		wantJSON string
		wantErr  bool
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
			// The server strips a single trailing slash from the publish URL
			// before it posts, this client does not. The URL has to survive the
			// round trip exactly as stored, otherwise an edit would rewrite the
			// value the user entered.
			name: "trailing slash in publish url is not normalized",
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
		{
			// The publish URL and the API key are required upstream, so neither
			// carries omitempty and an empty value must survive the round trip
			// as an empty key. A dropped key is silent data loss on update,
			// because the config sent back to the server is rebuilt from this
			// struct, and a missing publish URL makes the provider throw
			// instead of merely failing the request.
			name: "empty fields are preserved",
			data: []byte(
				`{"id":4,"name":"Empty Pinglet","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty Pinglet\",\"pingletPublishUrl\":\"\",\"pingletApiKey\":\"\",\"type\":\"pinglet\"}"}`,
			),

			want: notification.Pinglet{
				Base: notification.Base{
					ID:            4,
					Name:          "Empty Pinglet",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PingletDetails: notification.PingletDetails{
					PublishURL: "",
					APIKey:     "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Empty Pinglet","pingletPublishUrl":"","pingletApiKey":"","type":"pinglet","userId":1}`,
		},
		{
			name:    "missing config field",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false}`),
			wantErr: true,
		},
		{
			name:    "invalid config json",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"not-json"}`),
			wantErr: true,
		},
		{
			name: "invalid config detail type",
			data: []byte(
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"pingletPublishUrl\":123,\"type\":\"pinglet\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pinglet := notification.Pinglet{}

			err := json.Unmarshal(tc.data, &pinglet)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pinglet)

			data, err := json.Marshal(pinglet)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
