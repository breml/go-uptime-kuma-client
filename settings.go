package kuma

import (
	"context"
)

// GetSettings retrieves the server settings.
func (c *Client) GetSettings(ctx context.Context) (map[string]any, error) {
	resp, err := c.syncEmit(ctx, "getSettings")
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// SetSettings updates the server settings.
func (c *Client) SetSettings(ctx context.Context, settings map[string]any, password string) error {
	_, err := c.syncEmit(ctx, "setSettings", settings, password)
	return err
}
