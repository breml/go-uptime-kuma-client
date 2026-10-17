package kuma

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	engineio "github.com/maldikhan/go.socket.io/engine.io/v4/client"
	enginepolling "github.com/maldikhan/go.socket.io/engine.io/v4/client/transport/polling"
	enginews "github.com/maldikhan/go.socket.io/engine.io/v4/client/transport/websocket"
	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"golang.org/x/net/websocket"
)

// errWebSocketNotConnected mirrors the sentinel the socket.io library's own
// WebSocket implementation reports for a connection that was never dialed.
var errWebSocketNotConnected = errors.New("socket connection is not initialized")

// defaultPollingInterval is the interval the socket.io library polls with until
// the server's handshake states its own.
const defaultPollingInterval = 10 * time.Second

// newSocketIOClient builds the socket.io client for baseURL with transports
// that carry the cookies in jar: the long-polling requests through httpClient,
// which holds the same jar, and the WebSocket handshake through
// cookieWebSocket.
//
// The library's own transports take neither a cookie jar nor a header, and an
// Uptime Kuma 3.x server authenticates a socket.io connection by the session
// cookie on its handshake alone, see AuthSchemeHTTP. Both transports are built
// the way the library builds them by default otherwise.
func newSocketIOClient(
	baseURL string,
	logger socketio.Logger,
	httpClient *http.Client,
	jar http.CookieJar,
) (*socketio.Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	if parsed.Path == "" {
		parsed.Path = "/socket.io/"
	}

	wsTransport, err := enginews.NewTransport(
		enginews.WithLogger(logger),
		enginews.WithWebSocket(&cookieWebSocket{jar: jar}),
	)
	if err != nil {
		return nil, fmt.Errorf("create websocket transport: %w", err)
	}

	pollingTransport, err := enginepolling.NewTransport(
		enginepolling.WithLogger(logger),
		enginepolling.WithDefaultPinger(time.NewTicker(defaultPollingInterval)),
		enginepolling.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("create polling transport: %w", err)
	}

	engineClient, err := engineio.NewClient(
		engineio.WithURL(parsed),
		engineio.WithLogger(logger),
		engineio.WithSupportedTransports([]engineio.Transport{wsTransport, pollingTransport}),
	)
	if err != nil {
		return nil, fmt.Errorf("create engine.io client: %w", err)
	}

	client, err := socketio.NewClient(
		socketio.WithEngineIOClient(engineClient),
		socketio.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %w", err)
	}

	return client, nil
}

// cookieWebSocket is the socket.io library's WebSocket implementation with
// one addition: the handshake request carries the cookies the jar holds for
// the target, which is how the session cookie reaches the server on an
// upgraded connection.
type cookieWebSocket struct {
	jar  http.CookieJar
	conn *websocket.Conn
}

// Dial opens the WebSocket connection with the jar's cookies on the handshake.
func (w *cookieWebSocket) Dial(ctx context.Context, target *url.URL, origin *url.URL) error {
	config, err := websocket.NewConfig(target.String(), origin.String())
	if err != nil {
		return fmt.Errorf("create websocket config: %w", err)
	}

	if cookie := cookieHeader(w.jar, target); cookie != "" {
		config.Header.Set("Cookie", cookie)
	}

	conn, err := config.DialContext(ctx)
	if err != nil {
		return fmt.Errorf("dial websocket: %w", err)
	}

	w.conn = conn

	return nil
}

// Send writes a message frame to the WebSocket connection.
func (w *cookieWebSocket) Send(v []byte) error {
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
func (w *cookieWebSocket) Receive(v *[]byte) error {
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
func (w *cookieWebSocket) Close() error {
	if w.conn == nil {
		return nil
	}

	err := w.conn.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

// cookieHeader renders the cookies jar holds for target as a Cookie header
// value. The jar indexes by the http scheme, so a ws or wss target is looked up
// under the scheme its handshake is actually made over.
func cookieHeader(jar http.CookieJar, target *url.URL) string {
	if jar == nil {
		return ""
	}

	lookup := *target

	switch lookup.Scheme {
	case "ws":
		lookup.Scheme = "http"

	case "wss":
		lookup.Scheme = "https"

	default:
	}

	cookies := jar.Cookies(&lookup)

	parts := make([]string, 0, len(cookies))
	for _, cookie := range cookies {
		parts = append(parts, cookie.Name+"="+cookie.Value)
	}

	return strings.Join(parts, "; ")
}
