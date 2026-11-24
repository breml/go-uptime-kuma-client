package kuma_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
)

func newTLSAuthFake(t *testing.T) (server *httptest.Server) {
	t.Helper()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	fake.setHTTPAuth()

	server = httptest.NewTLSServer(fake)

	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	return server
}

func TestNewTLS(t *testing.T) {
	server := newTLSAuthFake(t)

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

	t.Run("rejects self-signed by default", func(t *testing.T) {
		require.Error(t, newClient())
	})

	t.Run("insecure skip verify", func(t *testing.T) {
		require.NoError(t, newClient(kuma.WithInsecureSkipVerify()))
	})

	t.Run("custom root CA", func(t *testing.T) {
		require.NoError(t, newClient(kuma.WithTLSConfig(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12})))
	})
}
