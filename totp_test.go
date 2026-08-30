package kuma_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// rfc6238Secret is the base32 form of the RFC 6238 test key
// "12345678901234567890".
const rfc6238Secret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

// TestTOTPCodeAtRFC6238 checks the generated codes against the SHA1 test
// vectors of RFC 6238, appendix B. Uptime Kuma verifies with the same
// parameters, so a code this agrees with is a code the server accepts.
func TestTOTPCodeAtRFC6238(t *testing.T) {
	tests := []struct {
		name string
		unix int64
		want string
	}{
		{name: "1970", unix: 59, want: "287082"},
		{name: "2005", unix: 1111111109, want: "081804"},
		{name: "2005 next step", unix: 1111111111, want: "050471"},
		{name: "2009", unix: 1234567890, want: "005924"},
		{name: "2033", unix: 2000000000, want: "279037"},
		{name: "2603", unix: 20000000000, want: "353130"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := kuma.TOTPCodeAt(rfc6238Secret, time.Unix(test.unix, 0).UTC())

			require.NoError(t, err)
			require.Equal(t, test.want, code)
		})
	}
}

// TestNormalizeTOTPSecret covers the shapes a secret reaches the client in.
// Uptime Kuma strips the base32 padding from the secret it puts in the
// otpauth:// URI, so restoring it is what makes such a secret usable at all.
func TestNormalizeTOTPSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		want    string
		wantErr bool
	}{
		{name: "already normalized", secret: rfc6238Secret, want: rfc6238Secret},
		{name: "lower case", secret: "gezdgnbvgy3tqojqgezdgnbvgy3tqojq", want: rfc6238Secret},
		{name: "spaced", secret: "GEZD GNBV GY3T QOJQ GEZD GNBV GY3T QOJQ", want: rfc6238Secret},
		{name: "hyphenated", secret: "GEZD-GNBV-GY3T-QOJQ-GEZD-GNBV-GY3T-QOJQ", want: rfc6238Secret},
		{
			name:   "unpadded",
			secret: "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQGEZA",
			want:   "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQGEZA====",
		},
		{
			name:   "padded",
			secret: "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQGEZA====",
			want:   "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQGEZA====",
		},
		{name: "empty", secret: "", wantErr: true},
		{name: "blank", secret: "   ", wantErr: true},
		{name: "not base32", secret: "not-a-secret!", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := kuma.NormalizeTOTPSecret(test.secret)

			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

// TestNormalizeTOTPSecretServerShape covers the secret Uptime Kuma actually
// hands out: the 64 random alphanumeric characters of genSecret, base32
// encoded and stripped of the one padding character that leaves. 103 is what
// makes this the interesting case - a length that is not a multiple of 8, so
// the padding has to be restored before the secret decodes at all.
func TestNormalizeTOTPSecretServerShape(t *testing.T) {
	const serverSecret = "MRXUIZ2RKI4TMY2JGFGDMWJZONETK5KZGE2WE3SGHBHEUUKIOVMU" +
		"OSCSHEZUY3BSJF2GUUKXNNVFCWRXJ44DKN3TOBIHK2JSN52DG5Y"

	require.NotZero(t, len(serverSecret)%8, "the fixture has to need the padding restored")

	normalized, err := kuma.NormalizeTOTPSecret(serverSecret)

	require.NoError(t, err)
	require.Zero(t, len(normalized)%8, "the decoder needs the padding back")

	code, err := kuma.TOTPCodeAt(normalized, time.Unix(1234567890, 0).UTC())

	require.NoError(t, err)
	require.Len(t, code, 6)
}
