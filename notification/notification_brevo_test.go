package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationBrevo_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Brevo
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Brevo Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Brevo Alert\",\"brevoApiKey\":\"test-api-key\",\"brevoToEmail\":\"recipient@example.com\",\"brevoFromEmail\":\"sender@example.com\",\"brevoFromName\":\"Uptime Kuma\",\"brevoSubject\":\"Alert Notification\",\"brevoCcEmail\":\"cc@example.com\",\"brevoBccEmail\":\"bcc@example.com\",\"type\":\"brevo\"}"}`,
			),

			want: notification.Brevo{
				Base: notification.Base{
					ID:            1,
					Name:          "My Brevo Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				BrevoDetails: notification.BrevoDetails{
					APIKey:    "test-api-key",
					ToEmail:   "recipient@example.com",
					FromEmail: "sender@example.com",
					FromName:  "Uptime Kuma",
					Subject:   "Alert Notification",
					CCEmail:   "cc@example.com",
					BCCEmail:  "bcc@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"brevoBccEmail":"bcc@example.com","brevoCcEmail":"cc@example.com","brevoApiKey":"test-api-key","brevoFromEmail":"sender@example.com","brevoFromName":"Uptime Kuma","brevoSubject":"Alert Notification","brevoToEmail":"recipient@example.com","id":1,"isDefault":true,"name":"My Brevo Alert","type":"brevo","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Brevo","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Brevo\",\"brevoApiKey\":\"minimal-key\",\"brevoToEmail\":\"to@example.com\",\"brevoFromEmail\":\"from@example.com\",\"type\":\"brevo\"}"}`,
			),

			want: notification.Brevo{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Brevo",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BrevoDetails: notification.BrevoDetails{
					APIKey:    "minimal-key",
					ToEmail:   "to@example.com",
					FromEmail: "from@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"brevoBccEmail":"","brevoCcEmail":"","brevoApiKey":"minimal-key","brevoFromEmail":"from@example.com","brevoFromName":"","brevoSubject":"","brevoToEmail":"to@example.com","id":2,"isDefault":false,"name":"Simple Brevo","type":"brevo","userId":1}`,
		},
		{
			name: "with multiple CC and BCC recipients",
			data: []byte(
				`{"id":3,"name":"Brevo Multi CC","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Brevo Multi CC\",\"brevoApiKey\":\"api-key\",\"brevoToEmail\":\"main@example.com\",\"brevoFromEmail\":\"sender@example.com\",\"brevoFromName\":\"Alert System\",\"brevoSubject\":\"System Alert\",\"brevoCcEmail\":\"cc1@example.com, cc2@example.com\",\"brevoBccEmail\":\"bcc1@example.com, bcc2@example.com\",\"type\":\"brevo\"}"}`,
			),

			want: notification.Brevo{
				Base: notification.Base{
					ID:            3,
					Name:          "Brevo Multi CC",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BrevoDetails: notification.BrevoDetails{
					APIKey:    "api-key",
					ToEmail:   "main@example.com",
					FromEmail: "sender@example.com",
					FromName:  "Alert System",
					Subject:   "System Alert",
					CCEmail:   "cc1@example.com, cc2@example.com",
					BCCEmail:  "bcc1@example.com, bcc2@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"brevoBccEmail":"bcc1@example.com, bcc2@example.com","brevoCcEmail":"cc1@example.com, cc2@example.com","brevoApiKey":"api-key","brevoFromEmail":"sender@example.com","brevoFromName":"Alert System","brevoSubject":"System Alert","brevoToEmail":"main@example.com","id":3,"isDefault":false,"name":"Brevo Multi CC","type":"brevo","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			brevo := notification.Brevo{}

			err := json.Unmarshal(tc.data, &brevo)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, brevo)

			data, err := json.Marshal(brevo)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
