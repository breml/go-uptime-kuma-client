package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTeltonika_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Teltonika
		wantJSON string
	}{
		{
			name: "single phone number with TLS validation",
			data: []byte(
				`{"id":1,"name":"My Teltonika Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Teltonika Alert\",\"teltonikaUrl\":\"https://192.168.100.1\",\"teltonikaUsername\":\"admin\",\"teltonikaPassword\":\"secret-password\",\"teltonikaModem\":\"1-1\",\"teltonikaPhoneNumber\":\"+336xxxxxxxx\",\"teltonikaUnsafeTls\":false,\"type\":\"Teltonika\"}"}`,
			),

			want: notification.Teltonika{
				Base: notification.Base{
					ID:            1,
					Name:          "My Teltonika Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TeltonikaDetails: notification.TeltonikaDetails{
					URL:         "https://192.168.100.1",
					Username:    "admin",
					Password:    "secret-password",
					Modem:       "1-1",
					PhoneNumber: "+336xxxxxxxx",
					UnsafeTLS:   false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Teltonika Alert","teltonikaModem":"1-1","teltonikaPassword":"secret-password","teltonikaPhoneNumber":"+336xxxxxxxx","teltonikaUnsafeTls":false,"teltonikaUrl":"https://192.168.100.1","teltonikaUsername":"admin","type":"Teltonika","userId":1}`,
		},
		{
			name: "multiple recipients with unsafe TLS",
			data: []byte(
				`{"id":2,"name":"Teltonika Self-Signed","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Teltonika Self-Signed\",\"teltonikaUrl\":\"http://teltonika.example.com:8080\",\"teltonikaUsername\":\"opsuser\",\"teltonikaPassword\":\"another-secret\",\"teltonikaModem\":\"2-1\",\"teltonikaPhoneNumber\":\"+336xxxxxxxx,+496xxxxxxxx\",\"teltonikaUnsafeTls\":true,\"type\":\"Teltonika\"}"}`,
			),

			want: notification.Teltonika{
				Base: notification.Base{
					ID:            2,
					Name:          "Teltonika Self-Signed",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TeltonikaDetails: notification.TeltonikaDetails{
					URL:         "http://teltonika.example.com:8080",
					Username:    "opsuser",
					Password:    "another-secret",
					Modem:       "2-1",
					PhoneNumber: "+336xxxxxxxx,+496xxxxxxxx",
					UnsafeTLS:   true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Teltonika Self-Signed","teltonikaModem":"2-1","teltonikaPassword":"another-secret","teltonikaPhoneNumber":"+336xxxxxxxx,+496xxxxxxxx","teltonikaUnsafeTls":true,"teltonikaUrl":"http://teltonika.example.com:8080","teltonikaUsername":"opsuser","type":"Teltonika","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teltonika := notification.Teltonika{}

			err := json.Unmarshal(tc.data, &teltonika)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, teltonika)

			data, err := json.Marshal(teltonika)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
