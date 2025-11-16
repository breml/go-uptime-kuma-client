package statuspage_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/statuspage"
)

func TestStatusPage_MarshalUnmarshal(t *testing.T) {
	sendURLTrue := true
	sendURLFalse := false

	tests := []struct {
		name string
		sp   statuspage.StatusPage
	}{
		{
			name: "complete status page",
			sp: statuspage.StatusPage{
				ID:                    1,
				Slug:                  "test-page",
				Title:                 "Test Status Page",
				Description:           "This is a test status page",
				Icon:                  "/icon.svg",
				Theme:                 "light",
				Published:             true,
				ShowTags:              true,
				DomainNameList:        []string{"status.example.com"},
				GoogleAnalyticsID:     "UA-123456-1",
				CustomCSS:             "body { background: #fff; }",
				FooterText:            "© 2024 Example Inc.",
				ShowPoweredBy:         false,
				ShowCertificateExpiry: true,
				PublicGroupList: []statuspage.PublicGroup{
					{
						ID:     10,
						Name:   "Web Services",
						Weight: 1,
						MonitorList: []statuspage.PublicMonitor{
							{ID: 100, SendURL: &sendURLTrue},
							{ID: 101, SendURL: &sendURLFalse},
						},
					},
					{
						ID:     20,
						Name:   "Databases",
						Weight: 2,
						MonitorList: []statuspage.PublicMonitor{
							{ID: 200, SendURL: nil},
						},
					},
				},
			},
		},
		{
			name: "minimal status page",
			sp: statuspage.StatusPage{
				Slug:  "minimal",
				Title: "Minimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.sp)
			require.NoError(t, err)

			var got statuspage.StatusPage
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)

			require.Equal(t, tt.sp, got)
		})
	}
}

func TestPublicGroup_MarshalUnmarshal(t *testing.T) {
	sendURL := true

	tests := []struct {
		name  string
		group statuspage.PublicGroup
	}{
		{
			name: "group with monitors",
			group: statuspage.PublicGroup{
				ID:     1,
				Name:   "Test Group",
				Weight: 5,
				MonitorList: []statuspage.PublicMonitor{
					{ID: 10, SendURL: &sendURL},
					{ID: 20, SendURL: nil},
				},
			},
		},
		{
			name: "empty group",
			group: statuspage.PublicGroup{
				Name: "Empty Group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.group)
			require.NoError(t, err)

			var got statuspage.PublicGroup
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)

			require.Equal(t, tt.group, got)
		})
	}
}

func TestPublicMonitor_MarshalUnmarshal(t *testing.T) {
	sendURLTrue := true
	sendURLFalse := false

	tests := []struct {
		name    string
		monitor statuspage.PublicMonitor
	}{
		{
			name:    "monitor with sendUrl true",
			monitor: statuspage.PublicMonitor{ID: 1, SendURL: &sendURLTrue},
		},
		{
			name:    "monitor with sendUrl false",
			monitor: statuspage.PublicMonitor{ID: 2, SendURL: &sendURLFalse},
		},
		{
			name:    "monitor without sendUrl",
			monitor: statuspage.PublicMonitor{ID: 3, SendURL: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.monitor)
			require.NoError(t, err)

			var got statuspage.PublicMonitor
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)

			require.Equal(t, tt.monitor, got)
		})
	}
}
