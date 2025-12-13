package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationMatrix_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Matrix
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Matrix Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Matrix Alert\",\"matrixHomeserverUrl\":\"https://matrix.example.com\",\"matrixInternalRoomId\":\"!roomid:example.com\",\"matrixAccessToken\":\"syt_token_example\",\"type\":\"matrix\"}"}`),

			want: notification.Matrix{
				Base: notification.Base{
					ID:            1,
					Name:          "My Matrix Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				MatrixDetails: notification.MatrixDetails{
					HomeserverURL:  "https://matrix.example.com",
					InternalRoomID: "!roomid:example.com",
					AccessToken:    "syt_token_example",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Matrix Alert","matrixHomeserverUrl":"https://matrix.example.com","matrixInternalRoomId":"!roomid:example.com","matrixAccessToken":"syt_token_example","type":"matrix","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Matrix","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Matrix\",\"matrixHomeserverUrl\":\"https://matrix.org\",\"matrixInternalRoomId\":\"!room:matrix.org\",\"matrixAccessToken\":\"token\",\"type\":\"matrix\"}"}`),

			want: notification.Matrix{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Matrix",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				MatrixDetails: notification.MatrixDetails{
					HomeserverURL:  "https://matrix.org",
					InternalRoomID: "!room:matrix.org",
					AccessToken:    "token",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Matrix","matrixHomeserverUrl":"https://matrix.org","matrixInternalRoomId":"!room:matrix.org","matrixAccessToken":"token","type":"matrix","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matrix := notification.Matrix{}

			err := json.Unmarshal(tc.data, &matrix)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, matrix)

			data, err := json.Marshal(matrix)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
