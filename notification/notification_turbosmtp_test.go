package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationTurboSMTP_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.TurboSMTP
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My TurboSMTP Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My TurboSMTP Alert\",\"turbosmtpBccEmail\":\"bcc@example.com\",\"turbosmtpCcEmail\":\"cc1@example.com,cc2@example.com\",\"turbosmtpConsumerKey\":\"test_consumer_key\",\"turbosmtpConsumerSecret\":\"test_consumer_secret\",\"turbosmtpFromEmail\":\"alerts@example.com\",\"turbosmtpRegion\":\"eu\",\"turbosmtpSubject\":\"Uptime Kuma Alert\",\"turbosmtpToEmail\":\"ops@example.com,oncall@example.com\",\"type\":\"TurboSMTP\"}"}`,
			),

			want: notification.TurboSMTP{
				Base: notification.Base{
					ID:            1,
					Name:          "My TurboSMTP Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				TurboSMTPDetails: notification.TurboSMTPDetails{
					ConsumerKey:    "test_consumer_key",
					ConsumerSecret: "test_consumer_secret",
					Region:         notification.TurboSMTPRegionEU,
					FromEmail:      "alerts@example.com",
					ToEmail:        "ops@example.com,oncall@example.com",
					CcEmail:        new("cc1@example.com,cc2@example.com"),
					BccEmail:       new("bcc@example.com"),
					Subject:        new("Uptime Kuma Alert"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My TurboSMTP Alert","turbosmtpBccEmail":"bcc@example.com","turbosmtpCcEmail":"cc1@example.com,cc2@example.com","turbosmtpConsumerKey":"test_consumer_key","turbosmtpConsumerSecret":"test_consumer_secret","turbosmtpFromEmail":"alerts@example.com","turbosmtpRegion":"eu","turbosmtpSubject":"Uptime Kuma Alert","turbosmtpToEmail":"ops@example.com,oncall@example.com","type":"TurboSMTP","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple TurboSMTP","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple TurboSMTP\",\"turbosmtpConsumerKey\":\"test_consumer_key\",\"turbosmtpConsumerSecret\":\"test_consumer_secret\",\"turbosmtpFromEmail\":\"alerts@example.com\",\"turbosmtpToEmail\":\"ops@example.com\",\"type\":\"TurboSMTP\"}"}`,
			),

			want: notification.TurboSMTP{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple TurboSMTP",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TurboSMTPDetails: notification.TurboSMTPDetails{
					ConsumerKey:    "test_consumer_key",
					ConsumerSecret: "test_consumer_secret",
					FromEmail:      "alerts@example.com",
					ToEmail:        "ops@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple TurboSMTP","turbosmtpConsumerKey":"test_consumer_key","turbosmtpConsumerSecret":"test_consumer_secret","turbosmtpFromEmail":"alerts@example.com","turbosmtpToEmail":"ops@example.com","type":"TurboSMTP","userId":1}`,
		},
		{
			// The consumer key, the consumer secret and both addresses are
			// required upstream, so none of them carries omitempty and an empty
			// value must survive the round trip as an empty key. A dropped key
			// is silent data loss on update, because the config sent back to
			// the server is rebuilt from this struct.
			name: "empty fields are preserved",
			data: []byte(
				`{"id":3,"name":"Empty TurboSMTP","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty TurboSMTP\",\"turbosmtpConsumerKey\":\"\",\"turbosmtpConsumerSecret\":\"\",\"turbosmtpFromEmail\":\"\",\"turbosmtpToEmail\":\"\",\"type\":\"TurboSMTP\"}"}`,
			),

			want: notification.TurboSMTP{
				Base: notification.Base{
					ID:            3,
					Name:          "Empty TurboSMTP",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				TurboSMTPDetails: notification.TurboSMTPDetails{
					ConsumerKey:    "",
					ConsumerSecret: "",
					FromEmail:      "",
					ToEmail:        "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Empty TurboSMTP","turbosmtpConsumerKey":"","turbosmtpConsumerSecret":"","turbosmtpFromEmail":"","turbosmtpToEmail":"","type":"TurboSMTP","userId":1}`,
		},
		{
			name:    "missing config field",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false}`),
			wantErr: true,
		},
		{
			name:    "invalid config json",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"not-json"}`),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			turbosmtp := notification.TurboSMTP{}

			err := json.Unmarshal(tc.data, &turbosmtp)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, turbosmtp)

			data, err := json.Marshal(turbosmtp)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestTurboSMTPRegion_String(t *testing.T) {
	tests := []struct {
		region notification.TurboSMTPRegion
		want   string
	}{
		{notification.TurboSMTPRegionUS, "us"},
		{notification.TurboSMTPRegionEU, "eu"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			require.Equal(t, tc.want, tc.region.String())
		})
	}
}
