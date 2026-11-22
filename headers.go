package kuma

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"
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

// WithHeaders adds HTTP headers to every request the client makes: the auth
// scheme probe, the HTTP login of a 3.x server, the socket.io long-polling
// requests and the WebSocket handshake. It is for what sits in front of the
// server, such as an access proxy that wants a header of its own
// (Cloudflare Access takes CF-Access-Client-Id and CF-Access-Client-Secret).
//
// The headers are copied, so changing h afterwards changes nothing. Calling
// the option more than once adds up; a header named again replaces the
// values it had.
//
// Headers the client sets itself for a request take precedence, so a
// configured header cannot break the protocol, with one exception: a
// configured Cookie is appended to the cookies of the client's own jar rather
// than replacing them, so the session cookie of a 3.x server still gets
// through. An Authorization header conflicts with WithAPIKey, which New
// reports as an error.
func WithHeaders(h http.Header) Option {
	return func(c *Client) {
		if c.headers == nil {
			c.headers = http.Header{}
		}

		for name, values := range h {
			c.headers[http.CanonicalHeaderKey(name)] = append([]string(nil), values...)
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
// turned on in its settings. It does not tell a key it does not know from a
// key it ignores because that setting is off: either is reported as
// ErrAPIKeyRejected. A 2.x server has no API key login on the socket.io API,
// so New fails with ErrAPIKeyUnsupported against one, whether the scheme was
// detected or pinned with WithAuthScheme.
//
// The key is the only credential such a client has: New rejects it together
// with a username and password, WithSessionToken, WithTOTPSecret or
// WithTOTPCode, because the server would authenticate the connection by one
// of them and silently ignore the rest. An empty key configures nothing.
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

// ErrAPIKeyRejected is returned when an Uptime Kuma 3.x server asks a client
// configured with WithAPIKey for a login instead of authenticating it by the
// key. The server answers both causes the same way: the key is invalid,
// expired or revoked, or "API and MCP" (apiEnabled) is turned off in its
// settings.
var ErrAPIKeyRejected = errors.New(
	`api key rejected: the key is invalid, expired or revoked, or "API and MCP" ` +
		"(apiEnabled) is disabled in the server's settings",
)

// requestHeaders assembles the headers New configured: the default
// User-Agent, what WithHeaders added, the User-Agent of WithUserAgent, and the
// Authorization of WithAPIKey.
func (c *Client) requestHeaders() (http.Header, error) {
	headers := http.Header{}
	headers.Set("User-Agent", defaultUserAgent())

	for name, values := range c.headers {
		headers[name] = append([]string(nil), values...)
	}

	if c.userAgent != "" {
		headers.Set("User-Agent", c.userAgent)
	}

	if c.apiKey != "" {
		if headers.Get(authorizationHeaderName) != "" {
			return nil, errors.New("headers: WithAPIKey sets the Authorization header that WithHeaders sets too")
		}

		headers.Set(authorizationHeaderName, "Bearer "+c.apiKey)
	}

	return headers, nil
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

// headerRoundTripper adds the configured headers to every request that goes
// out through the client's shared http.Client, see WithHeaders.
type headerRoundTripper struct {
	base    http.RoundTripper
	headers http.Header
}

// RoundTrip sends req with the configured headers added to a copy of it, as a
// RoundTripper must not change the request it is given.
func (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	applyHeaders(clone.Header, t.headers)

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
