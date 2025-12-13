package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMTP_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SMTP
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My SMTP Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMTP Alert\",\"smtpHost\":\"smtp.gmail.com\",\"smtpPort\":587,\"smtpSecure\":false,\"smtpIgnoreTLSError\":false,\"smtpFrom\":\"noreply@example.com\",\"smtpTo\":\"alerts@example.com\",\"customSubject\":\"Alert: {{ monitorJSON['name'] }}\",\"customBody\":\"Status: {{ msg }}\",\"htmlBody\":true,\"type\":\"smtp\"}"}`),

			want: notification.SMTP{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMTP Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMTPDetails: notification.SMTPDetails{
					Host:           "smtp.gmail.com",
					Port:           587,
					Secure:         false,
					IgnoreTLSError: false,
					From:           "noreply@example.com",
					To:             "alerts@example.com",
					CustomSubject:  "Alert: {{ monitorJSON['name'] }}",
					CustomBody:     "Status: {{ msg }}",
					HTMLBody:       true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"customBody":"Status: {{ msg }}","customSubject":"Alert: {{ monitorJSON['name'] }}","htmlBody":true,"id":1,"isDefault":true,"name":"My SMTP Alert","smtpCC":"","smtpBCC":"","smtpDkimDomain":"","smtpDkimHashAlgo":"","smtpDkimKeySelector":"","smtpDkimPrivateKey":"","smtpDkimheaderFieldNames":"","smtpDkimskipFields":"","smtpFrom":"noreply@example.com","smtpHost":"smtp.gmail.com","smtpIgnoreTLSError":false,"smtpPassword":"","smtpPort":587,"smtpSecure":false,"smtpTo":"alerts@example.com","smtpUsername":"","type":"smtp","userId":1}`,
		},
		{
			name: "with authentication and DKIM",
			data: []byte(`{"id":2,"name":"SMTP with Auth","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMTP with Auth\",\"smtpHost\":\"smtp.example.com\",\"smtpPort\":587,\"smtpSecure\":true,\"smtpIgnoreTLSError\":true,\"smtpUsername\":\"user@example.com\",\"smtpPassword\":\"secret\",\"smtpFrom\":\"sender@example.com\",\"smtpTo\":\"recipient@example.com\",\"smtpCC\":\"cc@example.com\",\"smtpBCC\":\"bcc@example.com\",\"smtpDkimDomain\":\"example.com\",\"smtpDkimKeySelector\":\"default\",\"smtpDkimPrivateKey\":\"-----BEGIN RSA PRIVATE KEY-----\\n...\",\"smtpDkimHashAlgo\":\"sha256\",\"htmlBody\":false,\"type\":\"smtp\"}"}`),

			want: notification.SMTP{
				Base: notification.Base{
					ID:            2,
					Name:          "SMTP with Auth",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMTPDetails: notification.SMTPDetails{
					Host:            "smtp.example.com",
					Port:            587,
					Secure:          true,
					IgnoreTLSError:  true,
					Username:        "user@example.com",
					Password:        "secret",
					From:            "sender@example.com",
					To:              "recipient@example.com",
					CC:              "cc@example.com",
					BCC:             "bcc@example.com",
					DkimDomain:      "example.com",
					DkimKeySelector: "default",
					DkimPrivateKey:  "-----BEGIN RSA PRIVATE KEY-----\n...",
					DkimHashAlgo:    "sha256",
					HTMLBody:        false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"customBody":"","customSubject":"","htmlBody":false,"id":2,"isDefault":false,"name":"SMTP with Auth","smtpBCC":"bcc@example.com","smtpCC":"cc@example.com","smtpDkimDomain":"example.com","smtpDkimHashAlgo":"sha256","smtpDkimKeySelector":"default","smtpDkimPrivateKey":"-----BEGIN RSA PRIVATE KEY-----\n...","smtpDkimheaderFieldNames":"","smtpDkimskipFields":"","smtpFrom":"sender@example.com","smtpHost":"smtp.example.com","smtpIgnoreTLSError":true,"smtpPassword":"secret","smtpPort":587,"smtpSecure":true,"smtpTo":"recipient@example.com","smtpUsername":"user@example.com","type":"smtp","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smtp := notification.SMTP{}

			err := json.Unmarshal(tc.data, &smtp)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smtp)

			data, err := json.Marshal(smtp)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
