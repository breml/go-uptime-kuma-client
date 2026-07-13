package kuma

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	engineioclient "github.com/maldikhan/go.socket.io/engine.io/v4/client"
	enginepolling "github.com/maldikhan/go.socket.io/engine.io/v4/client/transport/polling"
	enginews "github.com/maldikhan/go.socket.io/engine.io/v4/client/transport/websocket"
	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"

	"golang.org/x/net/websocket"
)

// errWebSocketNotConnected mirrors the unexported sentinel error used by the
// default socket.io WebSocket implementation.
var errWebSocketNotConnected = errors.New("socket connection is not initialized")

// setUserAgent sets the User-Agent header on req, unless userAgent is empty.
func setUserAgent(req *http.Request, userAgent string) {
	if userAgent == "" {
		return
	}

	req.Header.Set("User-Agent", userAgent)
}

// userAgentRoundTripper injects a User-Agent header into every request that
// does not already carry one.
type userAgentRoundTripper struct {
	userAgent string
}

// RoundTrip implements http.RoundTripper.
func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req = req.Clone(req.Context())
		req.Header.Set("User-Agent", rt.userAgent)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("round trip: %w", err)
	}

	return resp, nil
}

// userAgentWebSocket wraps golang.org/x/net/websocket to attach a custom
// User-Agent header during the WebSocket handshake. The socket.io client's
// default WebSocket implementation (ws_native.WebSocketConnection) does not
// expose any header configuration, so this type re-implements the same
// enginews.WebSocket interface with header support.
type userAgentWebSocket struct {
	userAgent string
	conn      *websocket.Conn
}

// Dial opens the WebSocket connection, attaching the configured User-Agent
// header to the handshake request.
func (w *userAgentWebSocket) Dial(ctx context.Context, target *url.URL, origin *url.URL) error {
	config, err := websocket.NewConfig(target.String(), origin.String())
	if err != nil {
		return fmt.Errorf("create websocket config: %w", err)
	}

	if w.userAgent != "" {
		config.Header.Set("User-Agent", w.userAgent)
	}

	conn, err := config.DialContext(ctx)
	if err != nil {
		return fmt.Errorf("dial websocket: %w", err)
	}

	w.conn = conn

	return nil
}

// Send writes a message frame to the WebSocket connection.
func (w *userAgentWebSocket) Send(v []byte) error {
	if w.conn == nil {
		return errWebSocketNotConnected
	}

	err := websocket.Message.Send(w.conn, string(v))
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

// Receive reads the next message frame from the WebSocket connection.
func (w *userAgentWebSocket) Receive(v *[]byte) error {
	if w.conn == nil {
		return errWebSocketNotConnected
	}

	err := websocket.Message.Receive(w.conn, v)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	return nil
}

// Close closes the WebSocket connection.
func (w *userAgentWebSocket) Close() error {
	if w.conn == nil {
		return nil
	}

	err := w.conn.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

// newEngineIOClientOption builds the socket.io ClientOption used to connect
// to baseURL. When userAgent is empty, it falls back to the library's
// default transports (socketio.WithRawURL) unmodified. Otherwise, it builds
// an Engine.IO client with WebSocket and polling transports that attach
// userAgent to every request.
func newEngineIOClientOption(baseURL string, logger socketio.Logger, userAgent string) (socketio.ClientOption, error) {
	if userAgent == "" {
		return socketio.WithRawURL(baseURL), nil
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	if parsedURL.Path == "" {
		parsedURL.Path = "/socket.io/"
	}

	wsTransport, err := enginews.NewTransport(
		enginews.WithLogger(logger),
		enginews.WithWebSocket(&userAgentWebSocket{userAgent: userAgent}),
	)
	if err != nil {
		return nil, fmt.Errorf("create websocket transport: %w", err)
	}

	pollingTransport, err := enginepolling.NewTransport(
		enginepolling.WithLogger(logger),
		enginepolling.WithHTTPClient(&http.Client{
			Transport: &userAgentRoundTripper{userAgent: userAgent},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create polling transport: %w", err)
	}

	engineioClient, err := engineioclient.NewClient(
		engineioclient.WithURL(parsedURL),
		engineioclient.WithLogger(logger),
		engineioclient.WithSupportedTransports([]engineioclient.Transport{wsTransport, pollingTransport}),
	)
	if err != nil {
		return nil, fmt.Errorf("create engine.io client: %w", err)
	}

	return socketio.WithEngineIOClient(engineioClient), nil
}
