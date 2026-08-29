package kuma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// readyEvents are the update events kuma.New() waits for before it returns. The
// payloads only have to unmarshal into the registered handlers' argument types.
//
//nolint:gochecknoglobals // fixed table, shared by the fake server below.
var readyEvents = []string{
	`42["monitorList",{}]`,
	`42["maintenanceList",{}]`,
	`42["notificationList",[]]`,
	`42["statusPageList",{}]`,
	`42["proxyList",[]]`,
	`42["dockerHostList",[]]`,
	`42["apiKeyList",[]]`,
}

// fakeAckCreatedID is the ID the fake server reports for a created
// notification.
const fakeAckCreatedID = 4242

// fakeAckServer is a minimal socket.io server over HTTP long-polling that
// completes the handshake, CONNECT, login and ready phases, so kuma.New()
// returns a usable client. Unlike fakeSocketIOServer it then keeps serving:
// every event the client emits afterwards is acknowledged after ackDelay,
// which lets a test place the ack relative to the caller's deadline.
//
// It never emits an update event for those commands, so a caller using
// syncEmitWithUpdateEvent keeps waiting after its ack arrived — the state in
// which the deadline and the ack collide.
type fakeAckServer struct {
	messages chan []byte

	mu       sync.Mutex
	ackDelay time.Duration
}

func (s *fakeAckServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	switch {
	case r.Method == http.MethodGet && sid == "":
		type engineIOHandshake struct {
			Sid          string   `json:"sid"`
			Upgrades     []string `json:"upgrades"`
			PingInterval int      `json:"pingInterval"`
			PingTimeout  int      `json:"pingTimeout"`
			MaxPayload   int      `json:"maxPayload"`
		}
		data, err := json.Marshal(engineIOHandshake{
			Sid:          "ack-test-sid",
			Upgrades:     []string{}, // no WebSocket upgrade keeps the client on polling
			PingInterval: 50,
			PingTimeout:  5000,
			MaxPayload:   1000000,
		})
		if err != nil {
			http.Error(w, "marshal handshake", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(append([]byte("0"), data...))

	case r.Method == http.MethodPost && sid != "":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusInternalServerError)
			return
		}

		s.handleClientMessage(body)
		w.WriteHeader(http.StatusOK)

	case r.Method == http.MethodGet && sid != "":
		select {
		case msg := <-s.messages:
			_, _ = w.Write(msg)

		case <-r.Context().Done():
		}

	default:
	}
}

func (s *fakeAckServer) setAckDelay(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ackDelay = d
}

func (s *fakeAckServer) currentAckDelay() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ackDelay
}

func (s *fakeAckServer) handleClientMessage(body []byte) {
	if len(body) < 2 {
		return
	}

	socketIOData := body[1:] // strip engine.io "4" (Message) prefix

	switch socketIOData[0] {
	case '0': // socket.io CONNECT
		s.messages <- []byte("40")

	case '2': // socket.io EVENT
		i := 1
		for i < len(socketIOData) && socketIOData[i] >= '0' && socketIOData[i] <= '9' {
			i++
		}

		ackID := string(socketIOData[1:i])
		payload := string(socketIOData[i:])

		if strings.Contains(payload, `"login"`) {
			s.messages <- fmt.Appendf(nil, `43%s[{"ok":true,"msg":"Logged in successfully."}]`, ackID)

			for _, event := range readyEvents {
				s.messages <- []byte(event)
			}

			return
		}

		// The server assigns the ID of a created resource in the ack, which
		// is the value a caller loses if a successful ack is discarded.
		ack := `{"ok":true,"msg":"ok"}`
		if strings.Contains(payload, `"addNotification"`) {
			ack = fmt.Sprintf(`{"ok":true,"msg":"ok","id":%d}`, fakeAckCreatedID)
		}

		// Acknowledge out of band: blocking here would stall the client's
		// send path instead of only delaying the ack.
		delay := s.currentAckDelay()
		go func() {
			time.Sleep(delay)
			s.messages <- fmt.Appendf(nil, `43%s[%s]`, ackID, ack)
		}()

	default:
	}
}

// TestAckDeliveryAroundDeadline sweeps the ack of an in-flight command across
// the caller's deadline, so that early, colliding and late acks are all
// exercised against the same call.
//
// The defect this guards against is a "send on closed channel" panic: the ack
// callback runs on a goroutine owned by the socket.io client, so an ack that
// arrives once the caller has returned finds the channel closed by its
// `defer close(res)`, and the panic is raised where no caller can recover from
// it. That kills the process, not just the request.
//
// What it catches, and what it does not: the late-ack iterations at the end of
// the sweep panic every time if `close(res)` is reintroduced, so the invariant
// "the ack channel is buffered and never closed" is genuinely enforced here.
// The original code guarded the send with `if ctx.Err() != nil` and crashed
// only when the context expired between that check and the send. That window
// is a few instructions wide, and it was verified that this test does not
// reproduce it — real transport latency swamps it. So a green run is not
// evidence that a guard-based implementation is safe; it is evidence that the
// channel is still handled the way syncEmit and syncEmitWithUpdateEvent
// document.
func TestAckDeliveryAroundDeadline(t *testing.T) {
	const (
		callTimeout = 30 * time.Millisecond
		iterations  = 24
		step        = 500 * time.Microsecond
	)

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
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
		kuma.WithConnectTimeout(10*time.Second),
	)
	require.NoError(t, err)

	// Start clearly before the deadline and end clearly after it.
	base := callTimeout - (iterations/2)*step

	t.Run("sync_emit_with_update_event", func(t *testing.T) {
		for i := range iterations {
			fake.setAckDelay(base + time.Duration(i)*step)

			callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
			// The fake acks the command but never emits notificationList, so
			// the call keeps waiting for its update event and always runs into
			// the deadline — with the ack landing somewhere around it.
			err := kumaClient.DeleteNotification(callCtx, 1)
			callCancel()

			// Both outcomes of the race keep wrapping the context error: the
			// ack that lost it produces a bare deadline error, the ack that
			// won it produces ErrUpdateEventTimeout wrapping the same
			// deadline error. Which one a given delay yields depends on
			// scheduling, so only the shared part is asserted here;
			// TestAwaitAckAndUpdateEvent pins the ack-wins case down
			// deterministically.
			require.ErrorIs(t, err, context.DeadlineExceeded,
				"ack delay %s should leave the call waiting on its update event", fake.currentAckDelay())
		}
	})

	t.Run("sync_emit", func(t *testing.T) {
		for i := range iterations {
			fake.setAckDelay(base + time.Duration(i)*step)

			callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
			msg, err := kumaClient.TestNotification(callCtx, notification.Generic{
				Base:           notification.Base{Name: "ack race"},
				GenericDetails: notification.GenericDetails{},
				TypeName:       "generic",
			})
			callCancel()

			// Whether the ack or the deadline won is a matter of scheduling and
			// either outcome is correct. A panic, a hang, or a garbled response
			// is not.
			if err != nil {
				require.ErrorIs(t, err, context.DeadlineExceeded)

				continue
			}

			require.Equal(t, "ok", msg)
		}
	})
}

// TestUpdateEventTimeoutKeepsSuccessfulAck covers the case the sweep above can
// only hit by chance: the server acknowledges the command well before the
// deadline and never emits the update event that follows it.
//
// The command was applied, so reporting a bare context error would tell the
// caller the write failed for a write the server performed — and for a create,
// it would also drop the ID the server assigned, leaving a retry to produce a
// duplicate.
func TestUpdateEventTimeoutKeepsSuccessfulAck(t *testing.T) {
	// Long enough that the ack, which the fake sends without delay, reliably
	// arrives first; the call then runs into the deadline waiting for the
	// update event the fake never emits. It travels over HTTP long-polling, so
	// this is a generous margin rather than a guarantee: raise callTimeout if
	// a loaded machine ever makes it flake.
	const callTimeout = 500 * time.Millisecond

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
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
		kuma.WithConnectTimeout(10*time.Second),
	)
	require.NoError(t, err)

	t.Run("create_reports_the_id_the_server_assigned", func(t *testing.T) {
		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		id, err := kumaClient.CreateNotification(callCtx, notification.Generic{
			Base:           notification.Base{Name: "update event timeout"},
			GenericDetails: notification.GenericDetails{},
			TypeName:       "generic",
		})

		require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
		require.ErrorIs(t, err, context.DeadlineExceeded,
			"the sentinel must not hide the timeout from callers matching on it")
		require.Equal(t, int64(fakeAckCreatedID), id,
			"the notification exists on the server, so its ID must survive the missing update event")
	})

	t.Run("delete_is_distinguishable_from_a_plain_timeout", func(t *testing.T) {
		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		err := kumaClient.DeleteNotification(callCtx, 1)

		require.ErrorIs(t, err, kuma.ErrUpdateEventTimeout)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("an_unacknowledged_command_stays_a_plain_timeout", func(t *testing.T) {
		// No ack at all before the deadline: nothing is known about the write,
		// so the sentinel must not be reported.
		fake.setAckDelay(2 * callTimeout)
		t.Cleanup(func() { fake.setAckDelay(0) })

		callCtx, callCancel := context.WithTimeout(ctx, callTimeout)
		defer callCancel()

		err := kumaClient.DeleteNotification(callCtx, 1)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NotErrorIs(t, err, kuma.ErrUpdateEventTimeout,
			"without an ack the client cannot claim the command was applied")
	})
}
