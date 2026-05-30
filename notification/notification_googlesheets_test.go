package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationGoogleSheets_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.GoogleSheets
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Google Sheets Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Google Sheets Alert\",\"googleSheetsWebhookUrl\":\"https://script.google.com/macros/s/AAAA/exec\",\"type\":\"GoogleSheets\"}"}`,
			),

			want: notification.GoogleSheets{
				Base: notification.Base{
					ID:            1,
					Name:          "My Google Sheets Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				GoogleSheetsDetails: notification.GoogleSheetsDetails{
					WebhookURL: "https://script.google.com/macros/s/AAAA/exec",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"googleSheetsWebhookUrl":"https://script.google.com/macros/s/AAAA/exec","id":1,"isDefault":true,"name":"My Google Sheets Alert","type":"GoogleSheets","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Google Sheets","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Google Sheets\",\"googleSheetsWebhookUrl\":\"https://script.google.com/macros/s/BBBB/exec\",\"type\":\"GoogleSheets\"}"}`,
			),

			want: notification.GoogleSheets{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Google Sheets",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				GoogleSheetsDetails: notification.GoogleSheetsDetails{
					WebhookURL: "https://script.google.com/macros/s/BBBB/exec",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"googleSheetsWebhookUrl":"https://script.google.com/macros/s/BBBB/exec","id":2,"isDefault":false,"name":"Simple Google Sheets","type":"GoogleSheets","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			googleSheets := notification.GoogleSheets{}

			err := json.Unmarshal(tc.data, &googleSheets)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, googleSheets)

			data, err := json.Marshal(googleSheets)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
