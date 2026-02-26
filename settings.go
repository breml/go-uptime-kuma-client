package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/settings"
)

// GetSettings retrieves the server settings.
func (c *Client) GetSettings(ctx context.Context) (*settings.Settings, error) {
	resp, err := c.syncEmit(ctx, "getSettings")
	if err != nil {
		return nil, err
	}

	var s settings.Settings
	err = convertToStruct(resp.Data, &s)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	return &s, nil
}

// SetSettings updates the server settings.
func (c *Client) SetSettings(ctx context.Context, s settings.Settings, password string) error {
	settingsMap, err := structToMap(s)
	if err != nil {
		return fmt.Errorf("set settings: %w", err)
	}

	_, err = c.syncEmit(ctx, "setSettings", settingsMap, password)
	return err
}
