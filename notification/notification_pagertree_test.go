package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPagerTree_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.PagerTree
		wantJSON string
	}{
		{
			name: "success with auto-resolve",
			data: []byte(
				`{"id":1,"name":"Test PagerTree","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"PagerTree\",\"pagertreeIntegrationUrl\":\"https://api.pagertree.com/api/v2/events\",\"pagertreeUrgency\":\"high\",\"pagertreeAutoResolve\":\"resolve\"}"}`,
			),

			want: notification.PagerTree{
				Base: notification.Base{
					ID:        1,
					Name:      "Test PagerTree",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				PagerTreeDetails: notification.PagerTreeDetails{
					IntegrationURL: "https://api.pagertree.com/api/v2/events",
					Urgency:        "high",
					AutoResolve:    "resolve",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"Test PagerTree","pagertreeAutoResolve":"resolve","pagertreeIntegrationUrl":"https://api.pagertree.com/api/v2/events","pagertreeUrgency":"high","type":"PagerTree","userId":42}`,
		},
		{
			name: "success without auto-resolve",
			data: []byte(
				`{"id":2,"name":"Test PagerTree No Resolve","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"PagerTree\",\"pagertreeIntegrationUrl\":\"https://events.pagertree.com/api/v2/events\",\"pagertreeUrgency\":\"medium\",\"pagertreeAutoResolve\":\"\"}"}`,
			),

			want: notification.PagerTree{
				Base: notification.Base{
					ID:            2,
					Name:          "Test PagerTree No Resolve",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				PagerTreeDetails: notification.PagerTreeDetails{
					IntegrationURL: "https://events.pagertree.com/api/v2/events",
					Urgency:        "medium",
					AutoResolve:    "",
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":true,"name":"Test PagerTree No Resolve","pagertreeAutoResolve":"","pagertreeIntegrationUrl":"https://events.pagertree.com/api/v2/events","pagertreeUrgency":"medium","type":"PagerTree","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(
				`{"id":3,"name":"Test PagerTree Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"PagerTree\",\"pagertreeIntegrationUrl\":\"\",\"pagertreeUrgency\":\"\",\"pagertreeAutoResolve\":\"\"}"}`,
			),

			want: notification.PagerTree{
				Base: notification.Base{
					ID:            3,
					Name:          "Test PagerTree Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PagerTreeDetails: notification.PagerTreeDetails{
					IntegrationURL: "",
					Urgency:        "",
					AutoResolve:    "",
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Test PagerTree Minimal","pagertreeAutoResolve":"","pagertreeIntegrationUrl":"","pagertreeUrgency":"","type":"PagerTree","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pagertree := notification.PagerTree{}

			err := json.Unmarshal(tc.data, &pagertree)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pagertree)

			data, err := json.Marshal(pagertree)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
