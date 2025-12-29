package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationLine_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Line
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My LINE Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My LINE Alert\",\"lineChannelAccessToken\":\"channel-access-token-123\",\"lineUserID\":\"U1234567890abcdef1234567890abcdef\",\"type\":\"line\"}"}`,
			),

			want: notification.Line{
				Base: notification.Base{
					ID:            1,
					Name:          "My LINE Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				LineDetails: notification.LineDetails{
					ChannelAccessToken: "channel-access-token-123",
					UserID:             "U1234567890abcdef1234567890abcdef",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"lineChannelAccessToken":"channel-access-token-123","lineUserID":"U1234567890abcdef1234567890abcdef","name":"My LINE Alert","type":"line","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple LINE","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple LINE\",\"lineChannelAccessToken\":\"token-abc\",\"lineUserID\":\"U9876543210fedcba9876543210fedcba\",\"type\":\"line\"}"}`,
			),

			want: notification.Line{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple LINE",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LineDetails: notification.LineDetails{
					ChannelAccessToken: "token-abc",
					UserID:             "U9876543210fedcba9876543210fedcba",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"lineChannelAccessToken":"token-abc","lineUserID":"U9876543210fedcba9876543210fedcba","name":"Simple LINE","type":"line","userId":1}`,
		},
		{
			name: "with different credentials",
			data: []byte(
				`{"id":3,"name":"LINE Production","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"LINE Production\",\"lineChannelAccessToken\":\"prod-token-xyz\",\"lineUserID\":\"U5555555555555555555555555555555555\",\"type\":\"line\"}"}`,
			),

			want: notification.Line{
				Base: notification.Base{
					ID:            3,
					Name:          "LINE Production",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				LineDetails: notification.LineDetails{
					ChannelAccessToken: "prod-token-xyz",
					UserID:             "U5555555555555555555555555555555555",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"lineChannelAccessToken":"prod-token-xyz","lineUserID":"U5555555555555555555555555555555555","name":"LINE Production","type":"line","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line := notification.Line{}

			err := json.Unmarshal(tc.data, &line)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, line)

			data, err := json.Marshal(line)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
