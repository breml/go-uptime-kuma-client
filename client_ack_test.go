package kuma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// requiredReadyEvents are the update events kuma.New() and Client.Resync wait
// for before they return. The payloads only have to unmarshal into the
// registered handlers' argument types.
//
//nolint:gochecknoglobals // fixed table, shared by the fake server below.
var requiredReadyEvents = []string{
	`42["monitorList",{}]`,
	`42["notificationList",[]]`,
	`42["statusPageList",{}]`,
}

// optionalReadyEvents are the best-effort update events. A server that never
// emits them (an older version, or a reverse proxy dropping them) must not
// block either call, see kuma.WithReadyEvents.
//
//nolint:gochecknoglobals // fixed table, shared by the fake server below.
var optionalReadyEvents = []string{
	`42["maintenanceList",{}]`,
	`42["proxyList",[]]`,
	`42["dockerHostList",[]]`,
	`42["apiKeyList",[]]`,
}

// fakeAckCreatedID is the ID the fake server reports for a created
// notification.
const fakeAckCreatedID = 4242

// fakeSessionToken is the JWT the fake server hands out at login, and the one it
// expects to see again in a loginByToken.
const fakeSessionToken = "fake-session-token"

// fakePresetToken is a token a test hands to WithSessionToken. It is distinct
// from the one the fake server issues at login, so that a test can tell the
// token that was presented from the one that was handed out.
const fakePresetToken = "preset-session-token"

// fakeAckServer is a minimal socket.io server over HTTP long-polling that
// completes the handshake, CONNECT, login and ready phases, so kuma.New()
// returns a usable client. Unlike fakeSocketIOServer it then keeps serving:
// every event the client emits afterwards is acknowledged after ackDelay,
// which lets a test place the ack relative to the caller's deadline.
//
// It never emits an update event for those commands, so a caller using
// syncEmitWithUpdateEvent keeps waiting after its ack arrived — the state in
// which the deadline and the ack collide.
type fakeAckServer struct {
	messages chan []byte

	mu       sync.Mutex
	ackDelay time.Duration
	// omitLoginToken drops the token from the login ack, as a server that
	// hands out none would.
	omitLoginToken bool
	// suppressLists answers a loginByToken without resending the lists, which
	// leaves a Resync waiting.
	suppressLists bool
	// suppressOptionalLists resends only the required lists, as a server that
	// never emits the best-effort ones does.
	suppressOptionalLists bool
	// resyncPayload is the last loginByToken frame received.
	resyncPayload string
	// rejectCreds answers a login the way a server rejecting the username and
	// password does. credsRejectionMsg is the message that rejection carries;
	// the empty string sends the translation key a 2.x server sends, while a
	// server from before the login messages were translated sends the message
	// itself.
	rejectCreds       bool
	credsRejectionMsg string
	// rejectTokenMsg answers a loginByToken the way a server rejecting the
	// session token does, with this as the reason. The empty string accepts.
	rejectTokenMsg string
	// tokenFrames are the tokens the loginByToken frames carried, in order,
	// recorded whether the token is accepted or refused. firstTokenAt is when
	// the first of them arrived, to assert it came before any password login.
	tokenFrames  []string
	firstTokenAt time.Time
	// autoLogin answers the connect the way a server with authentication
	// disabled does: it logs the client in itself and never asks for a login.
	autoLogin bool
	// omitLoginRequired answers the connect without stating how it wants to be
	// authenticated, as a server too old to send the event does.
	omitLoginRequired bool
	// loginRequiredDelay holds the loginRequired event back, the way a server
	// that registers its handlers after an await does.
	loginRequiredDelay time.Duration
	// twoFactorRequired answers a login without a code the way a server does
	// for an account with two-factor authentication enabled. twoFactorSecret
	// is the secret the codes are verified against, replayGuard turns the
	// server's refusal to see an accepted code twice on, and acceptedCodes is
	// what that refusal reads.
	twoFactorRequired bool
	twoFactorSecret   string
	replayGuard       bool
	acceptedCodes     map[string]struct{}
	// rejectCodedLogin rejects a login carrying a code for a reason that is
	// not the code, which the client has to hand over as it is instead of
	// spending its retry on it.
	rejectCodedLogin bool
	// loginCodes are the one-time codes the logins carried, in order.
	loginCodes []string
	// setupRequired answers the connect with a setup event and rejects the
	// credentials until the setup ran, as a server without users does.
	setupRequired bool
	// setupDone records that the setup command ran, which is what makes the
	// credentials work.
	setupDone bool
	// loginRequiredAt and firstLoginAt are when the server said it wants a
	// login and when the first login arrived, to assert their order.
	loginRequiredAt time.Time
	firstLoginAt    time.Time
	// loginFrames counts the login frames received.
	loginFrames int
}

// loginCode extracts the one-time code a login frame carries.
func loginCode(payload string) string {
	var frame []json.RawMessage

	start := strings.Index(payload, "[")
	if start < 0 || json.Unmarshal([]byte(payload[start:]), &frame) != nil || len(frame) < 2 {
		return ""
	}

	var data struct {
		Token string `json:"token"`
	}

	if json.Unmarshal(frame[1], &data) != nil {
		return ""
	}

	return data.Token
}

// tokenFromFrame extracts the session token a loginByToken frame carries.
func tokenFromFrame(payload string) string {
	var frame []json.RawMessage

	start := strings.Index(payload, "[")
	if start < 0 || json.Unmarshal([]byte(payload[start:]), &frame) != nil || len(frame) < 2 {
		return ""
	}

	var token string

	if json.Unmarshal(frame[1], &token) != nil {
		return ""
	}

	return token
}

func (s *fakeAckServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	switch {
	case r.URL.Path == "/api/entry-page":
		// The database is set up; what is missing is the first user, which the
		// setup event asks for.
		_, _ = w.Write([]byte(`{"type":"entryPage","entryPage":null}`))

	case r.Method == http.MethodGet && sid == "":
		type engineIOHandshake struct {
			Sid          string   `json:"sid"`
			Upgrades     []string `json:"upgrades"`
			PingInterval int      `json:"pingInterval"`
			PingTimeout  int      `json:"pingTimeout"`
			MaxPayload   int      `json:"maxPayload"`
		}
		data, err := json.Marshal(engineIOHandshake{
			Sid:          "ack-test-sid",
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
		select {
		case msg := <-s.messages:
			_, _ = w.Write(msg)

		case <-r.Context().Done():
		}

	default:
	}
}

func (s *fakeAckServer) setAckDelay(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ackDelay = d
}

func (s *fakeAckServer) currentAckDelay() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ackDelay
}

func (s *fakeAckServer) setOmitLoginToken(omit bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.omitLoginToken = omit
}

func (s *fakeAckServer) setRejectCreds(reject bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rejectCreds = reject
}

// setCredsRejectionMsg rejects the credentials with a specific message, such
// as the untranslated "Incorrect username or password." a 1.x server sends.
func (s *fakeAckServer) setCredsRejectionMsg(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.credsRejectionMsg = msg
}

func (s *fakeAckServer) setRejectToken() {
	s.setRejectTokenWith("authInvalidToken")
}

// setRejectTokenWith refuses a loginByToken with a specific reason, such as the
// authUserInactiveOrDeleted the server answers for a deactivated user.
func (s *fakeAckServer) setRejectTokenWith(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rejectTokenMsg = msg
}

func (s *fakeAckServer) setAutoLogin(auto bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.autoLogin = auto
}

func (s *fakeAckServer) setLoginRequiredDelay(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.loginRequiredDelay = d
}

// setTwoFactorSecret turns two-factor authentication on and verifies the codes
// against the shared test secret. With replayGuard set it also refuses every
// code it has already accepted; the real server only remembers the last one it
// let through, in twofa_last_token, which is stricter than these tests need to
// tell apart.
func (s *fakeAckServer) setTwoFactorSecret(replayGuard bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.twoFactorRequired = true
	s.twoFactorSecret = rfc6238Secret
	s.replayGuard = replayGuard
	s.acceptedCodes = map[string]struct{}{}
}

// useCode marks a code as already accepted, the way another client logging in
// with the same account inside the same time step would.
func (s *fakeAckServer) useCode(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.acceptedCodes[code] = struct{}{}
}

// receivedLoginCodes returns the one-time codes the logins carried, in order.
func (s *fakeAckServer) receivedLoginCodes() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.loginCodes...)
}

func (s *fakeAckServer) setRejectCodedLogin(reject bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rejectCodedLogin = reject
}

func (s *fakeAckServer) setSetupRequired(required bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.setupRequired = required
}

func (s *fakeAckServer) runSetup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.setupDone = true
}

func (s *fakeAckServer) setTwoFactorRequired(required bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.twoFactorRequired = required
}

// firstLoginAfterLoginRequired reports whether the client held its login back
// until the server asked for one.
func (s *fakeAckServer) firstLoginAfterLoginRequired() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return !s.firstLoginAt.IsZero() &&
		!s.loginRequiredAt.IsZero() &&
		s.firstLoginAt.After(s.loginRequiredAt)
}

func (s *fakeAckServer) setOmitLoginRequired(omit bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.omitLoginRequired = omit
}

func (s *fakeAckServer) recordLogin(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.loginFrames++
	s.loginCodes = append(s.loginCodes, code)

	if s.firstLoginAt.IsZero() {
		s.firstLoginAt = time.Now()
	}
}

func (s *fakeAckServer) loginFrameCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loginFrames
}

func (s *fakeAckServer) setSuppressLists(suppress bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.suppressLists = suppress
}

func (s *fakeAckServer) setSuppressOptionalLists(suppress bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.suppressOptionalLists = suppress
}

// tokenRejection returns the reason a loginByToken is to be refused with, or
// the empty string for a token the server accepts.
func (s *fakeAckServer) tokenRejection() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rejectTokenMsg
}

// recordToken records the token a loginByToken frame carried. It runs before
// the rejection is answered, so a refused frame is seen as well.
func (s *fakeAckServer) recordToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokenFrames = append(s.tokenFrames, token)

	if s.firstTokenAt.IsZero() {
		s.firstTokenAt = time.Now()
	}
}

// receivedTokens returns the tokens the loginByToken frames carried, in order.
func (s *fakeAckServer) receivedTokens() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.tokenFrames...)
}

// loginAfterToken reports whether the password login the client fell back to
// came after the token it presented first.
func (s *fakeAckServer) loginAfterToken() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return !s.firstLoginAt.IsZero() &&
		!s.firstTokenAt.IsZero() &&
		s.firstLoginAt.After(s.firstTokenAt)
}

func (s *fakeAckServer) recordResync(payload string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resyncPayload = payload

	return !s.suppressLists
}

func (s *fakeAckServer) lastResyncPayload() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.resyncPayload
}

// acceptsCodeLocked verifies a one-time code the way the server does, with the
// replay guard that refuses a code it has already let through. The caller holds
// the lock.
func (s *fakeAckServer) acceptsCodeLocked(code string) bool {
	if s.twoFactorSecret == "" {
		return false
	}

	if s.replayGuard {
		if _, used := s.acceptedCodes[code]; used {
			return false
		}
	}

	valid, err := totp.ValidateCustom(code, s.twoFactorSecret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil || !valid {
		return false
	}

	s.acceptedCodes[code] = struct{}{}

	return true
}

func (s *fakeAckServer) loginAck(code string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.rejectCreds || (s.setupRequired && !s.setupDone) {
		if s.credsRejectionMsg != "" {
			return fmt.Sprintf(`{"ok":false,"msg":%q}`, s.credsRejectionMsg)
		}

		return `{"ok":false,"msg":"authIncorrectCreds","msgi18n":true}`
	}

	// An ack asking for a one-time code carries neither an ok nor a message.
	if s.twoFactorRequired && code == "" {
		return `{"tokenRequired":true}`
	}

	if s.rejectCodedLogin && code != "" {
		return `{"ok":false,"msg":"authIncorrectCreds","msgi18n":true}`
	}

	if s.twoFactorRequired && !s.acceptsCodeLocked(code) {
		return `{"ok":false,"msg":"authInvalidToken","msgi18n":true}`
	}

	if s.omitLoginToken {
		return `{"ok":true,"msg":"Logged in successfully."}`
	}

	return fmt.Sprintf(`{"ok":true,"msg":"Logged in successfully.","token":%q}`, fakeSessionToken)
}

// sendAuthMode states how the connection is to be authenticated, which a real
// server does as the last thing it does for a new connection.
func (s *fakeAckServer) sendAuthMode() {
	s.mu.Lock()
	auto := s.autoLogin
	omit := s.omitLoginRequired
	delay := s.loginRequiredDelay
	needsSetup := s.setupRequired
	s.mu.Unlock()

	if needsSetup {
		// Both events go out in one engine.io payload, the way a real server
		// delivers them to a client that polls: it emits setup at the top of
		// its connection handler and the event below at the bottom, so a
		// single long poll picks up both.
		if !auto && !omit {
			s.messages <- []byte("42[\"setup\"]\x1e42[\"loginRequired\"]")

			return
		}

		s.messages <- []byte(`42["setup"]`)
	}

	switch {
	case auto:
		s.messages <- []byte(`42["autoLogin"]`)
		s.sendReadyEvents()

	case omit:

	default:
		// Held back out of band so the delay does not stall the request that
		// carried the connect.
		go func() {
			time.Sleep(delay)

			s.mu.Lock()
			s.loginRequiredAt = time.Now()
			s.mu.Unlock()

			s.messages <- []byte(`42["loginRequired"]`)
		}()
	}
}

func (s *fakeAckServer) sendReadyEvents() {
	s.mu.Lock()
	optionalSuppressed := s.suppressOptionalLists
	s.mu.Unlock()

	events := requiredReadyEvents
	if !optionalSuppressed {
		events = append(append([]string{}, requiredReadyEvents...), optionalReadyEvents...)
	}

	for _, event := range events {
		s.messages <- []byte(event)
	}
}

func (s *fakeAckServer) handleClientMessage(body []byte) {
	if len(body) < 2 {
		return
	}

	socketIOData := body[1:] // strip engine.io "4" (Message) prefix

	switch socketIOData[0] {
	case '0': // socket.io CONNECT
		s.messages <- []byte("40")
		s.sendAuthMode()

	case '2': // socket.io EVENT
		i := 1
		for i < len(socketIOData) && socketIOData[i] >= '0' && socketIOData[i] <= '9' {
			i++
		}

		ackID := string(socketIOData[1:i])
		payload := string(socketIOData[i:])

		// A resync logs in again, which is what makes the server resend the
		// lists. Checked before "login", which is a prefix of it.
		if strings.Contains(payload, `"loginByToken"`) {
			s.recordToken(tokenFromFrame(payload))

			if msg := s.tokenRejection(); msg != "" {
				s.messages <- fmt.Appendf(
					nil,
					`43%s[{"ok":false,"msg":%q,"msgi18n":true}]`,
					ackID,
					msg,
				)

				return
			}

			s.messages <- fmt.Appendf(nil, `43%s[{"ok":true}]`, ackID)

			if s.recordResync(payload) {
				s.sendReadyEvents()
			}

			return
		}

		if strings.Contains(payload, `"setup"`) {
			s.runSetup()
			s.messages <- fmt.Appendf(nil, `43%s[{"ok":true,"msg":"ok"}]`, ackID)

			return
		}

		if strings.Contains(payload, `"login"`) {
			code := loginCode(payload)
			s.recordLogin(code)

			ack := s.loginAck(code)
			s.messages <- fmt.Appendf(nil, `43%s[%s]`, ackID, ack)

			if strings.Contains(ack, `"ok":true`) {
				s.sendReadyEvents()
			}

			return
		}

		// The server assigns the ID of a created resource in the ack, which
		// is the value a caller loses if a successful ack is discarded.
		ack := `{"ok":true,"msg":"ok"}`
		if strings.Contains(payload, `"addNotification"`) {
			ack = fmt.Sprintf(`{"ok":true,"msg":"ok","id":%d}`, fakeAckCreatedID)
		}

		// Acknowledge out of band: blocking here would stall the client's
		// send path instead of only delaying the ack.
		delay := s.currentAckDelay()
		go func() {
			time.Sleep(delay)
			s.messages <- fmt.Appendf(nil, `43%s[%s]`, ackID, ack)
		}()

	default:
	}
}

// TestAckDeliveryAroundDeadline sweeps the ack of an in-flight command across
// the caller's deadline, so that early, colliding and late acks are all
// exercised against the same call.
//
// The defect this guards against is a "send on closed channel" panic: the ack
// callback runs on a goroutine owned by the socket.io client, so an ack that
// arrives once the caller has returned finds the channel closed by its
// `defer close(res)`, and the panic is raised where no caller can recover from
// it. That kills the process, not just the request.
//
// What it catches, and what it does not: the late-ack iterations at the end of
// the sweep panic every time if `close(res)` is reintroduced, so the invariant
// "the ack channel is buffered and never closed" is genuinely enforced here.
// The original code guarded the send with `if ctx.Err() != nil` and crashed
// only when the context expired between that check and the send. That window
// is a few instructions wide, and it was verified that this test does not
// reproduce it — real transport latency swamps it. So a green run is not
// evidence that a guard-based implementation is safe; it is evidence that the
// channel is still handled the way syncEmit and syncEmitWithUpdateEvent
// document.
func TestAckDeliveryAroundDeadline(t *testing.T) {
	const (
		callTimeout = 30 * time.Millisecond
		iterations  = 24
		step        = 500 * time.Microsecond
	)

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	server := httptest.NewServer(fake)

	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	kumaClient, err := kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(10*time.Second),
	)
	require.NoError(t, err)

	// Start clearly before the deadline and end clearly after it.
	base := callTimeout - (iterations/2)*step

	t.Run("sync_emit_with_update_event", func(t *testing.T) {
		for i := range iterations {
			fake.setAckDelay(base + time.Duration(i)*step)

			callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
			// The fake acks the command but never emits notificationList, so
			// the call keeps waiting for its update event and always runs into
			// the deadline — with the ack landing somewhere around it.
			err := kumaClient.DeleteNotification(callCtx, 1)
			callCancel()

			// Both outcomes of the race keep wrapping the context error: the
			// ack that lost it produces a bare deadline error, the ack that
			// won it produces ErrUpdateEventTimeout wrapping the same
			// deadline error. Which one a given delay yields depends on
			// scheduling, so only the shared part is asserted here;
			// TestAwaitAckAndUpdateEvent pins the ack-wins case down
			// deterministically.
			require.ErrorIs(t, err, context.DeadlineExceeded,
				"ack delay %s should leave the call waiting on its update event", fake.currentAckDelay())
		}
	})

	t.Run("sync_emit", func(t *testing.T) {
		for i := range iterations {
			fake.setAckDelay(base + time.Duration(i)*step)

			callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
			msg, err := kumaClient.TestNotification(callCtx, notification.Generic{
				Base:           notification.Base{Name: "ack race"},
				GenericDetails: notification.GenericDetails{},
				TypeName:       "generic",
			})
			callCancel()

			// Whether the ack or the deadline won is a matter of scheduling and
			// either outcome is correct. A panic, a hang, or a garbled response
			// is not.
			if err != nil {
				require.ErrorIs(t, err, context.DeadlineExceeded)

				continue
			}

			require.Equal(t, "ok", msg)
		}
	})
}

// TestUpdateEventTimeoutKeepsSuccessfulAck covers the case the sweep above can
// only hit by chance: the server acknowledges the command well before the
// deadline and never emits the update event that follows it.
//
// The command was applied, so reporting a bare context error would tell the
// caller the write failed for a write the server performed — and for a create,
// it would also drop the ID the server assigned, leaving a retry to produce a
// duplicate.
func TestUpdateEventTimeoutKeepsSuccessfulAck(t *testing.T) {
	// Long enough that the ack, which the fake sends without delay, reliably
	// arrives first; the call then runs into the deadline waiting for the
	// update event the fake never emits. It travels over HTTP long-polling, so
	// this is a generous margin rather than a guarantee: raise callTimeout if
	// a loaded machine ever makes it flake.
	const callTimeout = 500 * time.Millisecond

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	server := httptest.NewServer(fake)

	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	kumaClient, err := kuma.New(
		ctx,
		server.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(10*time.Second),
	)
	require.NoError(t, err)

	t.Run("create_reports_the_id_the_server_assigned", func(t *testing.T) {
		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		id, err := kumaClient.CreateNotification(callCtx, notification.Generic{
			Base:           notification.Base{Name: "update event timeout"},
			GenericDetails: notification.GenericDetails{},
			TypeName:       "generic",
		})

		require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
		require.ErrorIs(t, err, context.DeadlineExceeded,
			"the sentinel must not hide the timeout from callers matching on it")
		require.Equal(t, int64(fakeAckCreatedID), id,
			"the notification exists on the server, so its ID must survive the missing update event")
	})

	t.Run("the_id_survives_a_caller_that_propagates_only_the_error", func(t *testing.T) {
		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		// The reflex of every caller is to return the error and nothing else,
		// which drops the ID that sits next to it. Recovering the error has to
		// give that ID back, or the created notification is unreachable.
		_, err := func() (int64, error) {
			id, err := kumaClient.CreateNotification(callCtx, notification.Generic{
				Base:           notification.Base{Name: "propagated error"},
				GenericDetails: notification.GenericDetails{},
				TypeName:       "generic",
			})
			if err != nil {
				return 0, err
			}

			return id, nil
		}()

		var timeoutErr *kuma.UpdateEventTimeoutError
		require.ErrorAs(t, err, &timeoutErr)
		require.Equal(t, int64(fakeAckCreatedID), timeoutErr.ID)
		require.Equal(t, "addNotification", timeoutErr.Command)
	})

	t.Run("delete_is_distinguishable_from_a_plain_timeout", func(t *testing.T) {
		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		err := kumaClient.DeleteNotification(callCtx, 1)

		require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("an_unacknowledged_command_stays_a_plain_timeout", func(t *testing.T) {
		// No ack at all before the deadline: nothing is known about the write,
		// so the sentinel must not be reported.
		fake.setAckDelay(2 * callTimeout)
		t.Cleanup(func() { fake.setAckDelay(0) })

		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		err := kumaClient.DeleteNotification(callCtx, 1)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NotErrorIs(t, err, kuma.ErrUpdateEventTimeout,
			"without an ack the client cannot claim the command was applied")
	})
}
