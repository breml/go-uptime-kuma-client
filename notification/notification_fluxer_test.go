package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationFluxer_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Fluxer
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":1,"name":"My Fluxer Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Fluxer Alert\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz\",\"fluxerUsername\":\"Uptime Monitor\",\"fluxerPrefixMessage\":\"Alert:\",\"disableUrl\":true,\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            1,
					Name:          "My Fluxer Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL:    "https://fluxer.app/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz",
					Username:      "Uptime Monitor",
					PrefixMessage: "Alert:",
					DisableURL:    true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"disableUrl":true,"fluxerPrefixMessage":"Alert:","fluxerUsername":"Uptime Monitor","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz","id":1,"isDefault":true,"name":"My Fluxer Alert","type":"fluxer","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":2,"name":"Simple Fluxer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Fluxer\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/xyz/abc\",\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Fluxer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL: "https://fluxer.app/api/webhooks/xyz/abc",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"fluxerPrefixMessage":"","fluxerUsername":"","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/xyz/abc","id":2,"isDefault":false,"name":"Simple Fluxer","type":"fluxer","userId":1}`,
		},
		{
			name: "with minimalist message format",
			data: []byte(
				`{"id":3,"name":"Minimalist Fluxer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Minimalist Fluxer\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/min/webhook\",\"fluxerMessageFormat\":\"minimalist\",\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            3,
					Name:          "Minimalist Fluxer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL:    "https://fluxer.app/api/webhooks/min/webhook",
					MessageFormat: ptr.To("minimalist"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"fluxerMessageFormat":"minimalist","fluxerPrefixMessage":"","fluxerUsername":"","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/min/webhook","id":3,"isDefault":false,"name":"Minimalist Fluxer","type":"fluxer","userId":1}`,
		},
		{
			name: "with custom message format and template",
			data: []byte(
				`{"id":4,"name":"Custom Fluxer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Custom Fluxer\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/custom/webhook\",\"fluxerMessageFormat\":\"custom\",\"fluxerMessageTemplate\":\"Service {{name}} is {{status}}\",\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            4,
					Name:          "Custom Fluxer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL:      "https://fluxer.app/api/webhooks/custom/webhook",
					MessageFormat:   ptr.To("custom"),
					MessageTemplate: ptr.To("Service {{name}} is {{status}}"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"fluxerMessageFormat":"custom","fluxerMessageTemplate":"Service {{name}} is {{status}}","fluxerPrefixMessage":"","fluxerUsername":"","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/custom/webhook","id":4,"isDefault":false,"name":"Custom Fluxer","type":"fluxer","userId":1}`,
		},
		{
			name: "with legacy use message template",
			data: []byte(
				`{"id":5,"name":"Legacy Template Fluxer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Legacy Template Fluxer\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/legacy/webhook\",\"fluxerUseMessageTemplate\":true,\"fluxerMessageTemplate\":\"Alert: {{name}}\",\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            5,
					Name:          "Legacy Template Fluxer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL:         "https://fluxer.app/api/webhooks/legacy/webhook",
					UseMessageTemplate: ptr.To(true),
					MessageTemplate:    ptr.To("Alert: {{name}}"),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"fluxerMessageTemplate":"Alert: {{name}}","fluxerPrefixMessage":"","fluxerUseMessageTemplate":true,"fluxerUsername":"","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/legacy/webhook","id":5,"isDefault":false,"name":"Legacy Template Fluxer","type":"fluxer","userId":1}`,
		},
		{
			name: "pointer to false and empty string are serialized",
			data: []byte(
				`{"id":6,"name":"Explicit False Fluxer","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Explicit False Fluxer\",\"fluxerWebhookUrl\":\"https://fluxer.app/api/webhooks/explicit/webhook\",\"fluxerUseMessageTemplate\":false,\"fluxerMessageFormat\":\"\",\"type\":\"fluxer\"}"}`,
			),

			want: notification.Fluxer{
				Base: notification.Base{
					ID:            6,
					Name:          "Explicit False Fluxer",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				FluxerDetails: notification.FluxerDetails{
					WebhookURL:         "https://fluxer.app/api/webhooks/explicit/webhook",
					UseMessageTemplate: ptr.To(false),
					MessageFormat:      ptr.To(""),
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"disableUrl":false,"fluxerMessageFormat":"","fluxerPrefixMessage":"","fluxerUseMessageTemplate":false,"fluxerUsername":"","fluxerWebhookUrl":"https://fluxer.app/api/webhooks/explicit/webhook","id":6,"isDefault":false,"name":"Explicit False Fluxer","type":"fluxer","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fluxer := notification.Fluxer{}

			err := json.Unmarshal(tc.data, &fluxer)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, fluxer)

			data, err := json.Marshal(fluxer)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestNotificationFluxer_String(t *testing.T) {
	fluxer := notification.Fluxer{
		Base: notification.Base{
			ID:       1,
			Name:     "Test Fluxer",
			IsActive: true,
			UserID:   1,
		},
		FluxerDetails: notification.FluxerDetails{
			WebhookURL:         "https://fluxer.app/api/webhooks/123/abc",
			UseMessageTemplate: ptr.To(true),
			MessageFormat:      ptr.To("minimalist"),
			MessageTemplate:    ptr.To(""),
		},
	}

	str := fluxer.String()

	require.Contains(t, str, "Test Fluxer")
	require.Contains(t, str, "fluxer")
	require.Contains(t, str, "https://fluxer.app/api/webhooks/123/abc")
	require.Contains(t, str, "true")
	require.Contains(t, str, `"minimalist"`)
	require.NotContains(t, str, "0x")
}
