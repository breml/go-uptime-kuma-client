package kuma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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

// recordMonitorWrites makes the fake answer the add and editMonitor commands
// and record their monitor payloads, see answerMonitorWrite.
func (s *fakeAckServer) recordMonitorWrites() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.monitorWrites = map[string][]map[string]any{}
}

// recordedMonitorWrites returns the monitor payloads the command carried, in
// order.
func (s *fakeAckServer) recordedMonitorWrites(command string) []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.monitorWrites[command]
}

// answerMonitorWrite answers an add or editMonitor command the way Uptime
// Kuma does, with the updateMonitorIntoList broadcast first and the ack with
// the monitor ID after it, and records the monitor payload. It reports whether
// it answered the command; it answers none unless recordMonitorWrites ran.
func (s *fakeAckServer) answerMonitorWrite(session *fakeSession, ackID string, payload string) bool {
	var frame []json.RawMessage

	if json.Unmarshal([]byte(payload), &frame) != nil || len(frame) < 2 {
		return false
	}

	var command string

	if json.Unmarshal(frame[0], &command) != nil || (command != "add" && command != "editMonitor") {
		return false
	}

	var monitorData map[string]any

	if json.Unmarshal(frame[1], &monitorData) != nil {
		return false
	}

	s.mu.Lock()
	if s.monitorWrites == nil {
		s.mu.Unlock()

		return false
	}

	s.monitorWrites[command] = append(s.monitorWrites[command], monitorData)
	s.mu.Unlock()

	const id = fakeAckCreatedID

	s.sendTo(
		session,
		fmt.Appendf(nil, `42["updateMonitorIntoList",{"%d":{"id":%d,"name":"written","type":"pm2"}}]`, id, id),
	)
	s.sendTo(session, fmt.Appendf(nil, `43%s[{"ok":true,"msg":"ok","monitorID":%d}]`, ackID, id))

	return true
}

// TestMonitorWritesSendTheServerLinesShape covers what CreateMonitor and
// UpdateMonitor put on the wire for a pm2 monitor, per server line: the name
// as top-level system_service_name for 2.x, as config.systemServiceName for
// 3.x. TestAdaptMonitorPayload covers the rewrite itself; this one fails if a
// write path stops applying it.
func TestMonitorWritesSendTheServerLinesShape(t *testing.T) {
	tests := []struct {
		name   string
		fake   func(t *testing.T) (*fakeAckServer, string)
		scheme kuma.AuthScheme
		want   func(t *testing.T, payload map[string]any, processName string)
	}{
		{
			name:   "2.x",
			fake:   newAuthFake,
			scheme: kuma.AuthSchemeSocket,
			want: func(t *testing.T, payload map[string]any, processName string) {
				t.Helper()

				require.Equal(t, processName, payload["system_service_name"])
				require.NotContains(t, payload, "config")
			},
		},
		{
			name:   "3.x",
			fake:   newHTTPAuthFake,
			scheme: kuma.AuthSchemeHTTP,
			want: func(t *testing.T, payload map[string]any, processName string) {
				t.Helper()

				require.Equal(t, map[string]any{"systemServiceName": processName}, payload["config"])
				require.NotContains(t, payload, "system_service_name")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake, url := tc.fake(t)
			fake.recordMonitorWrites()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			kumaClient, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
			require.NoError(t, err)

			t.Cleanup(func() { _ = kumaClient.Disconnect() })

			require.Equal(t, tc.scheme, kumaClient.AuthScheme())

			pm2 := &monitor.PM2{
				Base:       monitor.Base{Name: "pm2", Interval: 60, RetryInterval: 60, IsActive: true},
				PM2Details: monitor.PM2Details{ProcessName: "worker 1"},
			}

			pm2.ID, err = kumaClient.CreateMonitor(ctx, pm2)
			require.NoError(t, err)

			pm2.ProcessName = "worker 2"

			err = kumaClient.UpdateMonitor(ctx, pm2)
			require.NoError(t, err)

			adds := fake.recordedMonitorWrites("add")
			require.Len(t, adds, 1)
			tc.want(t, adds[0], "worker 1")

			edits := fake.recordedMonitorWrites("editMonitor")
			require.Len(t, edits, 1)
			tc.want(t, edits[0], "worker 2")
		})
	}
}
