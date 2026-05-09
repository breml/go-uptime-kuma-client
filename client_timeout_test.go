package kuma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// fakeSocketIOServer implements a minimal socket.io server over HTTP long-polling.
// It handles the initial engine.io handshake, socket.io CONNECT, and login phases,
// but then deliberately blocks on subsequent GET requests without ever emitting
// ready events (e.g. apiKeyList). This simulates the Uptime Kuma 2.3.x behavior
// that triggers issue #271.
type fakeSocketIOServer struct {
	messages chan []byte
}

func (s *fakeSocketIOServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	switch {
	case r.Method == http.MethodGet && sid == "":
		// Initial engine.io handshake: respond with PacketOpen ("0") + JSON.
		type engineIOHandshake struct {
			Sid          string   `json:"sid"`
			Upgrades     []string `json:"upgrades"`
			PingInterval int      `json:"pingInterval"`
			PingTimeout  int      `json:"pingTimeout"`
			MaxPayload   int      `json:"maxPayload"`
		}
		data, _ := json.Marshal(engineIOHandshake{
			Sid:          "test-sid",
			Upgrades:     []string{}, // no WebSocket upgrade keeps the client on polling
			PingInterval: 50,
			PingTimeout:  5000,
			MaxPayload:   1000000,
		})
		_, _ = w.Write(append([]byte("0"), data...))

	case r.Method == http.MethodPost && sid != "":
		body, _ := io.ReadAll(r.Body)
		s.handleClientMessage(body)
		w.WriteHeader(http.StatusOK)

	case r.Method == http.MethodGet && sid != "":
		// Long-poll: deliver the next queued message or block until the
		// client disconnects (simulating a server that never emits ready events).
		select {
		case msg := <-s.messages:
			_, _ = w.Write(msg)

		case <-r.Context().Done():
		}

	default:
	}
}

// handleClientMessage processes socket.io messages received from the client via POST.
// The body format is: engine.io type byte ("4" for Message) + socket.io packet.
func (s *fakeSocketIOServer) handleClientMessage(body []byte) {
	if len(body) < 2 {
		return
	}

	socketIOData := body[1:] // strip engine.io "4" (Message) prefix

	switch socketIOData[0] {
	case '0': // socket.io CONNECT
		// Enqueue CONNECT ACK: engine.io "4" (Message) + socket.io "0" (Connect).
		s.messages <- []byte("40")

	case '2': // socket.io EVENT — the login request
		// Extract the ack ID: digits immediately after the "2" type byte.
		i := 1
		for i < len(socketIOData) && socketIOData[i] >= '0' && socketIOData[i] <= '9' {
			i++
		}

		ackID := string(socketIOData[1:i])
		// Enqueue login ACK: engine.io "4" + socket.io "3" (Ack) + ack ID + JSON payload.
		ack := fmt.Sprintf(`43%s[{"ok":true,"msg":"Logged in successfully."}]`, ackID)
		s.messages <- []byte(ack)

	default:
	}
}

func TestNewConnectTimeoutDuringReadyWait(t *testing.T) {
	// Regression test for issue #271: kuma.New() hangs indefinitely waiting for
	// ready events when WithConnectTimeout is set (Uptime Kuma 2.3.x).
	//
	// The fake server completes the handshake and login phases successfully, but
	// never sends the required ready events (e.g. apiKeyList). kuma.New() must
	// return a context.DeadlineExceeded error within the configured timeout.
	server := httptest.NewServer(&fakeSocketIOServer{
		messages: make(chan []byte, 10),
	})

	// Use an explicit cancellable context so that background goroutines
	// started by kuma.New() can be stopped before server.Close() waits
	// for all active connections to drain.
	ctx, cancel := context.WithCancel(t.Context())

	const timeout = 500 * time.Millisecond
	start := time.Now()

	_, err := kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(timeout),
	)

	elapsed := time.Since(start)

	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded, "error should wrap context.DeadlineExceeded, got: %v", err)
	require.Less(t, elapsed, 2*timeout, "New() should return within 2x timeout, took %s", elapsed)

	// Cancel the context to abort in-flight requests from background
	// goroutines, then force-close any remaining connections so that
	// server.Close() does not block waiting for active connections to drain.
	cancel()
	server.CloseClientConnections()
	server.Close()
}
