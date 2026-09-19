package kuma

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// awaitIterations is how often each case below is repeated. A channel that is
// ready and a context that is done are both ready in the same select, and
// select picks between ready cases at random, so a single run would pass for
// the wrong reason about half the time.
const awaitIterations = 100

// awaitSignals describes which signals are already delivered when
// awaitAckAndUpdateEvent is entered.
type awaitSignals struct {
	// ack is placed in the buffered ack channel if it is not nil.
	ack *ackResponse
	// updateEvent places a token in the channel the update event handler
	// writes to.
	updateEvent bool
	// confirm is the confirmation the caller supplied, nil for a caller that
	// returns on the first update event.
	confirm updateEventConfirm
}

// awaitWith drives awaitAckAndUpdateEvent with the signals already in the state
// the case under test describes.
func awaitWith(ctx context.Context, signals awaitSignals) (ackResponse, error) {
	res := make(chan ackResponse, 1)
	if signals.ack != nil {
		res <- *signals.ack
	}

	updated := make(chan struct{}, 1)
	if signals.updateEvent {
		updated <- struct{}{}
	}

	return awaitAckAndUpdateEvent(ctx, "addNotification", updated, res, signals.confirm)
}

func expiredContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Nanosecond)
	t.Cleanup(cancel)

	<-ctx.Done()

	return ctx
}

func canceledContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	cancel()

	return ctx
}

// TestAwaitAckAndUpdateEvent pins down what awaitAckAndUpdateEvent reports for
// every combination of signals that can be pending when the context is done.
// The socket level tests cannot do this: they depend on the ack and the
// deadline landing in a particular order on the wire, so the collision that
// matters here is only reached by chance.
func TestAwaitAckAndUpdateEvent(t *testing.T) {
	t.Parallel()

	okAck := ackResponse{OK: true, Msg: "ok", ID: 4242}
	rejectedAck := ackResponse{OK: false, Msg: "boom"}

	t.Run("both_signals_arrive_before_the_context_is_done", func(t *testing.T) {
		t.Parallel()

		for range awaitIterations {
			response, err := awaitWith(t.Context(), awaitSignals{ack: &okAck, updateEvent: true})

			require.NoError(t, err)
			require.Equal(t, okAck, response)
		}
	})

	t.Run("both_signals_are_ready_when_the_deadline_expires", func(t *testing.T) {
		t.Parallel()

		// Uptime Kuma broadcasts the list before it invokes the ack callback,
		// so this ordering is the common one, not an exotic one. Nothing is
		// missing here, and reporting ErrUpdateEventTimeout would claim a
		// stale cache for a call that completed.
		ctx := expiredContext(t)

		for range awaitIterations {
			response, err := awaitWith(ctx, awaitSignals{ack: &okAck, updateEvent: true})

			require.NoError(t, err)
			require.Equal(t, okAck, response)
		}
	})

	t.Run("update_event_missing_after_the_deadline", func(t *testing.T) {
		t.Parallel()

		ctx := expiredContext(t)

		for range awaitIterations {
			response, err := awaitWith(ctx, awaitSignals{ack: &okAck})

			require.ErrorIs(t, err, ErrUpdateEventTimeout)
			require.ErrorIs(t, err, context.DeadlineExceeded,
				"the sentinel must not hide the deadline from callers matching on it")
			require.Equal(t, okAck, response,
				"the ack is what carries the ID of a created resource")
		}
	})

	t.Run("the_error_carries_the_command_and_keeps_its_message", func(t *testing.T) {
		t.Parallel()

		ctx := expiredContext(t)

		for range awaitIterations {
			_, err := awaitWith(ctx, awaitSignals{ack: &okAck})

			var timeoutErr *UpdateEventTimeoutError
			require.ErrorAs(t, err, &timeoutErr)
			require.Equal(t, "addNotification", timeoutErr.Command)
			require.Zero(t, timeoutErr.ID,
				"the ID is filled in by the layer that knows which field of the ack carries it")
			require.EqualError(t, err,
				"addNotification: update event not received: context deadline exceeded")
		}
	})

	t.Run("update_event_missing_after_cancellation", func(t *testing.T) {
		t.Parallel()

		ctx := canceledContext(t)

		for range awaitIterations {
			_, err := awaitWith(ctx, awaitSignals{ack: &okAck})

			require.ErrorIs(t, err, ErrUpdateEventTimeout)
			require.ErrorIs(t, err, context.Canceled)
		}
	})

	t.Run("rejected_ack_at_the_deadline_reports_the_server_error", func(t *testing.T) {
		t.Parallel()

		ctx := expiredContext(t)

		for range awaitIterations {
			response, err := awaitWith(ctx, awaitSignals{ack: &rejectedAck})

			require.EqualError(t, err, "addNotification: boom")
			require.NotErrorIs(t, err, ErrUpdateEventTimeout,
				"the server rejected the command, so nothing was applied")
			require.Equal(t, ackResponse{}, response)
		}
	})

	t.Run("update_event_without_an_ack_stays_a_plain_timeout", func(t *testing.T) {
		t.Parallel()

		// The broadcast alone says something changed, but not what the server
		// made of this command and not which ID it assigned.
		ctx := expiredContext(t)

		for range awaitIterations {
			_, err := awaitWith(ctx, awaitSignals{updateEvent: true})

			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.NotErrorIs(t, err, ErrUpdateEventTimeout)
		}
	})

	t.Run("no_signal_at_all_stays_a_plain_timeout", func(t *testing.T) {
		t.Parallel()

		ctx := expiredContext(t)

		for range awaitIterations {
			_, err := awaitWith(ctx, awaitSignals{})

			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.NotErrorIs(t, err, ErrUpdateEventTimeout)
		}
	})
}

// TestAwaitAckAndUpdateEventConfirmation pins down what a confirmation changes:
// the wait no longer ends on the first broadcast, it ends when the state cache
// holds the write. The broadcasts carry a whole list and are matched by name, so
// with several writes in flight on one connection the first one a waiter sees is
// routinely a neighbour's, taken from a list that predates its own change.
func TestAwaitAckAndUpdateEventConfirmation(t *testing.T) {
	t.Parallel()

	okAck := ackResponse{OK: true, Msg: "ok", ID: 4242}

	t.Run("a_stale_broadcast_does_not_end_the_wait", func(t *testing.T) {
		t.Parallel()

		ctx := expiredContext(t)

		for range awaitIterations {
			// The update event is in hand, but the cache does not hold the
			// write: that is exactly the broadcast a neighbouring write
			// caused. Without the confirmation this returns success and leaves
			// the caller with a cache that is one broadcast behind.
			_, err := awaitWith(ctx, awaitSignals{
				ack:         &okAck,
				updateEvent: true,
				confirm:     func(ackResponse) bool { return false },
			})

			require.ErrorIs(t, err, ErrUpdateEventTimeout)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		}
	})

	t.Run("the_confirmed_write_succeeds_without_any_broadcast", func(t *testing.T) {
		t.Parallel()

		// A broadcast that arrived while the command was still in flight has
		// already refreshed the cache, so there is nothing left to wait for.
		ctx := expiredContext(t)

		for range awaitIterations {
			response, err := awaitWith(ctx, awaitSignals{
				ack:     &okAck,
				confirm: func(ackResponse) bool { return true },
			})

			require.NoError(t, err)
			require.Equal(t, okAck, response)
		}
	})

	t.Run("the_confirmation_sees_the_ack", func(t *testing.T) {
		t.Parallel()

		// The ID a create is confirmed by is the one the ack carries, so the
		// confirmation has to run after the ack and be handed it.
		for range awaitIterations {
			var seen []int64

			response, err := awaitWith(t.Context(), awaitSignals{
				ack:     &okAck,
				confirm: func(response ackResponse) bool { seen = append(seen, response.ID); return true },
			})

			require.NoError(t, err)
			require.Equal(t, okAck, response)
			require.Equal(t, []int64{okAck.ID}, seen)
		}
	})

	t.Run("a_later_broadcast_re_runs_the_confirmation", func(t *testing.T) {
		t.Parallel()

		// The first broadcast is the neighbour's and the second is this
		// write's own. Only re-reading the cache after every one of them gets
		// the waiter past the first.
		for range awaitIterations {
			updated := make(chan struct{}, 1)
			updated <- struct{}{}

			res := make(chan ackResponse, 1)
			res <- okAck

			calls := 0
			confirm := func(ackResponse) bool {
				calls++
				if calls == 1 {
					// Stand in for this write's own broadcast landing while
					// the neighbour's is still being looked at. Sent the way
					// the listener sends it, without blocking: the initial
					// token is still pending whenever the ack was taken
					// first.
					select {
					case updated <- struct{}{}:
					default:
					}

					return false
				}

				return true
			}

			response, err := awaitAckAndUpdateEvent(t.Context(), "addNotification", updated, res, confirm)

			require.NoError(t, err)
			require.Equal(t, okAck, response)
			require.Equal(t, 2, calls)
		}
	})
}
