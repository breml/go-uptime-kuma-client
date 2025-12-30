package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSPartner_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.SMSPartner
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My SMSPartner Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMSPartner Alert\",\"smspartnerApikey\":\"test-api-key\",\"smspartnerPhoneNumber\":\"33612345678\",\"smspartnerSenderName\":\"Uptime\",\"type\":\"SMSPartner\"}"}`,
			),

			want: notification.SMSPartner{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMSPartner Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSPartnerDetails: notification.SMSPartnerDetails{
					APIKey:      "test-api-key",
					PhoneNumber: "33612345678",
					SenderName:  "Uptime",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMSPartner Alert","smspartnerApikey":"test-api-key","smspartnerPhoneNumber":"33612345678","smspartnerSenderName":"Uptime","type":"SMSPartner","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple SMSPartner","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMSPartner\",\"smspartnerApikey\":\"minimal-key\",\"smspartnerPhoneNumber\":\"33687654321\",\"smspartnerSenderName\":\"Alert\",\"type\":\"SMSPartner\"}"}`,
			),

			want: notification.SMSPartner{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMSPartner",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSPartnerDetails: notification.SMSPartnerDetails{
					APIKey:      "minimal-key",
					PhoneNumber: "33687654321",
					SenderName:  "Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SMSPartner","smspartnerApikey":"minimal-key","smspartnerPhoneNumber":"33687654321","smspartnerSenderName":"Alert","type":"SMSPartner","userId":1}`,
		},
		{
			name: "with different phone number",
			data: []byte(
				`{"id":3,"name":"SMSPartner France","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSPartner France\",\"smspartnerApikey\":\"api-key-fr\",\"smspartnerPhoneNumber\":\"33698765432\",\"smspartnerSenderName\":\"Monitor\",\"type\":\"SMSPartner\"}"}`,
			),

			want: notification.SMSPartner{
				Base: notification.Base{
					ID:            3,
					Name:          "SMSPartner France",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSPartnerDetails: notification.SMSPartnerDetails{
					APIKey:      "api-key-fr",
					PhoneNumber: "33698765432",
					SenderName:  "Monitor",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SMSPartner France","smspartnerApikey":"api-key-fr","smspartnerPhoneNumber":"33698765432","smspartnerSenderName":"Monitor","type":"SMSPartner","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smspartner := notification.SMSPartner{}

			err := json.Unmarshal(tc.data, &smspartner)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smspartner)

			data, err := json.Marshal(smspartner)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
