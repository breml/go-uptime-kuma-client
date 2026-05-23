package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWhatsapp360messenger_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Whatsapp360messenger
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(
				`{"id":1,"name":"My 360messenger Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My 360messenger Alert\",\"Whatsapp360messengerAuthToken\":\"token123\",\"Whatsapp360messengerRecipient\":\"447488888888\",\"Whatsapp360messengerGroupIds\":[\"group1\",\"group2\"],\"Whatsapp360messengerUseTemplate\":true,\"Whatsapp360messengerTemplate\":\"Alert: {{ msg }}\",\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            1,
					Name:          "My 360messenger Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken:   "token123",
					Recipient:   "447488888888",
					GroupIDs:    []string{"group1", "group2"},
					UseTemplate: ptr.To(true),
					Template:    ptr.To("Alert: {{ msg }}"),
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"token123","Whatsapp360messengerGroupIds":["group1","group2"],"Whatsapp360messengerRecipient":"447488888888","Whatsapp360messengerTemplate":"Alert: {{ msg }}","Whatsapp360messengerUseTemplate":true,"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My 360messenger Alert","type":"Whatsapp360messenger","userId":1}`,
		},
		{
			name: "minimal configuration with recipient only",
			data: []byte(
				`{"id":2,"name":"Simple 360messenger","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple 360messenger\",\"Whatsapp360messengerAuthToken\":\"key\",\"Whatsapp360messengerRecipient\":\"447499999999\",\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple 360messenger",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken: "key",
					Recipient: "447499999999",
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"key","Whatsapp360messengerRecipient":"447499999999","active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple 360messenger","type":"Whatsapp360messenger","userId":1}`,
		},
		{
			name: "with group IDs only (no recipient)",
			data: []byte(
				`{"id":3,"name":"Group 360messenger","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Group 360messenger\",\"Whatsapp360messengerAuthToken\":\"grptoken\",\"Whatsapp360messengerRecipient\":\"\",\"Whatsapp360messengerGroupIds\":[\"grp1\",\"grp2\",\"grp3\"],\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            3,
					Name:          "Group 360messenger",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken: "grptoken",
					Recipient: "",
					GroupIDs:  []string{"grp1", "grp2", "grp3"},
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"grptoken","Whatsapp360messengerGroupIds":["grp1","grp2","grp3"],"Whatsapp360messengerRecipient":"","active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Group 360messenger","type":"Whatsapp360messenger","userId":1}`,
		},
		{
			name: "with legacy GroupId field",
			data: []byte(
				`{"id":4,"name":"Legacy 360messenger","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Legacy 360messenger\",\"Whatsapp360messengerAuthToken\":\"legtoken\",\"Whatsapp360messengerRecipient\":\"\",\"Whatsapp360messengerGroupId\":\"legacygrp\",\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            4,
					Name:          "Legacy 360messenger",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken: "legtoken",
					Recipient: "",
					GroupID:   ptr.To("legacygrp"),
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"legtoken","Whatsapp360messengerGroupId":"legacygrp","Whatsapp360messengerRecipient":"","active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Legacy 360messenger","type":"Whatsapp360messenger","userId":1}`,
		},
		{
			name: "with template disabled (pointer to false)",
			data: []byte(
				`{"id":5,"name":"No Template 360messenger","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"No Template 360messenger\",\"Whatsapp360messengerAuthToken\":\"tok\",\"Whatsapp360messengerRecipient\":\"447488888888\",\"Whatsapp360messengerUseTemplate\":false,\"Whatsapp360messengerTemplate\":\"\",\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            5,
					Name:          "No Template 360messenger",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken:   "tok",
					Recipient:   "447488888888",
					UseTemplate: ptr.To(false),
					Template:    ptr.To(""),
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"tok","Whatsapp360messengerRecipient":"447488888888","Whatsapp360messengerTemplate":"","Whatsapp360messengerUseTemplate":false,"active":true,"applyExisting":false,"id":5,"isDefault":false,"name":"No Template 360messenger","type":"Whatsapp360messenger","userId":1}`,
		},
		{
			name: "with explicitly empty GroupIDs slice serializes as empty array",
			data: []byte(
				`{"id":6,"name":"Empty Groups 360messenger","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty Groups 360messenger\",\"Whatsapp360messengerAuthToken\":\"tok\",\"Whatsapp360messengerRecipient\":\"447488888888\",\"Whatsapp360messengerGroupIds\":[],\"type\":\"Whatsapp360messenger\"}"}`,
			),

			want: notification.Whatsapp360messenger{
				Base: notification.Base{
					ID:            6,
					Name:          "Empty Groups 360messenger",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				Whatsapp360messengerDetails: notification.Whatsapp360messengerDetails{
					AuthToken: "tok",
					Recipient: "447488888888",
					GroupIDs:  []string{},
				},
			},
			wantJSON: `{"Whatsapp360messengerAuthToken":"tok","Whatsapp360messengerGroupIds":[],"Whatsapp360messengerRecipient":"447488888888","active":true,"applyExisting":false,"id":6,"isDefault":false,"name":"Empty Groups 360messenger","type":"Whatsapp360messenger","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wa := notification.Whatsapp360messenger{}

			err := json.Unmarshal(tc.data, &wa)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, wa)

			data, err := json.Marshal(wa)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
