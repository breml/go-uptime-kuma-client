package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationEvolution_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Evolution
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Evolution Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Evolution Alert\",\"evolutionApiUrl\":\"https://evolapicloud.com\",\"evolutionInstanceName\":\"myinstance\",\"evolutionAuthToken\":\"token123\",\"evolutionRecipient\":\"5511999999999\",\"type\":\"EvolutionApi\"}"}`,
			),

			want: notification.Evolution{
				Base: notification.Base{
					ID:            1,
					Name:          "My Evolution Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				EvolutionDetails: notification.EvolutionDetails{
					ApiUrl:       "https://evolapicloud.com",
					InstanceName: "myinstance",
					AuthToken:    "token123",
					Recipient:    "5511999999999",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"evolutionApiUrl":"https://evolapicloud.com","evolutionAuthToken":"token123","evolutionInstanceName":"myinstance","evolutionRecipient":"5511999999999","id":1,"isDefault":true,"name":"My Evolution Alert","type":"EvolutionApi","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Evolution","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Evolution\",\"evolutionApiUrl\":\"https://api.example.com\",\"evolutionInstanceName\":\"instance1\",\"evolutionAuthToken\":\"key\",\"evolutionRecipient\":\"551187654321\",\"type\":\"EvolutionApi\"}"}`,
			),

			want: notification.Evolution{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Evolution",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				EvolutionDetails: notification.EvolutionDetails{
					ApiUrl:       "https://api.example.com",
					InstanceName: "instance1",
					AuthToken:    "key",
					Recipient:    "551187654321",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"evolutionApiUrl":"https://api.example.com","evolutionAuthToken":"key","evolutionInstanceName":"instance1","evolutionRecipient":"551187654321","id":2,"isDefault":false,"name":"Simple Evolution","type":"EvolutionApi","userId":1}`,
		},
		{
			name: "with different recipient",
			data: []byte(
				`{"id":3,"name":"Evolution US","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Evolution US\",\"evolutionApiUrl\":\"https://custom.api.com\",\"evolutionInstanceName\":\"usinstance\",\"evolutionAuthToken\":\"ustoken\",\"evolutionRecipient\":\"12025551234\",\"type\":\"EvolutionApi\"}"}`,
			),

			want: notification.Evolution{
				Base: notification.Base{
					ID:            3,
					Name:          "Evolution US",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				EvolutionDetails: notification.EvolutionDetails{
					ApiUrl:       "https://custom.api.com",
					InstanceName: "usinstance",
					AuthToken:    "ustoken",
					Recipient:    "12025551234",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"evolutionApiUrl":"https://custom.api.com","evolutionAuthToken":"ustoken","evolutionInstanceName":"usinstance","evolutionRecipient":"12025551234","id":3,"isDefault":false,"name":"Evolution US","type":"EvolutionApi","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			evolution := notification.Evolution{}

			err := json.Unmarshal(tc.data, &evolution)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, evolution)

			data, err := json.Marshal(evolution)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
