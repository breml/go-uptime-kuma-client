package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWAHA_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.WAHA
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My WAHA Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My WAHA Alert\",\"wahaApiUrl\":\"https://waha.example.com\",\"wahaSession\":\"default\",\"wahaChatId\":\"5511999999999\",\"wahaApiKey\":\"test-api-key\",\"type\":\"waha\"}"}`),

			want: notification.WAHA{
				Base: notification.Base{
					ID:            1,
					Name:          "My WAHA Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WAHADetails: notification.WAHADetails{
					ApiURL:  "https://waha.example.com",
					Session: "default",
					ChatID:  "5511999999999",
					ApiKey:  "test-api-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My WAHA Alert","wahaApiKey":"test-api-key","wahaApiUrl":"https://waha.example.com","wahaChatId":"5511999999999","wahaSession":"default","type":"waha","userId":1}`,
		},
		{
			name: "minimal configuration without API key",
			data: []byte(`{"id":2,"name":"Simple WAHA","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple WAHA\",\"wahaApiUrl\":\"https://api.waha.cloud\",\"wahaSession\":\"main\",\"wahaChatId\":\"+5511987654321\",\"type\":\"waha\"}"}`),

			want: notification.WAHA{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple WAHA",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WAHADetails: notification.WAHADetails{
					ApiURL:  "https://api.waha.cloud",
					Session: "main",
					ChatID:  "+5511987654321",
					ApiKey:  "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple WAHA","wahaApiKey":"","wahaApiUrl":"https://api.waha.cloud","wahaChatId":"+5511987654321","wahaSession":"main","type":"waha","userId":1}`,
		},
		{
			name: "with different session and chat ID",
			data: []byte(`{"id":3,"name":"WAHA Multi Session","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"WAHA Multi Session\",\"wahaApiUrl\":\"https://custom.waha.io/\",\"wahaSession\":\"alerts\",\"wahaChatId\":\"120201234567@c.us\",\"wahaApiKey\":\"secure-key-123\",\"type\":\"waha\"}"}`),

			want: notification.WAHA{
				Base: notification.Base{
					ID:            3,
					Name:          "WAHA Multi Session",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WAHADetails: notification.WAHADetails{
					ApiURL:  "https://custom.waha.io/",
					Session: "alerts",
					ChatID:  "120201234567@c.us",
					ApiKey:  "secure-key-123",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"WAHA Multi Session","wahaApiKey":"secure-key-123","wahaApiUrl":"https://custom.waha.io/","wahaChatId":"120201234567@c.us","wahaSession":"alerts","type":"waha","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			waha := notification.WAHA{}

			err := json.Unmarshal(tc.data, &waha)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, waha)

			data, err := json.Marshal(waha)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
