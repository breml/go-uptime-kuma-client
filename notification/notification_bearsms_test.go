package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationBearSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.BearSMS
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My BearSMS Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My BearSMS Alert\",\"bearsmsUsername\":\"bear-user\",\"bearsmsHashKey\":\"hash-key-123\",\"bearsmsSenderId\":\"UptimeKuma\",\"bearsmsPhoneNumber\":\"972501234567\",\"type\":\"bearsms\"}"}`,
			),

			want: notification.BearSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My BearSMS Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				BearSMSDetails: notification.BearSMSDetails{
					Username:    "bear-user",
					HashKey:     "hash-key-123",
					SenderID:    new("UptimeKuma"),
					PhoneNumber: "972501234567",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"bearsmsUsername":"bear-user","bearsmsHashKey":"hash-key-123","bearsmsSenderId":"UptimeKuma","bearsmsPhoneNumber":"972501234567","id":1,"isDefault":true,"name":"My BearSMS Alert","type":"bearsms","userId":1}`,
		},
		{
			name: "minimal configuration without sender id",
			data: []byte(
				`{"id":2,"name":"Simple BearSMS","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple BearSMS\",\"bearsmsUsername\":\"user\",\"bearsmsHashKey\":\"key\",\"bearsmsPhoneNumber\":\"972509876543\",\"type\":\"bearsms\"}"}`,
			),

			want: notification.BearSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple BearSMS",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BearSMSDetails: notification.BearSMSDetails{
					Username:    "user",
					HashKey:     "key",
					PhoneNumber: "972509876543",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"bearsmsUsername":"user","bearsmsHashKey":"key","bearsmsPhoneNumber":"972509876543","id":2,"isDefault":false,"name":"Simple BearSMS","type":"bearsms","userId":1}`,
		},
		{
			name: "empty sender id is preserved",
			data: []byte(
				`{"id":3,"name":"BearSMS Empty Sender","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"BearSMS Empty Sender\",\"bearsmsUsername\":\"user\",\"bearsmsHashKey\":\"key\",\"bearsmsSenderId\":\"\",\"bearsmsPhoneNumber\":\"972501112222\",\"type\":\"bearsms\"}"}`,
			),

			want: notification.BearSMS{
				Base: notification.Base{
					ID:            3,
					Name:          "BearSMS Empty Sender",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BearSMSDetails: notification.BearSMSDetails{
					Username:    "user",
					HashKey:     "key",
					SenderID:    new(""),
					PhoneNumber: "972501112222",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"bearsmsUsername":"user","bearsmsHashKey":"key","bearsmsSenderId":"","bearsmsPhoneNumber":"972501112222","id":3,"isDefault":false,"name":"BearSMS Empty Sender","type":"bearsms","userId":1}`,
		},
		{
			// The username, the hash key and the phone number are required
			// upstream, so none of them carries omitempty and an empty value
			// must survive the round trip as an empty key. A dropped key is
			// silent data loss on update, because the config sent back to the
			// server is rebuilt from this struct.
			name: "empty fields are preserved",
			data: []byte(
				`{"id":4,"name":"Empty BearSMS","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty BearSMS\",\"bearsmsUsername\":\"\",\"bearsmsHashKey\":\"\",\"bearsmsPhoneNumber\":\"\",\"type\":\"bearsms\"}"}`,
			),

			want: notification.BearSMS{
				Base: notification.Base{
					ID:            4,
					Name:          "Empty BearSMS",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				BearSMSDetails: notification.BearSMSDetails{
					Username:    "",
					HashKey:     "",
					PhoneNumber: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"bearsmsUsername":"","bearsmsHashKey":"","bearsmsPhoneNumber":"","id":4,"isDefault":false,"name":"Empty BearSMS","type":"bearsms","userId":1}`,
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
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"bearsmsUsername\":123,\"type\":\"bearsms\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bearsms := notification.BearSMS{}

			err := json.Unmarshal(tc.data, &bearsms)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, bearsms)

			data, err := json.Marshal(bearsms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
