package kuma_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// newKumaContainer starts an Uptime Kuma of its own and returns its base URL
// and the auth scheme it speaks.
//
// The tests below turn two-factor authentication on for the admin account,
// which every other test in this package logs in as, so they cannot share the
// server TestMain brings up.
func newKumaContainer(t *testing.T) (baseURL string, scheme kuma.AuthScheme) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	repository, tag := kumaImage()

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       fmt.Sprintf("uptime-kuma-auth-%s", randomString(8)),
		Repository: repository,
		Tag:        tag,
		Env:        kumaContainerEnv(),
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(t, err)

	require.NoError(t, resource.Expire(600))

	t.Cleanup(func() {
		purgeErr := pool.Purge(resource)
		if purgeErr != nil {
			log.Printf("Could not purge resource: %v", purgeErr)
		}
	})

	baseURL = "http://localhost:" + resource.GetPort("3001/tcp")

	// The server needs a moment before it accepts connections, and the first
	// client is the one that runs the setup.
	err = pool.Retry(func() error {
		setupClient, setupErr := kuma.New(
			t.Context(),
			baseURL,
			testUsername, testPassword,
			kuma.WithAutosetup(),
			kuma.WithLogLevel(kuma.LogLevel(os.Getenv("SOCKETIO_LOG_LEVEL"))),
			kuma.WithConnectTimeout(10*time.Second),
		)
		if setupErr != nil {
			return setupErr
		}

		scheme = setupClient.AuthScheme()

		return setupClient.Disconnect()
	})
	require.NoError(t, err)

	return baseURL, scheme
}

// awaitSignInWindow lets the sign-in rate limit of a 3.x server, three per ten
// seconds per client address, recover before the next login. A 2.x server
// allows 20 per minute, which the tests stay inside without waiting.
func awaitSignInWindow(scheme kuma.AuthScheme) {
	if scheme == kuma.AuthSchemeHTTP {
		time.Sleep(11 * time.Second)
	}
}

// TestAuthenticationAgainstServer covers the authentication paths that only a
// real server can show: it turns two-factor authentication on for the admin
// account and then logs in the ways a client without a human at the keyboard
// has.
//
// Note the rate limits: a 2.x server allows 20 logins per minute and 30
// two-factor commands per minute, a 3.x server three sign-ins per ten seconds
// per client address, see kuma.RateLimitedError. The sequence below stays
// inside the first; against a 3.x server it waits the window out between the
// subtests that sign in, see awaitSignInWindow.
func TestAuthenticationAgainstServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// This one brings up a second Uptime Kuma of its own and waits out real
	// time steps, so it is gated the way the other end-to-end tests are rather
	// than run on every task test.
	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if !e2eTest {
		t.Skip(`skipping end to end test, "E2E_TEST" env var not set`)
	}

	baseURL, scheme := newKumaContainer(t)

	ctx, cancel := context.WithTimeout(t.Context(), 4*time.Minute)
	defer cancel()

	// The setup signs in twice: once to learn that the server wants a setup,
	// once after it.
	awaitSignInWindow(scheme)

	admin, err := kuma.New(ctx, baseURL, testUsername, testPassword, kuma.WithConnectTimeout(10*time.Second))
	require.NoError(t, err)

	defer func() { _ = admin.Disconnect() }()

	// The server checks the current password again for these commands, even
	// though the connection is already authenticated.
	secret, err := kuma.Enable2FA(admin, ctx, testPassword)
	if err != nil && scheme == kuma.AuthSchemeHTTP && strings.Contains(err.Error(), "status 500") {
		// The 3.0 nightly ships a better-auth that expects a column its own
		// 2FA migration does not create (better_auth_twoFactor has no
		// failedVerificationCount), so enabling two-factor authentication
		// fails on the server. Nothing below can run without it.
		t.Skipf("the server cannot enable two-factor authentication: %v", err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, secret)

	normalized, err := kuma.NormalizeTOTPSecret(secret)
	require.NoError(t, err)

	// A 3.x server wants a code from the new secret as proof before it turns
	// two-factor authentication on, a 2.x server the password again.
	code, err := kuma.TOTPCodeAt(normalized, time.Now())
	require.NoError(t, err)

	require.NoError(t, kuma.Confirm2FA(admin, ctx, testPassword, code))

	// Deferred rather than registered as cleanup: the test context is already
	// cancelled by the time cleanups run.
	defer func() { _ = kuma.Disable2FA(admin, ctx, testPassword) }()

	enabled, err := kuma.TwoFAStatus(admin, ctx)
	require.NoError(t, err)
	require.True(t, enabled)

	var sessionToken string

	awaitSignInWindow(scheme)

	t.Run("a login without a code is told what is missing", func(t *testing.T) {
		client, newErr := kuma.New(
			ctx,
			baseURL,
			testUsername,
			testPassword,
			kuma.WithConnectTimeout(10*time.Second),
		)

		require.ErrorIs(t, newErr, kuma.ErrTwoFactorRequired)
		require.Nil(t, client)
	})

	t.Run("the secret answers the request for a code", func(t *testing.T) {
		client, newErr := kuma.New(
			ctx,
			baseURL,
			testUsername,
			testPassword,
			kuma.WithTOTPSecret(normalized),
			kuma.WithConnectTimeout(10*time.Second),
		)

		require.NoError(t, newErr)
		require.NotNil(t, client)

		t.Cleanup(func() { _ = client.Disconnect() })

		sessionToken = client.SessionToken()
		require.NotEmpty(t, sessionToken, "a login answers with the token it can be repeated with")
	})

	awaitSignInWindow(scheme)

	t.Run("the session token needs neither password nor code", func(t *testing.T) {
		client, newErr := kuma.New(
			ctx,
			baseURL,
			"",
			"",
			kuma.WithSessionToken(sessionToken),
			kuma.WithConnectTimeout(10*time.Second),
		)

		require.NoError(t, newErr)
		require.NotNil(t, client)

		t.Cleanup(func() { _ = client.Disconnect() })

		require.Equal(t, sessionToken, client.SessionToken())

		// The lists come from a login, and the token is what repeats one.
		require.NoError(t, client.Resync(ctx))
	})

	t.Run("a code the server has seen is refused", func(t *testing.T) {
		if scheme == kuma.AuthSchemeHTTP {
			t.Skip("a 3.x server has no replay guard, see kuma.WithTOTPSecret")
		}

		// Line the two logins up inside one time step, which is what the
		// server's replay guard reacts to. Without the alignment the second
		// login could fall into the next step and succeed.
		waitForFreshTOTPStep(t)

		first, newErr := kuma.New(
			ctx,
			baseURL,
			testUsername,
			testPassword,
			kuma.WithTOTPSecret(normalized),
			kuma.WithConnectTimeout(10*time.Second),
		)
		require.NoError(t, newErr)

		t.Cleanup(func() { _ = first.Disconnect() })

		// A connect timeout shorter than what is left of the step keeps the
		// client from waiting the step out and retrying.
		second, newErr := kuma.New(
			ctx,
			baseURL,
			testUsername,
			testPassword,
			kuma.WithTOTPSecret(normalized),
			kuma.WithConnectTimeout(5*time.Second),
		)

		require.ErrorIs(t, newErr, kuma.ErrInvalidTOTPCode)
		require.Nil(t, second)
	})
}

// waitForFreshTOTPStep blocks until the current time step has ended, so that a
// pair of logins made right afterwards falls into one step the server has not
// seen a code from yet.
//
// Waiting for the boundary rather than for some headroom is what makes this
// reliable: an earlier login in this test has already used the current step's
// code, and the server refuses to see it again.
func waitForFreshTOTPStep(t *testing.T) {
	t.Helper()

	const step = 30 * time.Second

	time.Sleep(time.Until(time.Now().Truncate(step).Add(step)) + 100*time.Millisecond)
}
