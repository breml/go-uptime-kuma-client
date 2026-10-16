package kuma_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// newMonitorListClient connects to a fake server that completes the login with
// an empty monitor list and then hands every command to onCommand, and returns
// the connected client.
func newMonitorListClient(
	t *testing.T,
	onCommand func(fake *fakeSocketIOServer, command string, ackID string),
	opts ...kuma.Option,
) *kuma.Client {
	t.Helper()

	fake := &fakeSocketIOServer{
		messages: make(chan []byte, 64),
		eventsAfterLogin: []string{
			`42["monitorList",{}]`,
			`42["notificationList",[]]`,
			`42["statusPageList",{}]`,
		},
	}
	fake.onCommand = func(command string, ackID string) {
		onCommand(fake, command, ackID)
	}

	server := httptest.NewServer(fake)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(func() {
		cancel()
		server.CloseClientConnections()
		server.Close()
	})

	kumaClient, err := kuma.New(
		ctx, server.URL, "user", "pass",
		append([]kuma.Option{
			kuma.WithConnectTimeout(5 * time.Second),
			kuma.WithReadyGracePeriod(0),
		}, opts...)...,
	)
	require.NoError(t, err, "connecting to the fake server must succeed")

	// Registered after the connect, so the LIFO order disconnects the client
	// before the cleanup above closes the server under it.
	t.Cleanup(func() { _ = kumaClient.Disconnect() })

	return kumaClient
}

// okAck is the frame acknowledging the command with ackID as a success.
func okAck(ackID string) []byte {
	return fmt.Appendf(nil, `43%s[{"ok":true}]`, ackID)
}

// monitorListFrame is a monitorList broadcast holding one HTTP monitor.
func monitorListFrame(id int64, name string) []byte {
	return fmt.Appendf(nil, `42["monitorList",{"%d":{"id":%d,"name":%q,"type":"http"}}]`, id, id, name)
}

// TestGetMonitorsWaitsForTheListItRequested pins that GetMonitors returns the
// list its own request brought in, not whatever the cache held when the ack
// arrived. The fake acks first and broadcasts later - the client runs the ack
// on a goroutine of its own, so that order is one it can observe even against
// a server that broadcasts first - and a GetMonitors that returned on the ack
// would report the empty list from the login.
func TestGetMonitorsWaitsForTheListItRequested(t *testing.T) {
	t.Parallel()

	kumaClient := newMonitorListClient(t, func(fake *fakeSocketIOServer, command string, ackID string) {
		if command != "getMonitorList" {
			return
		}

		fake.messages <- okAck(ackID)

		go func() {
			time.Sleep(100 * time.Millisecond)

			fake.messages <- monitorListFrame(1, "fresh")
		}()
	})

	monitors, err := kumaClient.GetMonitors(t.Context())
	require.NoError(t, err)
	require.Len(t, monitors, 1, "the list the request brought in must be the one returned")
	require.Equal(t, "fresh", monitors[0].Name)
}

// TestGetMonitorsWithoutListIsAnError pins that an acknowledged request whose
// list never arrives fails with ErrUpdateEventTimeout, and returns no monitors
// rather than the cache it failed to refresh.
func TestGetMonitorsWithoutListIsAnError(t *testing.T) {
	t.Parallel()

	const budget = 200 * time.Millisecond

	kumaClient := newMonitorListClient(t, func(fake *fakeSocketIOServer, command string, ackID string) {
		if command == "getMonitorList" {
			fake.messages <- okAck(ackID)
		}
	}, kuma.WithOperationTimeout(budget))

	start := time.Now()
	monitors, err := kumaClient.GetMonitors(t.Context())
	elapsed := time.Since(start)

	require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
	require.ErrorIs(t, err, kuma.ErrOperationTimeout)
	require.Nil(t, monitors, "a list that was not refreshed must not be returned")
	require.Less(t, elapsed, 10*budget, "the wait must be bounded by the operation budget")
}

// TestGetMonitorsAreSerialized pins that a client never has two getMonitorList
// requests in flight. Uptime Kuma answers each with a whole-list broadcast in
// completion order, so a snapshot read before a monitor was created could
// otherwise land after one read later and drop that monitor from the cache.
func TestGetMonitorsAreSerialized(t *testing.T) {
	t.Parallel()

	const callers = 4

	var (
		inFlight    atomic.Int64
		maxInFlight atomic.Int64
		served      atomic.Int64
	)

	kumaClient := newMonitorListClient(t, func(fake *fakeSocketIOServer, command string, ackID string) {
		if command != "getMonitorList" {
			return
		}

		current := inFlight.Add(1)
		for {
			highest := maxInFlight.Load()
			if current <= highest || maxInFlight.CompareAndSwap(highest, current) {
				break
			}
		}

		go func() {
			// Long enough for the other callers to send theirs, were they not
			// held back.
			time.Sleep(50 * time.Millisecond)

			id := served.Add(1)

			// Released before the frames are queued: the next request can only
			// be sent once the client has read them.
			inFlight.Add(-1)

			fake.messages <- monitorListFrame(id, fmt.Sprintf("snapshot %d", id))
			fake.messages <- okAck(ackID)
		}()
	})

	var wg sync.WaitGroup
	for range callers {
		wg.Go(func() {
			_, err := kumaClient.GetMonitors(t.Context())
			require.NoError(t, err)
		})
	}

	wg.Wait()

	require.Equal(t, int64(callers), served.Load())
	require.Equal(t, int64(1), maxInFlight.Load(),
		"the server must never see two getMonitorList requests from one client at once")
}

// TestGetMonitorsQueueIsCancellable pins that a GetMonitors waiting behind
// another one gives up when its caller says so, rather than when the request
// ahead of it finishes.
func TestGetMonitorsQueueIsCancellable(t *testing.T) {
	t.Parallel()

	const (
		heldBudget    = 10 * time.Second
		callerTimeout = 100 * time.Millisecond
	)

	var received atomic.Int64

	// Leaves every request unanswered, so the first one holds the lock until
	// its budget expires.
	kumaClient := newMonitorListClient(t, func(_ *fakeSocketIOServer, command string, _ string) {
		if command == "getMonitorList" {
			received.Add(1)
		}
	}, kuma.WithOperationTimeout(heldBudget))

	go func() {
		_, _ = kumaClient.GetMonitors(t.Context())
	}()

	require.Eventually(t, func() bool {
		return received.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond, "the first request must reach the wire and hold the lock")

	ctx, cancel := context.WithTimeout(t.Context(), callerTimeout)
	defer cancel()

	start := time.Now()
	_, err := kumaClient.GetMonitors(ctx)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.NotErrorIs(t, err, kuma.ErrOperationTimeout,
		"the queue is bounded by the caller, not by the client's budget")
	require.Less(t, elapsed, heldBudget/2)
	require.Equal(t, int64(1), received.Load(), "the queued request must never have been sent")
}
