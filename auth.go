package kuma

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
// authentication enabled, which the client cannot answer yet: the server asks
// for a one-time code and the client has none to give.
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
// ErrInvalidSessionToken as well, so a caller that only cares that the token is
// gone keeps matching on that one.
//
// It is the token rejection a password cannot recover from: the server requires
// an active user for a password login too, so New reports it instead of falling
// back, see WithSessionToken.
var ErrUserInactive = errors.New("user inactive or deleted")

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
// The server records the last code it accepted and refuses to see it again, so
// two clients logging in with the same account inside one 30 second step
// cannot both succeed. The client retries once in the next step when the
// caller's deadline allows it, but a caller starting many clients at once is
// better served by WithSessionToken.
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

		c.totpCodeSet = true
		c.totpCode = func(_ context.Context) (string, error) {
			return totpCodeAt(normalized, time.Now())
		}
	}
}

// WithTOTPCode configures a callback that produces the one-time code, for a
// caller whose secret lives somewhere the client cannot read it, such as a
// hardware token or an external authenticator.
//
// It is called only when the server asks for a code, and once more for the
// single retry the client makes in the next time step, see ErrInvalidTOTPCode.
// WithTOTPSecret is this option with the code computed from a secret, and the
// two are mutually exclusive.
func WithTOTPCode(code func(ctx context.Context) (string, error)) Option {
	return func(c *Client) {
		if code == nil {
			return
		}

		c.totpCodeSet = true
		c.totpSources++
		c.totpCode = code
	}
}

// authBarrierWait bounds how long the client waits for the server to state how
// it wants to be authenticated, see authBarrier.
//
// The server emits loginRequired or autoLogin as the last thing it does for a
// new connection, so on a healthy server the wait is over in milliseconds. It
// only elapses in full for a server too old to send either event, or behind a
// proxy that drops it, and it must not be a hard failure for those.
const authBarrierWait = 500 * time.Millisecond

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
	setupRequired <-chan struct{}
}

// credentials are the ways New was given to authenticate, gathered from its
// arguments and options so the login path has a single thing to read.
type credentials struct {
	username string
	password string
	token    string
}

// loginError translates a rejected login ack into one of the package's sentinel
// errors, so a caller can tell the cases apart with errors.Is instead of
// matching on the message.
func loginError(command string, response ackResponse) error {
	switch response.Msg {
	case "authIncorrectCreds", "Incorrect username or password":
		return ErrInvalidCredentials

	case "authUserInactiveOrDeleted":
		return fmt.Errorf("%w: %w", ErrInvalidSessionToken, ErrUserInactive)

	case "authInvalidToken":
		// The server answers two different rejections with this one message: a
		// one-time code it does not accept for a login, and a session token it
		// does not accept for a loginByToken.
		if command == "loginByToken" {
			return ErrInvalidSessionToken
		}

		return ErrInvalidTOTPCode

	default:
		return fmt.Errorf("%s: %s", command, response.Msg)
	}
}

// await reports how the server wants the connection to be authenticated, or
// authModeUnknown once it is clear that it will not say.
func (b authBarrier) await(ctx context.Context) authMode {
	timer := time.NewTimer(authBarrierWait)
	defer timer.Stop()

	select {
	case <-b.autoLogin:
		return authModeAutoLogin

	case <-b.loginRequired:
		return authModeLoginRequired

	case <-b.setupRequired:
		return authModeSetup

	case <-timer.C:
		return authModeUnknown

	case <-ctx.Done():
		return authModeUnknown
	}
}

// isSet reports whether a username and password were given.
func (c credentials) isSet() bool {
	return c.username != "" && c.password != ""
}

// hasAny reports whether there is anything at all to authenticate with.
func (c credentials) hasAny() bool {
	return c.isSet() || c.token != ""
}

// authenticate logs the client in the way the server asked for, and records the
// session token the connection ends up authenticated with, whether the server
// handed it out or accepted the one the client presented.
//
// It returns nil without logging in for a server that authenticated the client
// itself, for one that wants to be set up first, and for one that never stated
// what it wants while the client has no credentials to offer anyway.
func (c *Client) authenticate(ctx context.Context, creds credentials, barrier authBarrier) error {
	switch barrier.await(ctx) {
	case authModeAutoLogin:
		c.setAutoLoggedIn()

		return nil

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
		creds.isSet()
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
		c.setSessionToken(response.Token)

		return nil
	}

	// An ack asking for a one-time code carries neither an ok nor a message,
	// so it has to be recognized before the message is looked at.
	if response.TokenRequired {
		return c.loginWithTOTP(ctx, username, password)
	}

	return loginError("login", response)
}

// loginWithTOTP answers the server's request for a one-time code.
//
// A code the server rejects is retried once in the next time step, because the
// server refuses a code it has accepted before and the next step is what
// produces a different one. Every other reason for the rejection - a wrong
// secret, a clock too far off - produces the same answer again, so a single
// retry is enough, and it is skipped altogether when the caller's deadline
// cannot cover the wait.
func (c *Client) loginWithTOTP(ctx context.Context, username string, password string) error {
	if !c.totpCodeSet {
		return ErrTwoFactorRequired
	}

	for retry := range 2 {
		if retry > 0 && !waitForNextTOTPStep(ctx) {
			break
		}

		code, err := c.totpCode(ctx)
		if err != nil {
			return fmt.Errorf("login: totp code: %w", err)
		}

		response, err := c.login(ctx, username, password, code)
		if err != nil {
			return err
		}

		if response.OK {
			c.setSessionToken(response.Token)

			return nil
		}

		if !errors.Is(loginError("login", response), ErrInvalidTOTPCode) {
			return loginError("login", response)
		}
	}

	return fmt.Errorf(
		"%w: the server also rejects a code it has accepted before, "+
			"so another login with this account may have just used this one",
		ErrInvalidTOTPCode,
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

// waitForNextTOTPStep waits until the current time step is over, and reports
// whether the wait completed inside the caller's budget. A deadline that
// cannot cover the wait is left alone rather than run into.
func waitForNextTOTPStep(ctx context.Context) bool {
	wait := time.Until(totpStepStart(time.Now()).Add(totpStep))

	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline && time.Until(deadline) <= wait {
		return false
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true

	case <-ctx.Done():
		return false
	}
}

// setupPending reports whether the server asked to be set up.
//
// The setup event and the ack of a login the server rejects because it has no
// users yet race each other, so the check is made after a moment rather than
// straight away.
func setupPending(setupRequired <-chan struct{}) bool {
	time.Sleep(10 * time.Millisecond)

	select {
	case <-setupRequired:
		return true

	default:
		return false
	}
}

// prepare2FA asks the server for a new two-factor secret and returns it.
//
// The server stores the secret right away, before it has been confirmed, and
// only starts requiring a code once save2FA ran. currentPassword is the
// password of the logged in account, which the server checks again for this
// command even though the connection is already authenticated.
func (c *Client) prepare2FA(ctx context.Context, currentPassword string) (string, error) {
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

// save2FA turns two-factor authentication on for the logged in account, using
// the secret prepare2FA handed out.
func (c *Client) save2FA(ctx context.Context, currentPassword string) error {
	_, err := c.syncEmit(ctx, "save2FA", currentPassword)

	return err
}

// disable2FA turns two-factor authentication off for the logged in account.
func (c *Client) disable2FA(ctx context.Context, currentPassword string) error {
	_, err := c.syncEmit(ctx, "disable2FA", currentPassword)

	return err
}

// twoFAStatus reports whether two-factor authentication is on for the logged
// in account.
func (c *Client) twoFAStatus(ctx context.Context) (bool, error) {
	response, err := c.syncEmit(ctx, "twoFAStatus")
	if err != nil {
		return false, err
	}

	return response.Status, nil
}
