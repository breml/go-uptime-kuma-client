package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationAmootSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.AmootSMS
		wantJSON string
		wantErr  bool
	}{
		{
			name: "plain",
			data: []byte(
				`{"id":1,"name":"My Amoot SMS Alert","active":true,"userId":1,"isDefault":true,"config":"{\"amootApiToken\":\"test-token\",\"amootLineNumber\":\"public\",\"amootMobiles\":\"9123456789, 09987654321\",\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Amoot SMS Alert\",\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My Amoot SMS Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				AmootSMSDetails: notification.AmootSMSDetails{
					APIToken:   "test-token",
					Mobiles:    "9123456789, 09987654321",
					LineNumber: new("public"),
				},
			},
			wantJSON: `{"active":true,"amootApiToken":"test-token","amootLineNumber":"public","amootMobiles":"9123456789, 09987654321","applyExisting":true,"id":1,"isDefault":true,"name":"My Amoot SMS Alert","type":"amootsms","userId":1}`,
		},
		{
			name: "pattern",
			data: []byte(
				`{"id":2,"name":"Amoot SMS Pattern","active":true,"userId":1,"isDefault":false,"config":"{\"amootApiToken\":\"test-token\",\"amootLineNumber\":\"public\",\"amootMobiles\":\"9123456789\",\"amootPatternCodeId\":1234,\"amootUsePattern\":true,\"applyExisting\":false,\"isDefault\":false,\"name\":\"Amoot SMS Pattern\",\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:       2,
					Name:     "Amoot SMS Pattern",
					IsActive: true,
					UserID:   1,
				},
				AmootSMSDetails: notification.AmootSMSDetails{
					APIToken:      "test-token",
					Mobiles:       "9123456789",
					UsePattern:    new(true),
					PatternCodeID: new(int64(1234)),
					LineNumber:    new("public"),
				},
			},
			wantJSON: `{"active":true,"amootApiToken":"test-token","amootLineNumber":"public","amootMobiles":"9123456789","amootPatternCodeId":1234,"amootUsePattern":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Amoot SMS Pattern","type":"amootsms","userId":1}`,
		},
		{
			name: "pattern with own line",
			data: []byte(
				`{"id":3,"name":"Amoot SMS Own Line","active":true,"userId":1,"isDefault":false,"config":"{\"amootApiToken\":\"test-token\",\"amootLineNumber\":\"50001234\",\"amootMobiles\":\"9123456789,09987654321\",\"amootPatternCodeId\":5678,\"amootUseOwnLine\":true,\"amootUsePattern\":true,\"applyExisting\":false,\"isDefault\":false,\"name\":\"Amoot SMS Own Line\",\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:       3,
					Name:     "Amoot SMS Own Line",
					IsActive: true,
					UserID:   1,
				},
				AmootSMSDetails: notification.AmootSMSDetails{
					APIToken:      "test-token",
					Mobiles:       "9123456789,09987654321",
					UsePattern:    new(true),
					PatternCodeID: new(int64(5678)),
					UseOwnLine:    new(true),
					LineNumber:    new("50001234"),
				},
			},
			wantJSON: `{"active":true,"amootApiToken":"test-token","amootLineNumber":"50001234","amootMobiles":"9123456789,09987654321","amootPatternCodeId":5678,"amootUseOwnLine":true,"amootUsePattern":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Amoot SMS Own Line","type":"amootsms","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":4,"name":"Simple Amoot SMS","active":true,"userId":1,"isDefault":false,"config":"{\"amootApiToken\":\"test-token\",\"amootMobiles\":\"9123456789\",\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Amoot SMS\",\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:       4,
					Name:     "Simple Amoot SMS",
					IsActive: true,
					UserID:   1,
				},
				AmootSMSDetails: notification.AmootSMSDetails{
					APIToken: "test-token",
					Mobiles:  "9123456789",
				},
			},
			wantJSON: `{"active":true,"amootApiToken":"test-token","amootMobiles":"9123456789","applyExisting":false,"id":4,"isDefault":false,"name":"Simple Amoot SMS","type":"amootsms","userId":1}`,
		},
		{
			name: "empty required fields are kept",
			data: []byte(
				`{"id":5,"name":"Amoot SMS Empty","active":false,"userId":1,"isDefault":false,"config":"{\"amootApiToken\":\"\",\"amootMobiles\":\"\",\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:     5,
					Name:   "Amoot SMS Empty",
					UserID: 1,
				},
			},
			wantJSON: `{"active":false,"amootApiToken":"","amootMobiles":"","applyExisting":false,"id":5,"isDefault":false,"name":"Amoot SMS Empty","type":"amootsms","userId":1}`,
		},
		{
			name: "explicit false and empty values are preserved",
			data: []byte(
				`{"id":6,"name":"Amoot SMS Unticked","active":true,"userId":1,"isDefault":false,"config":"{\"amootApiToken\":\"test-token\",\"amootLineNumber\":\"\",\"amootMobiles\":\"9123456789\",\"amootUseOwnLine\":false,\"amootUsePattern\":false,\"type\":\"amootsms\"}"}`,
			),

			want: notification.AmootSMS{
				Base: notification.Base{
					ID:       6,
					Name:     "Amoot SMS Unticked",
					IsActive: true,
					UserID:   1,
				},
				AmootSMSDetails: notification.AmootSMSDetails{
					APIToken:   "test-token",
					Mobiles:    "9123456789",
					UsePattern: new(false),
					UseOwnLine: new(false),
					LineNumber: new(""),
				},
			},
			wantJSON: `{"active":true,"amootApiToken":"test-token","amootLineNumber":"","amootMobiles":"9123456789","amootUseOwnLine":false,"amootUsePattern":false,"applyExisting":false,"id":6,"isDefault":false,"name":"Amoot SMS Unticked","type":"amootsms","userId":1}`,
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
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"amootMobiles\":123,\"type\":\"amootsms\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			amootSMS := notification.AmootSMS{}

			err := json.Unmarshal(tc.data, &amootSMS)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, amootSMS)

			data, err := json.Marshal(amootSMS)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
