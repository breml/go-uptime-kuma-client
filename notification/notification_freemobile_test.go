package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFreeMobile_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.FreeMobile
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My Free Mobile Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Free Mobile Alert\",\"freemobileUser\":\"12345678\",\"freemobilePass\":\"abcdef123456\",\"type\":\"FreeMobile\"}"}`),

			want: notification.FreeMobile{
				Base: notification.Base{
					ID:            1,
					Name:          "My Free Mobile Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FreeMobileDetails: notification.FreeMobileDetails{
					User: "12345678",
					Pass: "abcdef123456",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"freemobilePass":"abcdef123456","freemobileUser":"12345678","id":1,"isDefault":true,"name":"My Free Mobile Alert","type":"FreeMobile","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple Free Mobile","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Free Mobile\",\"freemobileUser\":\"87654321\",\"freemobilePass\":\"xyz789\",\"type\":\"FreeMobile\"}"}`),

			want: notification.FreeMobile{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Free Mobile",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FreeMobileDetails: notification.FreeMobileDetails{
					User: "87654321",
					Pass: "xyz789",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"freemobilePass":"xyz789","freemobileUser":"87654321","id":2,"isDefault":false,"name":"Simple Free Mobile","type":"FreeMobile","userId":1}`,
		},
		{
			name: "with different credentials",
			data: []byte(`{"id":3,"name":"Free Mobile Secondary","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Free Mobile Secondary\",\"freemobileUser\":\"99999999\",\"freemobilePass\":\"qwerty123456789\",\"type\":\"FreeMobile\"}"}`),

			want: notification.FreeMobile{
				Base: notification.Base{
					ID:            3,
					Name:          "Free Mobile Secondary",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FreeMobileDetails: notification.FreeMobileDetails{
					User: "99999999",
					Pass: "qwerty123456789",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"freemobilePass":"qwerty123456789","freemobileUser":"99999999","id":3,"isDefault":false,"name":"Free Mobile Secondary","type":"FreeMobile","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			freemobile := notification.FreeMobile{}

			err := json.Unmarshal(tc.data, &freemobile)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, freemobile)

			data, err := json.Marshal(freemobile)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
