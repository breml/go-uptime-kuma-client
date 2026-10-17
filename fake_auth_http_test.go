package kuma_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// This file holds the parts of fakeAckServer that make it behave like an
// Uptime Kuma 3.x server, see kuma.AuthSchemeHTTP: the engine.io sessions a
// reconnect produces, and the HTTP authentication endpoints of better-auth.

// fakeSession is one engine.io session of the fake: the long polls of one
// socket.io connection are served from its channel until it is replaced.
type fakeSession struct {
	sid      string
	messages chan []byte
	// cookie is the session cookie the handshake carried, or the empty
	// string.
	cookie string
}

// fakeHandshake records one engine.io handshake: when it was made and the
// Cookie header it carried.
type fakeHandshake struct {
	sid    string
	cookie string
	at     time.Time
}

// fakeTwoFactorCookie is the cookie a sign-in answered with a request for a
// code sets, the way better-auth names the pending login.
const fakeTwoFactorCookie = "better-auth.two_factor"

// fakeSessionCookie is the session cookie of a 3.x server.
const fakeSessionCookie = "better-auth.session_token"

// fakeRetryAfter is the wait the rate limiter states, in seconds.
const fakeRetryAfter = 7

// send queues a frame for the current session, which is the one the latest
// handshake started. It is for frames a test injects; frames the fake answers
// a command with go to the session the command came in on, see sendTo.
func (s *fakeAckServer) send(msg []byte) {
	s.mu.Lock()
	messages := s.messages
	s.mu.Unlock()

	messages <- msg
}

// sendTo queues a frame for one session.
func (*fakeAckServer) sendTo(session *fakeSession, msg []byte) {
	session.messages <- msg
}

// newSession starts the engine.io session a handshake asks for and makes it
// the current one. Earlier sessions stay served, the way a real server keeps
// every connection a client holds open: which one the client reads from is the
// client's business.
func (s *fakeAckServer) newSession(r *http.Request) *fakeSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sidCounter++

	session := &fakeSession{
		sid: fmt.Sprintf("ack-test-sid-%d", s.sidCounter),
	}

	cookie, err := r.Cookie(fakeSessionCookie)
	if err == nil {
		session.cookie = cookie.Value
	}

	if s.current == nil {
		// The first session adopts the channel the test created the fake
		// with, so that frames a test queued before connecting are seen.
		session.messages = s.messages
	} else {
		session.messages = make(chan []byte, cap(s.messages))
	}

	if s.sessions == nil {
		s.sessions = map[string]*fakeSession{}
	}

	s.sessions[session.sid] = session
	s.current = session
	s.messages = session.messages
	s.handshakes = append(s.handshakes, fakeHandshake{
		sid:    session.sid,
		cookie: r.Header.Get("Cookie"),
		at:     time.Now(),
	})

	return session
}

func (s *fakeAckServer) session(sid string) *fakeSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.sessions[sid]
}

// handshakeCount returns how many engine.io handshakes were made, which is how
// many connections the client opened.
func (s *fakeAckServer) handshakeCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.handshakes)
}

// handshakeCookies returns the Cookie header of every handshake, in order.
func (s *fakeAckServer) handshakeCookies() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	cookies := make([]string, 0, len(s.handshakes))
	for _, handshake := range s.handshakes {
		cookies = append(cookies, handshake.cookie)
	}

	return cookies
}

// receivedUnexpectedCommands returns the socket commands a 3.x server has no
// handler for that the client sent anyway.
func (s *fakeAckServer) receivedUnexpectedCommands() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.unexpectedCommands...)
}

func (s *fakeAckServer) recordUnexpected(command string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.unexpectedCommands = append(s.unexpectedCommands, command)
}

// setHTTPAuth switches the fake to the authentication of a 3.x server.
func (s *fakeAckServer) setHTTPAuth() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.httpAuth = true
}

func (s *fakeAckServer) isHTTPAuth() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.httpAuth
}

// probeCount returns how often the auth scheme probe was requested.
func (s *fakeAckServer) probeCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.probes
}

func (s *fakeAckServer) setProbeRateLimited(limited bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.probeRateLimited = limited
}

func (s *fakeAckServer) setSignInRateLimited(limited bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.signInRateLimited = limited
}

func (s *fakeAckServer) setStripCookie(strip bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stripCookie = strip
}

// issueSession hands out the session cookie a login answers with and records
// it as one the fake accepts from then on.
func (s *fakeAckServer) issueSession(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.issuedSessions == nil {
		s.issuedSessions = map[string]struct{}{}
	}

	s.issuedSessions[fakeSessionToken] = struct{}{}

	http.SetCookie(w, &http.Cookie{
		Name:     fakeSessionCookie,
		Value:    fakeSessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// acceptsSession reports whether a session cookie names a session the fake
// accepts: one a login handed out, or any other one unless the fake was told
// to refuse the token, see setRejectToken.
func (s *fakeAckServer) acceptsSession(cookie string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cookie == "" {
		return false
	}

	if _, issued := s.issuedSessions[cookie]; issued {
		return true
	}

	return s.rejectTokenMsg == ""
}

// unknownSession answers a poll for a session the fake does not have, the way
// engine.io does.
func unknownSession(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"code":1,"message":"Session ID unknown"}`))
}

// rateLimited answers a request the way better-auth's rate limiter does.
func rateLimited(w http.ResponseWriter) {
	w.Header().Set("X-Retry-After", strconv.Itoa(fakeRetryAfter))
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write([]byte(`{"message":"Too many requests. Please try again later."}`))
}

// betterAuthError answers a request the way better-auth refuses one, with the
// 401 every refusal of a login carries.
func betterAuthError(w http.ResponseWriter, code string, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = fmt.Fprintf(w, `{"code":%q,"message":%q}`, code, message)
}

// serveAuthHTTP answers the authentication endpoints of a 3.x server. In
// socket mode the fake has none of them: like a 2.x server it falls through
// to a catch-all, which serves the single-page application as HTML.
func (s *fakeAckServer) serveAuthHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.isHTTPAuth() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Uptime Kuma</body></html>"))

		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/api/auth/ok":
		s.mu.Lock()
		s.probes++
		limited := s.probeRateLimited
		s.mu.Unlock()

		if limited {
			rateLimited(w)

			return
		}

		_, _ = w.Write([]byte(`{"ok":true}`))

	case "/api/auth/sign-in/username":
		s.serveSignIn(w, r)

	case "/api/auth/two-factor/verify-totp":
		s.serveVerifyTOTP(w, r)

	case "/api/auth/get-session":
		s.serveGetSession(w, r)

	case "/api/setup":
		s.runSetup()
		_, _ = w.Write([]byte(`{"ok":true}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"NOT_FOUND","message":"Not Found"}`))
	}
}

// serveSignIn answers the sign-in of a 3.x server. It mirrors loginAck: the
// same switches decide the outcome, only the answer is an HTTP status and a
// cookie instead of an ack.
func (s *fakeAckServer) serveSignIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	_ = json.NewDecoder(r.Body).Decode(&body)

	s.recordLogin("")

	s.mu.Lock()
	limited := s.signInRateLimited
	rejected := s.rejectCreds || (s.setupRequired && !s.setupDone)
	twoFactor := s.twoFactorRequired
	omitToken := s.omitLoginToken
	s.mu.Unlock()

	switch {
	case limited:
		rateLimited(w)

	case rejected:
		betterAuthError(w, "INVALID_USERNAME_OR_PASSWORD", "Invalid username or password")

	case twoFactor:
		s.mu.Lock()
		if s.twoFactorCookies == nil {
			s.twoFactorCookies = map[string]struct{}{}
		}

		pending := fmt.Sprintf("pending-%d", len(s.twoFactorCookies)+1)
		s.twoFactorCookies[pending] = struct{}{}
		s.mu.Unlock()

		http.SetCookie(w, &http.Cookie{Name: fakeTwoFactorCookie, Value: pending, Path: "/", MaxAge: 600})
		_, _ = w.Write([]byte(`{"twoFactorRedirect":true}`))

	case omitToken:
		_, _ = w.Write([]byte(`{"token":"","user":{"username":"admin"}}`))

	default:
		s.issueSession(w)
		_, _ = fmt.Fprintf(w, `{"token":%q,"user":{"username":"admin"}}`, fakeSessionToken)
	}
}

// serveVerifyTOTP answers the one-time code of a 3.x server.
func (s *fakeAckServer) serveVerifyTOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}

	_ = json.NewDecoder(r.Body).Decode(&body)

	s.recordLogin(body.Code)

	pending, cookieErr := r.Cookie(fakeTwoFactorCookie)

	s.mu.Lock()

	known := false
	if cookieErr == nil {
		_, known = s.twoFactorCookies[pending.Value]
	}

	rejectCoded := s.rejectCodedLogin
	accepted := known && !rejectCoded && s.acceptsCodeLocked(body.Code)
	s.mu.Unlock()

	switch {
	case !known:
		betterAuthError(w, "INVALID_TWO_FACTOR_COOKIE", "Invalid two factor cookie")

	case rejectCoded:
		betterAuthError(w, "INVALID_USERNAME_OR_PASSWORD", "Invalid username or password")

	case !accepted:
		betterAuthError(w, "INVALID_CODE", "Invalid code")

	default:
		s.issueSession(w)
		_, _ = fmt.Fprintf(w, `{"token":%q,"user":{"username":"admin"}}`, fakeSessionToken)
	}
}

// serveGetSession answers whether the session cookie names a session the fake
// accepts, the way better-auth does: with the session, or with null.
func (s *fakeAckServer) serveGetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(fakeSessionCookie)
	if err != nil {
		_, _ = w.Write([]byte("null"))

		return
	}

	s.recordToken(cookie.Value)

	if !s.acceptsSession(cookie.Value) {
		_, _ = w.Write([]byte("null"))

		return
	}

	s.mu.Lock()
	twoFactor := s.twoFactorRequired
	s.mu.Unlock()

	_, _ = fmt.Fprintf(w, `{"session":{"token":%q},"user":{"username":"admin","twoFactorEnabled":%t}}`,
		cookie.Value, twoFactor)
}

// sendAuthModeHTTP states how a connection to a 3.x server is authenticated:
// from the session cookie on its handshake, which answers with the session
// event and every list, and otherwise with loginRequired. A server that wants
// to be set up has no user a cookie could name.
func (s *fakeAckServer) sendAuthModeHTTP(session *fakeSession) {
	s.mu.Lock()
	needsSetup := s.setupRequired && !s.setupDone
	auto := s.autoLogin
	strip := s.stripCookie
	delay := s.loginRequiredDelay
	suppressLists := s.suppressLists
	s.mu.Unlock()

	cookie := session.cookie
	if strip {
		cookie = ""
	}

	switch {
	case needsSetup:
		s.sendTo(session, []byte("42[\"setup\"]\x1e42[\"loginRequired\"]"))

	case auto || s.acceptsSession(cookie):
		s.sendTo(session, []byte(`42["session","admin"]`))

		// A server that authenticates the connection but resends nothing
		// leaves a resync waiting, see setSuppressLists.
		if !suppressLists {
			s.sendReadyEvents(session)
		}

	default:
		// Held back out of band so the delay does not stall the request that
		// carried the connect.
		go func() {
			time.Sleep(delay)

			s.mu.Lock()
			s.loginRequiredAt = time.Now()
			s.mu.Unlock()

			s.sendTo(session, []byte(`42["loginRequired"]`))
		}()
	}
}
