package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSignal_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Signal
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Signal Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Signal Alert\",\"signalURL\":\"http://localhost:9998\",\"signalNumber\":\"+1234567890\",\"signalRecipients\":\"+9876543210,+1112223333\",\"type\":\"signal\"}"}`),

			want: notification.Signal{
				Base: notification.Base{
					ID:            1,
					Name:          "My Signal Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SignalDetails: notification.SignalDetails{
					URL:        "http://localhost:9998",
					Number:     "+1234567890",
					Recipients: "+9876543210,+1112223333",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Signal Alert","signalURL":"http://localhost:9998","signalNumber":"+1234567890","signalRecipients":"+9876543210,+1112223333","type":"signal","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Signal","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Signal\",\"signalNumber\":\"+1234567890\",\"signalRecipients\":\"+9999999999\",\"type\":\"signal\"}"}`),

			want: notification.Signal{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Signal",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SignalDetails: notification.SignalDetails{
					Number:     "+1234567890",
					Recipients: "+9999999999",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Signal","signalURL":"","signalNumber":"+1234567890","signalRecipients":"+9999999999","type":"signal","userId":1}`,
		},
		{
			name: "multiple recipients",
			data: []byte(`{"id":3,"name":"Signal Multi","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Signal Multi\",\"signalURL\":\"http://signal-api:9998\",\"signalNumber\":\"+5555555555\",\"signalRecipients\":\"+1111111111,+2222222222,+3333333333\",\"type\":\"signal\"}"}`),

			want: notification.Signal{
				Base: notification.Base{
					ID:            3,
					Name:          "Signal Multi",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SignalDetails: notification.SignalDetails{
					URL:        "http://signal-api:9998",
					Number:     "+5555555555",
					Recipients: "+1111111111,+2222222222,+3333333333",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Signal Multi","signalURL":"http://signal-api:9998","signalNumber":"+5555555555","signalRecipients":"+1111111111,+2222222222,+3333333333","type":"signal","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signal := notification.Signal{}

			err := json.Unmarshal(tc.data, &signal)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, signal)

			data, err := json.Marshal(signal)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
