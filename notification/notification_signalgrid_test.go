package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSignalgrid_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Signalgrid
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Signalgrid Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Signalgrid Alert\",\"signalgridClientKey\":\"client-key-123\",\"signalgridChannel\":\"alerts\",\"type\":\"signalgrid\"}"}`,
			),

			want: notification.Signalgrid{
				Base: notification.Base{
					ID:            1,
					Name:          "My Signalgrid Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SignalgridDetails: notification.SignalgridDetails{
					ClientKey: "client-key-123",
					Channel:   "alerts",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Signalgrid Alert","signalgridChannel":"alerts","signalgridClientKey":"client-key-123","type":"signalgrid","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Signalgrid","active":false,"userId":1,"isDefault":false,"config":"{\"name\":\"Simple Signalgrid\",\"signalgridClientKey\":\"key-abc\",\"signalgridChannel\":\"ops\",\"type\":\"signalgrid\"}"}`,
			),

			want: notification.Signalgrid{
				Base: notification.Base{
					ID:       2,
					Name:     "Simple Signalgrid",
					IsActive: false,
					UserID:   1,
				},
				SignalgridDetails: notification.SignalgridDetails{
					ClientKey: "key-abc",
					Channel:   "ops",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Signalgrid","signalgridChannel":"ops","signalgridClientKey":"key-abc","type":"signalgrid","userId":1}`,
		},
		{
			// The client key and channel are required in the Uptime Kuma form,
			// so neither carries omitempty and an empty value survives the
			// round trip as an empty key.
			name: "empty required fields are kept",
			data: []byte(
				`{"id":3,"name":"Signalgrid Empty","active":false,"userId":1,"isDefault":false,"config":"{\"signalgridClientKey\":\"\",\"signalgridChannel\":\"\",\"type\":\"signalgrid\"}"}`,
			),

			want: notification.Signalgrid{
				Base: notification.Base{
					ID:     3,
					Name:   "Signalgrid Empty",
					UserID: 1,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Signalgrid Empty","signalgridChannel":"","signalgridClientKey":"","type":"signalgrid","userId":1}`,
		},
		{
			name:    "missing config field",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false}`),
			wantErr: true,
		},
		{
			name:    "invalid config json",
			data:    []byte(`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"not-json"}`),
			wantErr: true,
		},
		{
			name: "invalid config detail type",
			data: []byte(
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"signalgridClientKey\":123,\"type\":\"signalgrid\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signalgrid := notification.Signalgrid{}

			err := json.Unmarshal(tc.data, &signalgrid)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, signalgrid)

			data, err := json.Marshal(signalgrid)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
