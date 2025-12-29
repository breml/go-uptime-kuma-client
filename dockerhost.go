package kuma

import (
	"context"
	"errors"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/dockerhost"
)

// GetDockerHostList returns all Docker hosts for the authenticated user.
func (c *Client) GetDockerHostList(_ context.Context) []dockerhost.DockerHost {
	c.mu.Lock()
	defer c.mu.Unlock()

	hosts := make([]dockerhost.DockerHost, len(c.state.dockerHosts))
	copy(hosts, c.state.dockerHosts)

	return hosts
}

// GetDockerHost returns a specific Docker host by ID.
func (c *Client) GetDockerHost(_ context.Context, id int64) (*dockerhost.DockerHost, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, h := range c.state.dockerHosts {
		if h.GetID() == id {
			return &h, nil
		}
	}

	return nil, fmt.Errorf("get docker host: %w", ErrNotFound)
}

// CreateDockerHost creates a new Docker host and returns its ID.
func (c *Client) CreateDockerHost(ctx context.Context, config dockerhost.Config) (int64, error) {
	response, err := c.syncEmitWithUpdateEvent(ctx, "addDockerHost", "dockerHostList", config, nil)
	if err != nil {
		return 0, fmt.Errorf("create docker host: %w", err)
	}

	return response.ID, nil
}

// UpdateDockerHost updates an existing Docker host.
func (c *Client) UpdateDockerHost(ctx context.Context, config dockerhost.Config) error {
	if config.ID == 0 {
		return errors.New("update docker host: config must have ID set")
	}

	_, err := c.syncEmitWithUpdateEvent(ctx, "addDockerHost", "dockerHostList", config, config.ID)
	if err != nil {
		return fmt.Errorf("update docker host: %w", err)
	}

	return nil
}

// DeleteDockerHost deletes a Docker host by ID.
func (c *Client) DeleteDockerHost(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteDockerHost", "dockerHostList", id)
	if err != nil {
		return fmt.Errorf("delete docker host %d: %w", id, err)
	}

	return nil
}

// TestDockerHost tests a Docker host connection without creating the host.
func (c *Client) TestDockerHost(ctx context.Context, config dockerhost.Config) (*dockerhost.TestResult, error) {
	response, err := c.syncEmit(ctx, "testDockerHost", &config)
	if err != nil {
		return nil, fmt.Errorf("test docker host: %w", err)
	}

	// Convert the response to TestResult
	result := &dockerhost.TestResult{
		OK:  response.OK,
		Msg: response.Msg,
	}

	// Extract version if present
	if version, ok := response.Data["version"]; ok {
		switch v := version.(type) {
		case string:
			result.Version = v

		case map[string]any:
			if ver, versionOk := v["Version"].(string); versionOk {
				result.Version = ver
			}
		}
	}

	return result, nil
}
