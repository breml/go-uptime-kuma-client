package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWxPusher_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.WxPusher
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success with single spt",
			data: []byte(
				`{"id":1,"name":"My WxPusher Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My WxPusher Alert\",\"type\":\"WxPusher\",\"wxpusherSPT\":\"SPT_xxxxxxxxxxxx\"}"}`,
			),

			want: notification.WxPusher{
				Base: notification.Base{
					ID:            1,
					Name:          "My WxPusher Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WxPusherDetails: notification.WxPusherDetails{
					SPT: "SPT_xxxxxxxxxxxx",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My WxPusher Alert","type":"WxPusher","userId":1,"wxpusherSPT":"SPT_xxxxxxxxxxxx"}`,
		},
		{
			name: "comma separated spt list is preserved verbatim",
			data: []byte(
				`{"id":2,"name":"WxPusher Team","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"WxPusher Team\",\"type\":\"WxPusher\",\"wxpusherSPT\":\"SPT_aaaaaaaaaaaa, SPT_bbbbbbbbbbbb,SPT_cccccccccccc\"}"}`,
			),

			want: notification.WxPusher{
				Base: notification.Base{
					ID:            2,
					Name:          "WxPusher Team",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WxPusherDetails: notification.WxPusherDetails{
					SPT: "SPT_aaaaaaaaaaaa, SPT_bbbbbbbbbbbb,SPT_cccccccccccc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"WxPusher Team","type":"WxPusher","userId":1,"wxpusherSPT":"SPT_aaaaaaaaaaaa, SPT_bbbbbbbbbbbb,SPT_cccccccccccc"}`,
		},
		{
			name: "inactive notification",
			data: []byte(
				`{"id":3,"name":"Disabled WxPusher","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Disabled WxPusher\",\"type\":\"WxPusher\",\"wxpusherSPT\":\"SPT_dddddddddddd\"}"}`,
			),

			want: notification.WxPusher{
				Base: notification.Base{
					ID:            3,
					Name:          "Disabled WxPusher",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WxPusherDetails: notification.WxPusherDetails{
					SPT: "SPT_dddddddddddd",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Disabled WxPusher","type":"WxPusher","userId":1,"wxpusherSPT":"SPT_dddddddddddd"}`,
		},
		{
			name: "empty spt",
			data: []byte(
				`{"id":4,"name":"Empty WxPusher","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Empty WxPusher\",\"type\":\"WxPusher\",\"wxpusherSPT\":\"\"}"}`,
			),

			want: notification.WxPusher{
				Base: notification.Base{
					ID:            4,
					Name:          "Empty WxPusher",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WxPusherDetails: notification.WxPusherDetails{
					SPT: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Empty WxPusher","type":"WxPusher","userId":1,"wxpusherSPT":""}`,
		},
		{
			name: "spt missing from config is marshaled back as empty",
			data: []byte(
				`{"id":5,"name":"No SPT WxPusher","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"No SPT WxPusher\",\"type\":\"WxPusher\"}"}`,
			),

			want: notification.WxPusher{
				Base: notification.Base{
					ID:            5,
					Name:          "No SPT WxPusher",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WxPusherDetails: notification.WxPusherDetails{
					SPT: "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":5,"isDefault":false,"name":"No SPT WxPusher","type":"WxPusher","userId":1,"wxpusherSPT":""}`,
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
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"type\":\"WxPusher\",\"wxpusherSPT\":123}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wxpusher := notification.WxPusher{}

			err := json.Unmarshal(tc.data, &wxpusher)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, wxpusher)

			data, err := json.Marshal(wxpusher)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
