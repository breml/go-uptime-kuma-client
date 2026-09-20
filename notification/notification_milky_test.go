package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationMilky_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Milky
		wantJSON string
		wantErr  bool
	}{
		{
			name: "group message",
			data: []byte(
				`{"id":1,"name":"Test Milky Group","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"http://localhost:3000\",\"accessToken\":\"test-token\",\"msgType\":\"group\",\"recieverId\":\"123456789\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        1,
					Name:      "Test Milky Group",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "http://localhost:3000",
					AccessToken: "test-token",
					MsgType:     notification.MilkyMessageTypeGroup,
					ReceiverID:  "123456789",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"httpAddr":"http://localhost:3000","id":1,"isDefault":false,"msgType":"group","name":"Test Milky Group","recieverId":"123456789","type":"Milky","userId":42}`,
		},
		{
			name: "private message",
			data: []byte(
				`{"id":2,"name":"Test Milky Private","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"Milky\",\"httpAddr\":\"bot.example.com\",\"accessToken\":\"secret-token\",\"msgType\":\"private\",\"recieverId\":\"987654321\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        2,
					Name:      "Test Milky Private",
					IsActive:  true,
					UserID:    42,
					IsDefault: true,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "bot.example.com",
					AccessToken: "secret-token",
					MsgType:     notification.MilkyMessageTypePrivate,
					ReceiverID:  "987654321",
				},
			},
			wantJSON: `{"accessToken":"secret-token","active":true,"applyExisting":false,"httpAddr":"bot.example.com","id":2,"isDefault":true,"msgType":"private","name":"Test Milky Private","recieverId":"987654321","type":"Milky","userId":42}`,
		},
		{
			// The server appends "/api/" only if the address does not end with
			// a slash, so a trailing slash skips that segment and posts to an
			// endpoint that does not exist. This client does not normalize the
			// address, it has to survive the round trip exactly as stored,
			// otherwise an edit would rewrite the value the user entered.
			name: "trailing slash in http addr is not normalized",
			data: []byte(
				`{"id":3,"name":"Test Milky Trailing Slash","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"http://localhost:3000/\",\"accessToken\":\"test-token\",\"msgType\":\"group\",\"recieverId\":\"123456789\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        3,
					Name:      "Test Milky Trailing Slash",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "http://localhost:3000/",
					AccessToken: "test-token",
					MsgType:     notification.MilkyMessageTypeGroup,
					ReceiverID:  "123456789",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"httpAddr":"http://localhost:3000/","id":3,"isDefault":false,"msgType":"group","name":"Test Milky Trailing Slash","recieverId":"123456789","type":"Milky","userId":42}`,
		},
		{
			// The message type carries omitempty, because the form stores no
			// default. The server reads an absent and an empty message type
			// alike, so an empty value read from the server is written back as
			// an absent key.
			name: "empty message type is written back as an absent key",
			data: []byte(
				`{"id":4,"name":"Test Milky Empty Message Type","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"http://localhost:3000\",\"accessToken\":\"test-token\",\"msgType\":\"\",\"recieverId\":\"123456789\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        4,
					Name:      "Test Milky Empty Message Type",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "http://localhost:3000",
					AccessToken: "test-token",
					ReceiverID:  "123456789",
				},
			},
			wantJSON: `{"accessToken":"test-token","active":true,"applyExisting":false,"httpAddr":"http://localhost:3000","id":4,"isDefault":false,"name":"Test Milky Empty Message Type","recieverId":"123456789","type":"Milky","userId":42}`,
		},
		{
			// The address, the access token and the receiver ID are required
			// upstream, so none of them carries omitempty and an empty value
			// must survive the round trip as an empty key. A dropped key is
			// silent data loss on update, because the config sent back to the
			// server is rebuilt from this struct.
			name: "empty fields are preserved",
			data: []byte(
				`{"id":5,"name":"Test Milky Empty Fields","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"Milky\",\"httpAddr\":\"\",\"accessToken\":\"\",\"recieverId\":\"\"}"}`,
			),

			want: notification.Milky{
				Base: notification.Base{
					ID:        5,
					Name:      "Test Milky Empty Fields",
					IsActive:  false,
					UserID:    10,
					IsDefault: false,
				},
				MilkyDetails: notification.MilkyDetails{
					HTTPAddr:    "",
					AccessToken: "",
					ReceiverID:  "",
				},
			},
			wantJSON: `{"accessToken":"","active":false,"applyExisting":false,"httpAddr":"","id":5,"isDefault":false,"name":"Test Milky Empty Fields","recieverId":"","type":"Milky","userId":10}`,
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
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"httpAddr\":123,\"type\":\"Milky\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			milky := notification.Milky{}

			err := json.Unmarshal(tc.data, &milky)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, milky)

			data, err := json.Marshal(milky)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestMilkyMessageType_String(t *testing.T) {
	tests := []struct {
		messageType notification.MilkyMessageType
		want        string
	}{
		{notification.MilkyMessageTypeGroup, "group"},
		{notification.MilkyMessageTypePrivate, "private"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			require.Equal(t, tc.want, tc.messageType.String())
		})
	}
}
