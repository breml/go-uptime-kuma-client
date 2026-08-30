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

// newFakeAckClient connects a client to fake and returns it, with the server and
// the connection torn down when the test ends.
func newFakeAckClient(t *testing.T, fake *fakeAckServer, opts ...kuma.Option) *kuma.Client {
	t.Helper()

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
		append([]kuma.Option{kuma.WithConnectTimeout(10 * time.Second)}, opts...)...,
	)
	require.NoError(t, err)

	return kumaClient
}

// TestResync covers the way out of the stale cache an ErrUpdateEventTimeout
// leaves behind. The server has no command to re-request a single list, so
// Resync logs in again with the token from the initial login and waits for the
// lists that answers with.
func TestResync(t *testing.T) {
	t.Run("logs_in_again_with_the_token_from_the_first_login", func(t *testing.T) {
		fake := &fakeAckServer{messages: make(chan []byte, 32)}
		kumaClient := newFakeAckClient(t, fake)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		require.NoError(t, kumaClient.Resync(ctx))
		require.Contains(t, fake.lastResyncPayload(), fakeSessionToken,
			"the resync has to present the session the server handed out")
	})

	t.Run("waits_for_the_lists_and_names_the_ones_that_did_not_arrive", func(t *testing.T) {
		fake := &fakeAckServer{messages: make(chan []byte, 32)}
		kumaClient := newFakeAckClient(t, fake)

		// A server that acknowledges the login without resending the lists
		// leaves the cache exactly as stale as it was.
		fake.setSuppressLists(true)

		ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
		defer cancel()

		err := kumaClient.Resync(ctx)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.ErrorContains(t, err, "monitorList")
		require.ErrorContains(t, err, "notificationList")
		require.ErrorContains(t, err, "statusPageList")
		require.NotContains(t, err.Error(), "proxyList",
			"a best-effort list is not what a resync waits for")
	})

	t.Run("the_best_effort_lists_do_not_hold_it_up", func(t *testing.T) {
		// Long enough that a resync waiting it out is unmistakable in the
		// elapsed time below.
		const grace = time.Second

		fake := &fakeAckServer{messages: make(chan []byte, 32)}

		// A server that only ever emits the required lists, as older versions
		// and some reverse proxies do, has to resync just the same — and once
		// New has learned that it never sends them, without spending the grace
		// period on them again.
		fake.setSuppressOptionalLists(true)

		kumaClient := newFakeAckClient(t, fake, kuma.WithReadyGracePeriod(grace))

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		start := time.Now()

		require.NoError(t, kumaClient.Resync(ctx))
		require.Less(t, time.Since(start), grace/4,
			"the resync waited out the grace period for lists this server never sends")
	})

	t.Run("a_deadline_in_the_grace_window_does_not_fail_a_finished_resync", func(t *testing.T) {
		fake := &fakeAckServer{messages: make(chan []byte, 32)}
		kumaClient := newFakeAckClient(t, fake)

		fake.setSuppressOptionalLists(true)

		// Shorter than the default ready grace period, so the deadline is what
		// ends the wait for the best-effort lists. By then the required ones
		// have rebuilt the cache, which is all the caller asked for.
		ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
		defer cancel()

		require.NoError(t, kumaClient.Resync(ctx))
	})

	t.Run("without_a_session_token_it_reports_that_instead_of_hanging", func(t *testing.T) {
		// A server too old to state how it wants to be authenticated, given a
		// client with nothing to offer: New logs in nowhere, so there is no
		// token to resync with. An accepted login that carried no token cannot
		// produce this state, see TestNewLoginWithoutSessionToken.
		fake := &fakeAckServer{messages: make(chan []byte, 32)}
		fake.setOmitLoginRequired()

		server := httptest.NewServer(fake)

		newCtx, newCancel := context.WithTimeout(t.Context(), time.Minute)
		t.Cleanup(func() {
			newCancel()
			server.CloseClientConnections()
			server.Close()
		})

		kumaClient, err := kuma.New(
			newCtx, server.URL, "", "",
			kuma.WithConnectTimeout(10*time.Second),
			// Such a server sends no lists to a client it never logged in, so
			// the ready gate is what would be waited out, not the resync.
			kuma.WithReadyEvents(),
			kuma.WithReadyGracePeriod(0),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		err = kumaClient.Resync(ctx)

		require.ErrorContains(t, err, "no session token")
		require.Empty(t, fake.lastResyncPayload(),
			"without a token there is nothing to log in with, so nothing is emitted")
	})
}

// TestResyncAgainstServer runs the resync against a real Uptime Kuma instance,
// which is what proves the token is the one loginByToken accepts.
func TestResyncAgainstServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	id, err := client.CreateNotification(ctx, notification.Generic{
		Base:           notification.Base{Name: "resync"},
		GenericDetails: notification.GenericDetails{},
		TypeName:       "generic",
	})
	require.NoError(t, err)

	// Deferred rather than registered with t.Cleanup, which runs once
	// t.Context() is already cancelled.
	defer func() {
		require.NoError(t, client.DeleteNotification(ctx, id))
	}()

	require.NoError(t, client.Resync(ctx))

	notif, err := client.GetNotification(ctx, id)
	require.NoError(t, err, "the resync has to leave the cache complete, not empty")
	require.Equal(t, id, notif.GetID())
}
