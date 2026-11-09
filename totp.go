package kuma

import (
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// totpStep is the time step RFC 6238 uses and the one the server verifies
// with, see twoFAVerifyOptions in the server's server.js.
const totpStep = 30 * time.Second

// ErrInvalidTOTPSecret is returned when a TOTP secret is empty or not base32,
// so no one-time code can be derived from it. New reports it for the secret
// given to WithTOTPSecret, before any connection is made, and
// ValidateTOTPSecret reports it for a secret checked up front.
var ErrInvalidTOTPSecret = errors.New("invalid totp secret")

// ValidateTOTPSecret reports whether secret is one WithTOTPSecret accepts,
// wrapping ErrInvalidTOTPSecret when it is not. It applies the same rules, so
// a caller can reject a malformed secret where it is configured, for example
// in a configuration validator, rather than only once New runs.
func ValidateTOTPSecret(secret string) error {
	_, err := normalizeTOTPSecret(secret)

	return err
}

// normalizeTOTPSecret returns the secret in the padded upper-case base32 form
// base32.StdEncoding decodes, and reports a secret it cannot decode.
//
// The code generation tolerates a secret that is unpadded or lower case on its
// own, so neither of those is why this exists. What it is for is the separators
// a secret carries when it is copied out of a user interface - the server shows
// it in groups - which are stripped from the middle rather than only from the
// ends, and the decode below, which turns a mistyped secret into an error from
// New instead of a login the server rejects for reasons the caller cannot see.
func normalizeTOTPSecret(secret string) (string, error) {
	normalized := strings.ToUpper(strings.NewReplacer(" ", "", "-", "", "=", "").Replace(secret))
	if normalized == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidTOTPSecret)
	}

	if padding := len(normalized) % 8; padding != 0 {
		normalized += strings.Repeat("=", 8-padding)
	}

	_, err := base32.StdEncoding.DecodeString(normalized)
	if err != nil {
		return "", fmt.Errorf("%w: not base32: %w", ErrInvalidTOTPSecret, err)
	}

	return normalized, nil
}

// totpCodeAt returns the one-time code for secret at t, with the parameters
// the server verifies with: HMAC-SHA1, a 30 second step and six digits.
func totpCodeAt(secret string, t time.Time) (string, error) {
	code, err := totp.GenerateCodeCustom(secret, t, totp.ValidateOpts{
		Period:    uint(totpStep / time.Second),
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", fmt.Errorf("totp code: %w", err)
	}

	return code, nil
}

// totpStepStart returns the beginning of the time step t falls into.
//
// The server records the last code it accepted and refuses to see it again, so
// a code that was rejected for that reason is only worth replacing once the
// step it belongs to is over.
func totpStepStart(t time.Time) time.Time {
	return t.Truncate(totpStep)
}
