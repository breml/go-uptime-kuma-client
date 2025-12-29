package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationNostr_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		want     notification.Nostr
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My Nostr Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Nostr Alert\",\"sender\":\"nsec1sender\",\"recipients\":\"npub1recipient1\\nnpub1recipient2\",\"relays\":\"wss://relay1.example.com\\nwss://relay2.example.com\",\"type\":\"nostr\"}"}`,
			),

			want: notification.Nostr{
				Base: notification.Base{
					ID:            1,
					Name:          "My Nostr Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				NostrDetails: notification.NostrDetails{
					Sender:     "nsec1sender",
					Recipients: "npub1recipient1\nnpub1recipient2",
					Relays:     "wss://relay1.example.com\nwss://relay2.example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Nostr Alert","recipients":"npub1recipient1\nnpub1recipient2","relays":"wss://relay1.example.com\nwss://relay2.example.com","sender":"nsec1sender","type":"nostr","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple Nostr","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Nostr\",\"sender\":\"nsec1test\",\"recipients\":\"npub1single\",\"relays\":\"wss://relay.example.com\",\"type\":\"nostr\"}"}`,
			),

			want: notification.Nostr{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Nostr",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				NostrDetails: notification.NostrDetails{
					Sender:     "nsec1test",
					Recipients: "npub1single",
					Relays:     "wss://relay.example.com",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Nostr","recipients":"npub1single","relays":"wss://relay.example.com","sender":"nsec1test","type":"nostr","userId":1}`,
		},
		{
			name: "with multiple relays",
			data: []byte(
				`{"id":3,"name":"Nostr Multi-Relay","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Nostr Multi-Relay\",\"sender\":\"nsec1prod\",\"recipients\":\"npub1user1\\nnpub1user2\\nnpub1user3\",\"relays\":\"wss://relay1.example.com\\nwss://relay2.example.com\\nwss://relay3.example.com\",\"type\":\"nostr\"}"}`,
			),

			want: notification.Nostr{
				Base: notification.Base{
					ID:            3,
					Name:          "Nostr Multi-Relay",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				NostrDetails: notification.NostrDetails{
					Sender:     "nsec1prod",
					Recipients: "npub1user1\nnpub1user2\nnpub1user3",
					Relays:     "wss://relay1.example.com\nwss://relay2.example.com\nwss://relay3.example.com",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Nostr Multi-Relay","recipients":"npub1user1\nnpub1user2\nnpub1user3","relays":"wss://relay1.example.com\nwss://relay2.example.com\nwss://relay3.example.com","sender":"nsec1prod","type":"nostr","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nostr := notification.Nostr{}

			err := json.Unmarshal(tc.data, &nostr)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, nostr)

			data, err := json.Marshal(nostr)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
