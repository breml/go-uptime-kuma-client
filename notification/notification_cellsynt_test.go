package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationCellsynt_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Cellsynt
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":1,"name":"My Cellsynt Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Cellsynt Alert\",\"cellsyntLogin\":\"testuser\",\"cellsyntPassword\":\"testpass\",\"cellsyntDestination\":\"46701234567\",\"cellsyntOriginator\":\"Uptime\",\"cellsyntOriginatortype\":\"Numeric\",\"cellsyntAllowLongSMS\":true,\"type\":\"Cellsynt\"}"}`),

			want: notification.Cellsynt{
				Base: notification.Base{
					ID:            1,
					Name:          "My Cellsynt Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				CellsyntDetails: notification.CellsyntDetails{
					Login:          "testuser",
					Password:       "testpass",
					Destination:    "46701234567",
					Originator:     "Uptime",
					OriginatorType: "Numeric",
					AllowLongSMS:   true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"cellsyntAllowLongSMS":true,"cellsyntDestination":"46701234567","cellsyntLogin":"testuser","cellsyntOriginatortype":"Numeric","cellsyntOriginator":"Uptime","cellsyntPassword":"testpass","id":1,"isDefault":true,"name":"My Cellsynt Alert","type":"Cellsynt","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(`{"id":2,"name":"Simple Cellsynt","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Cellsynt\",\"cellsyntLogin\":\"user\",\"cellsyntPassword\":\"pass\",\"cellsyntDestination\":\"46700000000\",\"cellsyntOriginator\":\"Alert\",\"cellsyntOriginatortype\":\"Alphanumeric\",\"cellsyntAllowLongSMS\":false,\"type\":\"Cellsynt\"}"}`),

			want: notification.Cellsynt{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Cellsynt",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				CellsyntDetails: notification.CellsyntDetails{
					Login:          "user",
					Password:       "pass",
					Destination:    "46700000000",
					Originator:     "Alert",
					OriginatorType: "Alphanumeric",
					AllowLongSMS:   false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"cellsyntAllowLongSMS":false,"cellsyntDestination":"46700000000","cellsyntLogin":"user","cellsyntOriginatortype":"Alphanumeric","cellsyntOriginator":"Alert","cellsyntPassword":"pass","id":2,"isDefault":false,"name":"Simple Cellsynt","type":"Cellsynt","userId":1}`,
		},
		{
			name: "with long SMS allowed",
			data: []byte(`{"id":3,"name":"Cellsynt Long SMS","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Cellsynt Long SMS\",\"cellsyntLogin\":\"longuser\",\"cellsyntPassword\":\"longpass\",\"cellsyntDestination\":\"46799999999\",\"cellsyntOriginator\":\"LongMsg\",\"cellsyntOriginatortype\":\"Alphanumeric\",\"cellsyntAllowLongSMS\":true,\"type\":\"Cellsynt\"}"}`),

			want: notification.Cellsynt{
				Base: notification.Base{
					ID:            3,
					Name:          "Cellsynt Long SMS",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				CellsyntDetails: notification.CellsyntDetails{
					Login:          "longuser",
					Password:       "longpass",
					Destination:    "46799999999",
					Originator:     "LongMsg",
					OriginatorType: "Alphanumeric",
					AllowLongSMS:   true,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"cellsyntAllowLongSMS":true,"cellsyntDestination":"46799999999","cellsyntLogin":"longuser","cellsyntOriginatortype":"Alphanumeric","cellsyntOriginator":"LongMsg","cellsyntPassword":"longpass","id":3,"isDefault":false,"name":"Cellsynt Long SMS","type":"Cellsynt","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cellsynt := notification.Cellsynt{}

			err := json.Unmarshal(tc.data, &cellsynt)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, cellsynt)

			data, err := json.Marshal(cellsynt)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
