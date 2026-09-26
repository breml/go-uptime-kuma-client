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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// TestEndToEndManualMonitorStatus verifies that ManualStatus reaches the
// server and drives the heartbeats of a manual monitor. The server never
// returns manual_status, so the status is observed through the webhook
// notifications instead:
// 1. The monitor is created with status down, its first beat notifies down.
// 2. The monitor is edited to status up, the edit restarts the monitor and
// the down -> up transition notifies up.
func TestEndToEndManualMonitorStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if !e2eTest {
		t.Skip(`skipping end to end test, "E2E_TEST" env var not set`)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	payloads := make(chan WebhookPayload, 16)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Logf("Failed to decode webhook payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Logf("Webhook payload received: monitor=%d, status=%d, msg=%s",
			payload.Monitor.ID, payload.Heartbeat.Status, payload.Heartbeat.Msg)

		select {
		case payloads <- payload:
		default:
		}

		w.WriteHeader(http.StatusOK)
	})

	// Listen on all interfaces, so the Docker container can reach the server.
	server := httptest.NewUnstartedServer(mux)
	var err error
	server.Listener, err = net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err)
	server.Start()
	defer server.Close()

	gatewayIP := "172.17.0.1"
	gw := getDockerGatewayIP(t)
	if gw != "" {
		gatewayIP = gw
	}

	serverPort := server.Listener.Addr().(*net.TCPAddr).Port
	dockerAccessibleURL := "http://" + net.JoinHostPort(gatewayIP, strconv.Itoa(serverPort))

	notificationID, err := client.CreateNotification(ctx, notification.Webhook{
		Base: notification.Base{
			IsActive: true,
			Name:     "E2E Manual Monitor Webhook",
		},
		WebhookDetails: notification.WebhookDetails{
			WebhookURL:         fmt.Sprintf("%s/webhook", dockerAccessibleURL),
			WebhookContentType: "json",
		},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		//nolint:usetesting // t.Context() cannot be used in cleanup, since it is already done.
		err := client.DeleteNotification(context.Background(), notificationID)
		if err != nil {
			t.Logf("Failed to delete notification: %v", err)
		}
	})

	manualMonitor := monitor.Manual{
		Base: monitor.Base{
			Name:            "E2E Manual Monitor",
			Interval:        60,
			RetryInterval:   60,
			MaxRetries:      0,
			IsActive:        true,
			NotificationIDs: []int64{notificationID},
		},
		ManualDetails: monitor.ManualDetails{
			ManualStatus: new(monitor.ManualStatusDown),
		},
	}

	monitorID, err := client.CreateMonitor(ctx, &manualMonitor)
	require.NoError(t, err)

	t.Cleanup(func() {
		//nolint:usetesting // t.Context() cannot be used in cleanup, since it is already done.
		err := client.DeleteMonitor(context.Background(), monitorID)
		if err != nil {
			t.Logf("Failed to delete monitor: %v", err)
		}
	})

	waitForStatus := func(t *testing.T, wantStatus monitor.ManualStatus, wantMsg string) {
		t.Helper()

		for {
			select {
			case payload := <-payloads:
				if payload.Monitor.ID != monitorID {
					continue
				}

				require.Equal(t, int(wantStatus), payload.Heartbeat.Status)
				require.Equal(t, wantMsg, payload.Heartbeat.Msg)
				return

			case <-ctx.Done():
				t.Fatalf("timeout waiting for notification with status %d", wantStatus)
			}
		}
	}

	waitForStatus(t, monitor.ManualStatusDown, "Down")

	// The server stores the down heartbeat only after the notification has
	// been sent. If the edit restarts the monitor before that, the up
	// heartbeat counts as the first one, which does not notify.
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		t.Fatal("timeout waiting for the down heartbeat to be stored")
	}

	manualMonitor.ID = monitorID
	manualMonitor.ManualStatus = new(monitor.ManualStatusUp)
	err = client.UpdateMonitor(ctx, &manualMonitor)
	require.NoError(t, err)

	waitForStatus(t, monitor.ManualStatusUp, "Up")
}
