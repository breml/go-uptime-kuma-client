package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationResend_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Resend
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Resend Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Resend Alert\",\"resendApiKey\":\"re_xxxxxxxxxxxxxxxx\",\"resendFromEmail\":\"noreply@example.com\",\"resendFromName\":\"Uptime Kuma\",\"resendToEmail\":\"alerts@example.com\",\"resendSubject\":\"Uptime Alert\",\"type\":\"Resend\"}"}`,
			),

			want: notification.Resend{
				Base: notification.Base{
					ID:            1,
					Name:          "My Resend Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				ResendDetails: notification.ResendDetails{
					APIKey:    "re_xxxxxxxxxxxxxxxx",
					FromEmail: "noreply@example.com",
					FromName:  "Uptime Kuma",
					ToEmail:   "alerts@example.com",
					Subject:   "Uptime Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Resend Alert","resendApiKey":"re_xxxxxxxxxxxxxxxx","resendFromEmail":"noreply@example.com","resendFromName":"Uptime Kuma","resendSubject":"Uptime Alert","resendToEmail":"alerts@example.com","type":"Resend","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Resend","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Resend\",\"resendApiKey\":\"re_yyyyyyyyyyyyyyyy\",\"resendFromEmail\":\"sender@example.com\",\"resendToEmail\":\"user@example.com\",\"type\":\"Resend\"}"}`,
			),

			want: notification.Resend{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Resend",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ResendDetails: notification.ResendDetails{
					APIKey:    "re_yyyyyyyyyyyyyyyy",
					FromEmail: "sender@example.com",
					ToEmail:   "user@example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Resend","resendApiKey":"re_yyyyyyyyyyyyyyyy","resendFromEmail":"sender@example.com","resendToEmail":"user@example.com","type":"Resend","userId":1}`,
		},
		{
			name: "with multiple recipients",
			data: []byte(
				`{"id":3,"name":"Resend Multi To","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Resend Multi To\",\"resendApiKey\":\"re_zzzzzzzzzzzzzzzz\",\"resendFromEmail\":\"sender@example.com\",\"resendFromName\":\"Alert System\",\"resendToEmail\":\"to1@example.com, to2@example.com\",\"resendSubject\":\"System Alert\",\"type\":\"Resend\"}"}`,
			),

			want: notification.Resend{
				Base: notification.Base{
					ID:            3,
					Name:          "Resend Multi To",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				ResendDetails: notification.ResendDetails{
					APIKey:    "re_zzzzzzzzzzzzzzzz",
					FromEmail: "sender@example.com",
					FromName:  "Alert System",
					ToEmail:   "to1@example.com, to2@example.com",
					Subject:   "System Alert",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Resend Multi To","resendApiKey":"re_zzzzzzzzzzzzzzzz","resendFromEmail":"sender@example.com","resendFromName":"Alert System","resendSubject":"System Alert","resendToEmail":"to1@example.com, to2@example.com","type":"Resend","userId":1}`,
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resend := notification.Resend{}

			err := json.Unmarshal(tc.data, &resend)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, resend)

			data, err := json.Marshal(resend)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestNotificationResend_Type(t *testing.T) {
	require.Equal(t, "Resend", notification.Resend{}.Type())
	require.Equal(t, "Resend", notification.ResendDetails{}.Type())
}
