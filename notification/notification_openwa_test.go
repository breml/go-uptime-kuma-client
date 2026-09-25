package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOpenWa_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.OpenWa
		wantJSON string
	}{
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":1,"name":"My OpenWA Alert","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"My OpenWA Alert\",\"openwaApiUrl\":\"http://localhost:2785\",\"openwaApiKey\":\"secret-key\",\"openwaSession\":\"default\",\"openwaChatId\":\"00117612345678@c.us\",\"type\":\"openwa\"}"}`,
			),

			want: notification.OpenWa{
				Base: notification.Base{
					ID:        1,
					Name:      "My OpenWA Alert",
					IsActive:  true,
					UserID:    1,
					IsDefault: false,
				},
				OpenWaDetails: notification.OpenWaDetails{
					APIURL:  "http://localhost:2785",
					APIKey:  "secret-key",
					Session: "default",
					ChatID:  "00117612345678@c.us",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"My OpenWA Alert","openwaApiKey":"secret-key","openwaApiUrl":"http://localhost:2785","openwaChatId":"00117612345678@c.us","openwaSession":"default","type":"openwa","userId":1}`,
		},
		{
			name: "with custom message and multiple chat ids",
			data: []byte(
				`{"id":2,"name":"OpenWA Template","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"OpenWA Template\",\"openwaApiUrl\":\"https://wa.example.com\",\"openwaApiKey\":\"key123\",\"openwaSession\":\"alerts\",\"openwaChatId\":\"00117612345678@c.us, 123456789012345678@g.us,1234567890@lid\",\"openwaUseCustomMessage\":true,\"openwaCustomMessage\":\"Alert: {{ msg }}\",\"type\":\"openwa\"}"}`,
			),

			want: notification.OpenWa{
				Base: notification.Base{
					ID:            2,
					Name:          "OpenWA Template",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				OpenWaDetails: notification.OpenWaDetails{
					APIURL:           "https://wa.example.com",
					APIKey:           "key123",
					Session:          "alerts",
					ChatID:           "00117612345678@c.us, 123456789012345678@g.us,1234567890@lid",
					UseCustomMessage: new(true),
					CustomMessage:    new("Alert: {{ msg }}"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":2,"isDefault":true,"name":"OpenWA Template","openwaApiKey":"key123","openwaApiUrl":"https://wa.example.com","openwaChatId":"00117612345678@c.us, 123456789012345678@g.us,1234567890@lid","openwaCustomMessage":"Alert: {{ msg }}","openwaSession":"alerts","openwaUseCustomMessage":true,"type":"openwa","userId":1}`,
		},
		{
			name: "pointer to false and empty string are serialized",
			data: []byte(
				`{"id":3,"name":"OpenWA Explicit False","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"OpenWA Explicit False\",\"openwaApiUrl\":\"http://localhost:2785\",\"openwaApiKey\":\"secret-key\",\"openwaSession\":\"default\",\"openwaChatId\":\"00117612345678@c.us\",\"openwaUseCustomMessage\":false,\"openwaCustomMessage\":\"\",\"type\":\"openwa\"}"}`,
			),

			want: notification.OpenWa{
				Base: notification.Base{
					ID:        3,
					Name:      "OpenWA Explicit False",
					IsActive:  true,
					UserID:    1,
					IsDefault: false,
				},
				OpenWaDetails: notification.OpenWaDetails{
					APIURL:           "http://localhost:2785",
					APIKey:           "secret-key",
					Session:          "default",
					ChatID:           "00117612345678@c.us",
					UseCustomMessage: new(false),
					CustomMessage:    new(""),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"OpenWA Explicit False","openwaApiKey":"secret-key","openwaApiUrl":"http://localhost:2785","openwaChatId":"00117612345678@c.us","openwaCustomMessage":"","openwaSession":"default","openwaUseCustomMessage":false,"type":"openwa","userId":1}`,
		},
		{
			// The server strips trailing slashes from the API URL, this client
			// does not normalize it, otherwise an edit would rewrite the value
			// the user entered.
			name: "trailing slash in api url is not normalized",
			data: []byte(
				`{"id":4,"name":"OpenWA Trailing Slash","active":true,"userId":1,"isDefault":false,"config":"{\"openwaApiUrl\":\"http://localhost:2785/\",\"openwaApiKey\":\"secret-key\",\"openwaSession\":\"default\",\"openwaChatId\":\"00117612345678@c.us\",\"type\":\"openwa\"}"}`,
			),

			want: notification.OpenWa{
				Base: notification.Base{
					ID:       4,
					Name:     "OpenWA Trailing Slash",
					IsActive: true,
					UserID:   1,
				},
				OpenWaDetails: notification.OpenWaDetails{
					APIURL:  "http://localhost:2785/",
					APIKey:  "secret-key",
					Session: "default",
					ChatID:  "00117612345678@c.us",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"OpenWA Trailing Slash","openwaApiKey":"secret-key","openwaApiUrl":"http://localhost:2785/","openwaChatId":"00117612345678@c.us","openwaSession":"default","type":"openwa","userId":1}`,
		},
		{
			// The API URL, API key, session and chat ID are required upstream,
			// so none of them carries omitempty and an empty value survives the
			// round trip as an empty key.
			name: "empty required fields are kept",
			data: []byte(
				`{"id":5,"name":"OpenWA Empty","active":false,"userId":1,"isDefault":false,"config":"{\"openwaApiUrl\":\"\",\"openwaApiKey\":\"\",\"openwaSession\":\"\",\"openwaChatId\":\"\",\"type\":\"openwa\"}"}`,
			),

			want: notification.OpenWa{
				Base: notification.Base{
					ID:     5,
					Name:   "OpenWA Empty",
					UserID: 1,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":5,"isDefault":false,"name":"OpenWA Empty","openwaApiKey":"","openwaApiUrl":"","openwaChatId":"","openwaSession":"","type":"openwa","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			openwa := notification.OpenWa{}

			err := json.Unmarshal(tc.data, &openwa)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, openwa)

			data, err := json.Marshal(openwa)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
