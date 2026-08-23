package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFlowtriq_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Flowtriq
		wantJSON string
		wantErr  bool
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Flowtriq Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"flowtriqApiKey\":\"test_api_key\",\"flowtriqWebhookUrl\":\"https://app.flowtriq.com/api/webhooks/test\",\"isDefault\":true,\"name\":\"My Flowtriq Alert\",\"type\":\"Flowtriq\"}"}`,
			),

			want: notification.Flowtriq{
				Base: notification.Base{
					ID:            1,
					Name:          "My Flowtriq Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FlowtriqDetails: notification.FlowtriqDetails{
					WebhookURL: "https://app.flowtriq.com/api/webhooks/test",
					APIKey:     ptr.To("test_api_key"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"flowtriqApiKey":"test_api_key","flowtriqWebhookUrl":"https://app.flowtriq.com/api/webhooks/test","id":1,"isDefault":true,"name":"My Flowtriq Alert","type":"Flowtriq","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Flowtriq","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"flowtriqWebhookUrl\":\"https://app.flowtriq.com/api/webhooks/test\",\"isDefault\":false,\"name\":\"Simple Flowtriq\",\"type\":\"Flowtriq\"}"}`,
			),

			want: notification.Flowtriq{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Flowtriq",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FlowtriqDetails: notification.FlowtriqDetails{
					WebhookURL: "https://app.flowtriq.com/api/webhooks/test",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"flowtriqWebhookUrl":"https://app.flowtriq.com/api/webhooks/test","id":2,"isDefault":false,"name":"Simple Flowtriq","type":"Flowtriq","userId":1}`,
		},
		{
			name: "empty api key",
			data: []byte(
				`{"id":3,"name":"Cleared Flowtriq","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"flowtriqApiKey\":\"\",\"flowtriqWebhookUrl\":\"https://app.flowtriq.com/api/webhooks/test\",\"isDefault\":false,\"name\":\"Cleared Flowtriq\",\"type\":\"Flowtriq\"}"}`,
			),

			want: notification.Flowtriq{
				Base: notification.Base{
					ID:            3,
					Name:          "Cleared Flowtriq",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FlowtriqDetails: notification.FlowtriqDetails{
					WebhookURL: "https://app.flowtriq.com/api/webhooks/test",
					APIKey:     ptr.To(""),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"flowtriqApiKey":"","flowtriqWebhookUrl":"https://app.flowtriq.com/api/webhooks/test","id":3,"isDefault":false,"name":"Cleared Flowtriq","type":"Flowtriq","userId":1}`,
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
				`{"id":1,"name":"x","active":true,"userId":1,"isDefault":false,"config":"{\"flowtriqWebhookUrl\":123,\"type\":\"Flowtriq\"}"}`,
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flowtriq := notification.Flowtriq{}

			err := json.Unmarshal(tc.data, &flowtriq)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, flowtriq)

			data, err := json.Marshal(flowtriq)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
