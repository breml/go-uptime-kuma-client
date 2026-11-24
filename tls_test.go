package kuma_test

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

func newTLSAuthFake(t *testing.T) (fake *fakeAckServer, server *httptest.Server) {
	t.Helper()

	fake = &fakeAckServer{messages: make(chan []byte, 32)}
	fake.setHTTPAuth()

	server = httptest.NewTLSServer(fake)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return fake, server
}

func TestNewTLS(t *testing.T) {
	fake, server := newTLSAuthFake(t)

	pool := x509.NewCertPool()
	pool.AddCert(server.Certificate())

	newClient := func(opts ...kuma.Option) error {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		c, err := kuma.New(ctx, server.URL, "user", "password", opts...)
		if err == nil {
			_ = c.Disconnect()
		}

		return err
	}

	sum := sha256.Sum256(server.Certificate().Raw)
	fingerprint := hex.EncodeToString(sum[:])

	var colons []string
	for i := 0; i < len(fingerprint); i += 2 {
		colons = append(colons, strings.ToUpper(fingerprint[i:i+2]))
	}

	colonFingerprint := strings.Join(colons, ":")
	rootCAs := kuma.WithTLSConfig(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12})

	t.Run("rejects self-signed by default", func(t *testing.T) {
		err := newClient()

		var verifyErr *tls.CertificateVerificationError

		require.ErrorAs(t, err, &verifyErr)

		var unknownAuthority x509.UnknownAuthorityError

		require.ErrorAs(t, verifyErr.Err, &unknownAuthority)
	})

	wrongPin := kuma.WithCertificateFingerprint(strings.Repeat("00", sha256.Size))
	pin := kuma.WithCertificateFingerprint(fingerprint)
	// Fails on its own: the certificate is not valid for this host name.
	wrongHost := kuma.WithTLSConfig(&tls.Config{
		RootCAs: pool, ServerName: "other.invalid", MinVersion: tls.VersionTLS12,
	})

	t.Run("custom root CA", func(t *testing.T) {
		require.NoError(t, newClient(rootCAs))
	})

	t.Run("pinned fingerprint", func(t *testing.T) {
		require.NoError(t, newClient(pin))
		require.NoError(t, newClient(kuma.WithCertificateFingerprint(colonFingerprint)))
		require.NoError(t, newClient(kuma.WithCertificateFingerprint("sha256 Fingerprint="+colonFingerprint)))
	})

	t.Run("pinned fingerprint of another certificate", func(t *testing.T) {
		err := newClient(wrongPin)

		require.ErrorIs(t, err, kuma.ErrCertificateFingerprintMismatch)
		require.ErrorContains(t, err, fingerprint, "the error names the fingerprint the server presented")
	})

	t.Run("a wrong pin is not bypassed by a trusting TLS config", func(t *testing.T) {
		require.ErrorIs(t, newClient(wrongPin, rootCAs), kuma.ErrCertificateFingerprintMismatch)
		require.ErrorIs(t, newClient(rootCAs, wrongPin), kuma.ErrCertificateFingerprintMismatch)
	})

	t.Run("malformed fingerprint", func(t *testing.T) {
		for _, bad := range []string{"abc", "zz", strings.Repeat("00", sha256.Size-1), strings.Repeat("00", sha256.Size+1)} {
			require.ErrorContains(t, newClient(kuma.WithCertificateFingerprint(bad)), "certificate fingerprint", bad)
		}
	})

	t.Run("empty fingerprint means no pin", func(t *testing.T) {
		var verifyErr *tls.CertificateVerificationError

		require.ErrorAs(t, newClient(kuma.WithCertificateFingerprint("")), &verifyErr)
	})

	t.Run("fingerprint replaces the chain and host name verification, in either order", func(t *testing.T) {
		var verifyErr *tls.CertificateVerificationError

		require.ErrorAs(t, newClient(wrongHost), &verifyErr, "the config alone rejects the host name")

		require.NoError(t, newClient(pin, wrongHost))
		require.NoError(t, newClient(wrongHost, pin))
	})

	t.Run("a later TLS config replaces an earlier one", func(t *testing.T) {
		var verifyErr *tls.CertificateVerificationError

		empty := kuma.WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})

		require.ErrorAs(t, newClient(rootCAs, empty), &verifyErr)
	})

	t.Run("the TLS config is cloned", func(t *testing.T) {
		cfg := &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}
		option := kuma.WithTLSConfig(cfg)
		cfg.RootCAs = nil

		require.NoError(t, newClient(option))
	})

	t.Run("the caller's VerifyConnection runs after the pin", func(t *testing.T) {
		errCaller := errors.New("rejected by the caller")
		calls := 0

		option := kuma.WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
			VerifyConnection: func(tls.ConnectionState) error {
				calls++

				return errCaller
			},
		})

		require.ErrorIs(t, newClient(pin, option), errCaller)
		require.Positive(t, calls)

		calls = 0

		require.ErrorIs(t, newClient(wrongPin, option), kuma.ErrCertificateFingerprintMismatch)
		require.Zero(t, calls, "a wrong pin is refused before the caller's check")
	})

	t.Run("resync reconnects over TLS", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		c, err := kuma.New(ctx, server.URL, "user", "password", pin)
		require.NoError(t, err)

		t.Cleanup(func() { _ = c.Disconnect() })

		before := fake.handshakeCount()

		require.NoError(t, c.Resync(ctx))
		require.Greater(t, fake.handshakeCount(), before)
	})
}

// TestWebSocketOverTLS covers the WebSocket upgrade, which does not go through
// the shared http.Client, see cookieWebSocket: the handshake must reach the
// server with the TLS config in place.
func TestWebSocketOverTLS(t *testing.T) {
	recorder := &headerRecorder{next: &fakeSocketIOServerWithWebSocket{}}
	server := httptest.NewTLSServer(recorder)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	pool := x509.NewCertPool()
	pool.AddCert(server.Certificate())

	sum := sha256.Sum256(server.Certificate().Raw)

	tests := map[string]kuma.Option{
		"custom root CA":     kuma.WithTLSConfig(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}),
		"pinned fingerprint": kuma.WithCertificateFingerprint(hex.EncodeToString(sum[:])),
	}

	for name, option := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			// The fake upgrades and then never sends the lists, so New fails
			// waiting for them; the handshake has been made by then.
			_, err := kuma.New(ctx, server.URL, "admin", "admin1", kuma.WithConnectTimeout(time.Second), option)
			require.Error(t, err)

			require.NotNil(t, recorder.websocketHeader(), "the client upgraded to WebSocket over TLS")
		})

		recorder.mu.Lock()
		recorder.websocket = nil
		recorder.mu.Unlock()
	}
}
