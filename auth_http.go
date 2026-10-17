package kuma

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// AuthScheme is the way an Uptime Kuma server authenticates a client, which
// changed between the major versions, see WithAuthScheme.
type AuthScheme int

const (
	// AuthSchemeAuto lets New detect the scheme, which is the default.
	AuthSchemeAuto AuthScheme = iota

	// AuthSchemeSocket is the login of Uptime Kuma 2.x: a login command on the
	// socket.io connection, answered with a session token that is presented
	// the same way.
	AuthSchemeSocket

	// AuthSchemeHTTP is the login of Uptime Kuma 3.x: a login over HTTP that
	// answers with a session cookie, which the socket.io connection then
	// carries on its handshake. The server authenticates the connection from
	// that cookie alone and has no login command any more.
	AuthSchemeHTTP
)

// String names the scheme.
func (s AuthScheme) String() string {
	switch s {
	case AuthSchemeAuto:
		return "auto"

	case AuthSchemeSocket:
		return "socket"

	case AuthSchemeHTTP:
		return "http"

	default:
		return "unknown"
	}
}

// WithAuthScheme pins the way the client authenticates instead of letting New
// detect it.
//
// New detects the scheme with one request to the server before it connects,
// see Client.AuthScheme, and the detection is right for every Uptime Kuma
// release. Pinning it is for a proxy in front of the server that answers that
// request itself, and for tests. A pinned scheme the server does not speak
// fails the login rather than the detection.
func WithAuthScheme(scheme AuthScheme) Option {
	return func(c *Client) {
		c.authScheme = scheme
	}
}

// ErrTwoFactorLocked is returned when the server refuses to verify any further
// one-time codes for a login because too many wrong ones were sent, see
// ErrInvalidTOTPCode. An Uptime Kuma 3.x server locks the login after a
// handful of wrong codes; the way out is to start a new login.
var ErrTwoFactorLocked = errors.New("two-factor authentication locked after too many wrong codes")

// RateLimitedError is the error a login refused by the server's rate limiter
// is reported with. It wraps ErrRateLimited, so errors.Is keeps matching on
// that, and carries how long the server asked the client to wait where the
// server said so.
//
// An Uptime Kuma 3.x server allows three logins per ten seconds per client
// address and ten of any other authentication request per minute, and answers
// the ones beyond that with a Retry-After. A 2.x server, see ErrRateLimited,
// states no wait.
type RateLimitedError struct {
	// RetryAfter is how long the server asked the client to wait before the
	// next attempt, or zero where it did not say.
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	if e.RetryAfter <= 0 {
		return ErrRateLimited.Error()
	}

	return fmt.Sprintf("%s, retry after %s", ErrRateLimited, e.RetryAfter)
}

func (*RateLimitedError) Unwrap() error {
	return ErrRateLimited
}

// sessionCookieName is the cookie an Uptime Kuma 3.x server hands out on a
// successful login and authenticates the socket.io connection from. It is
// better-auth's default, and the server configures no other.
const sessionCookieName = "better-auth.session_token"

// secureSessionCookieName is the name better-auth uses for the same cookie when
// it is configured for secure cookies, which an Uptime Kuma 3.x server is not,
// see WithSessionToken. A preset token is stored under both names, so a server
// that starts to use it does not refuse a token that is fine.
const secureSessionCookieName = "__Secure-" + sessionCookieName

// betterAuthProbeTimeout bounds the request that detects the auth scheme, on
// its own next to the caller's deadline, because New has not started the
// connect it could charge it to yet.
const betterAuthProbeTimeout = 15 * time.Second

// better-auth's error codes the login classifies, as its JSON error bodies
// name them.
const (
	betterAuthCodeInvalidCredentials = "INVALID_USERNAME_OR_PASSWORD"
	betterAuthCodeInvalidCode        = "INVALID_CODE"
	betterAuthCodeInvalidTwoFactor   = "INVALID_TWO_FACTOR_COOKIE"
	betterAuthCodeTooManyAttempts    = "TOO_MANY_ATTEMPTS"
	betterAuthCodeBannedUser         = "BANNED_USER"
)

// betterAuthClient speaks to the HTTP authentication endpoints of an Uptime
// Kuma 3.x server, see AuthSchemeHTTP. The cookies the server sets are kept
// in the jar the socket.io transports read, which is how a login here
// authenticates the connection made afterwards.
type betterAuthClient struct {
	base *url.URL
	http *http.Client
	jar  http.CookieJar
}

// betterAuthError is the JSON body better-auth answers a refused request with.
type betterAuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// signInResponse is the body of an accepted sign-in or one-time code.
type signInResponse struct {
	Token             string `json:"token"`
	TwoFactorRedirect bool   `json:"twoFactorRedirect"`
}

// sessionResponse is the body of a get-session, null when there is none.
type sessionResponse struct {
	User sessionUser `json:"user"`
}

// sessionUser is the account a session belongs to, as far as the client reads
// it.
type sessionUser struct {
	TwoFactorEnabled bool `json:"twoFactorEnabled"`
}

// setupResponse is the body of a setup, which the server answers like its
// socket.io acks.
type setupResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// twoFactorEnableResponse is the body of a two-factor/enable.
type twoFactorEnableResponse struct {
	TOTPURI string `json:"totpURI"`
}

// detectAuthScheme asks the server how it authenticates, see AuthSchemeHTTP.
//
// The probe is better-auth's health endpoint, which a 3.x server answers with
// JSON, as it answers every request under /api/auth/ - its refusals and its
// rate limiting included. A 2.x server has no such route and answers the
// request with the single-page application its catch-all serves for any
// unknown path, as HTML. The media type is therefore what tells them apart,
// not the status: a 429 from the rate limiter is as much a 3.x answer as the
// 200.
func (a *betterAuthClient) detectAuthScheme(ctx context.Context) (AuthScheme, error) {
	ctx, cancel := context.WithTimeout(ctx, betterAuthProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.endpoint("/api/auth/ok"), http.NoBody)
	if err != nil {
		return AuthSchemeAuto, fmt.Errorf("connect to server: detect auth scheme: %w", err)
	}

	resp, err := a.http.Do(req)
	if err != nil {
		return AuthSchemeAuto, fmt.Errorf("connect to server: detect auth scheme: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// The body is read so the connection can be reused, its content does not
	// matter.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))

	if isJSON(resp.Header.Get("Content-Type")) {
		return AuthSchemeHTTP, nil
	}

	return AuthSchemeSocket, nil
}

// isJSON reports whether a Content-Type names a JSON media type.
func isJSON(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}

// signIn logs in with a username and password and reports whether the server
// wants a one-time code before it hands out a session, see verifyTOTP. A
// session the server does hand out is in the jar afterwards.
func (a *betterAuthClient) signIn(ctx context.Context, username string, password string) (bool, error) {
	var response signInResponse

	err := a.post(ctx, "/api/auth/sign-in/username", map[string]any{
		"username":   username,
		"password":   password,
		"rememberMe": true,
	}, &response)
	if err != nil {
		return false, err
	}

	return response.TwoFactorRedirect, nil
}

// verifyTOTP answers the server's request for a one-time code. The sign-in
// that asked for it left a cookie in the jar that names the pending login.
func (a *betterAuthClient) verifyTOTP(ctx context.Context, code string) error {
	var response signInResponse

	return a.post(ctx, "/api/auth/two-factor/verify-totp", map[string]any{"code": code}, &response)
}

// getSession reports whether the session cookie in the jar names a session the
// server still accepts, and whether that session's account has two-factor
// authentication enabled.
func (a *betterAuthClient) getSession(ctx context.Context) (ok bool, twoFactor bool, err error) {
	body, err := a.do(ctx, http.MethodGet, "/api/auth/get-session", nil)
	if err != nil {
		return false, false, err
	}

	if bytes.Equal(bytes.TrimSpace(body), []byte("null")) || len(bytes.TrimSpace(body)) == 0 {
		return false, false, nil
	}

	var response sessionResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, false, fmt.Errorf("get-session: parse response: %w", err)
	}

	return true, response.User.TwoFactorEnabled, nil
}

// enableTwoFactor asks the server for a new two-factor secret and returns it.
// The server turns two-factor authentication on only once verifyTOTP confirmed
// a code for it.
func (a *betterAuthClient) enableTwoFactor(ctx context.Context, password string) (string, error) {
	var response twoFactorEnableResponse

	err := a.post(ctx, "/api/auth/two-factor/enable", map[string]any{"password": password}, &response)
	if err != nil {
		return "", err
	}

	uri, err := url.Parse(response.TOTPURI)
	if err != nil {
		return "", fmt.Errorf("two-factor/enable: parse otpauth uri: %w", err)
	}

	secret := uri.Query().Get("secret")
	if secret == "" {
		return "", errors.New("two-factor/enable: otpauth uri carries no secret")
	}

	return secret, nil
}

// disableTwoFactor turns two-factor authentication off for the logged in
// account.
func (a *betterAuthClient) disableTwoFactor(ctx context.Context, password string) error {
	return a.post(ctx, "/api/auth/two-factor/disable", map[string]any{"password": password}, nil)
}

// sessionCookie returns the session cookie the jar holds for the server, or
// the empty string.
func (a *betterAuthClient) sessionCookie() string {
	var value string

	for _, cookie := range a.jar.Cookies(a.base) {
		switch cookie.Name {
		case secureSessionCookieName:
			return cookie.Value

		case sessionCookieName:
			value = cookie.Value

		default:
		}
	}

	return value
}

// presetSessionCookie puts a session token from an earlier login into the jar,
// see WithSessionToken.
func (a *betterAuthClient) presetSessionCookie(token string) {
	a.jar.SetCookies(a.base, sessionCookies(token, 0))
}

// clearSessionCookie drops the session cookie from the jar, for a token the
// server refused and a password login is about to replace.
func (a *betterAuthClient) clearSessionCookie() {
	a.jar.SetCookies(a.base, sessionCookies("", -1))
}

// sessionCookies builds the session cookie under both of its names, see
// secureSessionCookieName. The jar sends a cookie by name, path and host alone;
// the attributes are what the server would set, with Secure only where the
// name demands it.
func sessionCookies(token string, maxAge int) []*http.Cookie {
	plain := &http.Cookie{ //nolint:gosec // an http server sets the plain one without Secure.
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	secure := *plain //nolint:gosec // the copy gets the Secure attribute below.
	secure.Name = secureSessionCookieName
	secure.Secure = true

	return []*http.Cookie{plain, &secure}
}

// endpoint returns the URL of a path on the server.
func (a *betterAuthClient) endpoint(path string) string {
	return a.base.String() + path
}

// post sends a JSON body and decodes the JSON answer into response, which may
// be nil for a caller that only needs to know the request was accepted.
func (a *betterAuthClient) post(ctx context.Context, path string, body map[string]any, response any) error {
	data, err := a.do(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}

	if response == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}

	err = json.Unmarshal(data, response)
	if err != nil {
		return fmt.Errorf("%s: parse response: %w", path, err)
	}

	return nil
}

// do performs a request against the server and returns the body of an
// accepted one. A refused request is reported as the package's sentinel
// errors where the server's answer is one the login classifies, see
// betterAuthStatusError.
func (a *betterAuthClient) do(ctx context.Context, method string, path string, body map[string]any) ([]byte, error) {
	var reader io.Reader = http.NoBody

	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s: encode request: %w", path, err)
		}

		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, a.endpoint(path), reader)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	// better-auth refuses a state-changing request of a logged in session
	// without an Origin it trusts, which the server configures as any. A
	// browser sends its own; the client sends the server's.
	req.Header.Set("Origin", a.base.Scheme+"://"+a.base.Host)

	resp, err := a.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%s: read response: %w", path, err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return data, nil
	}

	return nil, betterAuthStatusError(path, resp, data)
}

// betterAuthStatusError translates a refused request into one of the package's
// sentinel errors, so a caller can tell the cases apart with errors.Is instead
// of matching on the message.
//
// A 2.x server never answers these paths, so a refusal that is not JSON is
// what a proxy in front of the server, or a pinned scheme the server does not
// speak, produces, and is reported with the status alone.
func betterAuthStatusError(path string, resp *http.Response, data []byte) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		return &RateLimitedError{RetryAfter: retryAfter(resp.Header)}
	}

	var body betterAuthError

	if !isJSON(resp.Header.Get("Content-Type")) || json.Unmarshal(data, &body) != nil {
		return fmt.Errorf("%s: the server answered with status %d%s", path, resp.StatusCode, bodyExcerpt(data))
	}

	switch body.Code {
	case betterAuthCodeInvalidCredentials:
		return ErrInvalidCredentials

	case betterAuthCodeInvalidCode:
		return ErrInvalidTOTPCode

	case betterAuthCodeInvalidTwoFactor:
		return fmt.Errorf(
			"%w: the server no longer knows the login the code was for, it expires ten minutes after the login",
			ErrTwoFactorRequired,
		)

	case betterAuthCodeTooManyAttempts:
		return ErrTwoFactorLocked

	case betterAuthCodeBannedUser:
		return ErrUserInactive

	default:
	}

	// The setup endpoint and anything the server itself throws answer with
	// the ok and msg of a socket.io ack rather than better-auth's code and
	// message.
	var ack setupResponse
	if json.Unmarshal(data, &ack) == nil && ack.Msg != "" {
		return fmt.Errorf("%s: %s", path, ack.Msg)
	}

	switch {
	case body.Code != "" && body.Message != "":
		return fmt.Errorf("%s: the server reported %s: %s", path, body.Code, body.Message)

	case body.Message != "":
		return fmt.Errorf("%s: %s", path, body.Message)

	default:
		return fmt.Errorf("%s: the server answered with status %d%s", path, resp.StatusCode, bodyExcerpt(data))
	}
}

// bodyExcerpt renders the start of a response body for an error message, or
// nothing for an empty one. Markup is left out: a page says nothing useful.
func bodyExcerpt(data []byte) string {
	const limit = 200

	text := strings.TrimSpace(string(data))
	if text == "" || strings.HasPrefix(text, "<") {
		return ""
	}

	if len(text) > limit {
		text = text[:limit] + "..."
	}

	return ": " + text
}

// retryAfter reads how long a rate limited request was asked to wait.
// better-auth states it in seconds, in a header of its own next to the
// standard one.
func retryAfter(header http.Header) time.Duration {
	for _, name := range []string{"X-Retry-After", "Retry-After"} {
		seconds, err := strconv.Atoi(strings.TrimSpace(header.Get(name)))
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}

	return 0
}
