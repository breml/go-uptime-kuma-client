package kuma

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
