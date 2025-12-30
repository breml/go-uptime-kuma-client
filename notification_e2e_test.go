package kuma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// TestEndToEndMonitorFailureNotification tests the complete notification flow:
// 1. Create a web server with a monitored endpoint and a webhook receiver.
// 2. Configure Uptime Kuma with a monitor and webhook notification.
// 3. Monitor performs first check successfully (endpoint returns 200, UP state).
// 4. Test waits for first successful check signal via channel.
// 5. Endpoint fails on second check (returns 503, triggers UP -> DOWN transition).
// 6. Uptime Kuma sends failure notification to the webhook.
// 7. Test verifies failure webhook was called with status=0.
// 8. Endpoint recovers on third check (returns 200, triggers DOWN -> UP transition).
// 9. Uptime Kuma sends recovery notification to the webhook.
// 10. Test verifies recovery webhook was called with status=1.
func TestEndToEndMonitorFailureNotification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if !e2eTest {
		t.Skip(`skipping end to end test, "E2E_TEST" env var not set`)
	}

	// Use a 60-second timeout to fit within the 120-second container expiration.
	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	// Track webhook calls.
	var webhookMu sync.Mutex
	var webhookCalls []WebhookPayload
	failureNotificationReceived := make(chan struct{}, 1)
	recoveryNotificationReceived := make(chan struct{}, 1)

	// Track check count and control endpoint behavior.
	var endpointMu sync.Mutex
	checkCount := 0
	firstCheckReceived := make(chan struct{}, 1)

	// Create test HTTP server with two handlers.
	mux := http.NewServeMux()

	// Handler 1: Monitored endpoint with state sequence: UP -> DOWN -> UP.
	// Check 1: Success (200) - establish UP state.
	// Check 2: Failure (503) - trigger DOWN notification.
	// Check 3+: Success (200) - trigger recovery notification.
	mux.HandleFunc("/monitored-endpoint", func(w http.ResponseWriter, _ *http.Request) {
		endpointMu.Lock()
		checkCount++
		currentCheck := checkCount
		endpointMu.Unlock()

		switch currentCheck {
		case 1:
			// First check: return success.
			t.Log("Monitor check #1: Returning 200 (success - establishing UP state)")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))

			// Signal that first check completed.
			select {
			case firstCheckReceived <- struct{}{}:
			default:
			}

			return

		case 2:
			// Second check: return failure to trigger DOWN notification.
			t.Log("Monitor check #2: Returning 503 (failure - triggering DOWN notification)")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Service Unavailable"))
			return

		default:
			// Third+ checks: return success to trigger recovery notification.
			t.Logf("Monitor check #%d: Returning 200 (success - triggering recovery notification)", currentCheck)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
			return
		}
	})

	// Handler 2: Webhook notification receiver.
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Webhook called! Method: %s, Content-Type: %s", r.Method, r.Header.Get("Content-Type"))

		var payload WebhookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Logf("Failed to decode webhook payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Logf("Webhook payload received: status=%d, msg=%s", payload.Heartbeat.Status, payload.Msg)

		webhookMu.Lock()
		webhookCalls = append(webhookCalls, payload)
		webhookMu.Unlock()

		// Signal based on notification type (failure vs recovery)
		switch payload.Heartbeat.Status {
		case 0:
			// Status 0 = DOWN (failure notification)
			t.Log("Failure notification received (status=0)")
			select {
			case failureNotificationReceived <- struct{}{}:
			default:
			}

		case 1:
			// Status 1 = UP (recovery notification)
			t.Log("Recovery notification received (status=1)")
			select {
			case recoveryNotificationReceived <- struct{}{}:
			default:
			}

		default:
			// Ignore unknown status codes
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Create server that listens on all interfaces.
	// This is necessary for the Docker container to reach it.
	server := httptest.NewUnstartedServer(mux)
	var err error
	server.Listener, err = net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err)
	server.Start()
	defer server.Close()

	t.Logf("Test server started at %s", server.URL)

	// Extract port from server URL and create URL accessible from Docker container.
	// On Linux, the Docker container can reach the host via the bridge gateway IP.
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// Try to detect Docker bridge gateway IP.
	// Common values: 172.17.0.1 (default), 192.168.x.1 (custom), or host.docker.internal.
	gatewayIP := "172.17.0.1" // Default Docker bridge gateway.

	// Try to get the actual gateway IP.
	gw := getDockerGatewayIP(t)
	if gw != "" {
		gatewayIP = gw
	}

	dockerAccessibleURL := "http://" + net.JoinHostPort(gatewayIP, strconv.Itoa(serverPort))
	t.Logf("Docker-accessible URL: %s", dockerAccessibleURL)

	// Create webhook notification.
	webhookNotification := notification.Webhook{
		Base: notification.Base{
			ApplyExisting: true,
			IsDefault:     true,
			IsActive:      true,
			Name:          "E2E Test Webhook",
		},
		WebhookDetails: notification.WebhookDetails{
			WebhookURL:         fmt.Sprintf("%s/webhook", dockerAccessibleURL),
			WebhookContentType: "json",
		},
	}

	notificationID, err := client.CreateNotification(ctx, webhookNotification)
	require.NoError(t, err)
	require.Positive(t, notificationID)
	t.Logf("Created notification with ID: %d", notificationID)

	// Cleanup notification
	t.Cleanup(func() {
		//nolint:usetesting // t.Context() cannot be used in cleanup, since it is already done.
		err := client.DeleteNotification(context.Background(), notificationID)
		if err != nil {
			t.Logf("Failed to delete notification: %v", err)
		}
	})

	// Create HTTP monitor that checks the monitored endpoint.
	// Endpoint returns success on first check, then automatically fails.
	// Use a short interval (20 seconds minimum) to speed up the test.
	httpMonitor := monitor.HTTP{
		Base: monitor.Base{
			Name:            "E2E Test Monitor",
			Interval:        20, // Check every 20 seconds (minimum allowed), see https://github.com/louislam/uptime-kuma/issues/1645 .
			RetryInterval:   20,
			ResendInterval:  0,
			MaxRetries:      0, // Send notification on first failure (no retries).
			UpsideDown:      false,
			IsActive:        true,
			NotificationIDs: []int64{notificationID},
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 fmt.Sprintf("%s/monitored-endpoint", dockerAccessibleURL),
			Timeout:             5,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"}, // Expect 2xx, but will get 503 on failure.
			AuthMethod:          monitor.AuthMethodNone,
		},
	}

	monitorID, err := client.CreateMonitor(ctx, &httpMonitor)
	require.NoError(t, err)
	require.Positive(t, monitorID)
	t.Logf("Created monitor with ID: %d", monitorID)

	// Verify the monitor was created with the correct configuration.
	var retrievedMonitor monitor.HTTP
	err = client.GetMonitorAs(ctx, monitorID, &retrievedMonitor)
	require.NoError(t, err)
	t.Logf("Monitor configured - URL: %s, NotificationIDs: %v, IsActive: %v",
		retrievedMonitor.URL, retrievedMonitor.NotificationIDs, retrievedMonitor.IsActive)

	// Cleanup monitor.
	t.Cleanup(func() {
		//nolint:usetesting // t.Context() cannot be used in cleanup, since it is already done.
		err := client.DeleteMonitor(context.Background(), monitorID)
		if err != nil {
			t.Logf("Failed to delete monitor: %v", err)
		}
	})

	// Wait for the monitor to perform its first successful check.
	// This establishes the "UP" state.
	// The handler will automatically fail on subsequent checks.
	t.Log("Waiting for initial successful check...")

	select {
	case <-firstCheckReceived:
		t.Log("First successful check received - monitor is UP")

	case <-time.After(30 * time.Second):
		t.Fatal("Timeout waiting for first successful check")

	case <-ctx.Done():
		t.Fatal("Context cancelled while waiting for first check")
	}

	// Phase 1: Wait for failure notification.
	// The endpoint will fail on check #2 (which happens automatically).
	// Monitor checks every 20s, so next check should happen within 20s.
	// Adding 5s buffer for notification delivery.
	t.Log("Waiting for failure notification (monitor will detect failure on check #2)...")

	select {
	case <-failureNotificationReceived:
		t.Log("Failure notification received!")

	case <-time.After(25 * time.Second):
		t.Fatal("Timeout waiting for failure notification")

	case <-ctx.Done():
		t.Fatal("Context cancelled while waiting for failure notification")
	}

	// Verify failure notification.
	webhookMu.Lock()
	require.NotEmpty(t, webhookCalls, "Expected at least one webhook call")
	failurePayload := webhookCalls[len(webhookCalls)-1] // Get most recent
	webhookMu.Unlock()

	t.Logf("Failure payload: status=%d, msg=%s", failurePayload.Heartbeat.Status, failurePayload.Msg)
	require.NotEmpty(t, failurePayload.Heartbeat, "Expected heartbeat information in failure webhook")
	require.NotEmpty(t, failurePayload.Monitor, "Expected monitor information in failure webhook")
	require.Equal(t, monitorID, failurePayload.Monitor.ID, "Failure webhook should be about our monitor")
	require.Equal(t, "E2E Test Monitor", failurePayload.Monitor.Name, "Monitor name should match")
	require.Equal(t, 0, failurePayload.Heartbeat.Status, "Monitor should be down/failed (status=0)")

	// Phase 2: Wait for recovery notification.
	// The endpoint will recover on check #3 (which happens automatically).
	// Monitor checks every 20s, so next check should happen within 20s.
	// Adding 5s buffer for notification delivery.
	t.Log("Waiting for recovery notification (monitor will detect recovery on check #3)...")

	select {
	case <-recoveryNotificationReceived:
		t.Log("Recovery notification received!")

	case <-time.After(25 * time.Second):
		t.Fatal("Timeout waiting for recovery notification")

	case <-ctx.Done():
		t.Fatal("Context cancelled while waiting for recovery notification")
	}

	// Verify recovery notification.
	webhookMu.Lock()
	require.GreaterOrEqual(t, len(webhookCalls), 2, "Expected at least two webhook calls (failure + recovery)")
	recoveryPayload := webhookCalls[len(webhookCalls)-1] // Get most recent.
	webhookMu.Unlock()

	t.Logf("Recovery payload: status=%d, msg=%s", recoveryPayload.Heartbeat.Status, recoveryPayload.Msg)
	require.NotEmpty(t, recoveryPayload.Heartbeat, "Expected heartbeat information in recovery webhook")
	require.NotEmpty(t, recoveryPayload.Monitor, "Expected monitor information in recovery webhook")
	require.Equal(t, monitorID, recoveryPayload.Monitor.ID, "Recovery webhook should be about our monitor")
	require.Equal(t, "E2E Test Monitor", recoveryPayload.Monitor.Name, "Monitor name should match")
	require.Equal(t, 1, recoveryPayload.Heartbeat.Status, "Monitor should be up/recovered (status=1)")

	t.Log("End-to-end notification test completed successfully!")
	t.Log("✓ Verified failure notification (UP → DOWN)")
	t.Log("✓ Verified recovery notification (DOWN → UP)")
}

// WebhookPayload represents the structure of Uptime Kuma's webhook notification.
type WebhookPayload struct {
	Heartbeat HeartbeatInfo `json:"heartbeat"`
	Monitor   MonitorInfo   `json:"monitor"`
	Msg       string        `json:"msg"`
}

type HeartbeatInfo struct {
	MonitorID int64  `json:"monitorID"`
	Status    int    `json:"status"` // 0 = down, 1 = up
	Time      string `json:"time"`
	Msg       string `json:"msg"`
	Important bool   `json:"important"`
	Duration  int64  `json:"duration"`
}

type MonitorInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
}

// getDockerGatewayIP attempts to detect the Docker bridge gateway IP.
func getDockerGatewayIP(t *testing.T) string {
	t.Helper()

	// Try to get interfaces and find the Docker bridge IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		// Look for docker0 or similar bridge interfaces
		if iface.Name == "docker0" || strings.HasPrefix(iface.Name, "br-") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						t.Logf("Detected Docker gateway IP: %s", ipnet.IP.String())
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}
