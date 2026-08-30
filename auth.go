package kuma

import (
	"context"
	"errors"
	"fmt"
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
}

// loginError translates a rejected login ack into one of the package's sentinel
// errors, so a caller can tell the cases apart with errors.Is instead of
// matching on the message.
func loginError(command string, response ackResponse) error {
	switch response.Msg {
	case "authIncorrectCreds", "Incorrect username or password":
		return ErrInvalidCredentials

	case "authUserInactiveOrDeleted":
		return fmt.Errorf("%w: user inactive or deleted", ErrInvalidSessionToken)

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

// authenticate logs the client in the way the server asked for, and records the
// session token a successful login answers with.
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
		if !creds.isSet() {
			return ErrAuthRequired
		}

	case authModeUnknown:
		if !creds.isSet() {
			// A server that never said it wants a login is not asked for one.
			return nil
		}

	default:
		// authMode has no other values.
	}

	return c.loginWithPassword(ctx, creds.username, creds.password)
}

// loginWithPassword logs in with a username and password.
func (c *Client) loginWithPassword(ctx context.Context, username string, password string) error {
	response, err := c.emitAck(
		ctx,
		"login",
		map[string]any{"username": username, "password": password, "token": ""},
	)
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
		return ErrTwoFactorRequired
	}

	return loginError("login", response)
}

// setAutoLoggedIn records that the server logged the client in itself, see
// Client.Resync.
func (c *Client) setAutoLoggedIn() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.autoLoggedIn = true
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
