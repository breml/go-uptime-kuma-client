package kuma

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// barrierIterations is how often the cases below are repeated. Several of the
// barrier's channels are ready in the same select, and select picks between
// ready cases at random, so a single run would pass for the wrong reason about
// half the time.
const barrierIterations = 100

// closedChan is a barrier event that has been delivered.
func closedChan() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)

	return ch
}

// TestAuthBarrierSetupWins covers the server that is not set up yet: it emits
// setup before it states how it wants to be authenticated, so both events are
// in by the time the barrier is read. The setup has to win every time, or the
// client sends a login that cannot succeed and the answer to the same
// connection differs from run to run.
func TestAuthBarrierSetupWins(t *testing.T) {
	tests := map[string]authBarrier{
		"with loginRequired": {
			loginRequired: closedChan(),
			autoLogin:     make(chan struct{}),
			setupRequired: closedChan(),
		},
		"with autoLogin": {
			loginRequired: make(chan struct{}),
			autoLogin:     closedChan(),
			setupRequired: closedChan(),
		},
	}

	for name, barrier := range tests {
		t.Run(name, func(t *testing.T) {
			for range barrierIterations {
				require.Equal(t, authModeSetup, barrier.await(t.Context()))
			}
		})
	}
}

// TestAuthBarrierWithoutSetup covers the servers that ask for nothing of the
// kind, which must not be answered with a setup they never requested.
func TestAuthBarrierWithoutSetup(t *testing.T) {
	tests := map[string]struct {
		barrier authBarrier
		want    authMode
	}{
		"login required": {
			barrier: authBarrier{
				loginRequired: closedChan(),
				autoLogin:     make(chan struct{}),
				setupRequired: make(chan struct{}),
			},
			want: authModeLoginRequired,
		},
		"auto login": {
			barrier: authBarrier{
				loginRequired: make(chan struct{}),
				autoLogin:     closedChan(),
				setupRequired: make(chan struct{}),
			},
			want: authModeAutoLogin,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for range barrierIterations {
				require.Equal(t, test.want, test.barrier.await(t.Context()))
			}
		})
	}
}

// TestBarrierWaitWithin covers the budget the wait for a server that states
// nothing is allowed to spend: a server too old to send either event must not
// leave the login it does answer with no time at all.
func TestBarrierWaitWithin(t *testing.T) {
	t.Run("without deadline", func(t *testing.T) {
		require.Equal(t, authBarrierWait, barrierWaitWithin(t.Context()))
	})

	t.Run("ample deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		require.Equal(t, authBarrierWait, barrierWaitWithin(ctx))
	})

	t.Run("short deadline", func(t *testing.T) {
		budget := authBarrierWait
		ctx, cancel := context.WithTimeout(t.Context(), budget)
		defer cancel()

		wait := barrierWaitWithin(ctx)

		require.Positive(t, wait)
		require.LessOrEqual(t, wait, budget/2)
	})

	t.Run("expired deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()

		require.Equal(t, time.Duration(0), barrierWaitWithin(ctx))
	})
}

// TestAwaitTOTPRetryStepAlreadyOver covers the login whose round trip crossed a
// time step boundary: the code the server rejected belongs to a step that is
// over, so the next code differs already and there is nothing left to wait for.
func TestAwaitTOTPRetryStepAlreadyOver(t *testing.T) {
	rejectedStep := totpStepStart(time.Now().Add(-totpStep))

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := awaitTOTPRetry(ctx, rejectedStep)

	require.NoError(t, err)
	require.Less(t, time.Since(start), 50*time.Millisecond)
}

// TestAwaitTOTPRetryDeadlineTooShort covers the caller whose deadline cannot
// cover the wait for the next step, which is reported instead of run into.
func TestAwaitTOTPRetryDeadlineTooShort(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
	defer cancel()

	err := awaitTOTPRetry(ctx, totpStepStart(time.Now()))

	require.ErrorIs(t, err, ErrInvalidTOTPCode)
	require.Contains(t, err.Error(), "does not fit in the deadline")
}

// TestAwaitTOTPRetryCancelled covers the context that ends while the wait is
// running, which says so rather than blaming the code.
func TestAwaitTOTPRetryCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := awaitTOTPRetry(ctx, totpStepStart(time.Now()))

	require.ErrorIs(t, err, context.Canceled)
}
