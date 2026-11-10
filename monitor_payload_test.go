package kuma_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

// TestAdaptMonitorPayload covers the shape the name of a pm2 or
// system-service monitor takes on the wire: top-level system_service_name for
// a 2.x server, config.systemServiceName for a 3.x server.
func TestAdaptMonitorPayload(t *testing.T) {
	tests := []struct {
		name   string
		scheme kuma.AuthScheme
		mon    any
		want   map[string]any
	}{
		{
			name:   "2.x keeps system_service_name",
			scheme: kuma.AuthSchemeSocket,
			mon: monitor.SystemService{
				SystemServiceDetails: monitor.SystemServiceDetails{SystemServiceName: "nginx.service"},
			},
			want: map[string]any{"system_service_name": "nginx.service"},
		},
		{
			name:   "3.x moves the system-service name into config",
			scheme: kuma.AuthSchemeHTTP,
			mon: monitor.SystemService{
				SystemServiceDetails: monitor.SystemServiceDetails{SystemServiceName: "nginx.service"},
			},
			want: map[string]any{"config": map[string]any{"systemServiceName": "nginx.service"}},
		},
		{
			name:   "3.x moves the pm2 process name into config",
			scheme: kuma.AuthSchemeHTTP,
			mon:    monitor.PM2{PM2Details: monitor.PM2Details{ProcessName: "worker 1"}},
			want:   map[string]any{"config": map[string]any{"systemServiceName": "worker 1"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.mon)
			require.NoError(t, err)

			payload := map[string]any{}
			err = json.Unmarshal(data, &payload)
			require.NoError(t, err)

			kuma.AdaptMonitorPayload(tc.scheme, payload)

			for key, value := range tc.want {
				require.Equal(t, value, payload[key], key)
			}

			if tc.scheme == kuma.AuthSchemeHTTP {
				require.NotContains(t, payload, "system_service_name")
			} else {
				require.NotContains(t, payload, "config")
			}
		})
	}

	t.Run("3.x merges into an existing config", func(t *testing.T) {
		payload := map[string]any{
			"system_service_name": "nginx.service",
			"config":              map[string]any{"other": true},
		}

		kuma.AdaptMonitorPayload(kuma.AuthSchemeHTTP, payload)

		require.Equal(t, map[string]any{
			"config": map[string]any{"other": true, "systemServiceName": "nginx.service"},
		}, payload)
	})

	t.Run("3.x drops a null name", func(t *testing.T) {
		payload := map[string]any{"system_service_name": nil}

		kuma.AdaptMonitorPayload(kuma.AuthSchemeHTTP, payload)

		require.Empty(t, payload)
	})

	t.Run("3.x leaves other monitor types alone", func(t *testing.T) {
		payload := map[string]any{"type": "http", "url": "https://example.com"}

		kuma.AdaptMonitorPayload(kuma.AuthSchemeHTTP, payload)

		require.Equal(t, map[string]any{"type": "http", "url": "https://example.com"}, payload)
	})
}
