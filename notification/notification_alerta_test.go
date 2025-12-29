package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationAlerta_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Alerta
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Alerta Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Alerta Alert\",\"alertaApiEndpoint\":\"https://alerta.example.com/api/alerts\",\"alertaApiKey\":\"api_key_secret\",\"alertaEnvironment\":\"Production\",\"alertaAlertState\":\"critical\",\"alertaRecoverState\":\"cleared\",\"type\":\"alerta\"}"}`,
			),

			want: notification.Alerta{
				Base: notification.Base{
					ID:            1,
					Name:          "My Alerta Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				AlertaDetails: notification.AlertaDetails{
					APIEndpoint:  "https://alerta.example.com/api/alerts",
					APIKey:       "api_key_secret",
					Environment:  "Production",
					AlertState:   "critical",
					RecoverState: "cleared",
				},
			},
			wantJSON: `{"active":true,"alertaAlertState":"critical","alertaApiEndpoint":"https://alerta.example.com/api/alerts","alertaApiKey":"api_key_secret","alertaEnvironment":"Production","alertaRecoverState":"cleared","applyExisting":true,"id":1,"isDefault":true,"name":"My Alerta Alert","type":"alerta","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Alerta","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Alerta\",\"alertaApiEndpoint\":\"https://alerta.io/api/alerts\",\"alertaApiKey\":\"key123\",\"alertaEnvironment\":\"Dev\",\"alertaAlertState\":\"critical\",\"alertaRecoverState\":\"cleared\",\"type\":\"alerta\"}"}`,
			),

			want: notification.Alerta{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Alerta",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				AlertaDetails: notification.AlertaDetails{
					APIEndpoint:  "https://alerta.io/api/alerts",
					APIKey:       "key123",
					Environment:  "Dev",
					AlertState:   "critical",
					RecoverState: "cleared",
				},
			},
			wantJSON: `{"active":true,"alertaAlertState":"critical","alertaApiEndpoint":"https://alerta.io/api/alerts","alertaApiKey":"key123","alertaEnvironment":"Dev","alertaRecoverState":"cleared","applyExisting":false,"id":2,"isDefault":false,"name":"Simple Alerta","type":"alerta","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			alerta := notification.Alerta{}

			err := json.Unmarshal(tc.data, &alerta)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, alerta)

			data, err := json.Marshal(alerta)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
