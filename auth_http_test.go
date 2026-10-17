package kuma_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// newHTTPAuthFake starts a fakeAckServer that authenticates the way an Uptime
// Kuma 3.x server does, see kuma.AuthSchemeHTTP, and returns it with its URL.
func newHTTPAuthFake(t *testing.T) (fake *fakeAckServer, url string) {
	t.Helper()

	fake = &fakeAckServer{messages: make(chan []byte, 32)}
	fake.setHTTPAuth()

	server := httptest.NewServer(fake)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return fake, server.URL
}

// TestNewDetectsAuthScheme covers how New tells a 2.x server from a 3.x one
// before it connects: by what the server answers on better-auth's health
// endpoint, which a 2.x server has no route for.
func TestNewDetectsAuthScheme(t *testing.T) {
	t.Run("a 2.x server answers the probe with its page", func(t *testing.T) {
		fake, url := newAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, kuma.AuthSchemeSocket, client.AuthScheme())
		require.Equal(t, 1, fake.loginFrameCount(), "the login goes over the socket")
	})

	t.Run("a 3.x server answers the probe with JSON", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, kuma.AuthSchemeHTTP, client.AuthScheme())
		require.Equal(t, 1, fake.probeCount())
		require.Empty(t, fake.receivedUnexpectedCommands(), "a 3.x server gets no socket login")
	})

	t.Run("a rate limited probe is still a 3.x answer", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setProbeRateLimited(true)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, kuma.AuthSchemeHTTP, client.AuthScheme())
	})

	t.Run("a pinned scheme skips the probe", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(
			ctx,
			url,
			"admin",
			"admin1",
			kuma.WithAuthScheme(kuma.AuthSchemeHTTP),
			kuma.WithConnectTimeout(5*time.Second),
		)
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, 0, fake.probeCount())
		require.Equal(t, kuma.AuthSchemeHTTP, client.AuthScheme())
	})

	t.Run("a pinned scheme the server does not speak fails the login", func(t *testing.T) {
		_, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		// The socket login is sent to a server that has no handler for it, so
		// it is never answered and the connect timeout is what ends it.
		client, err := kuma.New(
			ctx,
			url,
			"admin",
			"admin1",
			kuma.WithAuthScheme(kuma.AuthSchemeSocket),
			kuma.WithConnectTimeout(time.Second),
		)

		require.Error(t, err)
		require.Nil(t, client)
	})
}

// TestNewHTTPLogin covers the password login against a 3.x server: it is an
// HTTP request, the session it answers with is a cookie, and that cookie is
// what the socket.io handshake carries.
func TestNewHTTPLogin(t *testing.T) {
	fake, url := newHTTPAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 1, fake.loginFrameCount(), "one sign-in")
	require.Equal(t, fakeSessionToken, client.SessionToken(), "the session cookie is the token")
	require.False(t, client.SessionTokenRejected())
	require.Equal(t, 1, fake.handshakeCount())
	require.Contains(t, fake.handshakeCookies()[0], fakeSessionCookie+"="+fakeSessionToken,
		"the handshake has to carry the session cookie, it is all the server authenticates from")
	require.Empty(t, fake.receivedUnexpectedCommands())

	// The lists arrived on the authenticated connection.
	notif, err := client.GetNotification(ctx, fakeSeededNotificationID)
	require.NoError(t, err)
	require.Equal(t, int64(fakeSeededNotificationID), notif.GetID())
}

// TestNewHTTPLoginRejections covers the answers of a 3.x server that refuse a
// login, which are reported with the same sentinels a 2.x server's are.
func TestNewHTTPLoginRejections(t *testing.T) {
	tests := map[string]struct {
		prepare  func(fake *fakeAckServer)
		username string
		password string
		opts     []kuma.Option
		wantErr  error
		wantMsg  string
	}{
		"invalid credentials": {
			prepare:  func(fake *fakeAckServer) { fake.setRejectCreds() },
			username: "admin", password: "wrong",
			wantErr: kuma.ErrInvalidCredentials,
		},
		"two-factor required without a code source": {
			prepare:  func(fake *fakeAckServer) { fake.setTwoFactorRequired(true) },
			username: "admin", password: "admin1",
			wantErr: kuma.ErrTwoFactorRequired,
		},
		"no credentials against a server that wants a login": {
			prepare: func(*fakeAckServer) {},
			wantErr: kuma.ErrAuthRequired,
		},
		"rate limited": {
			prepare:  func(fake *fakeAckServer) { fake.setSignInRateLimited(true) },
			username: "admin", password: "admin1",
			wantErr: kuma.ErrRateLimited,
			wantMsg: "retry after 7s",
		},
		"a rejected token without a password": {
			prepare: func(fake *fakeAckServer) { fake.setRejectToken() },
			opts:    []kuma.Option{kuma.WithSessionToken("stale-token")},
			wantErr: kuma.ErrInvalidSessionToken,
		},
		"a proxy that strips the cookie": {
			prepare:  func(fake *fakeAckServer) { fake.setStripCookie(true) },
			username: "admin", password: "admin1",
			wantMsg: "did not see the session cookie",
		},
		"a login accepted without a session cookie": {
			prepare:  func(fake *fakeAckServer) { fake.setOmitLoginToken(true) },
			username: "admin", password: "admin1",
			wantMsg: "set no session cookie",
		},
		"setup required without autosetup": {
			prepare:  func(fake *fakeAckServer) { fake.setSetupRequired() },
			username: "admin", password: "admin1",
			wantMsg: "autosetup is disabled",
		},
		"setup required without credentials": {
			prepare: func(fake *fakeAckServer) { fake.setSetupRequired() },
			opts:    []kuma.Option{kuma.WithAutosetup()},
			wantMsg: "needs a username and password",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fake, url := newHTTPAuthFake(t)
			test.prepare(fake)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			client, err := kuma.New(
				ctx,
				url,
				test.username,
				test.password,
				append([]kuma.Option{kuma.WithConnectTimeout(5 * time.Second)}, test.opts...)...,
			)

			require.Error(t, err)
			require.Nil(t, client)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
			}

			if test.wantMsg != "" {
				require.ErrorContains(t, err, test.wantMsg)
			}

			require.Empty(t, fake.receivedUnexpectedCommands(), "a 3.x server gets no socket login")
		})
	}
}

// TestNewHTTPRateLimitedReportsTheWait covers the wait the server states with
// its refusal, which a caller that retries needs.
func TestNewHTTPRateLimitedReportsTheWait(t *testing.T) {
	fake, url := newHTTPAuthFake(t)
	fake.setSignInRateLimited(true)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	_, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))

	var limited *kuma.RateLimitedError

	require.ErrorAs(t, err, &limited)
	require.Equal(t, fakeRetryAfter*time.Second, limited.RetryAfter)
}

// TestNewHTTPTwoFactor covers the account with two-factor authentication
// enabled on a 3.x server: the sign-in asks for a code, which is answered with
// a second request.
func TestNewHTTPTwoFactor(t *testing.T) {
	fake, url := newHTTPAuthFake(t)
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

	t.Cleanup(func() { _ = client.Disconnect() })

	codes := fake.receivedLoginCodes()
	require.Len(t, codes, 2, "the sign-in and the code")
	require.Empty(t, codes[0], "the sign-in carries no code")
	require.NotEmpty(t, codes[1])
	require.Equal(t, fakeSessionToken, client.SessionToken())
}

// TestNewHTTPTwoFactorWrongCode covers the code a 3.x server rejects. It has no
// replay guard, so the retry in the next time step only confirms the secret is
// wrong; with a connect timeout shorter than a step the client gives up before
// that.
func TestNewHTTPTwoFactorWrongCode(t *testing.T) {
	fake, url := newHTTPAuthFake(t)
	fake.setTwoFactorSecret(false)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(
		ctx,
		url,
		"admin",
		"admin1",
		kuma.WithTOTPSecret("JBSWY3DPEHPK3PXP"),
		kuma.WithConnectTimeout(time.Second),
	)

	require.ErrorIs(t, err, kuma.ErrInvalidTOTPCode)
	require.Nil(t, client)
}

// TestNewHTTPSessionToken covers the session token of a 3.x server: the cookie
// a login handed out, presented again by a later client. The server is asked
// whether it still accepts it before the connection is made.
func TestNewHTTPSessionToken(t *testing.T) {
	t.Run("an accepted token needs neither password nor code", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setTwoFactorSecret(false)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(
			ctx,
			url,
			"",
			"",
			kuma.WithSessionToken(fakePresetToken),
			kuma.WithConnectTimeout(5*time.Second),
		)
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, 0, fake.loginFrameCount())
		require.Equal(t, []string{fakePresetToken}, fake.receivedTokens(), "the server is asked about the token")
		require.Equal(t, fakePresetToken, client.SessionToken())
		require.False(t, client.SessionTokenRejected())
		require.Contains(t, fake.handshakeCookies()[0], fakeSessionCookie+"="+fakePresetToken)
	})

	t.Run("a rejected token falls back to the password", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)
		fake.setRejectToken()

		var callbackErr error

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(
			ctx,
			url,
			"admin",
			"admin1",
			kuma.WithSessionToken("stale-token"),
			kuma.WithSessionTokenRejectedCallback(func(err error) { callbackErr = err }),
			kuma.WithConnectTimeout(5*time.Second),
		)
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, []string{"stale-token"}, fake.receivedTokens(), "the token is tried first")
		require.Equal(t, 1, fake.loginFrameCount(), "the password login is the fallback")
		require.Equal(t, fakeSessionToken, client.SessionToken(), "the fresh cookie replaces the stale one")
		require.True(t, client.SessionTokenRejected())
		require.ErrorIs(t, callbackErr, kuma.ErrInvalidSessionToken)
		require.Contains(t, fake.handshakeCookies()[0], fakeSessionCookie+"="+fakeSessionToken,
			"the handshake carries the fresh cookie, not the stale one")
	})
}

// TestNewHTTPAutoLogin covers the 3.x server with authentication disabled: it
// authenticates the connection without a cookie and states a session where a
// 2.x server states autoLogin.
func TestNewHTTPAutoLogin(t *testing.T) {
	fake, url := newHTTPAuthFake(t)
	fake.setAutoLogin()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "", "", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 0, fake.loginFrameCount())
	require.Empty(t, client.SessionToken())

	// Such a server authenticates the next connection the same way, so a
	// resync needs no credentials here, unlike on a 2.x server.
	require.NoError(t, client.Resync(ctx))
	require.Equal(t, 2, fake.handshakeCount())
}

// TestNewHTTPAutosetup covers the server that is not set up yet: the setup is
// an HTTP request, and the connection made before it, which the server could
// not authenticate, is replaced by one it can.
func TestNewHTTPAutosetup(t *testing.T) {
	fake, url := newHTTPAuthFake(t)
	fake.setSetupRequired()

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

	t.Cleanup(func() { _ = client.Disconnect() })

	require.Equal(t, 2, fake.loginFrameCount(), "the sign-in the setup explains, and the one after it")
	require.Equal(t, 2, fake.handshakeCount(), "the connection made before the setup carries no session")
	require.Equal(t, fakeSessionToken, client.SessionToken())
	require.Empty(t, fake.receivedUnexpectedCommands(), "the setup command of a 2.x server is not sent")

	notif, err := client.GetNotification(ctx, fakeSeededNotificationID)
	require.NoError(t, err)
	require.Equal(t, int64(fakeSeededNotificationID), notif.GetID())
}

// TestResyncHTTP covers the resync against a 3.x server, which has no command
// that resends every list: the client connects again with its cookie, and the
// server answers the new connection the way it answered the first.
func TestResyncHTTP(t *testing.T) {
	t.Run("connects again with the same cookie", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.NoError(t, client.Resync(ctx))

		cookies := fake.handshakeCookies()
		require.Len(t, cookies, 2)
		require.Equal(t, cookies[0], cookies[1])
		require.Empty(t, fake.receivedUnexpectedCommands(), "no loginByToken on a 3.x server")

		notif, err := client.GetNotification(ctx, fakeSeededNotificationID)
		require.NoError(t, err, "the resync has to leave the cache complete, not empty")
		require.Equal(t, int64(fakeSeededNotificationID), notif.GetID())

		// The new connection is the one commands go out on.
		msg, err := client.TestNotification(ctx, notification.Generic{
			Base:           notification.Base{Name: "after resync"},
			GenericDetails: notification.GenericDetails{},
			TypeName:       "generic",
		})
		require.NoError(t, err)
		require.Equal(t, "ok", msg)
	})

	t.Run("a command in flight finishes before the connection is replaced", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		const ackDelay = 300 * time.Millisecond

		fake.setAckDelay(ackDelay)

		type result struct {
			msg   string
			err   error
			ackAt time.Time
		}

		done := make(chan result, 1)

		go func() {
			msg, err := client.TestNotification(ctx, notification.Generic{
				Base:           notification.Base{Name: "in flight"},
				GenericDetails: notification.GenericDetails{},
				TypeName:       "generic",
			})
			done <- result{msg: msg, err: err, ackAt: time.Now()}
		}()

		// Give the command a moment to be emitted, then resync underneath it.
		time.Sleep(50 * time.Millisecond)

		fake.setAckDelay(0)

		require.NoError(t, client.Resync(ctx))

		got := <-done
		require.NoError(t, got.err, "the resync must not cost the command its ack")
		require.Equal(t, "ok", got.msg)
		require.False(t, fake.lastHandshakeAt().Before(got.ackAt),
			"the new connection must be made after the command in flight was answered")
	})

	t.Run("a revoked cookie is reported and the old connection kept", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		fake.revokeSessions()

		err = client.Resync(ctx)
		require.ErrorIs(t, err, kuma.ErrInvalidSessionToken)

		msg, err := client.TestNotification(ctx, notification.Generic{
			Base:           notification.Base{Name: "still served"},
			GenericDetails: notification.GenericDetails{},
			TypeName:       "generic",
		})
		require.NoError(t, err, "the connection the server did authenticate stays in service")
		require.Equal(t, "ok", msg)
	})

	t.Run("waits for the lists and names the ones that did not arrive", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)

		t.Cleanup(func() { _ = client.Disconnect() })

		fake.setSuppressLists(true)

		resyncCtx, resyncCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer resyncCancel()

		err = client.Resync(resyncCtx)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.ErrorContains(t, err, "monitorList")
	})

	t.Run("after a disconnect it reports that instead of connecting", func(t *testing.T) {
		fake, url := newHTTPAuthFake(t)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
		require.NoError(t, err)
		require.NoError(t, client.Disconnect())

		err = client.Resync(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "disconnected")
		require.Equal(t, 1, fake.handshakeCount())
	})
}

// TestServerVersion covers the version the server states to an authenticated
// connection, which the info event after the login carries.
func TestServerVersion(t *testing.T) {
	fake, url := newAuthFake(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	client, err := kuma.New(ctx, url, "admin", "admin1", kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = client.Disconnect() })

	fake.send([]byte(`42["info",{"version":"2.5.6","latestVersion":"2.5.6"}]`))

	require.Eventually(t, func() bool {
		return client.ServerVersion() == "2.5.6"
	}, 5*time.Second, 10*time.Millisecond)
}
