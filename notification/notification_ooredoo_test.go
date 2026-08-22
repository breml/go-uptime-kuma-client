package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOoredoo_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Ooredoo
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Ooredoo Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Ooredoo Alert\",\"ooredooAccessKey\":\"test_access_key\",\"ooredooBearerToken\":\"test_bearer_token\",\"ooredooServerUrl\":\"https://o-papi1-lb01.ooredoo.mv/bulk_sms/v2\",\"ooredooToNumber\":\"7712345, 9607798765\",\"ooredooUsername\":\"test_user\",\"type\":\"Ooredoo\"}"}`,
			),

			want: notification.Ooredoo{
				Base: notification.Base{
					ID:            1,
					Name:          "My Ooredoo Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				OoredooDetails: notification.OoredooDetails{
					Username:    "test_user",
					AccessKey:   "test_access_key",
					BearerToken: "test_bearer_token",
					ToNumber:    "7712345, 9607798765",
					ServerURL:   ptr.To("https://o-papi1-lb01.ooredoo.mv/bulk_sms/v2"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Ooredoo Alert","ooredooAccessKey":"test_access_key","ooredooBearerToken":"test_bearer_token","ooredooServerUrl":"https://o-papi1-lb01.ooredoo.mv/bulk_sms/v2","ooredooToNumber":"7712345, 9607798765","ooredooUsername":"test_user","type":"Ooredoo","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Ooredoo","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Ooredoo\",\"ooredooAccessKey\":\"test_access_key\",\"ooredooBearerToken\":\"test_bearer_token\",\"ooredooToNumber\":\"7712345\",\"ooredooUsername\":\"test_user\",\"type\":\"Ooredoo\"}"}`,
			),

			want: notification.Ooredoo{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Ooredoo",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OoredooDetails: notification.OoredooDetails{
					Username:    "test_user",
					AccessKey:   "test_access_key",
					BearerToken: "test_bearer_token",
					ToNumber:    "7712345",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Ooredoo","ooredooAccessKey":"test_access_key","ooredooBearerToken":"test_bearer_token","ooredooToNumber":"7712345","ooredooUsername":"test_user","type":"Ooredoo","userId":1}`,
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ooredoo := notification.Ooredoo{}

			err := json.Unmarshal(tc.data, &ooredoo)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, ooredoo)

			data, err := json.Marshal(ooredoo)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
