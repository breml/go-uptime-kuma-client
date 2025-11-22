package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/proxy"
)

func TestProxyCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	t.Run("initial_state", func(t *testing.T) {
		proxies := client.GetProxyList(ctx)
		t.Logf("Initial proxies count: %d", len(proxies))
	})

	var httpProxyID int64
	t.Run("create_http_proxy_without_auth", func(t *testing.T) {
		initialProxies := client.GetProxyList(ctx)
		initialCount := len(initialProxies)

		config := proxy.Config{
			Protocol: "http",
			Host:     "proxy.example.com",
			Port:     8080,
			Auth:     false,
			Active:   true,
			Default:  false,
		}

		httpProxyID, err = client.CreateProxy(ctx, config)
		require.NoError(t, err)
		require.Greater(t, httpProxyID, int64(0))

		proxies := client.GetProxyList(ctx)
		require.Len(t, proxies, initialCount+1)

		createdProxy, err := client.GetProxy(ctx, httpProxyID)
		require.NoError(t, err)
		require.Equal(t, httpProxyID, createdProxy.GetID())
		require.Equal(t, "http", createdProxy.Protocol)
		require.Equal(t, "proxy.example.com", createdProxy.Host)
		require.Equal(t, 8080, createdProxy.Port)
		require.False(t, createdProxy.Auth)
		require.True(t, createdProxy.Active)
		require.False(t, createdProxy.Default)
	})

	var socksProxyID int64
	t.Run("create_socks5_proxy_with_auth", func(t *testing.T) {
		initialProxies := client.GetProxyList(ctx)
		initialCount := len(initialProxies)

		config := proxy.Config{
			Protocol: "socks5",
			Host:     "socks.example.com",
			Port:     1080,
			Auth:     true,
			Username: "testuser",
			Password: "testpass",
			Active:   true,
			Default:  false,
		}

		socksProxyID, err = client.CreateProxy(ctx, config)
		require.NoError(t, err)
		require.Greater(t, socksProxyID, int64(0))

		proxies := client.GetProxyList(ctx)
		require.Len(t, proxies, initialCount+1)

		createdProxy, err := client.GetProxy(ctx, socksProxyID)
		require.NoError(t, err)
		require.Equal(t, socksProxyID, createdProxy.GetID())
		require.Equal(t, "socks5", createdProxy.Protocol)
		require.Equal(t, "socks.example.com", createdProxy.Host)
		require.Equal(t, 1080, createdProxy.Port)
		require.True(t, createdProxy.Auth)
		require.Equal(t, "testuser", createdProxy.Username)
		require.Equal(t, "testpass", createdProxy.Password)
		require.True(t, createdProxy.Active)
		require.False(t, createdProxy.Default)
	})

	t.Run("update_proxy", func(t *testing.T) {
		config := proxy.Config{
			ID:       httpProxyID,
			Protocol: "https",
			Host:     "proxy-updated.example.com",
			Port:     3128,
			Auth:     true,
			Username: "admin",
			Password: "admin123",
			Active:   true,
			Default:  false,
		}

		err = client.UpdateProxy(ctx, config)
		require.NoError(t, err)

		updatedProxy, err := client.GetProxy(ctx, httpProxyID)
		require.NoError(t, err)
		require.Equal(t, "https", updatedProxy.Protocol)
		require.Equal(t, "proxy-updated.example.com", updatedProxy.Host)
		require.Equal(t, 3128, updatedProxy.Port)
		require.True(t, updatedProxy.Auth)
		require.Equal(t, "admin", updatedProxy.Username)
		require.Equal(t, "admin123", updatedProxy.Password)
	})

	t.Run("set_default_proxy", func(t *testing.T) {
		config := proxy.Config{
			ID:       httpProxyID,
			Protocol: "https",
			Host:     "proxy-updated.example.com",
			Port:     3128,
			Auth:     true,
			Username: "admin",
			Password: "admin123",
			Active:   true,
			Default:  true,
		}

		err = client.UpdateProxy(ctx, config)
		require.NoError(t, err)

		updatedProxy, err := client.GetProxy(ctx, httpProxyID)
		require.NoError(t, err)
		require.True(t, updatedProxy.Default)
	})

	t.Run("switch_default_proxy", func(t *testing.T) {
		config := proxy.Config{
			ID:       socksProxyID,
			Protocol: "socks5",
			Host:     "socks.example.com",
			Port:     1080,
			Auth:     true,
			Username: "testuser",
			Password: "testpass",
			Active:   true,
			Default:  true,
		}

		err = client.UpdateProxy(ctx, config)
		require.NoError(t, err)

		newDefaultProxy, err := client.GetProxy(ctx, socksProxyID)
		require.NoError(t, err)
		require.True(t, newDefaultProxy.Default)

		oldDefaultProxy, err := client.GetProxy(ctx, httpProxyID)
		require.NoError(t, err)
		require.False(t, oldDefaultProxy.Default)
	})

	t.Run("delete_proxy", func(t *testing.T) {
		preDeleteProxies := client.GetProxyList(ctx)
		preDeleteCount := len(preDeleteProxies)

		err := client.DeleteProxy(ctx, httpProxyID)
		require.NoError(t, err)

		proxies := client.GetProxyList(ctx)
		require.Len(t, proxies, preDeleteCount-1)

		_, err = client.GetProxy(ctx, httpProxyID)
		require.Error(t, err)
		require.ErrorIs(t, err, kuma.ErrNotFound)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteProxy(ctx, socksProxyID)
		require.NoError(t, err)
	})
}

func TestProxyWithMonitor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	var proxyID int64
	t.Run("create_proxy", func(t *testing.T) {
		config := proxy.Config{
			Protocol: "http",
			Host:     "monitor-proxy.example.com",
			Port:     8888,
			Auth:     false,
			Active:   true,
			Default:  false,
		}

		proxyID, err = client.CreateProxy(ctx, config)
		require.NoError(t, err)
		require.Greater(t, proxyID, int64(0))
	})

	var monitorID int64
	t.Run("create_monitor_with_proxy", func(t *testing.T) {
		httpMonitor := monitor.HTTP{
			Base: monitor.Base{
				Name:     "Test Monitor with Proxy",
				Interval: 60,
				ProxyID:  &proxyID,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                "https://example.com",
				Method:             "GET",
				MaxRedirects:       10,
				AcceptedStatusCodes: []string{"200-299"},
			},
		}

		monitorID, err = client.CreateMonitor(ctx, httpMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		createdMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.NotNil(t, createdMonitor.ProxyID)
		require.Equal(t, proxyID, *createdMonitor.ProxyID)
	})

	t.Run("update_monitor_remove_proxy", func(t *testing.T) {
		currentMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)

		current := monitor.HTTP{}
		err = currentMonitor.As(&current)
		require.NoError(t, err)

		current.ProxyID = nil

		err = client.UpdateMonitor(ctx, current)
		require.NoError(t, err)

		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Nil(t, updatedMonitor.ProxyID)
	})

	t.Run("update_monitor_add_proxy", func(t *testing.T) {
		currentMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)

		current := monitor.HTTP{}
		err = currentMonitor.As(&current)
		require.NoError(t, err)

		current.ProxyID = &proxyID

		err = client.UpdateMonitor(ctx, current)
		require.NoError(t, err)

		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.NotNil(t, updatedMonitor.ProxyID)
		require.Equal(t, proxyID, *updatedMonitor.ProxyID)
	})

	t.Run("delete_proxy_removes_from_monitor", func(t *testing.T) {
		err := client.DeleteProxy(ctx, proxyID)
		require.NoError(t, err)

		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Nil(t, updatedMonitor.ProxyID)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)
	})
}

func TestProxyApplyExisting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	var monitor1ID, monitor2ID int64
	t.Run("create_monitors_without_proxy", func(t *testing.T) {
		httpMonitor1 := monitor.HTTP{
			Base: monitor.Base{
				Name:     "Test Monitor 1 for ApplyExisting",
				Interval: 60,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                "https://example1.com",
				Method:             "GET",
				MaxRedirects:       10,
				AcceptedStatusCodes: []string{"200-299"},
			},
		}

		monitor1ID, err = client.CreateMonitor(ctx, httpMonitor1)
		require.NoError(t, err)

		httpMonitor2 := monitor.HTTP{
			Base: monitor.Base{
				Name:     "Test Monitor 2 for ApplyExisting",
				Interval: 60,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                "https://example2.com",
				Method:             "GET",
				MaxRedirects:       10,
				AcceptedStatusCodes: []string{"200-299"},
			},
		}

		monitor2ID, err = client.CreateMonitor(ctx, httpMonitor2)
		require.NoError(t, err)

		m1, err := client.GetMonitor(ctx, monitor1ID)
		require.NoError(t, err)
		require.Nil(t, m1.ProxyID)

		m2, err := client.GetMonitor(ctx, monitor2ID)
		require.NoError(t, err)
		require.Nil(t, m2.ProxyID)
	})

	var proxyID int64
	t.Run("create_proxy_with_apply_existing", func(t *testing.T) {
		config := proxy.Config{
			Protocol:      "http",
			Host:          "apply-existing.example.com",
			Port:          9000,
			Auth:          false,
			Active:        true,
			Default:       false,
			ApplyExisting: true,
		}

		proxyID, err = client.CreateProxy(ctx, config)
		require.NoError(t, err)
		require.Greater(t, proxyID, int64(0))

		m1, err := client.GetMonitor(ctx, monitor1ID)
		require.NoError(t, err)
		require.NotNil(t, m1.ProxyID)
		require.Equal(t, proxyID, *m1.ProxyID)

		m2, err := client.GetMonitor(ctx, monitor2ID)
		require.NoError(t, err)
		require.NotNil(t, m2.ProxyID)
		require.Equal(t, proxyID, *m2.ProxyID)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteMonitor(ctx, monitor1ID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, monitor2ID)
		require.NoError(t, err)

		err = client.DeleteProxy(ctx, proxyID)
		require.NoError(t, err)
	})
}

func TestProxyProtocols(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	protocols := []string{"http", "https", "socks", "socks4", "socks5", "socks5h"}

	for _, protocol := range protocols {
		t.Run("protocol_"+protocol, func(t *testing.T) {
			config := proxy.Config{
				Protocol: protocol,
				Host:     "test-" + protocol + ".example.com",
				Port:     9000,
				Auth:     false,
				Active:   true,
				Default:  false,
			}

			id, err := client.CreateProxy(ctx, config)
			require.NoError(t, err)
			require.Greater(t, id, int64(0))

			createdProxy, err := client.GetProxy(ctx, id)
			require.NoError(t, err)
			require.Equal(t, protocol, createdProxy.Protocol)

			err = client.DeleteProxy(ctx, id)
			require.NoError(t, err)
		})
	}
}
