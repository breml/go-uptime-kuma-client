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

// normalizeTOTPSecret returns the secret in the padded upper-case base32 form
// the code generation expects, and reports a secret it cannot decode.
//
// The server hands the secret out in the otpauth:// URI it shows when
// two-factor authentication is set up, with the base32 padding stripped, so
// restoring the padding is what makes such a secret usable at all. A secret
// copied out of a user interface may also carry spaces, hyphens and lower
// case, all of which are accepted.
func normalizeTOTPSecret(secret string) (string, error) {
	normalized := strings.ToUpper(strings.NewReplacer(" ", "", "-", "", "=", "").Replace(secret))
	if normalized == "" {
		return "", errors.New("totp secret: empty")
	}

	if padding := len(normalized) % 8; padding != 0 {
		normalized += strings.Repeat("=", 8-padding)
	}

	_, err := base32.StdEncoding.DecodeString(normalized)
	if err != nil {
		return "", fmt.Errorf("totp secret: not base32: %w", err)
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
