package kuma

import (
	"context"
	"errors"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/proxy"
)

// GetProxyList returns all proxies for the authenticated user.
//
// They are served from the local state cache, which the server keeps up to
// date. Resync rebuilds it if an update event was missed, see
// ErrUpdateEventTimeout.
func (c *Client) GetProxyList(_ context.Context) []proxy.Proxy {
	c.mu.Lock()
	defer c.mu.Unlock()

	proxies := make([]proxy.Proxy, len(c.state.proxies))
	copy(proxies, c.state.proxies)

	return proxies
}

// GetProxy returns a specific proxy by ID. Like GetProxyList it serves from the
// local state cache, so a proxy whose update event was missed is reported as
// ErrNotFound until Resync.
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
//
// An error wrapping ErrUpdateEventTimeout means the proxy was created and only
// the update event is missing; the returned ID identifies it, as does the ID an
// UpdateEventTimeoutError carries, so retrying the call would create a
// duplicate.
func (c *Client) CreateProxy(ctx context.Context, config proxy.Config) (int64, error) {
	response, err := c.syncEmitWithConfirmedUpdateEvent(
		ctx, "addProxy", "proxyList",
		func(response ackResponse) bool { return c.hasProxy(response.ID) },
		config, nil,
	)
	if err != nil {
		if errors.Is(err, ErrUpdateEventTimeout) {
			return response.ID, fmt.Errorf("create proxy: %w", withCreatedID(err, response.ID))
		}

		return 0, fmt.Errorf("create proxy: %w", err)
	}

	return response.ID, nil
}

// UpdateProxy updates an existing proxy.
//
// An error wrapping ErrUpdateEventTimeout means the proxy was updated and only
// the update event is missing.
//
// Unlike CreateProxy and DeleteProxy this waits for the broadcast by name alone,
// because an edit changes no property the cache can be checked against: the
// proxy is in the list before and after it. With several writes in flight the
// broadcast it observes may therefore be a neighbour's, and the cache can hold
// the pre-edit values until the next one arrives. A read that has to see the
// edit resyncs.
func (c *Client) UpdateProxy(ctx context.Context, config proxy.Config) error {
	if config.ID == 0 {
		return errors.New("update proxy: config must have ID set")
	}

	_, err := c.syncEmitWithUpdateEvent(ctx, "addProxy", "proxyList", config, config.ID)
	if err != nil {
		return fmt.Errorf("update proxy: %w", err)
	}

	return nil
}

// DeleteProxy deletes a proxy by ID.
//
// An error wrapping ErrUpdateEventTimeout means the proxy was deleted and only
// the update event is missing.
func (c *Client) DeleteProxy(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithConfirmedUpdateEvent(
		ctx, "deleteProxy", "proxyList",
		func(ackResponse) bool { return !c.hasProxy(id) },
		id,
	)
	if err != nil {
		return fmt.Errorf("delete proxy %d: %w", id, err)
	}

	return nil
}

// hasProxy reports whether the state cache holds a proxy with the given ID.
func (c *Client) hasProxy(id int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return cachedByID(c.state.proxies, id)
}
