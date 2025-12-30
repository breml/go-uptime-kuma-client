package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTechulusPush_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.TechulusPush
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Techulus Push Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Techulus Push Alert\",\"pushAPIKey\":\"test-api-key\",\"pushTitle\":\"Alert Title\",\"pushSound\":\"default\",\"pushChannel\":\"alerts\",\"pushTimeSensitive\":true,\"type\":\"PushByTechulus\"}"}`,
			),

			want: notification.TechulusPush{
				Base: notification.Base{
					ID:            1,
					Name:          "My Techulus Push Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TechulusPushDetails: notification.TechulusPushDetails{
					APIKey:        "test-api-key",
					Title:         "Alert Title",
					Sound:         "default",
					Channel:       "alerts",
					TimeSensitive: true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Techulus Push Alert","pushAPIKey":"test-api-key","pushChannel":"alerts","pushSound":"default","pushTimeSensitive":true,"pushTitle":"Alert Title","type":"PushByTechulus","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Techulus Push","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Techulus Push\",\"pushAPIKey\":\"simple-key\",\"pushTitle\":\"\",\"pushSound\":\"\",\"pushChannel\":\"\",\"pushTimeSensitive\":false,\"type\":\"PushByTechulus\"}"}`,
			),

			want: notification.TechulusPush{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Techulus Push",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TechulusPushDetails: notification.TechulusPushDetails{
					APIKey:        "simple-key",
					Title:         "",
					Sound:         "",
					Channel:       "",
					TimeSensitive: false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Techulus Push","pushAPIKey":"simple-key","pushChannel":"","pushSound":"","pushTimeSensitive":false,"pushTitle":"","type":"PushByTechulus","userId":1}`,
		},
		{
			name: "with custom sound and channel",
			data: []byte(
				`{"id":3,"name":"Techulus Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Techulus Custom\",\"pushAPIKey\":\"custom-key\",\"pushTitle\":\"Custom Alert\",\"pushSound\":\"bell\",\"pushChannel\":\"monitoring\",\"pushTimeSensitive\":true,\"type\":\"PushByTechulus\"}"}`,
			),

			want: notification.TechulusPush{
				Base: notification.Base{
					ID:            3,
					Name:          "Techulus Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TechulusPushDetails: notification.TechulusPushDetails{
					APIKey:        "custom-key",
					Title:         "Custom Alert",
					Sound:         "bell",
					Channel:       "monitoring",
					TimeSensitive: true,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Techulus Custom","pushAPIKey":"custom-key","pushChannel":"monitoring","pushSound":"bell","pushTimeSensitive":true,"pushTitle":"Custom Alert","type":"PushByTechulus","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			techuluspush := notification.TechulusPush{}

			err := json.Unmarshal(tc.data, &techuluspush)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, techuluspush)

			data, err := json.Marshal(techuluspush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
