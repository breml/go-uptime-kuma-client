package kuma

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestContextErr pins how a done context is described. The distinction is
// carried by the context's cause, which only operationContext ever sets, so
// the identity comparison in contextErr is what keeps a caller's own
// cancellation from being reported as the client's timeout.
func TestContextErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// ctx is the done context contextErr is asked about.
		ctx func(t *testing.T) context.Context
		// wantIs are the errors the result must report.
		wantIs []error
		// wantNotIs are the errors it must not.
		wantNotIs []error
	}{
		{
			name: "a cancelled caller context stays a cancellation",
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			wantIs:    []error{context.Canceled},
			wantNotIs: []error{ErrOperationTimeout},
		},
		{
			name: "a caller deadline stays a plain deadline",
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
				t.Cleanup(cancel)

				return ctx
			},
			wantIs:    []error{context.DeadlineExceeded},
			wantNotIs: []error{ErrOperationTimeout},
		},
		{
			name: "a caller's own cause does not displace its cancellation",
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithCancelCause(t.Context())
				cancel(errors.New("apply aborted"))

				return ctx
			},
			// context.Cause reports the cause of the first cancelled
			// ancestor, so a cause the caller recorded is in reach here.
			// Returning it would drop context.Canceled from the chain and
			// break the contract on ErrUpdateEventTimeout.
			wantIs:    []error{context.Canceled},
			wantNotIs: []error{ErrOperationTimeout},
		},
		{
			name: "the client's own budget names itself",
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				c := &Client{operationTimeout: time.Nanosecond}

				ctx, cancel := c.operationContext(t.Context(), "getMonitorList")
				t.Cleanup(cancel)

				<-ctx.Done()

				return ctx
			},
			// Still a deadline, so the checks a caller already has keep
			// reporting; ErrOperationTimeout is what says whose deadline.
			wantIs: []error{ErrOperationTimeout, context.DeadlineExceeded},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := contextErr(tt.ctx(t))
			require.Error(t, err)

			for _, want := range tt.wantIs {
				require.ErrorIs(t, err, want)
			}

			for _, notWant := range tt.wantNotIs {
				require.NotErrorIs(t, err, notWant)
			}
		})
	}
}

// TestOperationContextDisabled pins that a zero or negative budget hands the
// caller's context straight back, rather than deriving one that never expires.
func TestOperationContextDisabled(t *testing.T) {
	t.Parallel()

	for _, timeout := range []time.Duration{0, -time.Second} {
		t.Run(timeout.String(), func(t *testing.T) {
			t.Parallel()

			c := &Client{operationTimeout: timeout}

			ctx, cancel := c.operationContext(t.Context(), "getMonitorList")
			defer cancel()

			require.Same(t, t.Context(), ctx, "a disabled timeout must not derive a context")

			_, ok := ctx.Deadline()
			require.False(t, ok)
		})
	}
}

// TestOperationTimeoutErrorUnwrap pins the error the cause carries, which is
// what lets errors.Is report both the client's budget and a plain deadline.
func TestOperationTimeoutErrorUnwrap(t *testing.T) {
	t.Parallel()

	err := error(&operationTimeoutError{timeout: 30 * time.Second})

	require.ErrorIs(t, err, ErrOperationTimeout)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Contains(t, err.Error(), "30s")
	require.NotErrorIs(t, err, context.Canceled)
}

// TestOperationContextUnboundedCommand pins the exemption: the commands that
// wait on the server while it makes a request of its own get no budget, so the
// answer they exist to produce is never replaced by a timeout of the client's.
func TestOperationContextUnboundedCommand(t *testing.T) {
	t.Parallel()

	c := &Client{operationTimeout: time.Minute}

	for _, command := range unboundedCommands() {
		t.Run(command, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := c.operationContext(t.Context(), command)
			defer cancel()

			require.Same(t, t.Context(), ctx, "an exempt command must not be budgeted")

			_, ok := ctx.Deadline()
			require.False(t, ok)
		})
	}
}
