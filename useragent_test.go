package kuma_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// capturedRequest records the User-Agent header and the "transport" query
// parameter (e.g. "websocket", or "" for the initial polling handshake) of
// one captured request.
type capturedRequest struct {
	userAgent string
	transport string
}

// headerCapturingHandler wraps an http.Handler and records the User-Agent
// header of every request it sees before delegating to inner.
type headerCapturingHandler struct {
	inner http.Handler

	mu       sync.Mutex
	requests []capturedRequest
}

func (h *headerCapturingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.requests = append(h.requests, capturedRequest{
		userAgent: r.Header.Get("User-Agent"),
		transport: r.URL.Query().Get("transport"),
	})
	h.mu.Unlock()

	h.inner.ServeHTTP(w, r)
}

func (h *headerCapturingHandler) first() string {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.requests) == 0 {
		return ""
	}

	return h.requests[0].userAgent
}

// userAgentFor returns the User-Agent header of the first captured request
// whose "transport" query parameter matches transport, or "" if none match.
func (h *headerCapturingHandler) userAgentFor(transport string) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, req := range h.requests {
		if req.transport == transport {
			return req.userAgent
		}
	}

	return ""
}

// TestWithUserAgent_Polling verifies that WithUserAgent attaches the
// configured User-Agent header to the Socket.IO long-polling transport's
// HTTP requests (the initial engine.io handshake in this case).
func TestWithUserAgent_Polling(t *testing.T) {
	const timeout = 500 * time.Millisecond
	const userAgent = "terraform-provider-uptimekuma/test"

	capture := &headerCapturingHandler{
		inner: &fakeSocketIOServer{messages: make(chan []byte, 10)},
	}
	server := httptest.NewServer(capture)

	ctx, cancel := context.WithTimeout(t.Context(), 5*timeout)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	// The fake server never emits ready events, so New() always times out.
	// We only care that the handshake request carried the right header.
	_, _ = kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(timeout),
		kuma.WithUserAgent(userAgent),
	)

	require.Equal(t, userAgent, capture.first())
}

// TestWithUserAgent_PollingDefault verifies that without WithUserAgent, the
// polling transport falls back to Go's default net/http User-Agent instead
// of sending no header at all.
func TestWithUserAgent_PollingDefault(t *testing.T) {
	const timeout = 500 * time.Millisecond

	capture := &headerCapturingHandler{
		inner: &fakeSocketIOServer{messages: make(chan []byte, 10)},
	}
	server := httptest.NewServer(capture)

	ctx, cancel := context.WithTimeout(t.Context(), 5*timeout)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	_, _ = kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(timeout),
	)

	require.Contains(t, capture.first(), "Go-http-client",
		"expected the unconfigured default net/http User-Agent, got %q", capture.first())
}

// TestWithUserAgent_WebSocketHandshake verifies that WithUserAgent attaches
// the configured User-Agent header to the WebSocket upgrade request.
func TestWithUserAgent_WebSocketHandshake(t *testing.T) {
	const timeout = 500 * time.Millisecond
	const userAgent = "terraform-provider-uptimekuma/test"

	capture := &headerCapturingHandler{
		inner: &fakeSocketIOServerWithWebSocket{},
	}
	server := httptest.NewServer(capture)

	ctx, cancel := context.WithTimeout(t.Context(), 10*timeout)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	_, _ = kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(timeout),
		kuma.WithUserAgent(userAgent),
	)

	require.Equal(t, userAgent, capture.first(),
		"the initial polling handshake should carry the configured User-Agent")
	require.Equal(t, userAgent, capture.userAgentFor("websocket"),
		"the WebSocket upgrade request should carry the configured User-Agent")
}

// entryPageHandler serves a minimal /api/entry-page response indicating no
// database setup is required, so setupDatabase (invoked via WithAutosetup)
// returns immediately after a single request. Every other path is delegated
// to a fakeSocketIOServer so the subsequent Socket.IO connection attempt
// behaves the same way as in TestWithUserAgent_Polling (a clean timeout)
// instead of hitting an unhandled 404, which the underlying Socket.IO client
// does not appear to recover from within ConnectTimeout.
type entryPageHandler struct {
	socketIO *fakeSocketIOServer
}

func (h *entryPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/entry-page" {
		h.socketIO.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"type":"setup"}`))
}

// TestWithUserAgent_Autosetup verifies that WithUserAgent attaches the
// configured User-Agent header to the entry-page HTTP request issued during
// autosetup.
func TestWithUserAgent_Autosetup(t *testing.T) {
	const timeout = 500 * time.Millisecond
	const userAgent = "terraform-provider-uptimekuma/test"

	capture := &headerCapturingHandler{
		inner: &entryPageHandler{socketIO: &fakeSocketIOServer{messages: make(chan []byte, 10)}},
	}
	server := httptest.NewServer(capture)

	ctx, cancel := context.WithTimeout(t.Context(), 5*timeout)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	// The fake server never emits ready events (see fakeSocketIOServer), so
	// New() always times out after the autosetup phase. We only care that
	// the entry-page request carried the right header.
	_, _ = kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithAutosetup(),
		kuma.WithConnectTimeout(timeout),
		kuma.WithUserAgent(userAgent),
	)

	require.Equal(t, userAgent, capture.first())
}
