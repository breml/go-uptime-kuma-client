package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGoAlert_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.GoAlert
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My GoAlert Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My GoAlert Alert\",\"goAlertBaseURL\":\"https://goalert.example.com\",\"goAlertToken\":\"token123\",\"type\":\"GoAlert\"}"}`),

			want: notification.GoAlert{
				Base: notification.Base{
					ID:            1,
					Name:          "My GoAlert Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GoAlertDetails: notification.GoAlertDetails{
					BaseURL: "https://goalert.example.com",
					Token:   "token123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"goAlertBaseURL":"https://goalert.example.com","goAlertToken":"token123","id":1,"isDefault":true,"name":"My GoAlert Alert","type":"GoAlert","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple GoAlert","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple GoAlert\",\"goAlertBaseURL\":\"https://alerts.example.com\",\"goAlertToken\":\"abc123\",\"type\":\"GoAlert\"}"}`),

			want: notification.GoAlert{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple GoAlert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GoAlertDetails: notification.GoAlertDetails{
					BaseURL: "https://alerts.example.com",
					Token:   "abc123",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"goAlertBaseURL":"https://alerts.example.com","goAlertToken":"abc123","id":2,"isDefault":false,"name":"Simple GoAlert","type":"GoAlert","userId":1}`,
		},
		{
			name: "with different base URL",
			data: []byte(`{"id":3,"name":"GoAlert Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"GoAlert Custom\",\"goAlertBaseURL\":\"https://custom.domain.com/goalert\",\"goAlertToken\":\"customtoken789\",\"type\":\"GoAlert\"}"}`),

			want: notification.GoAlert{
				Base: notification.Base{
					ID:            3,
					Name:          "GoAlert Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GoAlertDetails: notification.GoAlertDetails{
					BaseURL: "https://custom.domain.com/goalert",
					Token:   "customtoken789",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"goAlertBaseURL":"https://custom.domain.com/goalert","goAlertToken":"customtoken789","id":3,"isDefault":false,"name":"GoAlert Custom","type":"GoAlert","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			goalert := notification.GoAlert{}

			err := json.Unmarshal(tc.data, &goalert)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, goalert)

			data, err := json.Marshal(goalert)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
