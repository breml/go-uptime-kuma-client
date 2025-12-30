package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationSpugPush_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.SpugPush
		wantJSON string
	}{
		{
			name: "success with template key",
			data: []byte(
				`{"id":1,"name":"My SpugPush Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My SpugPush Alert\",\"templateKey\":\"test-template-key\",\"type\":\"SpugPush\"}"}`,
			),

			want: notification.SpugPush{
				Base: notification.Base{
					ID:            1,
					Name:          "My SpugPush Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				SpugPushDetails: notification.SpugPushDetails{
					TemplateKey: "test-template-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My SpugPush Alert","templateKey":"test-template-key","type":"SpugPush","userId":1}`,
		},
		{
			name: "minimal configuration",
			data: []byte(
				`{"id":2,"name":"Simple SpugPush","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple SpugPush\",\"templateKey\":\"simple-key\",\"type\":\"SpugPush\"}"}`,
			),

			want: notification.SpugPush{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple SpugPush",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SpugPushDetails: notification.SpugPushDetails{
					TemplateKey: "simple-key",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple SpugPush","templateKey":"simple-key","type":"SpugPush","userId":1}`,
		},
		{
			name: "with complex template key",
			data: []byte(
				`{"id":3,"name":"SpugPush Complex","active":false,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"SpugPush Complex\",\"templateKey\":\"complex-template-abc123-xyz\",\"type\":\"SpugPush\"}"}`,
			),

			want: notification.SpugPush{
				Base: notification.Base{
					ID:            3,
					Name:          "SpugPush Complex",
					IsActive:      false,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				SpugPushDetails: notification.SpugPushDetails{
					TemplateKey: "complex-template-abc123-xyz",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"SpugPush Complex","templateKey":"complex-template-abc123-xyz","type":"SpugPush","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spugpush := notification.SpugPush{}

			err := json.Unmarshal(tc.data, &spugpush)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, spugpush)

			data, err := json.Marshal(spugpush)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
