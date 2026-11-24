package kuma_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

func TestParseFingerprint(t *testing.T) {
	const hexValue = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	colons := strings.ToUpper(strings.Join(splitPairs(hexValue), ":"))

	valid := map[string]string{
		"lowercase":           hexValue,
		"uppercase":           strings.ToUpper(hexValue),
		"colons":              colons,
		"surrounding spaces":  "  " + colons + "\n",
		"openssl 3 line":      "sha256 Fingerprint=" + colons,
		"openssl 1 line":      "SHA256 Fingerprint=" + colons + "\n",
		"prefix, spaced hex":  "sha256 Fingerprint = " + colons,
		"prefix with no case": "Sha256 FINGERPRINT=" + hexValue,
	}

	for name, input := range valid {
		t.Run("valid "+name, func(t *testing.T) {
			got, err := kuma.ParseFingerprint(input)
			require.NoError(t, err)
			require.Len(t, got, 32)
			require.Equal(t, byte(0x01), got[0])
		})
	}

	invalid := map[string]string{
		"empty":               "",
		"empty after prefix":  "sha256 Fingerprint=",
		"31 bytes":            hexValue[:62],
		"33 bytes":            hexValue + "00",
		"not hex":             strings.Repeat("zz", 32),
		"sha1 prefix":         "sha1 Fingerprint=" + colons,
		"unknown prefix":      "md5=" + colons,
		"space separated":     strings.Join(splitPairs(hexValue), " "),
		"odd length":          hexValue + "0",
		"double separator":    strings.Replace(colons, ":", "::", 1),
		"prefix, wrong label": "sha256 Digest=" + colons,
	}

	for name, input := range invalid {
		t.Run("invalid "+name, func(t *testing.T) {
			_, err := kuma.ParseFingerprint(input)
			require.ErrorContains(t, err, "certificate fingerprint")
		})
	}
}

func splitPairs(hexValue string) []string {
	pairs := make([]string, 0, len(hexValue)/2)

	for i := 0; i < len(hexValue); i += 2 {
		pairs = append(pairs, hexValue[i:i+2])
	}

	return pairs
}
