package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSevenIO_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SevenIO
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Seven.io Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Seven.io Alert\",\"sevenioApiKey\":\"test-api-key\",\"sevenioSender\":\"UptimeKuma\",\"sevenioTo\":\"49123456789\",\"type\":\"sevenio\"}"}`,
			),

			want: notification.SevenIO{
				Base: notification.Base{
					ID:            1,
					Name:          "My Seven.io Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SevenIODetails: notification.SevenIODetails{
					APIKey: "test-api-key",
					Sender: "UptimeKuma",
					To:     "49123456789",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Seven.io Alert","sevenioApiKey":"test-api-key","sevenioSender":"UptimeKuma","sevenioTo":"49123456789","type":"sevenio","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Seven.io","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Seven.io\",\"sevenioApiKey\":\"minimal-key\",\"sevenioSender\":\"\",\"sevenioTo\":\"49999999999\",\"type\":\"sevenio\"}"}`,
			),

			want: notification.SevenIO{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Seven.io",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SevenIODetails: notification.SevenIODetails{
					APIKey: "minimal-key",
					Sender: "",
					To:     "49999999999",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Seven.io","sevenioApiKey":"minimal-key","sevenioSender":"","sevenioTo":"49999999999","type":"sevenio","userId":1}`,
		},
		{
			name: "with custom sender",
			data: []byte(
				`{"id":3,"name":"Seven.io Custom","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Seven.io Custom\",\"sevenioApiKey\":\"custom-key\",\"sevenioSender\":\"CustomAlert\",\"sevenioTo\":\"358501234567\",\"type\":\"sevenio\"}"}`,
			),

			want: notification.SevenIO{
				Base: notification.Base{
					ID:            3,
					Name:          "Seven.io Custom",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SevenIODetails: notification.SevenIODetails{
					APIKey: "custom-key",
					Sender: "CustomAlert",
					To:     "358501234567",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Seven.io Custom","sevenioApiKey":"custom-key","sevenioSender":"CustomAlert","sevenioTo":"358501234567","type":"sevenio","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sevenio := notification.SevenIO{}

			err := json.Unmarshal(tc.data, &sevenio)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sevenio)

			data, err := json.Marshal(sevenio)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
