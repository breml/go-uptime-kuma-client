package kuma_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// recordSeparator delimits the packets of an engine.io long-polling payload.
const recordSeparator = 0x1e

// batchingSocketIOServer is fakeSocketIOServer, except that a long poll
// answers with every packet queued at the time, joined into one payload by the
// record separator. That is how an engine.io server flushes its write buffer to
// a client that polls: the ack of a login and the lists that follow it arrive
// in one response body.
type batchingSocketIOServer struct {
	*fakeSocketIOServer

	// handshakes counts the engine.io handshakes, each of which opens a
	// session of its own on the server.
	handshakes atomic.Int64
	// batches counts the payloads that carried more than one packet.
	batches atomic.Int64
}

// batchSettle is how long a long poll that has a packet waits for the packets
// enqueued right behind it, so they go out in the same payload.
const batchSettle = 20 * time.Millisecond

func (s *batchingSocketIOServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	if r.Method != http.MethodGet || sid == "" || r.URL.Path == "/api/auth/ok" {
		if r.Method == http.MethodGet && sid == "" && r.URL.Path != "/api/auth/ok" {
			s.handshakes.Add(1)
		}

		s.fakeSocketIOServer.ServeHTTP(w, r)

		return
	}

	var batch [][]byte

	select {
	case msg := <-s.messages:
		batch = append(batch, msg)

	case <-r.Context().Done():
		return
	}

	timer := time.NewTimer(batchSettle)
	defer timer.Stop()

	for collecting := true; collecting; {
		select {
		case msg := <-s.messages:
			batch = append(batch, msg)

		case <-timer.C:
			collecting = false
		}
	}

	if len(batch) > 1 {
		s.batches.Add(1)
	}

	_, _ = w.Write(bytes.Join(batch, []byte{recordSeparator}))
}

// allListEvents are the list events Uptime Kuma sends after a login, each with
// an empty list.
func allListEvents() []string {
	return []string{
		`42["monitorList",{}]`,
		`42["notificationList",[]]`,
		`42["statusPageList",{}]`,
		`42["maintenanceList",{}]`,
		`42["proxyList",[]]`,
		`42["dockerHostList",[]]`,
		`42["apiKeyList",[]]`,
	}
}

// newPollingServer starts server and returns its URL.
func newPollingServer(t *testing.T, server http.Handler) string {
	t.Helper()

	httpServer := httptest.NewServer(server)
	t.Cleanup(func() {
		httpServer.CloseClientConnections()
		httpServer.Close()
	})

	return httpServer.URL
}

// TestNewPollingBatchedPayload covers a client that stays on long polling and
// gets the login ack and every list in one payload. A client that reads a
// payload as one packet keeps the first and drops the rest, and New waits out
// its timeout for lists that already arrived (issue #405).
func TestNewPollingBatchedPayload(t *testing.T) {
	server := &batchingSocketIOServer{
		fakeSocketIOServer: &fakeSocketIOServer{
			messages:         make(chan []byte, 20),
			eventsAfterLogin: allListEvents(),
		},
	}
	url := newPollingServer(t, server)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	pollingClient, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(3*time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { _ = pollingClient.Disconnect() })

	require.Empty(t, pollingClient.MissingReadyEvents(), "every list of the batched payload must reach the cache")
	require.Positive(t, server.batches.Load(), "the server must have batched packets for the test to mean anything")
	require.EqualValues(t, 1, server.handshakes.Load(), "the client must open exactly one engine.io session")
}

// TestPollingBatchedCommandResponse covers a command answered on long polling
// with its ack and the broadcast it waits for in one payload, the way a write
// on a client that polls is answered.
func TestPollingBatchedCommandResponse(t *testing.T) {
	server := &batchingSocketIOServer{
		fakeSocketIOServer: &fakeSocketIOServer{
			messages:         make(chan []byte, 20),
			eventsAfterLogin: allListEvents(),
		},
	}
	server.onCommand = func(command string, ackID string) {
		if command != "getMonitorList" {
			return
		}

		server.messages <- fmt.Appendf(nil, `43%s[{"ok":true}]`, ackID)
		server.messages <- []byte(`42["monitorList",{"1":{"id":1,"name":"batched","type":"http"},` +
			`"2":{"id":2,"name":"payload","type":"http"}}]`)
	}

	url := newPollingServer(t, server)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	pollingClient, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(3*time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { _ = pollingClient.Disconnect() })

	batchesBefore := server.batches.Load()

	monitors, err := pollingClient.GetMonitors(ctx)
	require.NoError(t, err)
	require.Len(t, monitors, 2)
	require.Greater(t, server.batches.Load(), batchesBefore, "the ack and the broadcast must have shared a payload")
}

// TestNewPollingLongPoll covers the latency of a client on long polling against
// a server with Uptime Kuma's ping interval of 25 seconds. The server holds a
// poll open until it has something to send, so a client has to have a poll in
// flight at all times; one that polls once per ping interval reads the connect
// ack and every packet after it up to 25 seconds late.
func TestNewPollingLongPoll(t *testing.T) {
	server := &batchingSocketIOServer{
		fakeSocketIOServer: &fakeSocketIOServer{
			messages:         make(chan []byte, 20),
			eventsAfterLogin: allListEvents(),
			pingInterval:     25000,
		},
	}
	url := newPollingServer(t, server)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	start := time.Now()

	pollingClient, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { _ = pollingClient.Disconnect() })

	require.Less(t, time.Since(start), 2*time.Second, "New must not wait for a ping interval between polls")
	require.Empty(t, pollingClient.MissingReadyEvents())
	require.EqualValues(t, 1, server.handshakes.Load(), "the client must open exactly one engine.io session")
}

// TestNewPollingLargePayload covers a payload above the socket.io library's
// default limit of 4 MiB for a long-polling response. A large instance sends
// lists of that size after a login, and a client on long polling gets them in
// one payload.
func TestNewPollingLargePayload(t *testing.T) {
	const (
		notifications = 5000
		nameSize      = 1024
	)

	var list strings.Builder

	list.WriteString(`42["notificationList",[`)

	for i := range notifications {
		if i > 0 {
			list.WriteString(",")
		}

		fmt.Fprintf(
			&list,
			`{"id":%d,"name":"%s","active":true,"config":"{\"type\":\"webhook\"}"}`,
			i+1,
			strings.Repeat("n", nameSize),
		)
	}

	list.WriteString(`]]`)
	require.Greater(t, list.Len(), 4<<20, "the payload must exceed the library's default limit")

	events := allListEvents()
	events[1] = list.String()

	server := &batchingSocketIOServer{
		fakeSocketIOServer: &fakeSocketIOServer{
			messages:         make(chan []byte, 20),
			eventsAfterLogin: events,
		},
	}
	url := newPollingServer(t, server)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	pollingClient, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { _ = pollingClient.Disconnect() })

	require.Empty(t, pollingClient.MissingReadyEvents())
	require.Len(t, pollingClient.GetNotifications(ctx), notifications)
}
