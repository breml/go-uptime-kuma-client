package kuma

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ErrCertificateFingerprintMismatch is returned, wrapped, when the server
// presents a certificate other than the one pinned with
// WithCertificateFingerprint.
var ErrCertificateFingerprintMismatch = errors.New("server certificate does not match the pinned fingerprint")

// WithTLSConfig sets the TLS configuration for every connection the client
// makes: the HTTP requests (login, autosetup, long-polling) and the WebSocket
// handshake. Use it to trust a private CA (RootCAs) or to present a client
// certificate. The config is cloned, so later changes by the caller have no
// effect. A later WithTLSConfig replaces an earlier one;
// WithCertificateFingerprint is independent of it and combines with it in
// either order.
func WithTLSConfig(cfg *tls.Config) Option {
	cloned := cfg.Clone()

	return func(c *Client) {
		c.tlsConfig = cloned
	}
}

// WithCertificateFingerprint pins the server's certificate: the SHA-256
// fingerprint of its leaf certificate, in hex, upper or lower case, with
// optional ':' separators, as `openssl x509 -fingerprint -sha256` prints it.
// It allows connecting to a server with a self-signed certificate: the chain
// and the host name are not verified, instead every connection, the HTTP
// requests and the WebSocket, is refused unless the server presents exactly
// that certificate, see ErrCertificateFingerprintMismatch. It combines with
// WithTLSConfig in either order.
func WithCertificateFingerprint(fingerprint string) Option {
	return func(c *Client) {
		c.certFingerprint = fingerprint
	}
}

// buildTLSConfig returns the TLS configuration WithTLSConfig and
// WithCertificateFingerprint ask for, or nil for the default verification.
// It runs after all options are applied, so their order does not matter.
func (c *Client) buildTLSConfig() (*tls.Config, error) {
	if c.certFingerprint == "" {
		return c.tlsConfig, nil
	}

	want, err := parseFingerprint(c.certFingerprint)
	if err != nil {
		return nil, err
	}

	cfg := c.tlsConfig.Clone()
	if cfg == nil {
		cfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	next := cfg.VerifyConnection

	// The chain is not verified; the pin takes its place. VerifyConnection,
	// unlike VerifyPeerCertificate, also runs for resumed sessions.
	cfg.InsecureSkipVerify = true
	cfg.VerifyConnection = func(state tls.ConnectionState) error {
		if len(state.PeerCertificates) == 0 {
			return fmt.Errorf("%w: the server presented no certificate", ErrCertificateFingerprintMismatch)
		}

		got := sha256.Sum256(state.PeerCertificates[0].Raw)
		if subtle.ConstantTimeCompare(got[:], want) != 1 {
			return fmt.Errorf("%w: the server's is %s", ErrCertificateFingerprintMismatch, hex.EncodeToString(got[:]))
		}

		if next != nil {
			return next(state)
		}

		return nil
	}

	return cfg, nil
}

// parseFingerprint reads the SHA-256 fingerprint as `openssl x509
// -fingerprint -sha256` prints it, with its "sha256 Fingerprint=" prefix or
// without, or as plain hex.
func parseFingerprint(fingerprint string) ([]byte, error) {
	value := strings.TrimSpace(fingerprint)

	if algorithm, rest, found := strings.Cut(value, "="); found {
		name, label, _ := strings.Cut(strings.ToLower(strings.TrimSpace(algorithm)), " ")
		if name != "sha256" || label != "fingerprint" {
			return nil, errBadFingerprint
		}

		value = rest
	}

	value = strings.TrimSpace(value)

	// With separators, every byte is two digits: "A:BC" and "AB::CD" are typos.
	if groups := strings.Split(value, ":"); len(groups) > 1 {
		for _, group := range groups {
			if len(group) != 2 {
				return nil, errBadFingerprint
			}
		}
	}

	want, err := hex.DecodeString(strings.ReplaceAll(value, ":", ""))
	if err != nil || len(want) != sha256.Size {
		return nil, errBadFingerprint
	}

	return want, nil
}

var errBadFingerprint = fmt.Errorf(
	"certificate fingerprint: want the SHA-256 as %d hex bytes, optionally separated by ':'",
	sha256.Size,
)

// newHTTPClient builds the http.Client the probe, the login and the long
// polling share: the cookie jar, the configured headers and, for the TLS
// configuration the options ask for, a clone of the default transport.
func (c *Client) newHTTPClient(jar http.CookieJar, httpBase *url.URL) (*http.Client, error) {
	tlsConfig, err := c.buildTLSConfig()
	if err != nil {
		return nil, err
	}

	c.tlsConfig = tlsConfig

	base := http.DefaultTransport

	if tlsConfig != nil {
		transport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return nil, errors.New("tls: default http transport is not an *http.Transport")
		}

		transport = transport.Clone()
		transport.TLSClientConfig = tlsConfig.Clone()
		base = transport
	}

	return &http.Client{
		Jar:           jar,
		Transport:     &headerRoundTripper{base: base, headers: c.requestHeaderSet},
		CheckRedirect: checkRedirect(httpBase, c.redirectNote()),
	}, nil
}
