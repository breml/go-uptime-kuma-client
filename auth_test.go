package kuma_test

import (
	"context"
	"net/http/httptest"
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
