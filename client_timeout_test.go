package kuma_test

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
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
		data, err := json.Marshal(engineIOHandshake{
			Sid:          "test-sid",
			Upgrades:     []string{}, // no WebSocket upgrade keeps the client on polling
			PingInterval: 50,
			PingTimeout:  5000,
			MaxPayload:   1000000,
		})
		if err != nil {
			http.Error(w, "marshal handshake", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(append([]byte("0"), data...))

	case r.Method == http.MethodPost && sid != "":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusInternalServerError)
			return
		}

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

// fakeSocketIOServerWithWebSocket is like fakeSocketIOServer but advertises a
// WebSocket upgrade in the initial handshake and completes the full engine.io
// upgrade sequence (2probe → 3probe → upgrade packet "5") plus the socket.io
// CONNECT and login phases over WebSocket. It never sends ready events
// (apiKeyList, monitorList, etc.) so kuma.New() must return with a
// "missing events" error within the configured ConnectTimeout.
//
// This server is the minimal reproduction for issue #325: the socket.io
// polling connection opens ("New polling connection") but without the
// afterConnect goroutine fix the client deadlocks — afterConnect calls
// Send() which waits on waitUpgrade, while waitUpgrade only closes after
// messageLoop processes 3probe, which it can't do because it is blocked
// in afterConnect. The result is a hang for the full ConnectTimeout duration.
type fakeSocketIOServerWithWebSocket struct{}

func (*fakeSocketIOServerWithWebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")
	transport := r.URL.Query().Get("transport")

	switch {
	case transport == "websocket":
		handleWebSocketUpgrade(w, r)

	case r.Method == http.MethodGet && sid == "":
		// Initial engine.io handshake — offer WebSocket upgrade.
		type engineIOHandshake struct {
			Sid          string   `json:"sid"`
			Upgrades     []string `json:"upgrades"`
			PingInterval int      `json:"pingInterval"`
			PingTimeout  int      `json:"pingTimeout"`
			MaxPayload   int      `json:"maxPayload"`
		}
		data, err := json.Marshal(engineIOHandshake{
			Sid:          "test-sid",
			Upgrades:     []string{"websocket"},
			PingInterval: 50,
			PingTimeout:  5000,
			MaxPayload:   1000000,
		})
		if err != nil {
			http.Error(w, "marshal handshake", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(append([]byte("0"), data...))

	case r.Method == http.MethodPost && sid != "":
		// Drain any polling POST that arrives during the upgrade window.
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)

	case r.Method == http.MethodGet && sid != "":
		// Block polling GET until the client disconnects.
		<-r.Context().Done()

	default:
	}
}

// handleWebSocketUpgrade performs the WebSocket handshake and then drives the
// engine.io upgrade sequence followed by the socket.io CONNECT and login.
// After the login ACK it blocks without emitting any ready events.
func handleWebSocketUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()

	// Engine.io probe exchange.
	msg, err := wsReadText(conn)
	if err != nil || msg != "2probe" {
		return
	}

	err = wsSendText(conn, "3probe")
	if err != nil {
		return
	}

	// Upgrade completion packet "5".
	_, err = wsReadText(conn)
	if err != nil {
		return
	}

	// socket.io CONNECT "40" → ACK.
	_, err = wsReadText(conn)
	if err != nil {
		return
	}

	err = wsSendText(conn, `40{"sid":"ws-test-sid"}`)
	if err != nil {
		return
	}

	// socket.io login event "42N[\"login\",...]" → ACK.
	msg, err = wsReadText(conn)
	if err != nil || len(msg) < 3 || !strings.HasPrefix(msg, "42") {
		return
	}

	i := 2
	for i < len(msg) && msg[i] >= '0' && msg[i] <= '9' {
		i++
	}

	ackID := msg[2:i]
	ack := fmt.Sprintf(`43%s[{"ok":true,"msg":"Logged in successfully."}]`, ackID)
	err = wsSendText(conn, ack)
	if err != nil {
		return
	}

	// Block until the client disconnects (never send ready events).
	_, _ = wsReadText(conn)
}

// wsHandshake upgrades the HTTP connection to WebSocket using http.Hijacker.
func wsHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	key := r.Header.Get("Sec-Websocket-Key")
	h := sha1.New()
	_, _ = io.WriteString(h, key+"258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("ResponseWriter is not a Hijacker")
	}

	conn, rw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}

	resp := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	_, err = rw.WriteString(resp)
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = rw.Flush()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// wsSendText sends an unmasked WebSocket text frame (server → client).
// Handles payloads up to 125 bytes, sufficient for all test messages.
func wsSendText(conn net.Conn, msg string) error {
	data := []byte(msg)
	frame := make([]byte, 2+len(data))
	frame[0] = 0x81 // FIN=1, opcode=text
	frame[1] = byte(len(data))
	copy(frame[2:], data)
	_, err := conn.Write(frame)
	return err
}

// wsReadText reads one masked WebSocket text frame (client → server).
// Handles payloads up to 125 bytes, sufficient for all test messages.
func wsReadText(conn net.Conn) (string, error) {
	header := make([]byte, 2)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return "", err
	}

	masked := header[1]&0x80 != 0
	n := int(header[1] & 0x7F)
	var mask [4]byte
	if masked {
		_, err = io.ReadFull(conn, mask[:])
		if err != nil {
			return "", err
		}
	}

	data := make([]byte, n)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return "", err
	}

	if masked {
		for i := range data {
			data[i] ^= mask[i%4]
		}
	}

	return string(data), nil
}

// TestNewConnectTimeoutDuringReadyWait_WebSocketUpgrade is the WebSocket-upgrade
// variant of TestNewConnectTimeoutDuringReadyWait and is the direct regression
// test for issue #325.
//
// The fake server completes the full engine.io WebSocket upgrade sequence
// (polling handshake → 2probe/3probe → upgrade "5" → socket.io CONNECT ACK
// → login ACK), but never emits ready events (apiKeyList etc.). kuma.New()
// must return with a "missing events" error — NOT with "connect to server:"
// — meaning it successfully passed the WebSocket upgrade phase.
//
// With the bug (afterConnect called inline from messageLoop): messageLoop
// blocks in Send() waiting for waitUpgrade, 3probe sits unprocessed in the
// buffer, waitUpgrade never closes → the select in kuma.New() never reaches
// the "connect" case. The error is "connect to server: context deadline
// exceeded" (stuck in the connect-wait phase), and the test fails because
// "missing events:" is absent from the error.
func TestNewConnectTimeoutDuringReadyWait_WebSocketUpgrade(t *testing.T) {
	const timeout = 500 * time.Millisecond

	server := httptest.NewServer(&fakeSocketIOServerWithWebSocket{})

	ctx, cancel := context.WithTimeout(t.Context(), 10*timeout)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	start := time.Now()

	_, err := kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(timeout),
	)

	elapsed := time.Since(start)

	require.Error(t, err)
	require.Less(t, elapsed, 2*timeout,
		"New() should not deadlock — elapsed %s exceeds 2× timeout", elapsed)
	require.ErrorIs(t, err, context.DeadlineExceeded,
		"error should wrap context.DeadlineExceeded, got: %v", err)
	// This assertion is the regression guard: with the bug, afterConnect deadlocks
	// so kuma.New() never reaches the ready-events phase and the error is
	// "connect to server: ...", not "missing events: ...".
	require.ErrorContains(t, err, "missing events:",
		"error should report missing events (not a connect-phase hang)")
	require.ErrorContains(t, err, "apiKeyList")
}

func TestNewConnectTimeoutDuringReadyWait(t *testing.T) {
	// Regression test for issue #271: kuma.New() hangs indefinitely waiting for
	// ready events when WithConnectTimeout is set (Uptime Kuma 2.3.x).
	//
	// The fake server completes the handshake and login phases successfully, but
	// never sends the required ready events (e.g. apiKeyList). kuma.New() must
	// return a context.DeadlineExceeded error within the configured timeout.
	const timeout = 500 * time.Millisecond

	server := httptest.NewServer(&fakeSocketIOServer{
		messages: make(chan []byte, 10),
	})

	// Bound the test itself so it fails fast if the regression reappears and
	// kuma.New() blocks indefinitely rather than honouring the timeout.
	ctx, cancel := context.WithTimeout(t.Context(), 5*timeout)

	// Always cancel the context and close the server, even when a require.*
	// assertion short-circuits the test via FailNow.
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

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
	require.ErrorContains(t, err, "missing events:", "error should list the still-missing ready events")
	require.ErrorContains(t, err, "apiKeyList", "apiKeyList should appear in the missing events list")
	require.Less(t, elapsed, 2*timeout, "New() should return within 2x timeout, took %s", elapsed)
}
