package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPlivo_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Plivo
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Plivo Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Plivo Alert\",\"plivoAnswerUrl\":\"https://example.com/answer.xml\",\"plivoAuthID\":\"MAXXXXXXXXXXXXXXXXXX\",\"plivoAuthToken\":\"test_auth_token\",\"plivoFromNumber\":\"+15559876543\",\"plivoMessageType\":\"call\",\"plivoToNumber\":\"+15551234567\",\"type\":\"plivo\"}"}`,
			),

			want: notification.Plivo{
				Base: notification.Base{
					ID:            1,
					Name:          "My Plivo Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				PlivoDetails: notification.PlivoDetails{
					AuthID:      "MAXXXXXXXXXXXXXXXXXX",
					AuthToken:   "test_auth_token",
					FromNumber:  "+15559876543",
					ToNumber:    "+15551234567",
					MessageType: notification.PlivoMessageTypeCall,
					AnswerURL:   ptr.To("https://example.com/answer.xml"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Plivo Alert","plivoAnswerUrl":"https://example.com/answer.xml","plivoAuthID":"MAXXXXXXXXXXXXXXXXXX","plivoAuthToken":"test_auth_token","plivoFromNumber":"+15559876543","plivoMessageType":"call","plivoToNumber":"+15551234567","type":"plivo","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Plivo","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Plivo\",\"plivoAuthID\":\"MAXXXXXXXXXXXXXXXXXX\",\"plivoAuthToken\":\"test_auth_token\",\"plivoFromNumber\":\"+15559876543\",\"plivoMessageType\":\"sms\",\"plivoToNumber\":\"+15551234567\",\"type\":\"plivo\"}"}`,
			),

			want: notification.Plivo{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Plivo",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PlivoDetails: notification.PlivoDetails{
					AuthID:      "MAXXXXXXXXXXXXXXXXXX",
					AuthToken:   "test_auth_token",
					FromNumber:  "+15559876543",
					ToNumber:    "+15551234567",
					MessageType: notification.PlivoMessageTypeSMS,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Plivo","plivoAuthID":"MAXXXXXXXXXXXXXXXXXX","plivoAuthToken":"test_auth_token","plivoFromNumber":"+15559876543","plivoMessageType":"sms","plivoToNumber":"+15551234567","type":"plivo","userId":1}`,
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
			plivo := notification.Plivo{}

			err := json.Unmarshal(tc.data, &plivo)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, plivo)

			data, err := json.Marshal(plivo)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
