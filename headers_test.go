package kuma_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

const fakeAPIKey = "uk1_fake-api-key"

// TestHeadersOnEveryHTTPRequest covers WithHeaders and WithUserAgent on the
// requests that go through the client's shared http.Client: the auth scheme
// probe, the HTTP login of a 3.x server and the socket.io long polling.
func TestHeadersOnEveryHTTPRequest(t *testing.T) {
	fake, url := newHTTPAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	configured := http.Header{}
	configured.Set("CF-Access-Client-Id", "client-id")
	configured.Set("Cookie", "proxy=1")
	configured.Set("Accept", "text/html")
	configured.Set("User-Agent", "overridden-by-with-user-agent")

	client, err := kuma.New(ctx, url, "admin", "admin1",
		kuma.WithConnectTimeout(5*time.Second),
		kuma.WithHeaders(configured),
		kuma.WithUserAgent("custom-agent/1.0"),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	// Changing the header after the option was applied changes nothing.
	configured.Set("CF-Access-Client-Id", "changed")

	require.NoError(t, client.Resync(ctx))

	requests := fake.receivedRequests()

	var sawProbe, sawSignIn, sawPolling bool

	for _, request := range requests {
		require.Equal(t, "client-id", request.header.Get("CF-Access-Client-Id"), request.path)
		require.Equal(t, "custom-agent/1.0", request.header.Get("User-Agent"), request.path)
		require.Contains(t, request.header.Get("Cookie"), "proxy=1", request.path)

		switch {
		case request.path == "/api/auth/ok":
			sawProbe = true

		case request.path == "/api/auth/sign-in/username":
			sawSignIn = true

			require.Equal(t, "application/json", request.header.Get("Accept"),
				"a header the client sets itself takes precedence")

		case request.transport == "polling":
			sawPolling = true

			require.Equal(t, "text/html", request.header.Get("Accept"))
			require.Contains(t, request.header.Get("Cookie"), "better-auth.session_token="+fakeSessionToken,
				"the configured cookie is added to the jar's, not put in its place")

		default:
		}
	}

	require.True(t, sawProbe, "the probe carries the headers")
	require.True(t, sawSignIn, "the sign-in carries the headers")
	require.True(t, sawPolling, "the long polling carries the headers")
}

// TestDefaultUserAgent covers the User-Agent a client sends unless it is told
// otherwise.
func TestDefaultUserAgent(t *testing.T) {
	fake, url := newHTTPAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	requests := fake.receivedRequests()
	require.NotEmpty(t, requests)

	for _, request := range requests {
		require.True(t, strings.HasPrefix(request.header.Get("User-Agent"), "go-uptime-kuma-client/"),
			"%s: %q", request.path, request.header.Get("User-Agent"))
		require.Empty(t, request.header.Get("Authorization"), request.path)
	}
}

// TestModuleVersion covers where the version of the default User-Agent comes
// from.
func TestModuleVersion(t *testing.T) {
	const module = "github.com/breml/go-uptime-kuma-client"

	tests := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{
			name: "no build information",
			want: "dev",
		},
		{
			name: "a dependency of the program",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/program", Version: "v1.2.3"},
				Deps: []*debug.Module{
					{Path: "example.com/other", Version: "v9.9.9"},
					{Path: module, Version: "v0.8.0"},
				},
			},
			want: "v0.8.0",
		},
		{
			name: "a dependency replaced by another version",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/program"},
				Deps: []*debug.Module{
					{
						Path:    module,
						Version: "v0.8.0",
						Replace: &debug.Module{Path: "example.com/fork", Version: "v0.8.1"},
					},
				},
			},
			want: "v0.8.1",
		},
		{
			name: "a dependency replaced by a local directory",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/program"},
				Deps: []*debug.Module{
					{Path: module, Version: "v0.8.0", Replace: &debug.Module{Path: "../go-uptime-kuma-client"}},
				},
			},
			want: "dev",
		},
		{
			name: "the module itself",
			info: &debug.BuildInfo{Main: debug.Module{Path: module, Version: "(devel)"}},
			want: "dev",
		},
		{
			name: "not among the dependencies",
			info: &debug.BuildInfo{Main: debug.Module{Path: "example.com/program", Version: "v1.2.3"}},
			want: "dev",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, kuma.ModuleVersion(tc.info))
		})
	}
}

// headerRecorder records the headers of the WebSocket handshake a server it
// wraps is asked for.
type headerRecorder struct {
	next http.Handler

	mu        sync.Mutex
	websocket http.Header
}

func (r *headerRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("transport") == "websocket" {
		r.mu.Lock()
		r.websocket = req.Header.Clone()
		r.mu.Unlock()
	}

	r.next.ServeHTTP(w, req)
}

func (r *headerRecorder) websocketHeader() http.Header {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.websocket
}

// TestHeadersOnWebSocketHandshake covers the configured headers on the
// WebSocket handshake, which the socket.io client makes without the shared
// http.Client, see cookieWebSocket.
func TestHeadersOnWebSocketHandshake(t *testing.T) {
	recorder := &headerRecorder{next: &fakeSocketIOServerWithWebSocket{}}
	server := httptest.NewServer(recorder)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	configured := http.Header{}
	configured.Set("CF-Access-Client-Id", "client-id")
	configured.Set("Cookie", "proxy=1")

	// The fake upgrades and then never sends the lists, so New fails waiting
	// for them; the handshake has been made by then.
	_, err := kuma.New(ctx, server.URL, "admin", "admin1",
		kuma.WithConnectTimeout(time.Second),
		kuma.WithHeaders(configured),
		kuma.WithUserAgent("custom-agent/1.0"),
	)
	require.Error(t, err)

	header := recorder.websocketHeader()
	require.NotNil(t, header, "the client upgraded to WebSocket")
	require.Equal(t, "client-id", header.Get("CF-Access-Client-Id"))
	require.Equal(t, "custom-agent/1.0", header.Get("User-Agent"))
	require.Equal(t, "proxy=1", header.Get("Cookie"))
}

// TestNewWithAPIKey covers the API key login of a 3.x server, see
// kuma.WithAPIKey: the key on the handshake instead of a login.
func TestNewWithAPIKey(t *testing.T) {
	t.Run("the server authenticates the connection by the key", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setAPIKey(fakeAPIKey)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, kuma.AuthSchemeHTTP, client.AuthScheme())
		require.Empty(t, client.SessionToken(), "an API key login has no session")
		require.Zero(t, fake.loginFrameCount(), "no sign-in is made")
		require.Empty(t, fake.receivedUnexpectedCommands())

		for _, request := range fake.receivedRequests() {
			require.Equal(t, "Bearer "+fakeAPIKey, request.header.Get("Authorization"), request.path)
		}

		// Resync connects again, presenting the key once more.
		require.NoError(t, client.Resync(ctx))
		require.Equal(t, 2, fake.handshakeCount())

		// A key the server stops accepting fails the next resync, and the
		// connection in service stays.
		fake.setAPIKey("")
		require.ErrorIs(t, client.Resync(ctx), kuma.ErrAPIKeyRejected)
	})

	t.Run("a key the server refuses", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setAPIKey("another-key")

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.ErrorIs(t, err, kuma.ErrAPIKeyRejected)
		require.ErrorContains(t, err, "apiEnabled", "the message names the setting as a cause")
		require.Nil(t, client)
	})

	t.Run("a 2.x server has no API key login", func(t *testing.T) {
		fake, url := newAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.ErrorIs(t, err, kuma.ErrAPIKeyUnsupported)
		require.Nil(t, client)
		require.Zero(t, fake.handshakeCount(), "the client does not connect")
	})

	t.Run("a pinned socket scheme has no API key login", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setAPIKey(fakeAPIKey)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAuthScheme(kuma.AuthSchemeSocket),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.ErrorIs(t, err, kuma.ErrAPIKeyUnsupported)
		require.Nil(t, client)
		require.Zero(t, fake.probeCount(), "a pinned scheme is not detected")
	})

	t.Run("an empty key configures nothing", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAPIKey(""),
		)
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, fakeSessionToken, client.SessionToken())
		require.Equal(t, 1, fake.loginFrameCount())
	})
}

// TestNewWithAPIKeyRejectsOtherCredentials covers the credentials an API key
// cannot be combined with: New fails before it makes a single request.
func TestNewWithAPIKeyRejectsOtherCredentials(t *testing.T) {
	authorization := http.Header{}
	authorization.Set("Authorization", "Basic Zm9vOmJhcg==")

	tests := []struct {
		name     string
		username string
		password string
		opts     []kuma.Option
		wantErr  string
	}{
		{
			name:     "username and password",
			username: "admin",
			password: "admin1",
			wantErr:  "WithAPIKey cannot be combined",
		},
		{
			name:    "session token",
			opts:    []kuma.Option{kuma.WithSessionToken(fakePresetToken)},
			wantErr: "WithAPIKey cannot be combined",
		},
		{
			name:    "totp secret",
			opts:    []kuma.Option{kuma.WithTOTPSecret("JBSWY3DPEHPK3PXP")},
			wantErr: "WithAPIKey cannot be combined",
		},
		{
			name: "totp code",
			opts: []kuma.Option{kuma.WithTOTPCode(func(context.Context) (string, error) {
				return "123456", nil
			})},
			wantErr: "WithAPIKey cannot be combined",
		},
		{
			name:    "an Authorization header",
			opts:    []kuma.Option{kuma.WithHeaders(authorization)},
			wantErr: "Authorization",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake, url := newHTTPAuthFake(t)
			fake.setAPIKey(fakeAPIKey)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			opts := append([]kuma.Option{
				kuma.WithConnectTimeout(5 * time.Second),
				kuma.WithAPIKey(fakeAPIKey),
			}, tc.opts...)

			client, err := kuma.New(ctx, url, tc.username, tc.password, opts...)
			require.ErrorContains(t, err, tc.wantErr)
			require.Nil(t, client)
			require.Empty(t, fake.receivedRequests(), "nothing is sent to the server")
		})
	}
}
