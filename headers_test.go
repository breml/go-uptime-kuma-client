package kuma_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"runtime/debug"
	"strconv"
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

	// A lower-case name in a literal is canonicalised like any other.
	configured := http.Header{"cookie": {"proxy=1"}}
	configured.Set("CF-Access-Client-Id", "client-id")
	configured.Set("User-Agent", "overridden-by-with-user-agent")

	withHeaders := kuma.WithHeaders(configured)

	// The option copied the headers when it was called, so changing them
	// before New runs it changes nothing.
	configured.Set("CF-Access-Client-Id", "changed")

	client, err := kuma.New(ctx, url, "admin", "admin1",
		kuma.WithConnectTimeout(5*time.Second),
		withHeaders,
		kuma.WithUserAgent("custom-agent/1.0"),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

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

		case request.transport == "polling":
			sawPolling = true

			require.Contains(t, request.header.Get("Cookie"), "better-auth.session_token="+fakeSessionToken,
				"the configured cookie is added to the jar's, not put in its place")

		default:
		}
	}

	require.True(t, sawProbe, "the probe carries the headers")
	require.True(t, sawSignIn, "the sign-in carries the headers")
	require.True(t, sawPolling, "the long polling carries the headers")
}

// TestWithHeadersAddsUp covers WithHeaders given more than once, and one
// header literal holding two spellings of a name.
func TestWithHeadersAddsUp(t *testing.T) {
	fake, url := newHTTPAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1",
		kuma.WithConnectTimeout(5*time.Second),
		kuma.WithHeaders(http.Header{"X-A": {"a1"}, "X-B": {"b1"}}),
		kuma.WithHeaders(http.Header{"X-B": {"b2"}, "X-C": {"c2"}}),
		kuma.WithHeaders(http.Header{"x-spelled": {"lower"}, "X-Spelled": {"canonical"}}),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	requests := fake.receivedRequests()
	require.NotEmpty(t, requests)

	for _, request := range requests {
		require.Equal(t, []string{"a1"}, request.header.Values("X-A"), request.path)
		require.Equal(t, []string{"b2"}, request.header.Values("X-B"), "a repeated name replaces")
		require.Equal(t, []string{"c2"}, request.header.Values("X-C"), request.path)
		require.Equal(t, []string{"canonical", "lower"}, request.header.Values("X-Spelled"),
			"two spellings are kept, in the sorted order of the spellings")
	}
}

// TestNewRejectsInvalidHeaders covers the headers and keys New refuses before
// it sends anything.
func TestNewRejectsInvalidHeaders(t *testing.T) {
	tests := []struct {
		name    string
		opts    []kuma.Option
		wantErr string
	}{
		{name: "Host", opts: []kuma.Option{kuma.WithHeaders(http.Header{"Host": {"x"}})}, wantErr: "Host is set by"},
		{
			name:    "Origin",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"origin": {"x"}})},
			wantErr: "Origin is set by",
		},
		{
			name:    "Accept-Encoding",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Accept-Encoding": {"gzip"}})},
			wantErr: "Accept-Encoding is set by",
		},
		{
			name:    "Accept",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Accept": {"x"}})},
			wantErr: "Accept is set by",
		},
		{name: "TE", opts: []kuma.Option{kuma.WithHeaders(http.Header{"TE": {"x"}})}, wantErr: "Te is set by"},
		{
			name:    "Sec-WebSocket-Extensions",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Sec-WebSocket-Extensions": {"x"}})},
			wantErr: "Sec-Websocket-Extensions is set by",
		},
		{
			name:    "invalid name",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"X Bad": {"x"}})},
			wantErr: "not a valid header name",
		},
		{
			name:    "invalid value",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"X-Bad": {"a\r\nX-Injected: 1"}})},
			wantErr: "X-Bad: not a valid header value",
		},
		{
			name:    "a Bearer token in Authorization",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Authorization": {"bearer abc"}})},
			wantErr: "use WithAPIKey",
		},
		{
			name:    "an invalid User-Agent",
			opts:    []kuma.Option{kuma.WithUserAgent("agent\n")},
			wantErr: "WithUserAgent",
		},
		{
			name:    "an API key with the Bearer prefix",
			opts:    []kuma.Option{kuma.WithAPIKey("BEARER " + fakeAPIKey)},
			wantErr: "without the Bearer prefix",
		},
		{
			name:    "an API key with whitespace",
			opts:    []kuma.Option{kuma.WithAPIKey("uk1 key")},
			wantErr: "whitespace or control characters",
		},
		{
			name:    "an API key with a control character",
			opts:    []kuma.Option{kuma.WithAPIKey("uk1_key\x00")},
			wantErr: "whitespace or control characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake, url := newHTTPAuthFake(t)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			client, err := kuma.New(ctx, url, "", "",
				append([]kuma.Option{kuma.WithConnectTimeout(5 * time.Second)}, tc.opts...)...)
			require.ErrorContains(t, err, tc.wantErr)
			require.Nil(t, client)
			require.Empty(t, fake.receivedRequests(), "nothing is sent to the server")
		})
	}
}

// TestBasicAuthorizationHeader covers an Authorization header for a reverse
// proxy in front of the server that wants HTTP basic auth.
func TestBasicAuthorizationHeader(t *testing.T) {
	fake, url := newHTTPAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1",
		kuma.WithConnectTimeout(5*time.Second),
		kuma.WithHeaders(http.Header{"Authorization": {"Basic Zm9vOmJhcg=="}}),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	requests := fake.receivedRequests()
	require.NotEmpty(t, requests)

	for _, request := range requests {
		require.Equal(t, "Basic Zm9vOmJhcg==", request.header.Get("Authorization"), request.path)
	}
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
// wraps is asked for. With setCookie it answers the auth scheme probe with that
// cookie, so the jar holds a cookie of the server's by the time the client
// upgrades.
type headerRecorder struct {
	next      http.Handler
	setCookie *http.Cookie

	mu        sync.Mutex
	websocket http.Header
}

func (r *headerRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("transport") == "websocket" {
		r.mu.Lock()
		r.websocket = req.Header.Clone()
		r.mu.Unlock()
	}

	if req.URL.Path == "/api/auth/ok" && r.setCookie != nil {
		http.SetCookie(w, r.setCookie)
	}

	r.next.ServeHTTP(w, req)
}

func (r *headerRecorder) websocketHeader() http.Header {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.websocket
}

// newWebSocketRecorder starts a fake that upgrades to WebSocket behind a
// headerRecorder.
func newWebSocketRecorder(t *testing.T, setCookie *http.Cookie) (recorder *headerRecorder, url string) {
	t.Helper()

	recorder = &headerRecorder{next: &fakeSocketIOServerWithWebSocket{}, setCookie: setCookie}
	server := httptest.NewServer(recorder)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return recorder, server.URL
}

// TestHeadersOnWebSocketHandshake covers the configured headers on the
// WebSocket handshake, which the socket.io client makes without the shared
// http.Client, see cookieWebSocket.
func TestHeadersOnWebSocketHandshake(t *testing.T) {
	t.Run("the configured headers and cookies, next to the jar's", func(t *testing.T) {
		recorder, url := newWebSocketRecorder(t, &http.Cookie{Name: "server-session", Value: "abc", Path: "/"})

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		configured := http.Header{}
		configured.Set("CF-Access-Client-Id", "client-id")
		configured.Set("Cookie", "proxy=1")

		// The fake upgrades and then never sends the lists, so New fails
		// waiting for them; the handshake has been made by then.
		_, err := kuma.New(ctx, url, "admin", "admin1",
			kuma.WithConnectTimeout(time.Second),
			kuma.WithHeaders(configured),
			kuma.WithUserAgent("custom-agent/1.0"),
		)
		require.Error(t, err)

		header := recorder.websocketHeader()
		require.NotNil(t, header, "the client upgraded to WebSocket")
		require.Equal(t, "client-id", header.Get("CF-Access-Client-Id"))
		require.Equal(t, "custom-agent/1.0", header.Get("User-Agent"))
		require.Equal(t, "server-session=abc; proxy=1", header.Get("Cookie"),
			"the configured cookie is added to the jar's, not put in its place")
	})

	t.Run("the API key", func(t *testing.T) {
		recorder, url := newWebSocketRecorder(t, nil)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		// The fake speaks no 3.x auth, so the scheme is pinned to get the key
		// past New's check; the handshake is made before New fails.
		_, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(time.Second),
			kuma.WithAuthScheme(kuma.AuthSchemeHTTP),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.Error(t, err)

		header := recorder.websocketHeader()
		require.NotNil(t, header, "the client upgraded to WebSocket")
		require.Equal(t, "Bearer "+fakeAPIKey, header.Get("Authorization"))
	})
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

		_, err = client.GetSettings(ctx)
		require.NoError(t, err, "the connection in service still answers commands")
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

	t.Run("a server that needs setup", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setAPIKey(fakeAPIKey)
		fake.setSetupRequired()

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.ErrorIs(t, err, kuma.ErrAPIKeySetupRequired)
		require.NotErrorIs(t, err, kuma.ErrAPIKeyRejected)
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
		require.ErrorContains(t, err, "status 200", "the error carries what the scheme was detected from")
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
		require.ErrorContains(t, err, "pinned with WithAuthScheme")
		require.Nil(t, client)
		require.Zero(t, fake.probeCount(), "a pinned scheme is not detected")
	})

	t.Run("a pinned 3.x scheme meets a 2.x server with authentication disabled", func(t *testing.T) {
		fake, url := newAuthFake(t)
		fake.setAutoLogin()

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "", "",
			kuma.WithConnectTimeout(5*time.Second),
			kuma.WithAuthScheme(kuma.AuthSchemeHTTP),
			kuma.WithAPIKey(fakeAPIKey),
		)
		require.ErrorIs(t, err, kuma.ErrAPIKeyUnsupported)
		require.ErrorContains(t, err, "pinned with WithAuthScheme")
		require.Nil(t, client)
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
			name:    "autosetup",
			opts:    []kuma.Option{kuma.WithAutosetup()},
			wantErr: "WithAPIKey cannot be combined with WithAutosetup",
		},
		{
			name:    "an Authorization header",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Authorization": {"Basic Zm9vOmJhcg=="}})},
			wantErr: "WithAPIKey sets the Authorization header",
		},
		{
			name:    "a lower-case authorization header",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"authorization": {"Basic Zm9vOmJhcg=="}})},
			wantErr: "WithAPIKey sets the Authorization header",
		},
		{
			name:    "an empty Authorization header",
			opts:    []kuma.Option{kuma.WithHeaders(http.Header{"Authorization": {""}})},
			wantErr: "WithAPIKey sets the Authorization header",
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

// TestSameOrigin covers how a URL is told to belong to the server, which
// decides whether a request gets the configured headers and whether a
// redirect is followed.
func TestSameOrigin(t *testing.T) {
	tests := []struct {
		base, target string
		same         bool
	}{
		{base: "http://kuma.example.com", target: "http://kuma.example.com:80/socket.io/?EIO=4", same: true},
		{base: "https://Kuma.Example.com", target: "wss://kuma.example.com:443/socket.io/", same: true},
		{base: "http://kuma.example.com:3001", target: "ws://kuma.example.com:3001/socket.io/", same: true},
		{base: "http://kuma.example.com", target: "https://kuma.example.com/api/auth/ok"},
		{base: "http://kuma.example.com", target: "wss://kuma.example.com/socket.io/"},
		{base: "https://kuma.example.com", target: "http://kuma.example.com"},
		{base: "http://kuma.example.com:3001", target: "https://kuma.example.com:3001"},
		{base: "http://kuma.example.com:3001", target: "https://kuma.example.com"},
		{base: "http://kuma.example.com", target: "https://other.example.com"},
		{base: "http://kuma.example.com:3001", target: "http://kuma.example.com:3002"},
		{base: "http://kuma.example.com", target: "http://evil.example.com"},
		{base: "http://kuma.example.com", target: "http://kuma.example.com.evil.example"},
	}

	for _, tc := range tests {
		t.Run(tc.base+" "+tc.target, func(t *testing.T) {
			base, err := neturl.Parse(tc.base)
			require.NoError(t, err)

			target, err := neturl.Parse(tc.target)
			require.NoError(t, err)

			require.Equal(t, tc.same, kuma.SameOrigin(base, target))
		})
	}
}

// TestCheckRedirect covers the redirects the client follows: within the
// server's origin, nothing else.
func TestCheckRedirect(t *testing.T) {
	request := func(target string) *http.Request {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
		require.NoError(t, err)

		return req
	}

	tests := []struct {
		name      string
		base      string
		target    string
		sensitive bool
		wantErr   []string
	}{
		{
			name:   "within the origin",
			base:   "http://kuma.example.com",
			target: "http://kuma.example.com/dashboard",
		},
		{
			name:   "to another host",
			base:   "http://kuma.example.com",
			target: "http://elsewhere.example.com/api/auth/ok",
			wantErr: []string{
				"http://elsewhere.example.com/api/auth/ok",
				"use http://elsewhere.example.com as the base URL",
			},
		},
		{
			name:    "the upgrade to https",
			base:    "http://kuma.example.com",
			target:  "https://kuma.example.com/api/auth/ok",
			wantErr: []string{"upgrades to https", "use https://kuma.example.com as the base URL"},
		},
		{
			name:    "the upgrade keeps the base path",
			base:    "http://kuma.example.com/kuma",
			target:  "https://kuma.example.com/kuma/api/auth/ok",
			wantErr: []string{"upgrades to https", "use https://kuma.example.com/kuma as the base URL"},
		},
		{
			name:    "to a path that is not the request's",
			base:    "http://kuma.example.com",
			target:  "https://sso.example.com/sign_in",
			wantErr: []string{"connect to the server on https://sso.example.com"},
		},
		{
			name:    "the downgrade to http",
			base:    "https://kuma.example.com",
			target:  "http://kuma.example.com/api/auth/ok",
			wantErr: []string{"use http://kuma.example.com as the base URL"},
		},
		{
			name:    "a port change",
			base:    "http://kuma.example.com:3001",
			target:  "http://kuma.example.com:3002/api/auth/ok",
			wantErr: []string{"use http://kuma.example.com:3002 as the base URL"},
		},
		{
			name:      "a client with credentials is told why",
			base:      "http://kuma.example.com",
			target:    "http://elsewhere.example.com/api/auth/ok",
			sensitive: true,
			wantErr:   []string{"configured headers and API key stay with the server"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base, err := neturl.Parse(tc.base)
			require.NoError(t, err)

			via := []*http.Request{request(tc.base + "/api/auth/ok")}

			note := ""
			if tc.sensitive {
				note = kuma.SensitiveRedirectNote
			}

			err = kuma.CheckRedirect(base, note)(request(tc.target), via)
			if len(tc.wantErr) == 0 {
				require.NoError(t, err)

				return
			}

			require.ErrorContains(t, err, tc.base+"/api/auth/ok")

			for _, want := range tc.wantErr {
				require.ErrorContains(t, err, want)
			}

			if !tc.sensitive {
				require.NotContains(t, err.Error(), "API key")
			}
		})
	}

	t.Run("no password in the error", func(t *testing.T) {
		base, err := neturl.Parse("http://admin:s3cret@kuma.example.com")
		require.NoError(t, err)

		via := []*http.Request{request("http://admin:s3cret@kuma.example.com/api/auth/ok")}

		err = kuma.CheckRedirect(base, "")(request("https://admin:s3cret@kuma.example.com/api/auth/ok"), via)
		require.ErrorContains(t, err, "upgrades to https")
		require.NotContains(t, err.Error(), "s3cret")
	})

	t.Run("no more than ten redirects", func(t *testing.T) {
		base, err := neturl.Parse("http://kuma.example.com")
		require.NoError(t, err)

		check := kuma.CheckRedirect(base, "")

		via := make([]*http.Request, 0, 10)
		for range 9 {
			via = append(via, request("http://kuma.example.com/hop"))
		}

		require.NoError(t, check(request("http://kuma.example.com/hop"), via))

		via = append(via, request("http://kuma.example.com/hop"))
		require.ErrorContains(t, check(request("http://kuma.example.com/hop"), via), "stopped after 10 redirects")
	})
}

// TestBaseURLUserinfo covers a base URL with userinfo, whose Basic
// Authorization net/http sets before the configured headers are added.
func TestBaseURLUserinfo(t *testing.T) {
	for name, opt := range map[string]kuma.Option{
		"API key":              kuma.WithAPIKey(fakeAPIKey),
		"Authorization header": kuma.WithHeaders(http.Header{"Authorization": {"Basic Zm9vOmJhcg=="}}),
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newHTTPAuthFake(t)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			withUserinfo := strings.Replace(url, "http://", "http://user:secret@", 1)

			client, err := kuma.New(ctx, withUserinfo, "", "", kuma.WithConnectTimeout(5*time.Second), opt)
			require.ErrorContains(t, err, "userinfo")
			require.Nil(t, client)
			require.Empty(t, fake.receivedRequests(), "nothing is sent to the server")
		})
	}
}

// TestRedirectOffTheServersOrigin covers a redirect to another server: the
// client refuses to follow it, so the other server gets neither the API key
// nor a configured header, while the server the client was configured for
// gets all of them.
func TestRedirectOffTheServersOrigin(t *testing.T) {
	var mu sync.Mutex

	var originHeaders, otherHeaders []http.Header

	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		otherHeaders = append(otherHeaders, r.Header.Clone())
		mu.Unlock()

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html></html>"))
	}))

	t.Cleanup(other.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		originHeaders = append(originHeaders, r.Header.Clone())
		mu.Unlock()

		http.Redirect(w, r, other.URL+r.URL.RequestURI(), http.StatusFound)
	}))

	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	configured := http.Header{}
	configured.Set("CF-Access-Client-Secret", "proxy-secret")
	configured.Set("Cookie", "proxy=1")

	_, err := kuma.New(ctx, server.URL, "", "",
		kuma.WithConnectTimeout(2*time.Second),
		kuma.WithAPIKey(fakeAPIKey),
		kuma.WithHeaders(configured),
		kuma.WithUserAgent("custom-agent/1.0"),
	)
	require.ErrorContains(t, err, "leaves the server's origin")
	require.ErrorContains(t, err, other.URL+"/api/auth/ok")

	mu.Lock()
	defer mu.Unlock()

	require.NotEmpty(t, originHeaders)

	for _, header := range originHeaders {
		require.Equal(t, "Bearer "+fakeAPIKey, header.Get("Authorization"))
		require.Equal(t, "proxy-secret", header.Get("CF-Access-Client-Secret"))
		require.Equal(t, "proxy=1", header.Get("Cookie"))
		require.Equal(t, "custom-agent/1.0", header.Get("User-Agent"))
	}

	require.Empty(t, otherHeaders, "the other server sees no request at all")
}

// TestHeadersOffTheServersOrigin covers the round tripper itself for a
// request that leaves the server's origin without a redirect the client
// could refuse: it gets the User-Agent and nothing else.
func TestHeadersOffTheServersOrigin(t *testing.T) {
	base, err := neturl.Parse("http://kuma.example.com")
	require.NoError(t, err)

	configured := http.Header{}
	configured.Set("Authorization", "Bearer "+fakeAPIKey)
	configured.Set("CF-Access-Client-Secret", "proxy-secret")
	configured.Set("User-Agent", "custom-agent/1.0")

	for target, wantSecrets := range map[string]bool{
		"http://kuma.example.com/api/auth/ok":    true,
		"https://kuma.example.com/api/auth/ok":   false,
		"http://other.example.com/api/auth/ok":   false,
		"http://kuma.example.com:81/api/auth/ok": false,
	} {
		t.Run(target, func(t *testing.T) {
			req, reqErr := http.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
			require.NoError(t, reqErr)

			header := kuma.ApplyRequestHeaders(base, configured, req)

			require.Equal(t, "custom-agent/1.0", header.Get("User-Agent"))

			if wantSecrets {
				require.Equal(t, "Bearer "+fakeAPIKey, header.Get("Authorization"))
				require.Equal(t, "proxy-secret", header.Get("CF-Access-Client-Secret"))

				return
			}

			require.Empty(t, header.Get("Authorization"))
			require.Empty(t, header.Get("CF-Access-Client-Secret"))
		})
	}
}

// TestProbeStatus covers which answers of the auth scheme probe New reads as a
// 2.x server, and which it reports instead of guessing: JSON means 3.x only on
// a 2xx, and JSON with any other status is what a proxy or gateway in front of
// the server answers.
func TestProbeStatus(t *testing.T) {
	tests := []struct {
		status      int
		contentType string
		wantErr     bool
	}{
		{status: http.StatusOK, contentType: "text/html"},
		{status: http.StatusNotFound, contentType: "text/html"},
		{status: http.StatusTooManyRequests, contentType: "text/html"},
		{status: http.StatusUnauthorized, contentType: "text/html", wantErr: true},
		{status: http.StatusForbidden, contentType: "text/html", wantErr: true},
		{status: http.StatusFound, contentType: "text/html", wantErr: true},
		{status: http.StatusBadGateway, contentType: "text/html", wantErr: true},
		{status: http.StatusUnauthorized, contentType: "application/json", wantErr: true},
		{status: http.StatusForbidden, contentType: "application/json", wantErr: true},
		{status: http.StatusNotFound, contentType: "application/json", wantErr: true},
		{status: http.StatusTooManyRequests, contentType: "application/json", wantErr: true},
		{status: http.StatusBadGateway, contentType: "application/json", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(strconv.Itoa(tc.status)+" "+tc.contentType, func(t *testing.T) {
			fake := &fakeAckServer{messages: make(chan []byte, 32)}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/auth/ok" {
					fake.ServeHTTP(w, r)

					return
				}

				// A 302 without a Location is not followed.
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte("{}"))
			}))

			t.Cleanup(func() {
				server.CloseClientConnections()
				server.Close()
			})

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			client, err := kuma.New(ctx, server.URL, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))

			if tc.wantErr {
				require.ErrorContains(t, err, "status "+strconv.Itoa(tc.status))
				require.ErrorContains(t, err, "Content-Type "+strconv.Quote(tc.contentType))
				require.ErrorContains(t, err, server.URL+"/api/auth/ok")
				require.ErrorContains(t, err, "WithHeaders")
				require.Nil(t, client)

				return
			}

			require.NoError(t, err)

			t.Cleanup(func() { _ = client.Disconnect() })

			require.Equal(t, kuma.AuthSchemeSocket, client.AuthScheme())
		})
	}
}

// TestProbeRedirectedWithinTheOrigin covers a probe that a proxy in front of
// the server redirects to its sign-in page on the same host: neither Uptime
// Kuma line redirects the probe, so New reports it instead of reading the
// sign-in page as a 2.x server.
func TestProbeRedirectedWithinTheOrigin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth/ok" {
			http.Redirect(w, r, "/oauth2/sign_in", http.StatusFound)

			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html></html>"))
	}))

	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, server.URL, "", "",
		kuma.WithConnectTimeout(2*time.Second),
		kuma.WithAPIKey(fakeAPIKey),
	)
	require.ErrorContains(t, err, "was redirected to "+server.URL+"/oauth2/sign_in")
	require.NotErrorIs(t, err, kuma.ErrAPIKeyUnsupported)
	require.Nil(t, client)
}

// TestProbeErrorHidesPassword covers a base URL with userinfo, for HTTP basic
// auth to a reverse proxy: an error that names the probe's URL does not show
// the password.
func TestProbeErrorHidesPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadGateway)
	}))

	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	withUserinfo := strings.Replace(server.URL, "http://", "http://admin:s3cret@", 1)

	client, err := kuma.New(ctx, withUserinfo, "admin", "admin1", kuma.WithConnectTimeout(2*time.Second))
	require.ErrorContains(t, err, "status 502")
	require.ErrorContains(t, err, "admin:xxxxx@")
	require.NotContains(t, err.Error(), "s3cret")
	require.Nil(t, client)
}
