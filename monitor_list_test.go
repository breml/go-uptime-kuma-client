package kuma_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/monitor"
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

// monitorAddedFrames are the frames answering the add command with ackID: the
// updateMonitorIntoList broadcast merging the new monitor, then the ack with its
// ID, in the order Uptime Kuma sends them.
func monitorAddedFrames(ackID string, id int64) [][]byte {
	return [][]byte{
		fmt.Appendf(nil, `42["updateMonitorIntoList",{"%d":{"id":%d,"name":"created","type":"http"}}]`, id, id),
		fmt.Appendf(nil, `43%s[{"ok":true,"monitorID":%d}]`, ackID, id),
	}
}

// TestGetMonitorsWaitsForTheListItRequested pins that GetMonitors returns the
// list its own request brought in, not whatever the cache held when the ack
// arrived. The fake acks first and broadcasts later, an order the client can
// observe even from a server that broadcasts first, because it runs acks on a
// goroutine of their own. Returning on the ack would report the empty list from
// the login.
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
// rather than the cache it failed to refresh. The failed request must release
// the lock, so the next one goes through.
func TestGetMonitorsWithoutListIsAnError(t *testing.T) {
	t.Parallel()

	const budget = 200 * time.Millisecond

	var sendList atomic.Bool

	kumaClient := newMonitorListClient(t, func(fake *fakeSocketIOServer, command string, ackID string) {
		if command != "getMonitorList" {
			return
		}

		if sendList.Load() {
			fake.messages <- monitorListFrame(1, "fresh")
		}

		fake.messages <- okAck(ackID)
	}, kuma.WithOperationTimeout(budget))

	start := time.Now()
	monitors, err := kumaClient.GetMonitors(t.Context())
	elapsed := time.Since(start)

	require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
	require.ErrorIs(t, err, kuma.ErrOperationTimeout)
	require.Nil(t, monitors, "a list that was not refreshed must not be returned")
	require.Less(t, elapsed, 10*budget, "the wait must be bounded by the operation budget")

	sendList.Store(true)

	monitors, err = kumaClient.GetMonitors(t.Context())
	require.NoError(t, err, "a failed request must not keep the lock")
	require.Len(t, monitors, 1)
}

// TestGetMonitorsAreSerialized pins that a client never has two getMonitorList
// requests in flight, and that each caller gets the list its own request
// brought in.
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

	type result struct {
		monitors []monitor.Base
		err      error
	}

	results := make(chan result, callers)

	for range callers {
		go func() {
			monitors, err := kumaClient.GetMonitors(t.Context())
			results <- result{monitors: monitors, err: err}
		}()
	}

	names := make([]string, 0, callers)

	for range callers {
		res := <-results
		require.NoError(t, res.err)
		require.Len(t, res.monitors, 1)

		names = append(names, res.monitors[0].Name)
	}

	require.Equal(t, int64(callers), served.Load())
	require.Equal(t, int64(1), maxInFlight.Load(),
		"the server must never see two getMonitorList requests from one client at once")
	require.ElementsMatch(t, []string{"snapshot 1", "snapshot 2", "snapshot 3", "snapshot 4"}, names,
		"every caller must get the list its own request brought in")
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

// heldCommands records the commands a fake server receives and holds their
// answers back until the test sends them.
type heldCommands struct {
	mu       sync.Mutex
	received []string
	pending  chan heldCommand
}

type heldCommand struct {
	fake    *fakeSocketIOServer
	command string
	ackID   string
}

func newHeldCommands() *heldCommands {
	return &heldCommands{pending: make(chan heldCommand, 16)}
}

func (h *heldCommands) onCommand(fake *fakeSocketIOServer, command string, ackID string) {
	h.mu.Lock()
	h.received = append(h.received, command)
	h.mu.Unlock()

	h.pending <- heldCommand{fake: fake, command: command, ackID: ackID}
}

func (h *heldCommands) commands() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return slices.Clone(h.received)
}

// next returns the oldest command not yet taken, which must be command.
func (h *heldCommands) next(t *testing.T, command string) heldCommand {
	t.Helper()

	select {
	case held := <-h.pending:
		require.Equal(t, command, held.command)

		return held

	case <-time.After(5 * time.Second):
		require.FailNow(t, "command never reached the server", command)

		return heldCommand{}
	}
}

func (c heldCommand) answer(frames ...[]byte) {
	for _, frame := range frames {
		c.fake.messages <- frame
	}
}

// heldBackWindow is how long a test watches for a command the client must hold
// back, before concluding it did.
const heldBackWindow = 100 * time.Millisecond

// TestGetMonitorsHoldsBackMonitorWrites pins that a monitor write waits for a
// GetMonitors in flight. Sent alongside, the server could insert the monitor,
// merge it into the cache and confirm the create before the list GetMonitors
// asked for arrives, and that list, read before the insert, would then drop the
// new monitor from the cache.
func TestGetMonitorsHoldsBackMonitorWrites(t *testing.T) {
	t.Parallel()

	held := newHeldCommands()
	kumaClient := newMonitorListClient(t, held.onCommand)

	listDone := make(chan error, 1)

	go func() {
		_, err := kumaClient.GetMonitors(t.Context())
		listDone <- err
	}()

	list := held.next(t, "getMonitorList")

	createDone := make(chan error, 1)

	go func() {
		_, err := kumaClient.CreateMonitor(t.Context(), operationTimeoutMonitor())
		createDone <- err
	}()

	time.Sleep(heldBackWindow)
	require.Equal(t, []string{"getMonitorList"}, held.commands(),
		"the create must not reach the server while GetMonitors is in flight")

	list.answer(monitorListFrame(1, "listed"), okAck(list.ackID))
	require.NoError(t, <-listDone)

	add := held.next(t, "add")
	add.answer(monitorAddedFrames(add.ackID, 2)...)
	require.NoError(t, <-createDone)
}

// TestMonitorWritesHoldBackGetMonitors pins the other order: a GetMonitors
// waits for the monitor writes in flight, so the list it asks for is read after
// them.
func TestMonitorWritesHoldBackGetMonitors(t *testing.T) {
	t.Parallel()

	held := newHeldCommands()
	kumaClient := newMonitorListClient(t, held.onCommand)

	createDone := make(chan error, 1)

	go func() {
		_, err := kumaClient.CreateMonitor(t.Context(), operationTimeoutMonitor())
		createDone <- err
	}()

	add := held.next(t, "add")

	listDone := make(chan error, 1)

	go func() {
		_, err := kumaClient.GetMonitors(t.Context())
		listDone <- err
	}()

	time.Sleep(heldBackWindow)
	require.Equal(t, []string{"add"}, held.commands(),
		"GetMonitors must not reach the server while a create is in flight")

	add.answer(monitorAddedFrames(add.ackID, 2)...)
	require.NoError(t, <-createDone)

	list := held.next(t, "getMonitorList")
	list.answer(monitorListFrame(2, "created"), okAck(list.ackID))
	require.NoError(t, <-listDone)
}

// TestMonitorWritesStayConcurrent pins that the lock GetMonitors takes does not
// serialize the monitor writes among themselves: both creates reach the server
// before either is answered.
func TestMonitorWritesStayConcurrent(t *testing.T) {
	t.Parallel()

	held := newHeldCommands()
	kumaClient := newMonitorListClient(t, held.onCommand)

	createDone := make(chan error, 2)

	for range 2 {
		go func() {
			_, err := kumaClient.CreateMonitor(t.Context(), operationTimeoutMonitor())
			createDone <- err
		}()
	}

	first := held.next(t, "add")
	second := held.next(t, "add")

	first.answer(monitorAddedFrames(first.ackID, 2)...)
	second.answer(monitorAddedFrames(second.ackID, 3)...)

	require.NoError(t, <-createDone)
	require.NoError(t, <-createDone)
}

// TestMonitorWriteQueueIsCancellable pins that a monitor write waiting behind a
// GetMonitors gives up when its caller says so, without being sent.
func TestMonitorWriteQueueIsCancellable(t *testing.T) {
	t.Parallel()

	const callerTimeout = 100 * time.Millisecond

	held := newHeldCommands()
	kumaClient := newMonitorListClient(t, held.onCommand, kuma.WithOperationTimeout(10*time.Second))

	go func() {
		_, _ = kumaClient.GetMonitors(t.Context())
	}()

	held.next(t, "getMonitorList")

	ctx, cancel := context.WithTimeout(t.Context(), callerTimeout)
	defer cancel()

	_, err := kumaClient.CreateMonitor(ctx, operationTimeoutMonitor())
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.NotErrorIs(t, err, kuma.ErrOperationTimeout,
		"the queue is bounded by the caller, not by the client's budget")
	require.Equal(t, []string{"getMonitorList"}, held.commands(), "the queued create must never have been sent")
}
