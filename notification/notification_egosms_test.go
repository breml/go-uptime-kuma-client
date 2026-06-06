package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationEgoSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.EgoSMS
		wantJSON string
		wantErr  bool
	}{
		{
			name: "full with custom sender",
			data: []byte(
				`{"id":1,"name":"My EgoSMS Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My EgoSMS Alert\",\"egosmsUsername\":\"myuser\",\"egosmsPassword\":\"mypassword\",\"egosmsSender\":\"MYAPP\",\"egosmsPhoneNumber\":\"2567XXXXXXXX\",\"type\":\"egosms\"}"}`,
			),

			want: notification.EgoSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My EgoSMS Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				EgoSMSDetails: notification.EgoSMSDetails{
					Username:    "myuser",
					Password:    "mypassword",
					Sender:      ptr.To("MYAPP"),
					PhoneNumber: "2567XXXXXXXX",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"egosmsPassword":"mypassword","egosmsPhoneNumber":"2567XXXXXXXX","egosmsSender":"MYAPP","egosmsUsername":"myuser","id":1,"isDefault":true,"name":"My EgoSMS Alert","type":"egosms","userId":1}`,
		},
		{
			name: "minimal without sender",
			data: []byte(
				`{"id":2,"name":"Simple EgoSMS","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple EgoSMS\",\"egosmsUsername\":\"myuser\",\"egosmsPassword\":\"mypassword\",\"egosmsPhoneNumber\":\"2567XXXXXXXX\",\"type\":\"egosms\"}"}`,
			),

			want: notification.EgoSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple EgoSMS",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				EgoSMSDetails: notification.EgoSMSDetails{
					Username:    "myuser",
					Password:    "mypassword",
					PhoneNumber: "2567XXXXXXXX",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"egosmsPassword":"mypassword","egosmsPhoneNumber":"2567XXXXXXXX","egosmsUsername":"myuser","id":2,"isDefault":false,"name":"Simple EgoSMS","type":"egosms","userId":1}`,
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
			egosms := notification.EgoSMS{}

			err := json.Unmarshal(tc.data, &egosms)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, egosms)

			data, err := json.Marshal(egosms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
