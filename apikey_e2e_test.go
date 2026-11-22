package kuma_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// TestAPIKeyAgainstServer covers WithAPIKey against a real server: an Uptime
// Kuma 3.x server authenticates the socket.io connection by the key once "API
// and MCP" is turned on, a 2.x server has no such login.
//
// It brings up an Uptime Kuma of its own, because it changes the server's
// settings, and is gated the way the other end-to-end tests are.
func TestAPIKeyAgainstServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if !e2eTest {
		t.Skip(`skipping end to end test, "E2E_TEST" env var not set`)
	}

	baseURL, scheme := newKumaContainer(t)

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()

	newWithKey := func(key string) (*kuma.Client, error) {
		return kuma.New(ctx, baseURL, "", "",
			kuma.WithAPIKey(key),
			kuma.WithLogLevel(kuma.LogLevel(os.Getenv("SOCKETIO_LOG_LEVEL"))),
			kuma.WithConnectTimeout(10*time.Second),
		)
	}

	if scheme != kuma.AuthSchemeHTTP {
		keyClient, err := newWithKey("uk1_not-a-key-a-2x-server-would-take")
		require.ErrorIs(t, err, kuma.ErrAPIKeyUnsupported)
		require.Nil(t, keyClient)

		return
	}

	// The setup signs in twice, see TestAuthenticationAgainstServer.
	awaitSignInWindow(scheme)

	admin, err := kuma.New(ctx, baseURL, testUsername, testPassword, kuma.WithConnectTimeout(10*time.Second))
	require.NoError(t, err)

	defer func() { _ = admin.Disconnect() }()

	key, err := kuma.AddAPIKey(admin, ctx, "e2e-"+randomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, key)

	t.Run("refused while the API is disabled", func(t *testing.T) {
		keyClient, err := newWithKey(key)
		require.ErrorIs(t, err, kuma.ErrAPIKeyRejected)
		require.Nil(t, keyClient)
	})

	require.NoError(t, kuma.SetAPIEnabled(admin, ctx, true, testPassword))

	t.Run("authenticates the connection by the key", func(t *testing.T) {
		keyClient, err := newWithKey(key)
		require.NoError(t, err)

		defer func() { _ = keyClient.Disconnect() }()

		require.Equal(t, kuma.AuthSchemeHTTP, keyClient.AuthScheme())
		require.Empty(t, keyClient.SessionToken(), "an API key login has no session")

		_, err = keyClient.GetSettings(ctx)
		require.NoError(t, err)

		require.NoError(t, keyClient.Resync(ctx), "a resync connects again with the key")

		_, err = keyClient.GetSettings(ctx)
		require.NoError(t, err)
	})

	t.Run("the key and the configured headers reach the server on every path", func(t *testing.T) {
		proxyURL, requests := newRecordingProxy(t, baseURL)

		headers := http.Header{}
		headers.Set("X-E2E-Proxy", "seen")

		keyClient, err := kuma.New(ctx, proxyURL, "", "",
			kuma.WithAPIKey(key),
			kuma.WithHeaders(headers),
			kuma.WithUserAgent("go-uptime-kuma-client-e2e/1.0"),
			kuma.WithConnectTimeout(10*time.Second),
		)
		require.NoError(t, err)
		require.NoError(t, keyClient.Disconnect())

		seen := map[string]bool{}

		for _, request := range requests() {
			kind := request.path
			if request.transport != "" {
				kind = request.transport
			}

			seen[kind] = true

			require.Equal(t, "seen", request.header.Get("X-E2E-Proxy"), kind)
			require.Equal(t, "go-uptime-kuma-client-e2e/1.0", request.header.Get("User-Agent"), kind)
			require.Equal(t, "Bearer "+key, request.header.Get("Authorization"), kind)
		}

		require.True(t, seen["/api/auth/ok"], "the probe went through the proxy")
		require.True(t, seen["polling"], "the long polling went through the proxy")
		require.True(t, seen["websocket"], "the WebSocket upgrade went through the proxy")
	})

	t.Run("an unknown key is refused", func(t *testing.T) {
		keyClient, err := newWithKey("uk0_unknown-key")
		require.ErrorIs(t, err, kuma.ErrAPIKeyRejected)
		require.Nil(t, keyClient)
	})
}

// newRecordingProxy starts a reverse proxy to target that records the path,
// the engine.io transport and the headers of every request it forwards, the
// WebSocket upgrade included.
func newRecordingProxy(t *testing.T, target string) (proxyURL string, requests func() []fakeRequest) {
	t.Helper()

	upstream, err := url.Parse(target)
	require.NoError(t, err)

	proxy := httputil.NewSingleHostReverseProxy(upstream)

	var (
		mu       sync.Mutex
		recorded []fakeRequest
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		recorded = append(recorded, fakeRequest{
			path:      r.URL.Path,
			transport: r.URL.Query().Get("transport"),
			header:    r.Header.Clone(),
		})
		mu.Unlock()

		proxy.ServeHTTP(w, r)
	}))

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return server.URL, func() []fakeRequest {
		mu.Lock()
		defer mu.Unlock()

		return append([]fakeRequest(nil), recorded...)
	}
}
