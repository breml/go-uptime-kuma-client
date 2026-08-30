package kuma_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// newAuthFake starts a fakeAckServer and returns it together with its URL, so
// each test can set the auth behaviour it needs before connecting.
func newAuthFake(t *testing.T) (fake *fakeAckServer, url string) {
	t.Helper()

	fake = &fakeAckServer{messages: make(chan []byte, 32)}
	server := httptest.NewServer(fake)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return fake, server.URL
}

// TestNewTwoFactorRequired covers the account with two-factor authentication
// enabled. The server answers such a login with an ack that carries neither an
// ok nor a message, which without the typed error surfaces as an error with an
// empty message.
func TestNewTwoFactorRequired(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setTwoFactorRequired(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))

	require.ErrorIs(t, err, kuma.ErrTwoFactorRequired)
	require.Nil(t, client)
	require.Equal(t, 1, fake.loginFrameCount())
}

// TestNewInvalidCredentials covers the server that rejects the username and
// password, without ever asking to be set up.
func TestNewInvalidCredentials(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setRejectCreds(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "wrong", kuma.WithConnectTimeout(5*time.Second))

	require.ErrorIs(t, err, kuma.ErrInvalidCredentials)
	require.Nil(t, client)
}

// TestNewAuthRequiredWithoutCredentials covers the caller that offers no
// credentials to a server that asks for a login. Before the client observed
// that request it had nothing to react to and waited out its connect timeout
// instead.
func TestNewAuthRequiredWithoutCredentials(t *testing.T) {
	_, url := newAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	start := time.Now()
	client, err := kuma.New(ctx, url, "", "", kuma.WithConnectTimeout(5*time.Second))
	elapsed := time.Since(start)

	require.ErrorIs(t, err, kuma.ErrAuthRequired)
	require.Nil(t, client)
	require.Less(t, elapsed, 3*time.Second, "must not wait out the connect timeout")
}

// TestNewAutoLogin covers the server with authentication disabled: it logs the
// client in itself and broadcasts the lists unprompted, so the client must not
// send a login at all.
func TestNewAutoLogin(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setAutoLogin(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "", "", kuma.WithConnectTimeout(5*time.Second))

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 0, fake.loginFrameCount())

	// The server hands out no session token on this path, and has no command
	// that resends every list, so a resync cannot work.
	err = client.Resync(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "authentication disabled")
}

// TestNewWithoutLoginRequiredEvent covers the server that never states how it
// wants to be authenticated, as a version too old to send the event or a proxy
// dropping it does. The client falls back to logging in anyway.
func TestNewWithoutLoginRequiredEvent(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setOmitLoginRequired(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 1, fake.loginFrameCount())
}

// TestNewLoginWaitsForLoginRequired covers the ordering the barrier exists for:
// the server registers its login handler after an await and only then says it
// wants a login, so a login sent before that can be dropped without a trace.
func TestNewLoginWaitsForLoginRequired(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setLoginRequiredDelay(200 * time.Millisecond)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.True(
		t,
		fake.firstLoginAfterLoginRequired(),
		"the login must not be sent before the server asked for one",
	)
}

// TestNewSetupFallbackAfterRejectedCredentials covers the server that is not
// set up yet: it has no user to match the credentials against, so it rejects
// them and asks to be set up instead. Running the setup is what makes the same
// credentials work, so the rejection must not end the connect.
//
// This is the regression guard for classifying that rejection by the message
// the server sent rather than by matching on a formatted error string.
func TestNewSetupFallbackAfterRejectedCredentials(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setSetupRequired(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithAutosetup(),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })
}

// TestNewRejectsHalfCredentials covers the caller that sets only one of the two
// credentials, which used to reach the server as a login that could not
// succeed.
func TestNewRejectsHalfCredentials(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	client, err := kuma.New(ctx, "http://127.0.0.1:1", "admin", "")

	require.Error(t, err)
	require.Contains(t, err.Error(), "have to be set together")
	require.Nil(t, client)
}

// TestResyncRejectedToken covers the session token the server no longer
// accepts, which is what a password change leaves behind.
func TestResyncRejectedToken(t *testing.T) {
	fake, url := newAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	fake.setRejectToken(true)

	resyncCtx, resyncCancel := context.WithTimeout(ctx, 2*time.Second)
	defer resyncCancel()

	err = client.Resync(resyncCtx)

	require.ErrorIs(t, err, kuma.ErrInvalidSessionToken)
}

// TestNewWithTOTPSecret covers the account with two-factor authentication
// enabled against a client that holds the shared secret: the server asks for a
// one-time code and the client answers it without a human at the keyboard.
func TestNewWithTOTPSecret(t *testing.T) {
	// The code is compared against one generated after New returned, so both
	// have to fall into the same step for the comparison to mean anything.
	requireTOTPStepHeadroom(t, 5*time.Second)

	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(false)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPSecret(rfc6238Secret),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	// The code is generated only once the server asks for one, because the
	// server evaluates the fields of a login independently and would verify a
	// code it never asked for against a secret that does not exist.
	codes := fake.receivedLoginCodes()
	require.Len(t, codes, 2)
	require.Empty(t, codes[0])

	want, err := kuma.TOTPCodeAt(rfc6238Secret, time.Now())
	require.NoError(t, err)
	require.Equal(t, want, codes[1])
}

// TestNewWithTOTPCode covers the caller whose secret lives somewhere the client
// cannot read it, so it hands over a callback instead.
func TestNewWithTOTPCode(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(false)

	calls := 0

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPCode(func(context.Context) (string, error) {
			calls++

			return kuma.TOTPCodeAt(rfc6238Secret, time.Now())
		}),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 1, calls, "the code is asked for only once the server wants one")
}

// TestNewTOTPReplayGuardTightDeadline covers the code the server has accepted
// before, which it refuses to see again. Recovering means waiting out the
// current time step, so a caller whose deadline cannot cover that wait has to
// be told promptly instead of being run into it.
func TestNewTOTPReplayGuardTightDeadline(t *testing.T) {
	// The deadline has to be shorter than the wait for the next step for this
	// to cover anything, and how long that wait is depends on the time of day,
	// so the step is given room for both before either is measured.
	const tightDeadline = 3 * time.Second

	requireTOTPStepHeadroom(t, 2*tightDeadline)

	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(true)

	// Use up the code of the current step, the way another client logging in
	// with the same account just before would have.
	code, err := kuma.TOTPCodeAt(rfc6238Secret, time.Now())
	require.NoError(t, err)
	fake.useCode(code)

	ctx, cancel := context.WithTimeout(t.Context(), tightDeadline)
	defer cancel()

	start := time.Now()
	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPSecret(rfc6238Secret),
		kuma.WithConnectTimeout(3*time.Second),
	)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, kuma.ErrInvalidTOTPCode)
	require.Nil(t, client)
	require.Less(t, elapsed, tightDeadline, "must not run into the deadline waiting for the next step")
}

// requireTOTPStepHeadroom leaves at least want of the current time step, by
// sleeping out a step that is nearly over.
//
// A test that says something about the wait for the next step otherwise says
// it only for the time of day it happens to run at: near the end of a step
// that wait is short, and a deadline meant to be too small to cover it covers
// it after all.
func requireTOTPStepHeadroom(t *testing.T, want time.Duration) {
	t.Helper()

	const step = 30 * time.Second

	remaining := time.Until(time.Now().Truncate(step).Add(step))
	if remaining > want {
		return
	}

	time.Sleep(remaining + 100*time.Millisecond)
}

// TestNewTOTPSecretRejected covers the secret the client cannot decode, which
// fails the connect rather than the login it would be needed for.
func TestNewTOTPSecretRejected(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	client, err := kuma.New(ctx, "http://127.0.0.1:1", "admin", "admin1", kuma.WithTOTPSecret("!!!"))

	require.Error(t, err)
	require.Contains(t, err.Error(), "totp secret")
	require.Nil(t, client)
}

// TestNewTOTPOptionsConflict covers configuring both ways of producing a code.
func TestNewTOTPOptionsConflict(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		"http://127.0.0.1:1",
		"admin",
		"admin1",
		kuma.WithTOTPSecret(rfc6238Secret),
		kuma.WithTOTPCode(func(context.Context) (string, error) { return "000000", nil }),
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "at most one")
	require.Nil(t, client)
}

// TestNewTOTPCodeNil covers the nil callback, which is a caller who meant to
// configure a code source and did not. It fails the connect rather than
// leaving the client to report at login time that nothing was configured.
func TestNewTOTPCodeNil(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	client, err := kuma.New(ctx, "http://127.0.0.1:1", "admin", "admin1", kuma.WithTOTPCode(nil))

	require.Error(t, err)
	require.Contains(t, err.Error(), "nil callback")
	require.Nil(t, client)
}

// TestNewTOTPOptionsConflictReportsTheConflictFirst covers both options given
// with a secret that also fails to decode: the conflict is the caller's first
// mistake, so fixing the secret must not be what uncovers it.
func TestNewTOTPOptionsConflictReportsTheConflictFirst(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		"http://127.0.0.1:1",
		"admin",
		"admin1",
		kuma.WithTOTPSecret("!!!"),
		kuma.WithTOTPCode(func(context.Context) (string, error) { return "000000", nil }),
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "at most one")
	require.Nil(t, client)
}

// TestNewTOTPEmptyCode covers the callback that produces nothing, such as a
// hardware token prompt the user dismissed. The server reads an empty code as
// no code and asks again, in an ack carrying neither an ok nor a message, so
// sending it would buy an error with nothing in it.
func TestNewTOTPEmptyCode(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(false)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPCode(func(context.Context) (string, error) { return "", nil }),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.ErrorIs(t, err, kuma.ErrTwoFactorRequired)
	require.Nil(t, client)
	require.Empty(t, fake.receivedLoginCodes()[1:], "an empty code is not worth sending")
}

// TestNewTOTPCodeFails covers the callback that reports an error, which is the
// caller's own failure and has to reach them as one rather than as a rejected
// code.
func TestNewTOTPCodeFails(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(false)

	wantErr := errors.New("hardware token unplugged")

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPCode(func(context.Context) (string, error) { return "", wantErr }),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.ErrorIs(t, err, wantErr)
	require.NotErrorIs(t, err, kuma.ErrInvalidTOTPCode)
	require.Nil(t, client)
}

// TestNewTOTPRejectedForAnotherReason covers the coded login the server
// rejects for something other than the code. Only a rejected code is worth the
// retry, so anything else has to reach the caller as itself rather than as the
// replay guard the retry would end in.
func TestNewTOTPRejectedForAnotherReason(t *testing.T) {
	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(false)
	fake.setRejectCodedLogin(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPSecret(rfc6238Secret),
		kuma.WithConnectTimeout(5*time.Second),
	)

	require.ErrorIs(t, err, kuma.ErrInvalidCredentials)
	require.NotErrorIs(t, err, kuma.ErrInvalidTOTPCode)
	require.Nil(t, client)
	require.Len(t, fake.receivedLoginCodes(), 2, "a rejection the retry cannot change is not retried")
}

// TestNewTOTPWaitCancelled covers the caller who gives up while the client is
// waiting out the time step: the cancellation is what happened, and reporting
// it as a rejected code would send them looking for a second client that used
// the same one.
func TestNewTOTPWaitCancelled(t *testing.T) {
	requireTOTPStepHeadroom(t, 5*time.Second)

	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(true)

	code, err := kuma.TOTPCodeAt(rfc6238Secret, time.Now())
	require.NoError(t, err)
	fake.useCode(code)

	// No deadline, so the client commits to the wait instead of skipping it.
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go func() {
		for range 100 {
			if len(fake.receivedLoginCodes()) > 1 {
				break
			}

			time.Sleep(20 * time.Millisecond)
		}

		cancel()
	}()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithTOTPSecret(rfc6238Secret))

	require.ErrorIs(t, err, context.Canceled)
	require.NotErrorIs(t, err, kuma.ErrInvalidTOTPCode)
	require.Nil(t, client)
}

// TestNewTOTPReplayGuardRecovers covers the same rejection with a deadline
// generous enough to wait the time step out: the next step produces a code the
// server has not seen, so the login goes through without the caller doing
// anything.
//
// Waiting for the step boundary takes up to 30 seconds, so this runs only when
// E2E_TEST is set, even though it drives the fake rather than a real server.
func TestNewTOTPReplayGuardRecovers(t *testing.T) {
	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if !e2eTest {
		t.Skip("skipping test that waits out a TOTP time step")
	}

	fake, url := newAuthFake(t)
	fake.setTwoFactorSecret(true)

	code, err := kuma.TOTPCodeAt(rfc6238Secret, time.Now())
	require.NoError(t, err)
	fake.useCode(code)

	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPSecret(rfc6238Secret),
		kuma.WithConnectTimeout(60*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() { _ = client.Disconnect() })

	codes := fake.receivedLoginCodes()
	require.Len(t, codes, 3, "the login without a code, the rejected one and the retry")
	require.NotEqual(t, codes[1], codes[2], "the retry has to fall into a later time step")
}
