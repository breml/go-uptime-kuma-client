package statuspage_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/statuspage"
)

func TestIncident_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		incident statuspage.Incident
	}{
		{
			name: "complete incident",
			incident: statuspage.Incident{
				ID:      1,
				Title:   "Service Degradation",
				Content: "We are experiencing issues with our API service.",
				Style:   "warning",
				Pin:     true,
			},
		},
		{
			name: "minimal incident without ID",
			incident: statuspage.Incident{
				Title:   "Maintenance Window",
				Content: "Scheduled maintenance in progress.",
				Style:   "info",
				Pin:     false,
			},
		},
		{
			name: "danger incident",
			incident: statuspage.Incident{
				ID:      2,
				Title:   "Major Outage",
				Content: "All services are currently down.",
				Style:   "danger",
				Pin:     true,
			},
		},
		{
			name: "primary incident",
			incident: statuspage.Incident{
				Title:   "New Feature Release",
				Content: "We've just released a new feature!",
				Style:   "primary",
				Pin:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.incident)
			require.NoError(t, err)

			var got statuspage.Incident
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)

			require.Equal(t, tt.incident, got)
		})
	}
}

func TestIncident_OmitEmptyID(t *testing.T) {
	incident := statuspage.Incident{
		Title:   "Test",
		Content: "Test content",
		Style:   "info",
		Pin:     false,
	}

	data, err := json.Marshal(incident)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// ID should not be present when it's 0
	_, hasID := result["id"]
	require.False(t, hasID, "ID field should be omitted when zero")
}
