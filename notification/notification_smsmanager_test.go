package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSMSManager_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.SMSManager
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My SMSManager Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SMSManager Alert\",\"smsmanagerApiKey\":\"test-api-key\",\"numbers\":\"420777123456\",\"messageType\":\"1\",\"type\":\"SMSManager\"}"}`),

			want: notification.SMSManager{
				Base: notification.Base{
					ID:            1,
					Name:          "My SMSManager Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SMSManagerDetails: notification.SMSManagerDetails{
					APIKey:      "test-api-key",
					Numbers:     "420777123456",
					MessageType: "1",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"messageType":"1","name":"My SMSManager Alert","numbers":"420777123456","smsmanagerApiKey":"test-api-key","type":"SMSManager","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple SMSManager","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SMSManager\",\"smsmanagerApiKey\":\"minimal-key\",\"numbers\":\"420123456789\",\"messageType\":\"0\",\"type\":\"SMSManager\"}"}`),

			want: notification.SMSManager{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SMSManager",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSManagerDetails: notification.SMSManagerDetails{
					APIKey:      "minimal-key",
					Numbers:     "420123456789",
					MessageType: "0",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"messageType":"0","name":"Simple SMSManager","numbers":"420123456789","smsmanagerApiKey":"minimal-key","type":"SMSManager","userId":1}`,
		},
		{
			name: "with different gateway type",
			data: []byte(`{"id":3,"name":"SMSManager Gateway 2","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SMSManager Gateway 2\",\"smsmanagerApiKey\":\"api-key-2\",\"numbers\":\"420987654321\",\"messageType\":\"2\",\"type\":\"SMSManager\"}"}`),

			want: notification.SMSManager{
				Base: notification.Base{
					ID:            3,
					Name:          "SMSManager Gateway 2",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SMSManagerDetails: notification.SMSManagerDetails{
					APIKey:      "api-key-2",
					Numbers:     "420987654321",
					MessageType: "2",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"messageType":"2","name":"SMSManager Gateway 2","numbers":"420987654321","smsmanagerApiKey":"api-key-2","type":"SMSManager","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smsmanager := notification.SMSManager{}

			err := json.Unmarshal(tc.data, &smsmanager)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smsmanager)

			data, err := json.Marshal(smsmanager)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
