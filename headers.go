package kuma

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/net/http/httpguts"
)

// modulePath is this module's import path, which the default User-Agent looks
// its own version up by, see defaultUserAgent.
const modulePath = "github.com/breml/go-uptime-kuma-client"

// userAgentProduct is the product token of the default User-Agent.
const userAgentProduct = "go-uptime-kuma-client"

// devVersion is the version the default User-Agent names when the build
// carries none for this module: a build of the module itself, a replace
// directive pointing at a local checkout, or a binary built without module
// information.
const devVersion = "dev"

// cookieHeaderName is the header the cookie jar fills, which configured
// headers add to rather than replace, see WithHeaders.
const cookieHeaderName = "Cookie"

// authorizationHeaderName is the header WithAPIKey sets.
const authorizationHeaderName = "Authorization"

// userAgentHeaderName is the header WithUserAgent sets.
const userAgentHeaderName = "User-Agent"

// WithHeaders adds HTTP headers to every request the client makes: the auth
// scheme probe, the HTTP login of a 3.x server, the socket.io long-polling
// requests and the WebSocket handshake. It is for what sits in front of the
// server, such as an access proxy that wants a header of its own
// (Cloudflare Access takes CF-Access-Client-Id and CF-Access-Client-Secret).
//
// The headers go to the server's origin (scheme, host and port of the base
// URL) only. The client follows no redirect that leaves that origin, with or
// without configured headers, an upgrade of the same host from http to https
// included: a redirect turns the POSTs of the login and the long polling into
// GETs without a body, and the WebSocket handshake follows none. New reports
// such a redirect with the base URL to use instead, as far as the redirect
// tells it.
//
// The headers are copied when the option is called, so changing h afterwards
// changes nothing. Names are canonicalised; values that one h holds under two
// spellings of a name are kept, in the sorted order of the spellings. Calling
// the option more than once adds up; a header named again replaces the values
// it had.
//
// New rejects what would break the protocol or the transport: a name or value
// that is not a valid HTTP header field, and the headers the client, net/http
// or the WebSocket handshake set themselves (Host, Origin, Connection,
// Upgrade, Accept, Accept-Encoding, Content-Type, Content-Length,
// Transfer-Encoding, TE, Trailer and Sec-WebSocket-*). A configured Cookie is
// appended to the cookies of the client's own jar rather than replacing them,
// so the session cookie of a 3.x server still gets through. Authorization is
// accepted for a scheme other than Bearer, for a reverse proxy in front of the
// server that wants HTTP basic auth; a Bearer token belongs in WithAPIKey,
// and an Authorization header next to WithAPIKey is an error.
func WithHeaders(h http.Header) Option {
	names := make([]string, 0, len(h))
	for name := range h {
		names = append(names, name)
	}

	slices.Sort(names)

	copied := http.Header{}

	for _, name := range names {
		canonical := http.CanonicalHeaderKey(name)
		copied[canonical] = append(copied[canonical], h[name]...)
	}

	return func(c *Client) {
		if c.headers == nil {
			c.headers = http.Header{}
		}

		for name, values := range copied {
			c.headers[name] = slices.Clone(values)
		}
	}
}

// WithUserAgent sets the User-Agent of every request the client makes, see
// WithHeaders. It defaults to "go-uptime-kuma-client/<version>", with the
// version of this module the program was built with, or "dev" where the build
// does not record one. It takes precedence over a User-Agent given to
// WithHeaders.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithAPIKey authenticates with an API key of an Uptime Kuma 3.x server
// instead of a login, presented as "Authorization: Bearer <key>" on every
// request the client makes, see WithHeaders.
//
// The server reads the key from the socket.io handshake and authenticates
// the connection with it, provided "API and MCP" (the apiEnabled setting) is
// turned on in its settings. New takes a session the server answers the
// handshake with as the key being accepted. It does not tell a key it does
// not know from a key it ignores because that setting is off: either is
// reported as ErrAPIKeyRejected. A 3.x server with authentication disabled
// ignores the key altogether and opens a session for any connection, so there
// even an invalid key "succeeds".
//
// A 2.x server has no API key login on the socket.io API. New fails with
// ErrAPIKeyUnsupported when the scheme is detected as, or pinned to,
// AuthSchemeSocket, and when a pinned AuthSchemeHTTP meets a 2.x server with
// authentication disabled, which logs the client in by itself (autoLogin).
// Pinning AuthSchemeHTTP against a 2.x server that wants a login surfaces as
// ErrAPIKeyRejected instead, because that server asks for a login the way a
// 3.x server refusing the key does.
//
// A 3.x server that is not set up yet fails New with ErrAPIKeySetupRequired;
// a 2.x one fails with ErrAPIKeyUnsupported, as every 2.x server does.
//
// The key is the only credential such a client has: New rejects it together
// with a username and password, WithSessionToken, WithTOTPSecret or
// WithTOTPCode, because the server would authenticate the connection by one
// of them and silently ignore the rest. It rejects WithAutosetup too: a server
// that is not set up has no API keys, and setting it up takes a username and
// password. The key is passed without the "Bearer " prefix, and New rejects a
// key with whitespace or control characters in it.
//
// WithAPIKey("") is a no-op: the client authenticates exactly as if the option
// had not been given.
//
// No session is created, so Client.SessionToken returns the empty string;
// Client.Resync connects again and presents the key once more. The server
// authenticates the connection as the user who created the key, but some of
// its commands that act on that user's account (two-factor authentication,
// creating API keys) want a login session and fail on such a connection.
//
// Security: the key is a bearer credential that is valid until it expires or
// is revoked in the server's settings. Store it the way a password would be
// stored.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// ErrAPIKeyUnsupported is returned by New for a client configured with
// WithAPIKey against an Uptime Kuma 2.x server, see AuthSchemeSocket, which
// authenticates the socket.io API by a login only.
var ErrAPIKeyUnsupported = errors.New(
	"api key: Uptime Kuma authenticates the socket.io API with an API key only from 3.0 on",
)

// ErrAPIKeySetupRequired is returned by New when a 3.x server that is not set
// up yet answers a client configured with WithAPIKey: such a server has no API
// keys, and its setup takes a username and password.
var ErrAPIKeySetupRequired = errors.New(
	"api key: the server needs setup, which an API key cannot run: " +
		"set the server up with a username and password first",
)

// ErrAPIKeyRejected is returned when an Uptime Kuma 3.x server asks a client
// configured with WithAPIKey for a login instead of authenticating it by the
// key. The server answers both causes the same way: the key is invalid,
// expired or revoked, or "API and MCP" (apiEnabled) is turned off in its
// settings.
var ErrAPIKeyRejected = errors.New(
	`api key rejected: the key is invalid, expired or revoked, or "API and MCP" ` +
		"(apiEnabled) is disabled in the server's settings",
)

// buildRequestHeaders assembles the headers of every request, see
// requestHeaders, and binds them to the origin of base, the server's HTTP
// base URL.
//
// A base URL with userinfo is rejected next to an API key or a configured
// Authorization header: net/http turns the userinfo into a Basic
// Authorization of its own before the headers are added, so the key would
// reach the server on the WebSocket handshake but not on the HTTP requests.
func (c *Client) buildRequestHeaders(base *url.URL) error {
	if base.User != nil && (c.apiKey != "" || len(c.headers.Values(authorizationHeaderName)) > 0) {
		return errors.New(
			"base url: userinfo (user:password@) sets an Authorization header of its own, " +
				"which cannot be combined with WithAPIKey or an Authorization header from WithHeaders",
		)
	}

	headers, err := c.requestHeaders()
	if err != nil {
		return err
	}

	c.requestHeaderSet = newRequestHeaderSet(base, headers)

	return nil
}

// requestHeaders assembles the headers New configured: the default
// User-Agent, what WithHeaders added, the User-Agent of WithUserAgent, and the
// Authorization of WithAPIKey.
//
// It rejects what WithHeaders and WithAPIKey document as invalid.
func (c *Client) requestHeaders() (http.Header, error) {
	err := validateHeaders(c.headers)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set(userAgentHeaderName, defaultUserAgent())

	for name, values := range c.headers {
		headers[name] = slices.Clone(values)
	}

	if c.userAgent != "" {
		if !httpguts.ValidHeaderFieldValue(c.userAgent) {
			return nil, errors.New("headers: WithUserAgent: not a valid header value")
		}

		headers.Set(userAgentHeaderName, c.userAgent)
	}

	if c.apiKey != "" {
		err = validateAPIKey(c.apiKey)
		if err != nil {
			return nil, err
		}

		if len(headers.Values(authorizationHeaderName)) > 0 {
			return nil, errors.New("headers: WithAPIKey sets the Authorization header that WithHeaders sets too")
		}

		headers.Set(authorizationHeaderName, "Bearer "+c.apiKey)
	}

	return headers, nil
}

// reservedHeaders are the headers the client, net/http or the WebSocket
// handshake set themselves, which WithHeaders may not set, see there. The
// Sec-WebSocket-* headers are matched by their prefix.
func reservedHeaders() []string {
	return []string{
		"Host",
		"Origin",
		"Connection",
		"Upgrade",
		"Accept",
		"Accept-Encoding",
		"Content-Type",
		"Content-Length",
		"Transfer-Encoding",
		"Te",
		"Trailer",
	}
}

// validateHeaders rejects the headers WithHeaders documents as invalid.
func validateHeaders(headers http.Header) error {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}

	slices.Sort(names)

	for _, name := range names {
		if !httpguts.ValidHeaderFieldName(name) {
			return fmt.Errorf("headers: %q is not a valid header name", name)
		}

		if slices.Contains(reservedHeaders(), name) ||
			strings.HasPrefix(strings.ToLower(name), "sec-websocket-") {
			return fmt.Errorf("headers: %s is set by the client or its transport and cannot be configured", name)
		}

		for _, value := range headers[name] {
			if !httpguts.ValidHeaderFieldValue(value) {
				return fmt.Errorf("headers: %s: not a valid header value", name)
			}

			if name == authorizationHeaderName && hasBearerPrefix(value) {
				return errors.New("headers: Authorization: a Bearer token is an API key, use WithAPIKey")
			}
		}
	}

	return nil
}

// validateAPIKey rejects a key WithAPIKey documents as invalid.
func validateAPIKey(key string) error {
	if hasBearerPrefix(key) {
		return errors.New("api key: pass the key without the Bearer prefix")
	}

	if strings.IndexFunc(key, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) >= 0 {
		return errors.New("api key: the key contains whitespace or control characters")
	}

	if !httpguts.ValidHeaderFieldValue("Bearer " + key) {
		return errors.New("api key: not a valid header value")
	}

	return nil
}

// hasBearerPrefix reports whether value starts with the Bearer scheme, in any
// case.
func hasBearerPrefix(value string) bool {
	const prefix = "bearer "

	return len(value) >= len(prefix) && strings.EqualFold(value[:len(prefix)], prefix)
}

// defaultUserAgent is the User-Agent a client sends unless WithUserAgent or
// WithHeaders configures another.
func defaultUserAgent() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		info = nil
	}

	return userAgentProduct + "/" + moduleVersion(info)
}

// moduleVersion is the version of this module the build information records,
// or devVersion where it records none: in a build of the module itself, which
// reports "(devel)", and for a dependency replaced by a local directory, which
// carries no version. A nil info is a build without module information.
func moduleVersion(info *debug.BuildInfo) string {
	if info == nil {
		return devVersion
	}

	module := &info.Main
	if module.Path != modulePath {
		module = nil

		for _, dep := range info.Deps {
			if dep.Path == modulePath {
				module = dep

				break
			}
		}
	}

	if module == nil {
		return devVersion
	}

	if module.Replace != nil {
		module = module.Replace
	}

	if module.Version == "" || module.Version == "(devel)" {
		return devVersion
	}

	return module.Version
}

// requestHeaderSet is the headers New assembled, see requestHeaders, bound to
// the origin of the server they are meant for, see sameOrigin.
//
// The configured headers carry credentials (the API key, an access proxy's
// secret), so they go to that origin alone: net/http strips Authorization
// from a redirect that leaves the original host, and a round tripper that
// added it back on every request would hand the key to wherever the server,
// or anything in between, redirects. A request to another origin gets the
// User-Agent and nothing else.
type requestHeaderSet struct {
	base    *url.URL
	headers http.Header
}

// newRequestHeaderSet binds headers to the origin of base.
func newRequestHeaderSet(base *url.URL, headers http.Header) *requestHeaderSet {
	return &requestHeaderSet{base: base, headers: headers}
}

// apply adds the headers meant for target's origin to header, see
// applyHeaders.
func (s *requestHeaderSet) apply(header http.Header, target *url.URL) {
	if s == nil {
		return
	}

	if sameOrigin(s.base, target) {
		applyHeaders(header, s.headers)

		return
	}

	if header.Get(userAgentHeaderName) == "" && s.headers.Get(userAgentHeaderName) != "" {
		header.Set(userAgentHeaderName, s.headers.Get(userAgentHeaderName))
	}
}

// sameOrigin reports whether target belongs to the server at base: the same
// scheme, host and port, see origin.
func sameOrigin(base *url.URL, target *url.URL) bool {
	return origin(base) == origin(target)
}

// checkRedirect is the redirect policy of the client's http.Client: it
// follows a redirect only within the server's origin, see sameOrigin, and no
// more than ten of them, as net/http does by default. A non-empty note is
// added to the error, see redirectNote.
func checkRedirect(base *url.URL, note string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRedirects)
		}

		if sameOrigin(base, req.URL) {
			return nil
		}

		from := via[len(via)-1].URL
		suggestion := suggestBaseURL(base, via[0].URL, req.URL)

		msg := fmt.Sprintf("redirect from %s to %s leaves the server's origin; %s",
			from.Redacted(), req.URL.Redacted(), suggestion)

		if strings.EqualFold(from.Hostname(), req.URL.Hostname()) &&
			from.Scheme == "http" && req.URL.Scheme == "https" {
			msg = fmt.Sprintf("redirect from %s to %s upgrades to https; %s",
				from.Redacted(), req.URL.Redacted(), suggestion)
		}

		if note != "" {
			msg += ", " + note
		}

		return errors.New(msg)
	}
}

// suggestBaseURL tells what to use instead of base after a redirect from the
// client's first request to target: target without the part of the path the
// client appended to base, or, for a redirect that does not keep that part,
// target's origin. The URLs it prints carry no userinfo.
func suggestBaseURL(base *url.URL, first *url.URL, target *url.URL) string {
	origin := target.Scheme + "://" + target.Host

	basePath := strings.TrimSuffix(base.Path, "/")
	if !strings.HasPrefix(first.Path, basePath) {
		return "connect to the server on " + origin
	}

	appended := strings.TrimPrefix(first.Path, basePath)
	if !strings.HasSuffix(target.Path, appended) {
		return "connect to the server on " + origin
	}

	return "use " + origin + strings.TrimSuffix(strings.TrimSuffix(target.Path, appended), "/") + " as the base URL"
}

// sensitiveRedirectNote is what a refused redirect adds for a client with
// configured headers or an API key, see redirectNote.
const sensitiveRedirectNote = "the redirect is not followed so that the configured headers and API key " +
	"stay with the server"

// redirectNote is the note checkRedirect adds to a refused redirect: why it
// matters for a client with configured headers or an API key, and nothing for
// one with neither.
func (c *Client) redirectNote() string {
	if len(c.headers) > 0 || c.apiKey != "" {
		return sensitiveRedirectNote
	}

	return ""
}

// maxRedirects is how many redirects net/http follows by default.
const maxRedirects = 10

// origin renders the scheme, host and port of u the way two URLs of one
// server compare equal: a ws or wss scheme as the http or https it is made
// over, host names in lower case, and the default port spelled out.
func origin(u *url.URL) string {
	scheme := strings.ToLower(u.Scheme)

	switch scheme {
	case "ws":
		scheme = "http"

	case "wss":
		scheme = "https"

	default:
	}

	port := u.Port()
	if port == "" {
		switch scheme {
		case "http":
			port = "80"

		case "https":
			port = "443"

		default:
		}
	}

	return scheme + "://" + net.JoinHostPort(strings.ToLower(u.Hostname()), port)
}

// headerRoundTripper adds the configured headers to the requests that go out
// through the client's shared http.Client to the server, see WithHeaders and
// requestHeaderSet.
type headerRoundTripper struct {
	base    http.RoundTripper
	headers *requestHeaderSet
}

// RoundTrip sends req with the configured headers added to a copy of it, as a
// RoundTripper must not change the request it is given.
func (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	t.headers.apply(clone.Header, clone.URL)

	return t.base.RoundTrip(clone) //nolint:wrapcheck // a transport passes its errors on as they are.
}

// applyHeaders adds configured to target: a header target carries already is
// left as it is, except for Cookie, whose configured values are appended to
// the ones the cookie jar put there.
func applyHeaders(target http.Header, configured http.Header) {
	for name, values := range configured {
		if name == cookieHeaderName {
			cookies := nonEmpty(append(target.Values(cookieHeaderName), values...))
			if len(cookies) > 0 {
				target.Set(cookieHeaderName, strings.Join(cookies, "; "))
			}

			continue
		}

		if len(target.Values(name)) > 0 {
			continue
		}

		target[name] = append([]string(nil), values...)
	}
}

// nonEmpty returns values without the empty strings.
func nonEmpty(values []string) []string {
	result := make([]string, 0, len(values))

	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
