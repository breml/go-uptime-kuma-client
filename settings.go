package kuma

import (
	"context"
	"fmt"

	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
)

func (c *Client) GetSettings(ctx context.Context) (map[string]any, error) {
	res := make(chan ackResponse)
	defer close(res)

	err := c.socketioClient.Emit("getSettings",
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("getSettings: %v", err)
	}

	select {
	case response := <-res:
		if !response.OK {
			return nil, fmt.Errorf("getSettings: %s", response.Msg)
		}

		return response.Data, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("getSettings: %v", ctx.Err())
	}
}

func (c *Client) SetSettings(ctx context.Context, settings map[string]any, password string) error {
	res := make(chan ackResponse)
	defer close(res)

	err := c.socketioClient.Emit("setSettings",
		settings,
		password,
		emit.WithAck(func(response ackResponse) {
			res <- response
		}),
	)
	if err != nil {
		return fmt.Errorf("setSettings: %v", err)
	}

	select {
	case response := <-res:
		if !response.OK {
			return fmt.Errorf("setSettings: %s", response.Msg)
		}

		return nil
	case <-ctx.Done():
		return fmt.Errorf("setSettings: %v", ctx.Err())
	}
}
