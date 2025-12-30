package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSendGrid_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SendGrid
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My SendGrid Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SendGrid Alert\",\"sendgridApiKey\":\"SG.xxxxxxxxxxxxxxxx\",\"sendgridToEmail\":\"alerts@example.com\",\"sendgridFromEmail\":\"noreply@example.com\",\"sendgridSubject\":\"Uptime Alert\",\"sendgridCcEmail\":\"cc@example.com\",\"sendgridBccEmail\":\"bcc@example.com\",\"type\":\"SendGrid\"}"}`,
			),

			want: notification.SendGrid{
				Base: notification.Base{
					ID:            1,
					Name:          "My SendGrid Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SendGridDetails: notification.SendGridDetails{
					APIKey:    "SG.xxxxxxxxxxxxxxxx",
					ToEmail:   "alerts@example.com",
					FromEmail: "noreply@example.com",
					Subject:   "Uptime Alert",
					CcEmail:   "cc@example.com",
					BccEmail:  "bcc@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SendGrid Alert","sendgridApiKey":"SG.xxxxxxxxxxxxxxxx","sendgridBccEmail":"bcc@example.com","sendgridCcEmail":"cc@example.com","sendgridFromEmail":"noreply@example.com","sendgridSubject":"Uptime Alert","sendgridToEmail":"alerts@example.com","type":"SendGrid","userId":1}`,
		},
		{
			name: "minimal without CC and BCC",
			data: []byte(
				`{"id":2,"name":"Simple SendGrid","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SendGrid\",\"sendgridApiKey\":\"SG.yyyyyyyyyyyyyyyy\",\"sendgridToEmail\":\"user@example.com\",\"sendgridFromEmail\":\"sender@example.com\",\"sendgridSubject\":\"\",\"sendgridCcEmail\":\"\",\"sendgridBccEmail\":\"\",\"type\":\"SendGrid\"}"}`,
			),

			want: notification.SendGrid{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SendGrid",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SendGridDetails: notification.SendGridDetails{
					APIKey:    "SG.yyyyyyyyyyyyyyyy",
					ToEmail:   "user@example.com",
					FromEmail: "sender@example.com",
					Subject:   "",
					CcEmail:   "",
					BccEmail:  "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SendGrid","sendgridApiKey":"SG.yyyyyyyyyyyyyyyy","sendgridBccEmail":"","sendgridCcEmail":"","sendgridFromEmail":"sender@example.com","sendgridSubject":"","sendgridToEmail":"user@example.com","type":"SendGrid","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sendgrid := notification.SendGrid{}

			err := json.Unmarshal(tc.data, &sendgrid)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sendgrid)

			data, err := json.Marshal(sendgrid)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
