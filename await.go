package kuma

import (
	"context"
	"fmt"
	"time"
)

// sendErrorLostGrace is how long emitOn waits, after a send failed, for the
// connection to be reported lost. The close handler runs once the transport's
// reader notices the drop, which against a restarted server is usually just
// after the send that hit it failed. It is only ever paid on a failed send.
const sendErrorLostGrace = 500 * time.Millisecond

// emitAckOn emits command on conn and returns the ack the server answered
// with, see Client.emitAck, which is the caller for the connection in service.
// It takes no lock, so a caller emitting on the connection in service holds the
// shared side of the connection lock, see lockConn.
func emitAckOn(ctx context.Context, conn *socketConn, command string, args ...any) (ackResponse, error) {
	// Buffered and never closed, so a late ack (e.g. after the context
	// expired) neither leaks the ack goroutine forever nor panics sending on
	// a closed channel.
	res := make(chan ackResponse, 1)

	args = append(args, ackOptions(ctx, res)...)

	err := emitOn(ctx, conn, command, args...)
	if err != nil {
		return ackResponse{}, fmt.Errorf("%s: %w", command, err)
	}

	select {
	case response := <-res:
		return response, nil

	case <-conn.lost:
		// An ack that came in just before the drop is still the answer.
		select {
		case response := <-res:
			return response, nil

		default:
			return ackResponse{}, fmt.Errorf("%s: %w", command, ErrConnectionLost)
		}

	case <-ctx.Done():
		return ackResponse{}, fmt.Errorf("%s: %w", command, contextErr(ctx))
	}
}

// emitOn sends command on conn and returns once it is on its way, conn is
// lost, or ctx is done, whichever comes first.
//
// The socket.io client sends synchronously, and over long polling a send is an
// HTTP request bounded by nothing but the connection's context, so against a
// server that stopped answering it does not return until the heartbeat gives
// the connection up. Waiting for it here keeps the caller's context and the
// operation timeout in charge. A send that is abandoned ends once the
// connection is lost or closed, which cancels that context.
func emitOn(ctx context.Context, conn *socketConn, command string, args ...any) error {
	if conn.isLost() {
		return ErrConnectionLost
	}

	sent := make(chan error, 1)

	go func() {
		sent <- conn.socket.Emit(command, args...)
	}()

	select {
	case err := <-sent:
		if err != nil && awaitLost(ctx, conn) {
			return fmt.Errorf("%w: %w", ErrConnectionLost, err)
		}

		return err

	case <-conn.lost:
		return ErrConnectionLost

	case <-ctx.Done():
		return contextErr(ctx)
	}
}

// awaitLost reports whether conn is lost, or is reported lost within
// sendErrorLostGrace, see emitOn.
func awaitLost(ctx context.Context, conn *socketConn) bool {
	timer := time.NewTimer(sendErrorLostGrace)
	defer timer.Stop()

	select {
	case <-conn.lost:
		return true

	case <-timer.C:
		return false

	case <-ctx.Done():
		return conn.isLost()
	}
}

// awaitAckAndUpdateEvent waits for the ack delivered on res and for the update
// event that confirm accepts, or, when confirm is nil, for the first update
// event signalled on updated.
//
// An ack that reports a failure is returned as the server's error even when ctx
// is done or lost is closed as well, because the rejection is the more useful
// answer and it is what the caller would have received had the ack arrived a
// moment earlier. An update event without an ack stays a plain context or
// connection-lost error: without the ack the
// client knows neither the outcome of the command nor the ID the server
// assigned, and a confirmation cannot run without the ID either.
func awaitAckAndUpdateEvent(
	ctx context.Context,
	command string,
	updated <-chan struct{},
	res <-chan ackResponse,
	lost <-chan struct{},
	confirm updateEventConfirm,
) (ackResponse, error) {
	var state awaitState

	for {
		if state.acked && confirmUpdate(confirm, state.response, state.sawUpdate) {
			return state.response, nil
		}

		select {
		case <-updated:
			state.sawUpdate = true

		case state.response = <-res:
			if !state.response.OK {
				return ackResponse{}, fmt.Errorf("%s: %s", command, state.response.Msg)
			}

			state.acked = true
			// Nil blocks forever, so the loop keeps waiting on the remaining
			// signals instead of spinning on a channel that is done.
			res = nil

		case <-lost:
			// The connection is gone, and with it any signal still to come.
			// What was delivered before the drop still counts.
			state.collectDelivered(res, updated)

			return state.result(command, confirm, ErrConnectionLost)

		case <-ctx.Done():
			// A signal that has already been delivered and the context being
			// done can become ready in the same select, which picks between
			// ready cases at random. Both signals are therefore collected
			// explicitly, instead of letting that coin flip decide whether
			// what they report is seen.
			state.collectDelivered(res, updated)

			return state.result(command, confirm, contextErr(ctx))
		}
	}
}

// awaitState is what awaitAckAndUpdateEvent has received so far: the ack, if
// acked, and whether an update event was signalled.
type awaitState struct {
	response  ackResponse
	acked     bool
	sawUpdate bool
}

// collectDelivered takes the ack and the update event that have already been
// delivered, without waiting for either. Receiving from a nil channel is never
// ready, so the default case covers an ack already taken.
func (s *awaitState) collectDelivered(res <-chan ackResponse, updated <-chan struct{}) {
	if !s.acked {
		select {
		case s.response = <-res:
			s.acked = true

		default:
		}
	}

	select {
	case <-updated:
		s.sawUpdate = true

	default:
	}
}

// result reports the outcome of a wait that cause ended, see
// resultOnWaitEnded.
func (s *awaitState) result(command string, confirm updateEventConfirm, cause error) (ackResponse, error) {
	settled := s.acked && confirmUpdate(confirm, s.response, s.sawUpdate)

	return resultOnWaitEnded(command, s.response, s.acked, settled, cause)
}

// confirmUpdate reports whether the write acknowledged by response is settled:
// visible in the state cache when the caller supplied a confirmation, and
// otherwise simply broadcast at least once.
func confirmUpdate(confirm updateEventConfirm, response ackResponse, sawUpdate bool) bool {
	if confirm == nil {
		return sawUpdate
	}

	return confirm(response)
}

// resultOnWaitEnded reports the outcome of a command whose wait ended before
// it was settled, from the signals that are in hand: acked tells whether the
// ack in response arrived, settled whether the write is visible in the state
// cache - the update event for a command that waits for one by name, the
// confirmation for a command that recognizes its own write. cause is what
// ended the wait: the context's error, see contextErr, or ErrConnectionLost.
func resultOnWaitEnded(
	command string,
	response ackResponse,
	acked bool,
	settled bool,
	cause error,
) (ackResponse, error) {
	switch {
	case acked && !response.OK:
		return ackResponse{}, fmt.Errorf("%s: %s", command, response.Msg)

	case acked && settled:
		// Both signals are in hand and the wait merely ended while they were
		// collected. Nothing is missing, so this is the same success the
		// caller's loop returns.
		return response, nil

	case acked:
		// The server applied the command, only the state cache has not caught
		// up with it. Returning the response lets the caller keep what the ack
		// carried, e.g. the ID of a created resource.
		return response, &UpdateEventTimeoutError{Command: command, Err: cause}

	default:
		return ackResponse{}, fmt.Errorf("%s: %w", command, cause)
	}
}
