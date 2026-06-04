package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWebpush_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Webpush
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Webpush Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Webpush Alert\",\"subscription\":{\"endpoint\":\"https://push.example.com/abc123\",\"keys\":{\"p256dh\":\"BGxi5eHcCn...\",\"auth\":\"abc\"}},\"type\":\"Webpush\"}"}`,
			),

			want: notification.Webpush{
				Base: notification.Base{
					ID:            1,
					Name:          "My Webpush Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WebpushDetails: notification.WebpushDetails{
					Subscription: notification.WebpushSubscription{
						Endpoint: "https://push.example.com/abc123",
						Keys: notification.WebpushSubscriptionKeys{
							P256dh: "BGxi5eHcCn...",
							Auth:   "abc",
						},
					},
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Webpush Alert","subscription":{"endpoint":"https://push.example.com/abc123","keys":{"p256dh":"BGxi5eHcCn...","auth":"abc"}},"type":"Webpush","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Webpush","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Webpush\",\"subscription\":{\"endpoint\":\"https://fcm.googleapis.com/fcm/send/xyz\",\"keys\":{\"p256dh\":\"BNc...\",\"auth\":\"xyz\"}},\"type\":\"Webpush\"}"}`,
			),

			want: notification.Webpush{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Webpush",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WebpushDetails: notification.WebpushDetails{
					Subscription: notification.WebpushSubscription{
						Endpoint: "https://fcm.googleapis.com/fcm/send/xyz",
						Keys: notification.WebpushSubscriptionKeys{
							P256dh: "BNc...",
							Auth:   "xyz",
						},
					},
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Webpush","subscription":{"endpoint":"https://fcm.googleapis.com/fcm/send/xyz","keys":{"p256dh":"BNc...","auth":"xyz"}},"type":"Webpush","userId":1}`,
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
			webpush := notification.Webpush{}

			err := json.Unmarshal(tc.data, &webpush)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, webpush)

			data, err := json.Marshal(webpush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestNotificationWebpush_Type(t *testing.T) {
	webpush := notification.Webpush{}
	require.Equal(t, "Webpush", webpush.Type())

	details := notification.WebpushDetails{}
	require.Equal(t, "Webpush", details.Type())
}

func TestNotificationWebpush_String(t *testing.T) {
	webpush := notification.Webpush{
		Base: notification.Base{
			ID:       1,
			Name:     "Test Webpush",
			IsActive: true,
			UserID:   1,
		},
		WebpushDetails: notification.WebpushDetails{
			Subscription: notification.WebpushSubscription{
				Endpoint: "https://push.example.com/abc123",
				Keys: notification.WebpushSubscriptionKeys{
					P256dh: "BGxi5eHcCn...",
					Auth:   "abc",
				},
			},
		},
	}

	str := webpush.String()
	require.Contains(t, str, "Test Webpush")
	require.Contains(t, str, "Webpush")
	require.Contains(t, str, "https://push.example.com/abc123")
	require.NotContains(t, str, "0x")
}
