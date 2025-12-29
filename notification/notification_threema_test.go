package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationThreema_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Threema
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Threema Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Threema Alert\",\"threemaSenderIdentity\":\"GATEWAY1\",\"threemaSecret\":\"test-secret-key\",\"threemaRecipient\":\"USERID123\",\"threemaRecipientType\":\"identity\",\"type\":\"threema\"}"}`,
			),

			want: notification.Threema{
				Base: notification.Base{
					ID:            1,
					Name:          "My Threema Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ThreemaDetails: notification.ThreemaDetails{
					SenderIdentity: "GATEWAY1",
					Secret:         "test-secret-key",
					Recipient:      "USERID123",
					RecipientType:  "identity",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Threema Alert","threemaRecipient":"USERID123","threemaRecipientType":"identity","threemaSenderIdentity":"GATEWAY1","threemaSecret":"test-secret-key","type":"threema","userId":1}`,
		},
		{
			name: "minimal configuration with phone recipient",
			data: []byte(
				`{"id":2,"name":"Simple Threema","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Threema\",\"threemaSenderIdentity\":\"GATEWAY2\",\"threemaSecret\":\"secret\",\"threemaRecipient\":\"+41791234567\",\"threemaRecipientType\":\"phone\",\"type\":\"threema\"}"}`,
			),

			want: notification.Threema{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Threema",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ThreemaDetails: notification.ThreemaDetails{
					SenderIdentity: "GATEWAY2",
					Secret:         "secret",
					Recipient:      "+41791234567",
					RecipientType:  "phone",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Threema","threemaRecipient":"+41791234567","threemaRecipientType":"phone","threemaSenderIdentity":"GATEWAY2","threemaSecret":"secret","type":"threema","userId":1}`,
		},
		{
			name: "with email recipient",
			data: []byte(
				`{"id":3,"name":"Threema Email Alert","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Threema Email Alert\",\"threemaSenderIdentity\":\"GATEWAY3\",\"threemaSecret\":\"api-key\",\"threemaRecipient\":\"user@example.com\",\"threemaRecipientType\":\"email\",\"type\":\"threema\"}"}`,
			),

			want: notification.Threema{
				Base: notification.Base{
					ID:            3,
					Name:          "Threema Email Alert",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ThreemaDetails: notification.ThreemaDetails{
					SenderIdentity: "GATEWAY3",
					Secret:         "api-key",
					Recipient:      "user@example.com",
					RecipientType:  "email",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Threema Email Alert","threemaRecipient":"user@example.com","threemaRecipientType":"email","threemaSenderIdentity":"GATEWAY3","threemaSecret":"api-key","type":"threema","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			threema := notification.Threema{}

			err := json.Unmarshal(tc.data, &threema)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, threema)

			data, err := json.Marshal(threema)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
