package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/proxy"
)

// GetProxyList returns all proxies for the authenticated user.
func (c *Client) GetProxyList(_ context.Context) []proxy.Proxy {
	c.mu.Lock()
	defer c.mu.Unlock()

	proxies := make([]proxy.Proxy, len(c.state.proxies))
	copy(proxies, c.state.proxies)

	return proxies
}

// GetProxy returns a specific proxy by ID.
func (c *Client) GetProxy(_ context.Context, id int64) (*proxy.Proxy, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, p := range c.state.proxies {
		if p.GetID() == id {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("get proxy: %w", ErrNotFound)
}

// CreateProxy creates a new proxy.
func (c *Client) CreateProxy(ctx context.Context, config proxy.Config) (int64, error) {
	response, err := c.syncEmitWithUpdateEvent(ctx, "addProxy", "proxyList", config, nil)
	if err != nil {
		return 0, fmt.Errorf("create proxy: %w", err)
	}

	return response.ID, nil
}

// UpdateProxy updates an existing proxy.
func (c *Client) UpdateProxy(ctx context.Context, config proxy.Config) error {
	if config.ID == 0 {
		return fmt.Errorf("update proxy: config must have ID set")
	}

	_, err := c.syncEmitWithUpdateEvent(ctx, "addProxy", "proxyList", config, config.ID)
	if err != nil {
		return fmt.Errorf("update proxy: %w", err)
	}

	return nil
}

// DeleteProxy deletes a proxy by ID.
func (c *Client) DeleteProxy(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteProxy", "proxyList", id)
	if err != nil {
		return fmt.Errorf("delete proxy %d: %w", id, err)
	}

	return nil
}
