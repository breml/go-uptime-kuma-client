package kuma

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ErrAuthRequired is returned when the server asks for a login and the client
// has no credentials to offer.
//
// The server states what it expects right after the connect: it emits
// loginRequired for a server that wants a login and autoLogin for one that has
// authentication disabled and logged the client in itself.
var ErrAuthRequired = errors.New("server requires authentication")

// ErrInvalidCredentials is returned when the server rejects the username and
// password. It is what the server reports as authIncorrectCreds, and as
// "Incorrect username or password" on versions before the login messages were
// translated.
var ErrInvalidCredentials = errors.New("incorrect username or password")

// ErrTwoFactorRequired is returned when the account has two-factor
// authentication enabled and the login could not answer the server's request
// for a code: neither WithTOTPSecret nor WithTOTPCode was configured, the
// configured source produced an empty code, or the server asked for a code
// again instead of answering the one it was given.
//
// A code source that fails outright reports its own error instead, because what
// went wrong there is not that a code was missing.
//
// It is the typed form of the server's tokenRequired answer to a login. That
// answer carries neither an ok nor a message, so without this it surfaces as an
// error with an empty message.
var ErrTwoFactorRequired = errors.New("two-factor authentication required")

// ErrInvalidTOTPCode is returned when the server rejects the one-time code of a
// login. Besides a wrong secret and a clock too far off, this is what the
// server answers for a code it has accepted once already: it records the last
// code it let through and refuses to see it again, so two logins inside the
// same 30 second step cannot both succeed.
var ErrInvalidTOTPCode = errors.New("invalid two-factor authentication code")

// ErrInvalidSessionToken is returned when the server rejects a session token,
// see Client.Resync. The token carries no expiry, so what invalidates it is a
// changed password or a user that was deactivated or deleted.
var ErrInvalidSessionToken = errors.New("session token rejected")

// ErrUserInactive is returned when the server rejects a session token because
// the account it names is deactivated or deleted. It wraps
// ErrInvalidSessionToken, so a caller that only cares that the token is gone
// keeps matching on that one.
//
// It is the token rejection a password cannot recover from: the server requires
// an active user for a password login too, so New reports it instead of falling
// back, see WithSessionToken.
var ErrUserInactive = fmt.Errorf("%w: user inactive or deleted", ErrInvalidSessionToken)

// ErrRateLimited is returned when the server refuses a login because too many
// were attempted in a short time. The server allows 20 per minute and answers
// the ones beyond that with an untranslated sentence rather than a key.
//
// A login that answers a one-time code costs two of those attempts, because the
// server asks for the code in the ack of a first login, so an account with
// two-factor authentication enabled reaches the limit in half as many logins.
var ErrRateLimited = errors.New("too many login attempts, rate limited by the server")

// WithSessionToken authenticates with a session token from an earlier login
// instead of a password, see Client.SessionToken.
//
// The token is what the server hands out on a successful login and what it
// accepts in place of one. It bypasses two-factor authentication entirely, so
// a client holding one needs neither the password nor the one-time code, which
// also makes it the way to start many clients at once against an account with
// two-factor authentication enabled, see WithTOTPSecret.
//
// If a username and password are given as well, the token is tried first and
// the password login is the fallback for a token the server refuses.
// Client.SessionTokenRejected reports that this happened, because the login
// succeeds either way and the stored token is dead. The fallback is skipped
// where the password cannot help: for a deactivated or deleted user, see
// ErrUserInactive, and for a failure that is not a rejection at all, such as a
// transport error. Without a username and password a rejected token fails New
// with ErrInvalidSessionToken.
//
// Security: the token is a bearer credential that does not expire. The server
// invalidates it only when the password changes or the user is deactivated,
// and it is a stronger credential than the password for an account with
// two-factor authentication, because it needs no second factor. Store it the
// way the password would be stored.
func WithSessionToken(token string) Option {
	return func(c *Client) {
		c.sessionTokenPreset = token
	}
}

// WithSessionTokenRejectedCallback registers a function New calls when the
// server refused the token WithSessionToken configured and the password login
// took over, with the error the token was refused with.
//
// It is the pushed form of Client.SessionTokenRejected, for a caller that
// persists tokens and would otherwise have to remember to ask: the fallback
// leaves New succeeding, so nothing about the returned client says that the
// stored token is dead. The callback runs on the goroutine that called New,
// before it returns, and a client whose token was accepted never calls it.
func WithSessionTokenRejectedCallback(callback func(err error)) Option {
	return func(c *Client) {
		c.sessionTokenRejectedCallback = callback
	}
}

// WithTOTPSecret configures the shared secret of an account with two-factor
// authentication enabled, so the client can answer the server's request for a
// one-time code itself and log in without a human at the keyboard.
//
// The secret is the base32 string from the otpauth:// URI Uptime Kuma shows
// when two-factor authentication is set up. Spaces, hyphens, lower case and
// the missing base32 padding the server strips are all accepted; a secret that
// cannot be decoded fails New rather than the login it would be needed for.
//
// The code is generated only once the server asks for one, so configuring a
// secret for an account without two-factor authentication changes nothing. It
// is generated then and not before because the server evaluates the fields of
// a login independently: a code sent to an account that has two-factor
// authentication disabled is verified against a secret that does not exist,
// after the login has already been answered.
//
// The server records the last code it accepted and refuses to see it again,
// see ErrInvalidTOTPCode, so two clients logging in with the same account
// inside one 30 second step cannot both succeed. The client retries once in
// the next step when the caller's deadline allows it, which settles a collision
// between two clients and no more: of three or more meeting in one step, the
// ones that lose the retry as well have spent both attempts, however generous
// the deadline. WithSessionToken is the answer that scales, because it needs no
// code and so meets no such guard.
//
// The server answers a wrong secret and a clock too far off with the same
// rejection, so those cost the same retry: a login that can never succeed
// takes up to one step to fail. WithConnectTimeout bounds it.
//
// WithTOTPSecret and WithTOTPCode are mutually exclusive.
func WithTOTPSecret(secret string) Option {
	return func(c *Client) {
		c.totpSources++

		normalized, err := normalizeTOTPSecret(secret)
		if err != nil {
			c.totpErr = err

			return
		}

		c.totpCode = func(_ context.Context) (string, error) {
			return totpCodeAt(normalized, time.Now())
		}
	}
}

// WithTOTPCode configures a callback that produces the one-time code, for a
// caller whose secret lives somewhere the client cannot read it, such as a
// hardware token or an external authenticator.
//
// It is called at most twice per login: once when the server asks for a code,
// and once more for the single retry the client makes in the next time step,
// which happens only if the server rejected the first code and the caller's
// deadline can cover the wait, see ErrInvalidTOTPCode. A callback that returns
// an empty code fails the login with ErrTwoFactorRequired rather than sending
// it, because the server reads an empty code as no code at all.
//
// A nil callback fails New, the way a secret WithTOTPSecret cannot decode
// does. WithTOTPSecret is this option with the code computed from a secret,
// and the two are mutually exclusive.
func WithTOTPCode(code func(ctx context.Context) (string, error)) Option {
	return func(c *Client) {
		c.totpSources++

		if code == nil {
			c.totpErr = errors.New("totp: WithTOTPCode was given a nil callback")

			return
		}

		c.totpCode = code
	}
}

// authBarrierWait bounds how long the client waits for the server to state how
// it wants to be authenticated, see authBarrier.
//
// The server emits loginRequired or autoLogin as the last thing it does for a
// new connection, so on a supported server the wait is over in milliseconds. It
// only elapses in full behind a proxy that drops the event, or against an
// Uptime Kuma 1.x server, which sends neither and is not supported, see the
// package documentation. It must not be a hard failure for those - which is
// also why it never takes more than half of what is left of the caller's
// deadline, see barrierWaitWithin.
const authBarrierWait = 500 * time.Millisecond

// setupEventGrace is how long setupPending waits for a setup event that has not
// arrived yet, to settle its race with the ack of the login the server rejects
// because it has no users to match against.
//
// It is a variable so the tests can widen it, see export_test.go.
//
//nolint:gochecknoglobals // widened by the tests, constant everywhere else.
var setupEventGrace = 10 * time.Millisecond

// authMode is how the server wants the connection to be authenticated, as the
// events it emits right after the connect state it.
type authMode int

const (
	// authModeUnknown means the server said nothing the client recognizes.
	authModeUnknown authMode = iota

	// authModeLoginRequired means the server wants a login.
	authModeLoginRequired

	// authModeAutoLogin means the server has authentication disabled and has
	// logged the client in itself.
	authModeAutoLogin

	// authModeSetup means the server is not set up yet and wants that first.
	authModeSetup

	// authModeSession means the server authenticated the connection itself:
	// from the session cookie on its handshake, or because it has
	// authentication disabled. It is what an Uptime Kuma 3.x server states in
	// place of a login ack, see AuthSchemeHTTP.
	authModeSession
)

// authBarrier reports how the server wants the connection to be authenticated.
//
// It doubles as the point from which the socket.io API is known to be live. The
// server registers its login and loginByToken handlers after an await in the
// connection handler and emits loginRequired or autoLogin only once every
// handler is in place, so a login emitted before that can land in the gap and
// be dropped silently, leaving the client to wait out its timeout on a
// connection that is otherwise healthy. Waiting for the barrier closes that
// window.
type authBarrier struct {
	loginRequired <-chan struct{}
	autoLogin     <-chan struct{}
	session       <-chan struct{}
	setupRequired <-chan struct{}
}

// credentials are the ways New was given to authenticate, gathered from its
// arguments and options so the login path has a single thing to read.
type credentials struct {
	username string
	password string
	token    string
}

// newCredentials gathers what New was given to authenticate with, and rejects a
// username without a password (or the other way round): the server has no login
// that takes one without the other, so it can only ever be a caller mistake.
func newCredentials(username string, password string, token string) (credentials, error) {
	if (username == "") != (password == "") {
		return credentials{}, errors.New("credentials: username and password have to be set together")
	}

	return credentials{username: username, password: password, token: token}, nil
}

// loginError translates a rejected login ack into one of the package's sentinel
// errors, so a caller can tell the cases apart with errors.Is instead of
// matching on the message.
func loginError(command string, response ackResponse) error {
	// A server from before the login messages were translated sends the
	// message itself, and sends it with a trailing period the translation key
	// that replaced it has no equivalent of, so it is matched by its prefix.
	// The rate limiter's answer is a sentence on every version, because it was
	// never given a key of its own.
	//
	// These two are the only untranslated messages that are classified. Uptime
	// Kuma 1.x, which answers every rejection this way, is not supported: its
	// remaining rejections - an invalid session token, a deactivated user, a
	// wrong one-time code - fall through to unclassifiedLoginError.
	if !response.MsgI18n {
		if strings.HasPrefix(response.Msg, "Incorrect username or password") {
			return ErrInvalidCredentials
		}

		if strings.HasPrefix(response.Msg, "Too frequently") {
			return ErrRateLimited
		}
	}

	switch response.Msg {
	case "authIncorrectCreds":
		return ErrInvalidCredentials

	case "authUserInactiveOrDeleted":
		return ErrUserInactive

	case "authInvalidToken":
		// The server answers two different rejections with this one message: a
		// one-time code it does not accept for a login, and a session token it
		// does not accept for a loginByToken.
		if command == "loginByToken" {
			return ErrInvalidSessionToken
		}

		return ErrInvalidTOTPCode

	default:
		return unclassifiedLoginError(command, response)
	}
}

// unclassifiedLoginError reports a rejection this client version has no
// sentinel for, without passing a translation key off as a sentence and without
// producing an error with nothing in it when the server named no reason at all.
func unclassifiedLoginError(command string, response ackResponse) error {
	switch {
	case response.Msg == "":
		return fmt.Errorf("%s: the server rejected it without saying why", command)

	case response.MsgI18n:
		return fmt.Errorf("%s: the server reported the untranslated reason %q", command, response.Msg)

	default:
		return fmt.Errorf("%s: %s", command, response.Msg)
	}
}

// await reports how the server wants the connection to be authenticated, or
// authModeUnknown once it is clear that it will not say.
func (b authBarrier) await(ctx context.Context) authMode {
	timer := time.NewTimer(barrierWaitWithin(ctx))
	defer timer.Stop()

	var mode authMode

	select {
	case <-b.setupRequired:
		return authModeSetup

	case <-b.autoLogin:
		mode = authModeAutoLogin

	case <-b.session:
		mode = authModeSession

	case <-b.loginRequired:
		mode = authModeLoginRequired

	case <-timer.C:
		return authModeUnknown

	case <-ctx.Done():
		return authModeUnknown
	}

	// A server that wants to be set up emits setup before it states how it
	// wants to be authenticated, so by the time either of those events is in,
	// a setup that is coming has been delivered too. Answering it here rather
	// than from the select above is what keeps a server that sent both from
	// being read one way or the other at random.
	if isRequested(b.setupRequired) {
		return authModeSetup
	}

	return mode
}

// isRequested reports whether an event of the barrier has been delivered,
// without waiting for one that has not.
func isRequested(event <-chan struct{}) bool {
	select {
	case <-event:
		return true

	default:
		return false
	}
}

// barrierWaitWithin is how long the barrier may wait for a server that states
// nothing: authBarrierWait, but never more than half of what is left of the
// caller's deadline, so that such a server cannot spend the whole budget
// before the login it does answer is even sent.
func barrierWaitWithin(ctx context.Context) time.Duration {
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		return authBarrierWait
	}

	return min(authBarrierWait, max(time.Until(deadline)/2, 0))
}

// hasPassword reports whether a username and password were given.
func (c credentials) hasPassword() bool {
	return c.username != "" && c.password != ""
}

// hasAny reports whether there is anything at all to authenticate with.
func (c credentials) hasAny() bool {
	return c.hasPassword() || c.token != ""
}

// authenticate logs the client in the way the server asked for, and records the
// session token the connection ends up authenticated with, whether the server
// handed it out or accepted the one the client presented.
//
// It returns nil without logging in for a server that wants to be set up first,
// and for one that never stated what it wants while the client has no
// credentials to offer anyway. A server that authenticated the client itself is
// still logged in to when credentials were given, because that is what produces
// a session token, see Resync.
func (c *Client) authenticate(ctx context.Context, creds credentials, barrier authBarrier) error {
	mode := barrier.await(ctx)

	// await answers authModeUnknown both for a server that stated nothing and
	// for a wait the caller cut short. Only the context can tell those apart,
	// and reading it here keeps a cancellation from being reported as a server
	// that wants no login - which for a client without credentials would be a
	// successful New, and for one with them a login emitted into a connection
	// that is already being torn down.
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("waiting for the server's auth mode: %w", err)
	}

	switch mode {
	case authModeAutoLogin, authModeSession:
		c.setAutoLoggedIn()

		if !creds.hasAny() {
			return nil
		}

		// The connection is authenticated already, but a server with
		// authentication disabled keeps its login handler in place and a login
		// is the only thing that hands out a session token, which Resync and
		// Client.SessionToken need. Credentials given explicitly are therefore
		// still used.

	case authModeSetup:
		// The caller runs the setup and logs in afterwards.
		return nil

	case authModeLoginRequired:
		if !creds.hasAny() {
			return ErrAuthRequired
		}

	case authModeUnknown:
		if !creds.hasAny() {
			// A server that never said it wants a login is not asked for one.
			return nil
		}

	default:
		// authMode has no other values.
	}

	if creds.token != "" {
		tokenErr := c.loginWithToken(ctx, creds.token)
		if tokenErr == nil {
			return nil
		}

		if !recoverableWithPassword(tokenErr, creds) {
			return tokenErr
		}

		// A token the server no longer accepts is what a password change
		// leaves behind, and the password is what recovers from it.
		c.socketioLogger.Warnf(
			"session token rejected (%v), falling back to the password login for %q",
			tokenErr,
			creds.username,
		)

		err := c.loginWithPassword(ctx, creds.username, creds.password)
		if err != nil {
			// The rejected token is what sent the login down this path, which
			// the password failure on its own does not say.
			return fmt.Errorf("%w, and the password login that followed: %w", tokenErr, err)
		}

		c.setSessionTokenRejected()

		if c.sessionTokenRejectedCallback != nil {
			c.sessionTokenRejectedCallback(tokenErr)
		}

		return nil
	}

	return c.loginWithPassword(ctx, creds.username, creds.password)
}

// recoverableWithPassword reports whether a rejected session token is worth
// following with a password login: only a token the server refuses, and only
// when there is a password to offer. A deactivated or deleted user is the
// rejection the server answers the same way for both credentials, and anything
// that is not a rejection - a transport failure, a message the client does not
// know - says nothing about the password either way.
func recoverableWithPassword(err error, creds credentials) bool {
	return errors.Is(err, ErrInvalidSessionToken) &&
		!errors.Is(err, ErrUserInactive) &&
		creds.hasPassword()
}

// authenticateHTTP logs in over HTTP, see AuthSchemeHTTP, before the socket.io
// connection is made: the session cookie it leaves in the jar is what
// authenticates that connection on its handshake. It records the session token
// the client ends up with, see Client.SessionToken.
//
// It returns nil without logging in when there is nothing to log in with; the
// server then states on the connection whether it wants a login at all, see
// confirmSession. A server that is not set up yet rejects the credentials it
// has no user for, and New reads that together with the setup event, so the
// rejection is returned as it is.
func (c *Client) authenticateHTTP(ctx context.Context, creds credentials) error {
	if creds.token != "" {
		tokenErr := c.loginWithCookie(ctx, creds.token)
		if tokenErr == nil {
			return nil
		}

		if !recoverableWithPassword(tokenErr, creds) {
			return tokenErr
		}

		// A token the server no longer accepts is what a password change,
		// or on this server simply its expiry, leaves behind, and the
		// password is what recovers from it.
		c.socketioLogger.Warnf(
			"session token rejected (%v), falling back to the password login for %q",
			tokenErr,
			creds.username,
		)

		c.auth.clearSessionCookie()

		err := c.signInHTTP(ctx, creds.username, creds.password)
		if err != nil {
			return fmt.Errorf("%w, and the password login that followed: %w", tokenErr, err)
		}

		c.setSessionTokenRejected()

		if c.sessionTokenRejectedCallback != nil {
			c.sessionTokenRejectedCallback(tokenErr)
		}

		return nil
	}

	if !creds.hasPassword() {
		return nil
	}

	return c.signInHTTP(ctx, creds.username, creds.password)
}

// loginWithCookie presents a session token from an earlier login to a 3.x
// server, see WithSessionToken, by asking the server whether it still names a
// session. The token is a signed cookie the server hands out, and it expires,
// so the answer is checked here rather than read off the connection: a
// connection the server refuses to authenticate looks the same for an expired
// token and for a proxy that strips the cookie, and only the first is one a
// password recovers from.
func (c *Client) loginWithCookie(ctx context.Context, token string) error {
	c.auth.presetSessionCookie(token)

	ok, _, err := c.auth.getSession(ctx)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidSessionToken
	}

	c.setSessionToken(token)

	return nil
}

// signInHTTP logs in with a username and password over HTTP, answers the
// server's request for a one-time code if it makes one, and records the
// session cookie the server handed out.
func (c *Client) signInHTTP(ctx context.Context, username string, password string) error {
	twoFactor, err := c.auth.signIn(ctx, username, password)
	if err != nil {
		return err
	}

	if twoFactor {
		err = c.loginWithTOTP(ctx, func(ctx context.Context, code string) error {
			return c.auth.verifyTOTP(ctx, code)
		})
		if err != nil {
			return err
		}
	}

	return c.keepSessionCookie()
}

// keepSessionCookie records the session cookie an accepted login left in the
// jar, see keepSessionToken.
func (c *Client) keepSessionCookie() error {
	token := c.auth.sessionCookie()
	if token == "" {
		return errors.New("login: the server accepted the login but set no session cookie")
	}

	c.setSessionToken(token)

	return nil
}

// confirmSession reads how the server authenticated a connection made after
// authenticateHTTP, see AuthSchemeHTTP. loginErr is what authenticateHTTP
// reported, which New carried on past only for a rejection that a setup
// event explains.
//
// A server that authenticated the connection without being handed a session
// cookie has authentication disabled and is treated like the autoLogin of a
// 2.x server. One that asks for a login although a login over HTTP succeeded
// did not see the cookie on the handshake, which is what a proxy in between
// that drops it looks like; that is reported for what it is, because neither
// the password nor the token is at fault.
func (c *Client) confirmSession(
	ctx context.Context,
	creds credentials,
	barrier authBarrier,
	loginErr error,
) error {
	mode := barrier.await(ctx)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("waiting for the server's auth mode: %w", err)
	}

	switch mode {
	case authModeSetup:
		// The caller runs the setup and logs in afterwards.
		return nil

	case authModeSession:
		if c.SessionToken() == "" {
			c.setAutoLoggedIn()
		}

		return nil

	case authModeLoginRequired:
		if loginErr != nil {
			return loginErr
		}

		if c.SessionToken() != "" {
			return errors.New(
				"the server accepted the login over HTTP but did not see the session cookie " +
					"on the socket.io connection, check that a proxy in between forwards cookies",
			)
		}

		if !creds.hasAny() {
			return ErrAuthRequired
		}

		return errors.New("the server asks for a login although none of the credentials given were used")

	case authModeAutoLogin, authModeUnknown:
		if loginErr != nil {
			return loginErr
		}

		if !creds.hasAny() {
			// A server that never said it wants a login is not asked for one.
			return nil
		}

		return errors.New("the server stated neither a session nor that it wants a login")

	default:
		return fmt.Errorf("auth mode: unknown: %d", mode)
	}
}

// loginWithToken presents a session token from an earlier login, see
// WithSessionToken.
func (c *Client) loginWithToken(ctx context.Context, token string) error {
	response, err := c.emitAck(ctx, "loginByToken", token)
	if err != nil {
		return err
	}

	if !response.OK {
		return loginError("loginByToken", response)
	}

	// The server answers a loginByToken without a token of its own, so the one
	// that worked is the one to keep, see Client.SessionToken.
	c.setSessionToken(token)

	return nil
}

// loginWithPassword logs in with a username and password, and answers the
// server's request for a one-time code if it makes one.
func (c *Client) loginWithPassword(ctx context.Context, username string, password string) error {
	response, err := c.login(ctx, username, password, "")
	if err != nil {
		return err
	}

	if response.OK {
		return c.keepSessionToken(response)
	}

	// An ack asking for a one-time code carries neither an ok nor a message,
	// so it has to be recognized before the message is looked at.
	if response.TokenRequired {
		return c.loginWithTOTP(ctx, func(ctx context.Context, code string) error {
			return c.loginWithCode(ctx, username, password, code)
		})
	}

	return loginError("login", response)
}

// loginWithCode makes one socket login with a one-time code and interprets
// the ack the way loginWithTOTP's attempts are interpreted.
func (c *Client) loginWithCode(ctx context.Context, username string, password string, code string) error {
	response, err := c.login(ctx, username, password, code)
	if err != nil {
		return err
	}

	if response.OK {
		return c.keepSessionToken(response)
	}

	if response.TokenRequired {
		return fmt.Errorf(
			"%w: the server asked for a code again instead of answering the one it was given",
			ErrTwoFactorRequired,
		)
	}

	return loginError("login", response)
}

// totpAttempt presents one one-time code to the server: on the socket for a
// 2.x server, over HTTP for a 3.x one. A code the server rejects comes back
// wrapping ErrInvalidTOTPCode, which is the only outcome another attempt could
// change.
type totpAttempt func(ctx context.Context, code string) error

// loginWithTOTP answers the server's request for a one-time code.
//
// A code the server rejects is retried once in the step after the one it was
// generated in, because a 2.x server refuses a code it has accepted before and
// a later step is what produces a different one. Every other reason for the
// rejection - a wrong secret, a clock too far off - produces the same answer
// again, so a single retry is enough. The wait is skipped when the round trip
// already crossed into the next step, and the retry is given up on when the
// caller's deadline cannot cover the wait. A 3.x server has no such replay
// guard, so against it the retry only ever confirms a wrong secret.
//
// The ways this fails are kept apart in the error, because they call for
// different fixes: no code source, a source that produced nothing usable, a
// deadline too short to outlast the replay guard, and a code rejected in two
// consecutive steps, which is what rules the replay guard out.
func (c *Client) loginWithTOTP(ctx context.Context, attempt totpAttempt) error {
	if c.totpCode == nil {
		return ErrTwoFactorRequired
	}

	var rejectedStep time.Time

	for retry := range 2 {
		if retry > 0 {
			waitErr := awaitTOTPRetry(ctx, rejectedStep)
			if waitErr != nil {
				return waitErr
			}
		}

		step, err := c.loginOnceWithTOTP(ctx, attempt)
		if !errors.Is(err, ErrInvalidTOTPCode) {
			return err
		}

		rejectedStep = step
	}

	return fmt.Errorf(
		"%w: rejected in two consecutive time steps, which rules out the server's replay "+
			"guard, so what is left to check is the secret and this host's clock",
		ErrInvalidTOTPCode,
	)
}

// loginOnceWithTOTP makes one attempt at a login with a one-time code.
//
// It also reports the time step the code it sent belongs to, which is what
// tells a retry whether there is still a step to wait out. The step is read
// after the code was produced, so a code generated just before a step boundary
// counts as belonging to the later one - the reading that at worst waits when
// it did not have to, rather than sending the same code twice.
func (c *Client) loginOnceWithTOTP(ctx context.Context, attempt totpAttempt) (time.Time, error) {
	code, err := c.totpCode(ctx)

	codeStep := totpStepStart(time.Now())

	if err != nil {
		return codeStep, fmt.Errorf("login: %w", err)
	}

	// The server reads an empty code as no code at all and asks again instead
	// of rejecting it, in an ack carrying neither an ok nor a message, so
	// sending one buys an error with nothing in it.
	if code == "" {
		return codeStep, fmt.Errorf(
			"%w: the configured code source produced an empty code",
			ErrTwoFactorRequired,
		)
	}

	return codeStep, attempt(ctx, code)
}

// keepSessionToken records the token an accepted login handed out.
//
// A success without one would leave the client authenticated but unable to
// resync, and would surface far from here as Resync reporting a client that was
// created without credentials, so it is reported where it happens instead.
func (c *Client) keepSessionToken(response ackResponse) error {
	if response.Token == "" {
		return errors.New("login: the server accepted the login but handed out no session token")
	}

	c.setSessionToken(response.Token)

	return nil
}

// awaitTOTPRetry waits out the time step the rejected code belongs to, so the
// attempt after it has a code the server has not seen before, and says why it
// could not. A login whose round trip already crossed into a later step has
// nothing left to wait for.
func awaitTOTPRetry(ctx context.Context, rejectedStep time.Time) error {
	if totpStepStart(time.Now()).After(rejectedStep) {
		return nil
	}

	waited, err := waitForNextTOTPStep(ctx)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	if waited {
		return nil
	}

	return fmt.Errorf(
		"%w: recovering from a code the server has seen before needs a wait of up to %s "+
			"for the next time step, which does not fit in the deadline",
		ErrInvalidTOTPCode,
		totpStep,
	)
}

// login emits the login command with an optional one-time code and returns the
// ack, whether it reports success or not.
func (c *Client) login(
	ctx context.Context,
	username string,
	password string,
	code string,
) (ackResponse, error) {
	return c.emitAck(
		ctx,
		"login",
		map[string]any{"username": username, "password": password, "token": code},
	)
}

// setAutoLoggedIn records that the server logged the client in itself, see
// Client.Resync.
func (c *Client) setAutoLoggedIn() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.autoLoggedIn = true
}

// waitForNextTOTPStep waits until the current time step is over and reports
// whether it got there. A deadline that cannot cover the wait is left alone
// rather than run into, which is the (false, nil) answer; a wait the context
// cuts short returns that context's error, so a cancellation on the way is not
// reported as a code the server disliked.
func waitForNextTOTPStep(ctx context.Context) (bool, error) {
	wait := time.Until(totpStepStart(time.Now()).Add(totpStep))

	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline && time.Until(deadline) <= wait {
		return false, nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true, nil

	case <-ctx.Done():
		return false, fmt.Errorf("waiting for the next totp time step: %w", ctx.Err())
	}
}

// setupPending reports whether the server asked to be set up.
//
// The setup event and the ack of a login the server rejects because it has no
// users yet race each other, so an event that has not arrived is given a moment
// rather than being read as absent straight away. One that is already in costs
// no wait, and a cancelled caller is not made to sit out the grace period.
func setupPending(ctx context.Context, setupRequired <-chan struct{}) bool {
	timer := time.NewTimer(setupEventGrace)
	defer timer.Stop()

	select {
	case <-setupRequired:
		return true

	case <-timer.C:
		return false

	case <-ctx.Done():
		return false
	}
}

// enable2FA asks the server for a new two-factor secret and returns it.
//
// The server stores the secret right away, before it has been confirmed, and
// only starts requiring a code once confirm2FA ran. currentPassword is the
// password of the logged in account, which the server checks again for this
// command even though the connection is already authenticated.
func (c *Client) enable2FA(ctx context.Context, currentPassword string) (string, error) {
	if c.scheme == AuthSchemeHTTP {
		return c.auth.enableTwoFactor(ctx, currentPassword)
	}

	response, err := c.syncEmit(ctx, "prepare2FA", currentPassword)
	if err != nil {
		return "", err
	}

	uri, err := url.Parse(response.URI)
	if err != nil {
		return "", fmt.Errorf("prepare2FA: parse otpauth uri: %w", err)
	}

	secret := uri.Query().Get("secret")
	if secret == "" {
		return "", errors.New("prepare2FA: otpauth uri carries no secret")
	}

	return secret, nil
}

// confirm2FA turns two-factor authentication on for the logged in account,
// using the secret enable2FA handed out. A 3.x server wants a code generated
// from that secret as proof, a 2.x server the password again.
func (c *Client) confirm2FA(ctx context.Context, currentPassword string, code string) error {
	if c.scheme == AuthSchemeHTTP {
		return c.auth.verifyTOTP(ctx, code)
	}

	_, err := c.syncEmit(ctx, "save2FA", currentPassword)

	return err
}

// disable2FA turns two-factor authentication off for the logged in account.
func (c *Client) disable2FA(ctx context.Context, currentPassword string) error {
	if c.scheme == AuthSchemeHTTP {
		return c.auth.disableTwoFactor(ctx, currentPassword)
	}

	_, err := c.syncEmit(ctx, "disable2FA", currentPassword)

	return err
}

// twoFAStatus reports whether two-factor authentication is on for the logged
// in account.
func (c *Client) twoFAStatus(ctx context.Context) (bool, error) {
	if c.scheme == AuthSchemeHTTP {
		_, enabled, err := c.auth.getSession(ctx)

		return enabled, err
	}

	response, err := c.syncEmit(ctx, "twoFAStatus")
	if err != nil {
		return false, err
	}

	return response.Status, nil
}
