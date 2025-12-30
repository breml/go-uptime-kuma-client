package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSEagle_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SMSEagle
		wantJSON string
	}{
		{
			name: "API v1 SMS to phone number with priority",
			data: []byte(
				`{"id":1,"name":"My SMSEagle Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMSEagle Alert\",\"smseagleUrl\":\"https://192.168.1.100\",\"smseagleToken\":\"access-token-abc123\",\"smseagleRecipientType\":\"smseagle-to\",\"smseagleRecipient\":\"441234567890\",\"smseagleMsgType\":\"smseagle-sms\",\"smseaglePriority\":5,\"smseagleEncoding\":false,\"smseagleDuration\":0,\"smseagleTtsModel\":0,\"smseagleApiType\":\"smseagle-apiv1\",\"type\":\"SMSEagle\"}"}`,
			),

			want: notification.SMSEagle{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMSEagle Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSEagleDetails: notification.SMSEagleDetails{
					URL:           "https://192.168.1.100",
					Token:         "access-token-abc123",
					RecipientType: "smseagle-to",
					Recipient:     "441234567890",
					MsgType:       "smseagle-sms",
					Priority:      5,
					Encoding:      false,
					Duration:      0,
					TtsModel:      0,
					APIType:       "smseagle-apiv1",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SMSEagle Alert","smseagleApiType":"smseagle-apiv1","smseagleDuration":0,"smseagleEncoding":false,"smseagleMsgType":"smseagle-sms","smseaglePriority":5,"smseagleRecipient":"441234567890","smseagleRecipientContact":"","smseagleRecipientGroup":"","smseagleRecipientTo":"","smseagleRecipientType":"smseagle-to","smseagleTtsModel":0,"smseagleToken":"access-token-abc123","smseagleUrl":"https://192.168.1.100","type":"SMSEagle","userId":1}`,
		},
		{
			name: "API v2 SMS with contacts and unicode encoding",
			data: []byte(
				`{"id":2,"name":"SMSEagle Contact Alert","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSEagle Contact Alert\",\"smseagleUrl\":\"https://smseagle.example.com/device\",\"smseagleToken\":\"api-token-xyz789\",\"smseagleRecipientType\":\"smseagle-contact\",\"smseagleRecipient\":\"\",\"smseagleRecipientContact\":\"123,456,789\",\"smseagleMsgType\":\"smseagle-sms\",\"smseaglePriority\":3,\"smseagleEncoding\":true,\"smseagleDuration\":0,\"smseagleTtsModel\":0,\"smseagleApiType\":\"smseagle-apiv2\",\"type\":\"SMSEagle\"}"}`,
			),

			want: notification.SMSEagle{
				Base: notification.Base{
					ID:            2,
					Name:          "SMSEagle Contact Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSEagleDetails: notification.SMSEagleDetails{
					URL:              "https://smseagle.example.com/device",
					Token:            "api-token-xyz789",
					RecipientType:    "smseagle-contact",
					Recipient:        "",
					RecipientContact: "123,456,789",
					MsgType:          "smseagle-sms",
					Priority:         3,
					Encoding:         true,
					Duration:         0,
					TtsModel:         0,
					APIType:          "smseagle-apiv2",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"SMSEagle Contact Alert","smseagleApiType":"smseagle-apiv2","smseagleDuration":0,"smseagleEncoding":true,"smseagleMsgType":"smseagle-sms","smseaglePriority":3,"smseagleRecipient":"","smseagleRecipientContact":"123,456,789","smseagleRecipientGroup":"","smseagleRecipientTo":"","smseagleRecipientType":"smseagle-contact","smseagleTtsModel":0,"smseagleToken":"api-token-xyz789","smseagleUrl":"https://smseagle.example.com/device","type":"SMSEagle","userId":1}`,
		},
		{
			name: "API v2 TTS advanced call to group",
			data: []byte(
				`{"id":3,"name":"SMSEagle Voice Alert","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSEagle Voice Alert\",\"smseagleUrl\":\"https://smseagle.local\",\"smseagleToken\":\"voice-token-def456\",\"smseagleRecipientType\":\"smseagle-group\",\"smseagleRecipient\":\"\",\"smseagleRecipientGroup\":\"10,20,30\",\"smseagleMsgType\":\"smseagle-tts-advanced\",\"smseaglePriority\":1,\"smseagleEncoding\":false,\"smseagleDuration\":20,\"smseagleTtsModel\":2,\"smseagleApiType\":\"smseagle-apiv2\",\"type\":\"SMSEagle\"}"}`,
			),

			want: notification.SMSEagle{
				Base: notification.Base{
					ID:            3,
					Name:          "SMSEagle Voice Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSEagleDetails: notification.SMSEagleDetails{
					URL:            "https://smseagle.local",
					Token:          "voice-token-def456",
					RecipientType:  "smseagle-group",
					Recipient:      "",
					RecipientGroup: "10,20,30",
					MsgType:        "smseagle-tts-advanced",
					Priority:       1,
					Encoding:       false,
					Duration:       20,
					TtsModel:       2,
					APIType:        "smseagle-apiv2",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"SMSEagle Voice Alert","smseagleApiType":"smseagle-apiv2","smseagleDuration":20,"smseagleEncoding":false,"smseagleMsgType":"smseagle-tts-advanced","smseaglePriority":1,"smseagleRecipient":"","smseagleRecipientContact":"","smseagleRecipientGroup":"10,20,30","smseagleRecipientTo":"","smseagleRecipientType":"smseagle-group","smseagleTtsModel":2,"smseagleToken":"voice-token-def456","smseagleUrl":"https://smseagle.local","type":"SMSEagle","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smseagle := notification.SMSEagle{}

			err := json.Unmarshal(tc.data, &smseagle)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smseagle)

			data, err := json.Marshal(smseagle)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
