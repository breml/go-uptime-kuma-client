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
	// updateEvent closes the channel the update event handler closes.
	updateEvent bool
}

// awaitWith drives awaitAckAndUpdateEvent with the signals already in the state
// the case under test describes.
func awaitWith(ctx context.Context, signals awaitSignals) (ackResponse, error) {
	res := make(chan ackResponse, 1)
	if signals.ack != nil {
		res <- *signals.ack
	}

	done := make(chan struct{})
	if signals.updateEvent {
		close(done)
	}

	return awaitAckAndUpdateEvent(ctx, "addNotification", done, res)
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
