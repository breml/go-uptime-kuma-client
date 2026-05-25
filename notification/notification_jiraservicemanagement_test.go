package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationJiraServiceManagement_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.JiraServiceManagement
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My JSM Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My JSM Alert\",\"jsmCloudId\":\"cloud-123\",\"jsmEmail\":\"user@example.com\",\"jsmApiToken\":\"token-123\",\"jsmPriority\":1,\"type\":\"JiraServiceManagement\"}"}`,
			),

			want: notification.JiraServiceManagement{
				Base: notification.Base{
					ID:            1,
					Name:          "My JSM Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				JiraServiceManagementDetails: notification.JiraServiceManagementDetails{
					CloudID:  "cloud-123",
					Email:    "user@example.com",
					APIToken: "token-123",
					Priority: 1,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My JSM Alert","jsmCloudId":"cloud-123","jsmEmail":"user@example.com","jsmApiToken":"token-123","jsmPriority":1,"type":"JiraServiceManagement","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple JSM","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple JSM\",\"jsmCloudId\":\"cloud-456\",\"jsmEmail\":\"admin@example.com\",\"jsmApiToken\":\"token-456\",\"type\":\"JiraServiceManagement\"}"}`,
			),

			want: notification.JiraServiceManagement{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple JSM",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				JiraServiceManagementDetails: notification.JiraServiceManagementDetails{
					CloudID:  "cloud-456",
					Email:    "admin@example.com",
					APIToken: "token-456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple JSM","jsmCloudId":"cloud-456","jsmEmail":"admin@example.com","jsmApiToken":"token-456","jsmPriority":0,"type":"JiraServiceManagement","userId":1}`,
		},
		{
			name: "with low priority",
			data: []byte(
				`{"id":3,"name":"Low Priority JSM","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Low Priority JSM\",\"jsmCloudId\":\"cloud-789\",\"jsmEmail\":\"ops@example.com\",\"jsmApiToken\":\"token-789\",\"jsmPriority\":5,\"type\":\"JiraServiceManagement\"}"}`,
			),

			want: notification.JiraServiceManagement{
				Base: notification.Base{
					ID:            3,
					Name:          "Low Priority JSM",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				JiraServiceManagementDetails: notification.JiraServiceManagementDetails{
					CloudID:  "cloud-789",
					Email:    "ops@example.com",
					APIToken: "token-789",
					Priority: 5,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Low Priority JSM","jsmCloudId":"cloud-789","jsmEmail":"ops@example.com","jsmApiToken":"token-789","jsmPriority":5,"type":"JiraServiceManagement","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsm := notification.JiraServiceManagement{}

			err := json.Unmarshal(tc.data, &jsm)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, jsm)

			data, err := json.Marshal(jsm)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
