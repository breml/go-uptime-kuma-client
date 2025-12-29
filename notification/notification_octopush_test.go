package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationOctopush_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Octopush
		wantJSON string
	}{
		{
			name: "success V2",
			data: []byte(
				`{"id":1,"name":"Test Octopush V2","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"octopush\",\"octopushVersion\":\"2\",\"octopushAPIKey\":\"test-api-key\",\"octopushLogin\":\"testuser\",\"octopushPhoneNumber\":\"+33612345678\",\"octopushSMSType\":\"sms_premium\",\"octopushSenderName\":\"AlertBot\"}"}`,
			),

			want: notification.Octopush{
				Base: notification.Base{
					ID:        1,
					Name:      "Test Octopush V2",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				OctopushDetails: notification.OctopushDetails{
					Version:     "2",
					APIKey:      "test-api-key",
					Login:       "testuser",
					PhoneNumber: "+33612345678",
					SMSType:     "sms_premium",
					SenderName:  "AlertBot",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"Test Octopush V2","octopushAPIKey":"test-api-key","octopushDMAPIKey":"","octopushDMLogin":"","octopushDMPhoneNumber":"","octopushDMSMSType":"","octopushDMSenderName":"","octopushLogin":"testuser","octopushPhoneNumber":"+33612345678","octopushSMSType":"sms_premium","octopushSenderName":"AlertBot","octopushVersion":"2","type":"octopush","userId":42}`,
		},
		{
			name: "success V1",
			data: []byte(
				`{"id":2,"name":"Test Octopush V1","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"octopush\",\"octopushVersion\":\"1\",\"octopushDMLogin\":\"dmuser\",\"octopushDMAPIKey\":\"dm-api-key\",\"octopushDMPhoneNumber\":\"+33698765432\",\"octopushDMSMSType\":\"sms_premium\",\"octopushDMSenderName\":\"SMSBot\"}"}`,
			),

			want: notification.Octopush{
				Base: notification.Base{
					ID:        2,
					Name:      "Test Octopush V1",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				OctopushDetails: notification.OctopushDetails{
					Version:       "1",
					DMLogin:       "dmuser",
					DMAPIKey:      "dm-api-key",
					DMPhoneNumber: "+33698765432",
					DMSMSType:     "sms_premium",
					DMSenderName:  "SMSBot",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Test Octopush V1","octopushAPIKey":"","octopushDMAPIKey":"dm-api-key","octopushDMLogin":"dmuser","octopushDMPhoneNumber":"+33698765432","octopushDMSMSType":"sms_premium","octopushDMSenderName":"SMSBot","octopushLogin":"","octopushPhoneNumber":"","octopushSMSType":"","octopushSenderName":"","octopushVersion":"1","type":"octopush","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test Octopush Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"octopush\",\"octopushVersion\":\"\"}"}`,
			),

			want: notification.Octopush{
				Base: notification.Base{
					ID:            3,
					Name:          "Test Octopush Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				OctopushDetails: notification.OctopushDetails{
					Version: "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Test Octopush Minimal","octopushAPIKey":"","octopushDMAPIKey":"","octopushDMLogin":"","octopushDMPhoneNumber":"","octopushDMSMSType":"","octopushDMSenderName":"","octopushLogin":"","octopushPhoneNumber":"","octopushSMSType":"","octopushSenderName":"","octopushVersion":"","type":"octopush","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			octopush := notification.Octopush{}

			err := json.Unmarshal(tc.data, &octopush)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, octopush)

			data, err := json.Marshal(octopush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
