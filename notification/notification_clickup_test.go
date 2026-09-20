package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationClickUp_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.ClickUp
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My ClickUp Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My ClickUp Alert\",\"clickupToken\":\"pk_12345_ABCDEF\",\"clickupWorkspaceId\":\"9001234567\",\"clickupChannelId\":\"channel-abc\",\"clickupDisableUrl\":true,\"type\":\"ClickUp\"}"}`,
			),

			want: notification.ClickUp{
				Base: notification.Base{
					ID:            1,
					Name:          "My ClickUp Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ClickUpDetails: notification.ClickUpDetails{
					Token:       "pk_12345_ABCDEF",
					WorkspaceID: "9001234567",
					ChannelID:   "channel-abc",
					DisableURL:  ptr.To(true),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My ClickUp Alert","clickupToken":"pk_12345_ABCDEF","clickupWorkspaceId":"9001234567","clickupChannelId":"channel-abc","clickupDisableUrl":true,"type":"ClickUp","userId":1}`,
		},
		{
			// The checkbox writes its key as soon as it has been toggled, so an
			// explicit false is a value the server stores and must survive the
			// round trip.
			name: "url explicitly enabled",
			data: []byte(
				`{"id":2,"name":"ClickUp With Address","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"ClickUp With Address\",\"clickupToken\":\"pk_67890_GHIJKL\",\"clickupWorkspaceId\":\"9007654321\",\"clickupChannelId\":\"channel-xyz\",\"clickupDisableUrl\":false,\"type\":\"ClickUp\"}"}`,
			),

			want: notification.ClickUp{
				Base: notification.Base{
					ID:            2,
					Name:          "ClickUp With Address",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ClickUpDetails: notification.ClickUpDetails{
					Token:       "pk_67890_GHIJKL",
					WorkspaceID: "9007654321",
					ChannelID:   "channel-xyz",
					DisableURL:  ptr.To(false),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"ClickUp With Address","clickupToken":"pk_67890_GHIJKL","clickupWorkspaceId":"9007654321","clickupChannelId":"channel-xyz","clickupDisableUrl":false,"type":"ClickUp","userId":1}`,
		},
		{
			// An untouched checkbox is absent from the stored config. Upstream
			// treats the absent key and an explicit false alike, so nothing
			// changes for the message, but the key must stay absent to keep the
			// config this client sends back faithful to what the server holds.
			name: "minimal configuration without disable url",
			data: []byte(
				`{"id":3,"name":"Simple ClickUp","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple ClickUp\",\"clickupToken\":\"pk_11111_MNOPQR\",\"clickupWorkspaceId\":\"9001111111\",\"clickupChannelId\":\"channel-min\",\"type\":\"ClickUp\"}"}`,
			),

			want: notification.ClickUp{
				Base: notification.Base{
					ID:            3,
					Name:          "Simple ClickUp",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ClickUpDetails: notification.ClickUpDetails{
					Token:       "pk_11111_MNOPQR",
					WorkspaceID: "9001111111",
					ChannelID:   "channel-min",
					DisableURL:  nil,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Simple ClickUp","clickupToken":"pk_11111_MNOPQR","clickupWorkspaceId":"9001111111","clickupChannelId":"channel-min","type":"ClickUp","userId":1}`,
		},
		{
			// The three ids are required upstream, so none of them carries
			// omitempty and an empty value must survive the round trip as an
			// empty key. A dropped key is silent data loss on update, because
			// the config sent back to the server is rebuilt from this struct.
			name: "empty fields are preserved",
			data: []byte(
				`{"id":4,"name":"Empty ClickUp","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty ClickUp\",\"clickupToken\":\"\",\"clickupWorkspaceId\":\"\",\"clickupChannelId\":\"\",\"type\":\"ClickUp\"}"}`,
			),

			want: notification.ClickUp{
				Base: notification.Base{
					ID:            4,
					Name:          "Empty ClickUp",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ClickUpDetails: notification.ClickUpDetails{
					Token:       "",
					WorkspaceID: "",
					ChannelID:   "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Empty ClickUp","clickupToken":"","clickupWorkspaceId":"","clickupChannelId":"","type":"ClickUp","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clickup := notification.ClickUp{}

			err := json.Unmarshal(tc.data, &clickup)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, clickup)

			data, err := json.Marshal(clickup)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
