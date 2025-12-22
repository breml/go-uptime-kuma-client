package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWhapi_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Whapi
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My Whapi Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Whapi Alert\",\"whapiApiUrl\":\"https://gate.whapi.cloud\",\"whapiAuthToken\":\"test-auth-token\",\"whapiRecipient\":\"5511999999999\",\"type\":\"whapi\"}"}`),

			want: notification.Whapi{
				Base: notification.Base{
					ID:            1,
					Name:          "My Whapi Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WhapiDetails: notification.WhapiDetails{
					ApiURL:    "https://gate.whapi.cloud",
					AuthToken: "test-auth-token",
					Recipient: "5511999999999",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Whapi Alert","whapiApiUrl":"https://gate.whapi.cloud","whapiAuthToken":"test-auth-token","whapiRecipient":"5511999999999","type":"whapi","userId":1}`,
		},
		{
			name: "minimal configuration with default API URL",
			data: []byte(`{"id":2,"name":"Simple Whapi","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Whapi\",\"whapiAuthToken\":\"simple-token\",\"whapiRecipient\":\"+5511987654321\",\"type\":\"whapi\"}"}`),

			want: notification.Whapi{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Whapi",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WhapiDetails: notification.WhapiDetails{
					ApiURL:    "",
					AuthToken: "simple-token",
					Recipient: "+5511987654321",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Whapi","whapiApiUrl":"","whapiAuthToken":"simple-token","whapiRecipient":"+5511987654321","type":"whapi","userId":1}`,
		},
		{
			name: "with custom API URL",
			data: []byte(`{"id":3,"name":"Whapi Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Whapi Custom\",\"whapiApiUrl\":\"https://custom.whapi.io/\",\"whapiAuthToken\":\"custom-token-123\",\"whapiRecipient\":\"12025551234\",\"type\":\"whapi\"}"}`),

			want: notification.Whapi{
				Base: notification.Base{
					ID:            3,
					Name:          "Whapi Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WhapiDetails: notification.WhapiDetails{
					ApiURL:    "https://custom.whapi.io/",
					AuthToken: "custom-token-123",
					Recipient: "12025551234",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Whapi Custom","whapiApiUrl":"https://custom.whapi.io/","whapiAuthToken":"custom-token-123","whapiRecipient":"12025551234","type":"whapi","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whapi := notification.Whapi{}

			err := json.Unmarshal(tc.data, &whapi)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, whapi)

			data, err := json.Marshal(whapi)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
