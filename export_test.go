package kuma

import "time"

// The two-factor management commands and the code generation are internal:
// they administer an account rather than authenticate a client, and there is
// no request for them as public API yet. The end-to-end tests need them to put
// a real server into the state they exercise, and those live in the external
// test package, so they are exported for the tests alone.
//
//nolint:gochecknoglobals // the standard way to reach internals from an external test package.
var (
	Prepare2FA  = (*Client).prepare2FA
	Save2FA     = (*Client).save2FA
	Disable2FA  = (*Client).disable2FA
	TwoFAStatus = (*Client).twoFAStatus

	NormalizeTOTPSecret = normalizeTOTPSecret
	TOTPCodeAt          = totpCodeAt
)

// SetSetupEventGrace widens the window setupPending gives a setup event that
// has not arrived yet, and returns the function that puts it back.
//
// The window settles a race whose timing a fake server cannot hit reliably, so
// a test that wants the late-setup path takes the timing out of the assertion
// rather than tuning a delay against it.
func SetSetupEventGrace(d time.Duration) func() {
	previous := setupEventGrace
	setupEventGrace = d

	return func() { setupEventGrace = previous }
}

// SetDefaultOperationTimeout shortens the budget New applies when the caller
// configures none, and returns the function that puts it back.
//
// It exists so a test can exercise the default itself rather than a value it
// passed in: without it, removing the default from New would leave every
// command unbounded and the suite green.
func SetDefaultOperationTimeout(d time.Duration) func() {
	previous := defaultOperationTimeout
	defaultOperationTimeout = d

	return func() { defaultOperationTimeout = previous }
}
