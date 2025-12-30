package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationWebhook_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.Webhook
		wantJSON string
	}{
		{
			name: "json content type",
			data: []byte(
				`{"id":1,"name":"Webhook JSON","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":true,\"isDefault\":false,\"name\":\"Webhook JSON\",\"webhookURL\":\"https://example.com/webhook\",\"webhookContentType\":\"json\",\"type\":\"webhook\"}"}`,
			),

			want: notification.Webhook{
				Base: notification.Base{
					ID:            1,
					Name:          "Webhook JSON",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: true,
				},
				WebhookDetails: notification.WebhookDetails{
					WebhookURL:         "https://example.com/webhook",
					WebhookContentType: "json",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":false,"name":"Webhook JSON","type":"webhook","userId":1,"webhookContentType":"json","webhookURL":"https://example.com/webhook"}`,
		},
		{
			name: "form-data content type",
			data: []byte(
				`{"id":2,"name":"Webhook Form","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Webhook Form\",\"webhookURL\":\"https://example.com/form-handler\",\"webhookContentType\":\"form-data\",\"type\":\"webhook\"}"}`,
			),

			want: notification.Webhook{
				Base: notification.Base{
					ID:            2,
					Name:          "Webhook Form",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WebhookDetails: notification.WebhookDetails{
					WebhookURL:         "https://example.com/form-handler",
					WebhookContentType: "form-data",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Webhook Form","type":"webhook","userId":1,"webhookContentType":"form-data","webhookURL":"https://example.com/form-handler"}`,
		},
		{
			name: "custom content type with body",
			data: []byte(
				`{"id":3,"name":"Webhook Custom","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Webhook Custom\",\"webhookURL\":\"https://api.example.com/alerts\",\"webhookContentType\":\"custom\",\"webhookCustomBody\":\"{\\\"title\\\": \\\"Alert - {{ monitorJSON['name'] }}\\\", \\\"message\\\": \\\"{{ msg }}\\\"}\",\"type\":\"webhook\"}"}`,
			),

			want: notification.Webhook{
				Base: notification.Base{
					ID:            3,
					Name:          "Webhook Custom",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WebhookDetails: notification.WebhookDetails{
					WebhookURL:         "https://api.example.com/alerts",
					WebhookContentType: "custom",
					WebhookCustomBody:  `{"title": "Alert - {{ monitorJSON['name'] }}", "message": "{{ msg }}"}`,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":3,"isDefault":false,"name":"Webhook Custom","type":"webhook","userId":1,"webhookContentType":"custom","webhookCustomBody":"{\"title\": \"Alert - {{ monitorJSON['name'] }}\", \"message\": \"{{ msg }}\"}","webhookURL":"https://api.example.com/alerts"}`,
		},
		{
			name: "with additional headers",
			data: []byte(
				`{"id":4,"name":"Webhook Headers","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Webhook Headers\",\"webhookURL\":\"https://api.example.com/notify\",\"webhookContentType\":\"json\",\"webhookAdditionalHeaders\":\"{\\\"Authorization\\\":\\\"Bearer secret-token\\\",\\\"X-App-ID\\\":\\\"uptime-kuma\\\"}\",\"type\":\"webhook\"}"}`,
			),

			want: notification.Webhook{
				Base: notification.Base{
					ID:            4,
					Name:          "Webhook Headers",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				WebhookDetails: notification.WebhookDetails{
					WebhookURL:         "https://api.example.com/notify",
					WebhookContentType: "json",
					WebhookAdditionalHeaders: notification.WebhookAdditionalHeaders{
						"Authorization": "Bearer secret-token",
						"X-App-ID":      "uptime-kuma",
					},
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":4,"isDefault":false,"name":"Webhook Headers","type":"webhook","userId":1,"webhookContentType":"json","webhookAdditionalHeaders":"{\"Authorization\":\"Bearer secret-token\",\"X-App-ID\":\"uptime-kuma\"}","webhookURL":"https://api.example.com/notify"}`,
		},
		{
			name: "all features combined",
			data: []byte(
				`{"id":5,"name":"Webhook Full","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"Webhook Full\",\"webhookURL\":\"https://alerts.example.com/v1/webhook\",\"webhookContentType\":\"custom\",\"webhookCustomBody\":\"{\\\"event\\\": \\\"{{ msg }}\\\", \\\"service\\\": \\\"{{ monitorJSON['name'] }}\\\"}\",\"webhookAdditionalHeaders\":\"{\\\"Authorization\\\":\\\"Bearer xyz123\\\",\\\"Content-Type\\\":\\\"application/json\\\",\\\"X-Custom-Header\\\":\\\"custom-value\\\"}\",\"type\":\"webhook\"}"}`,
			),

			want: notification.Webhook{
				Base: notification.Base{
					ID:            5,
					Name:          "Webhook Full",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				WebhookDetails: notification.WebhookDetails{
					WebhookURL:         "https://alerts.example.com/v1/webhook",
					WebhookContentType: "custom",
					WebhookCustomBody:  `{"event": "{{ msg }}", "service": "{{ monitorJSON['name'] }}"}`,
					WebhookAdditionalHeaders: notification.WebhookAdditionalHeaders{
						"Authorization":   "Bearer xyz123",
						"Content-Type":    "application/json",
						"X-Custom-Header": "custom-value",
					},
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":5,"isDefault":true,"name":"Webhook Full","type":"webhook","userId":1,"webhookContentType":"custom","webhookCustomBody":"{\"event\": \"{{ msg }}\", \"service\": \"{{ monitorJSON['name'] }}\"}","webhookAdditionalHeaders":"{\"Authorization\":\"Bearer xyz123\",\"Content-Type\":\"application/json\",\"X-Custom-Header\":\"custom-value\"}","webhookURL":"https://alerts.example.com/v1/webhook"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			webhook := notification.Webhook{}

			err := json.Unmarshal(tc.data, &webhook)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, webhook)

			data, err := json.Marshal(webhook)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestNotificationWebhook_Type(t *testing.T) {
	webhook := notification.Webhook{}
	require.Equal(t, "webhook", webhook.Type())

	details := notification.WebhookDetails{}
	require.Equal(t, "webhook", details.Type())
}

func TestNotificationWebhook_String(t *testing.T) {
	webhook := notification.Webhook{
		Base: notification.Base{
			ID:            1,
			Name:          "Test Webhook",
			IsActive:      true,
			UserID:        1,
			IsDefault:     false,
			ApplyExisting: true,
		},
		WebhookDetails: notification.WebhookDetails{
			WebhookURL:         "https://example.com/webhook",
			WebhookContentType: "json",
		},
	}

	str := webhook.String()
	require.Contains(t, str, "Test Webhook")
	require.Contains(t, str, "webhook")
	require.Contains(t, str, "https://example.com/webhook")
	require.Contains(t, str, "json")
}

func TestWebhookAdditionalHeaders_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		headers  notification.WebhookAdditionalHeaders
		wantJSON string
	}{
		{
			name: "single header",
			headers: notification.WebhookAdditionalHeaders{
				"Authorization": "Bearer token",
			},
			wantJSON: `"{\"Authorization\":\"Bearer token\"}"`,
		},
		{
			name: "multiple headers",
			headers: notification.WebhookAdditionalHeaders{
				"Authorization": "Bearer token",
				"X-Custom":      "value",
			},
			wantJSON: `"{\"Authorization\":\"Bearer token\",\"X-Custom\":\"value\"}"`,
		},
		{
			name:     "empty map",
			headers:  notification.WebhookAdditionalHeaders{},
			wantJSON: `"{}"`,
		},
		{
			name:     "nil map",
			headers:  nil,
			wantJSON: `null`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tc.headers)
			require.NoError(t, err)
			require.JSONEq(t, tc.wantJSON, string(data))

			// Test unmarshaling
			var headers notification.WebhookAdditionalHeaders
			err = json.Unmarshal(data, &headers)
			require.NoError(t, err)

			if tc.headers == nil {
				require.Nil(t, headers)
			} else {
				require.Equal(t, tc.headers, headers)
			}
		})
	}
}

func TestWebhookAdditionalHeaders_UnmarshalEmptyString(t *testing.T) {
	var headers notification.WebhookAdditionalHeaders
	err := json.Unmarshal([]byte(`""`), &headers)
	require.NoError(t, err)
	require.Nil(t, headers)
}
