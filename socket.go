package kuma

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"

	"github.com/breml/go-uptime-kuma-client/dockerhost"
	"github.com/breml/go-uptime-kuma-client/maintenance"
	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
	"github.com/breml/go-uptime-kuma-client/proxy"
	"github.com/breml/go-uptime-kuma-client/statuspage"
)

// handledEvents names the events registerHandlers applies to the state cache
// or the auth barrier itself, and which the catch-all handler therefore leaves
// alone.
func handledEvents() []string {
	return []string{
		"notificationList",
		"monitorList",
		"updateMonitorIntoList",
		"deleteMonitorFromList",
		"statusPageList",
		"maintenanceList",
		"proxyList",
		"dockerHostList",
	}
}

// socketConn is one socket.io connection of a Client, with everything that
// belongs to that connection rather than to the client: the socket itself,
// the cancel of the context its lifetime is bound to, and the events the
// server states its auth mode with on it.
//
// The socket.io library's client is single-use - closing it is final - so
// connecting again, which is what a resync on an Uptime Kuma 3.x server takes,
// means a new socketConn with every handler registered again, see
// Client.reconnectLocked.
type socketConn struct {
	socket *socketio.Client

	// cancel ends the context the socket was connected with, which is what
	// stops its transport. It is nil until dial ran.
	cancel context.CancelFunc

	connect       chan struct{}
	setupRequired chan struct{}
	loginRequired chan struct{}
	autoLogin     chan struct{}

	closeConnect       func()
	closeSetupRequired func()
	closeLoginRequired func()
	closeAutoLogin     func()

	// closed records that the connection was torn down or never came up.
	// Guarded by Client.mu, and read by every handler so that a late event of
	// a replaced connection cannot touch the state cache.
	closed bool
}

// barrier returns the auth barrier that reads this connection's events.
func (conn *socketConn) barrier() authBarrier {
	return authBarrier{
		loginRequired: conn.loginRequired,
		autoLogin:     conn.autoLogin,
		setupRequired: conn.setupRequired,
	}
}

// newSocket creates a socket.io client for the server with every handler
// registered, not yet connected.
func (c *Client) newSocket() (*socketConn, error) {
	socket, err := newSocketIOClient(c.baseURL, c.socketioLogger, c.httpClient, c.jar)
	if err != nil {
		return nil, err
	}

	conn := &socketConn{
		socket:        socket,
		connect:       make(chan struct{}),
		setupRequired: make(chan struct{}),
		loginRequired: make(chan struct{}),
		autoLogin:     make(chan struct{}),
	}

	conn.closeConnect = sync.OnceFunc(func() { close(conn.connect) })
	conn.closeSetupRequired = sync.OnceFunc(func() { close(conn.setupRequired) })
	conn.closeLoginRequired = sync.OnceFunc(func() { close(conn.loginRequired) })
	conn.closeAutoLogin = sync.OnceFunc(func() { close(conn.autoLogin) })

	c.registerHandlers(conn)

	return conn, nil
}

// registerHandlers registers the event handlers that feed the state cache and
// the auth barrier on conn.
//
// The socket.io library runs every handler on a goroutine of its own, so two
// events of one connection can be applied in either order. The handlers that
// touch the state go through applyEvent, which drops the events of a
// connection that was torn down.
func (c *Client) registerHandlers(conn *socketConn) {
	c.registerStateHandlers(conn)
	c.registerConnectionHandlers(conn)
}

// registerStateHandlers registers the handlers for the events that carry the
// lists the state cache is built from.
func (c *Client) registerStateHandlers(conn *socketConn) {
	socket := conn.socket

	socket.On("notificationList", func(notificationList []notification.Base) {
		c.applyEvent(conn, "notificationList", func(s *state) {
			s.notifications = notificationList
		})
	})

	socket.On("monitorList", func(monitorMap map[string]monitor.Base) {
		c.applyEvent(conn, "monitorList", func(s *state) {
			s.monitors = slices.Collect(maps.Values(monitorMap))
		})
	})

	// Uptime Kuma v2 sends updateMonitorIntoList for individual monitor
	// updates (add/edit/pause/resume).
	socket.On("updateMonitorIntoList", func(monitorMap map[string]monitor.Base) {
		c.applyEvent(conn, "updateMonitorIntoList", func(s *state) {
			for _, updatedMonitor := range monitorMap {
				s.mergeMonitor(updatedMonitor)
			}
		})
	})

	// Uptime Kuma v2 sends deleteMonitorFromList when a monitor is deleted.
	socket.On("deleteMonitorFromList", func(monitorID int64) {
		c.applyEvent(conn, "deleteMonitorFromList", func(s *state) {
			s.monitors = slices.DeleteFunc(s.monitors, func(m monitor.Base) bool {
				return m.ID == monitorID
			})
		})
	})

	socket.On("statusPageList", func(statusPageMap map[int64]statuspage.StatusPage) {
		c.applyEvent(conn, "statusPageList", func(s *state) {
			s.statusPages = statusPageMap
		})
	})

	socket.On("maintenanceList", func(maintenanceMap map[string]maintenance.Maintenance) {
		c.applyEvent(conn, "maintenanceList", func(s *state) {
			s.maintenances = slices.Collect(maps.Values(maintenanceMap))
		})
	})

	socket.On("proxyList", func(proxyList []proxy.Proxy) {
		c.applyEvent(conn, "proxyList", func(s *state) {
			s.proxies = proxyList
		})
	})

	socket.On("dockerHostList", func(dockerHostList []dockerhost.DockerHost) {
		c.applyEvent(conn, "dockerHostList", func(s *state) {
			s.dockerHosts = dockerHostList
		})
	})
}

// registerConnectionHandlers registers the handlers for the events that
// describe the connection itself: the connect, the server's auth mode and its
// version, and the catch-all that announces every other event.
func (c *Client) registerConnectionHandlers(conn *socketConn) {
	socket := conn.socket

	socket.On("connect", func() {
		conn.closeConnect()
	})

	// The server states how it wants the connection to be authenticated as
	// the last thing it does for it, so these are registered before the
	// connect and unconditionally, see authBarrier.
	socket.On("setup", func() {
		conn.closeSetupRequired()
	})

	socket.On("loginRequired", func() {
		conn.closeLoginRequired()
	})

	socket.On("autoLogin", func() {
		conn.closeAutoLogin()
	})

	// Every event that has no handler of its own is still announced on the
	// update signal, for the commands that wait for one by name. The events
	// handled above announce themselves once the state holds them, so they
	// are excluded here: a second announcement racing the handler would wake
	// a waiter before the cache is current.
	socket.OnAny(func(event string, _ []any) {
		if slices.Contains(handledEvents(), event) {
			return
		}

		c.mu.Lock()
		closed := conn.closed
		c.mu.Unlock()

		if closed {
			return
		}

		c.updates.Emit(context.Background(), event)
	})
}

// mergeMonitor replaces the cached monitor with the same ID, or appends one
// the cache does not hold yet.
func (s *state) mergeMonitor(updated monitor.Base) {
	for i, existing := range s.monitors {
		if existing.ID == updated.ID {
			s.monitors[i] = updated

			return
		}
	}

	s.monitors = append(s.monitors, updated)
}

// applyEvent applies an event's payload to the state cache and announces it on
// the update signal, unless the connection it arrived on has been torn down.
func (c *Client) applyEvent(conn *socketConn, event string, apply func(s *state)) {
	c.mu.Lock()

	if conn.closed {
		c.mu.Unlock()

		return
	}

	apply(&c.state)
	c.mu.Unlock()

	c.updates.Emit(context.Background(), event)
}

// dial connects conn and waits for the server's socket.io connect, bounded by
// ctx. The socket itself is connected with a context derived from the one New
// was given, because the socket.io client keeps reading from the connection
// until that context is done.
//
// On any failure the connection is left in a state teardown can handle: a dial
// that returned an error is only ever cancelled, never closed, because the
// socket.io client's Close blocks forever after a failed Connect.
func (c *Client) dial(ctx context.Context, conn *socketConn) error {
	connCtx, cancel := context.WithCancel(c.baseCtx)
	conn.cancel = cancel

	// client.Connect returns as soon as the transport is running and the
	// handshake has been requested; a non-nil error is the real transport
	// failure, which is surfaced instead of waiting out the timeout.
	connectErr := make(chan error, 1)
	go func() {
		connectErr <- conn.socket.Connect(connCtx)
	}()

	for {
		select {
		case <-conn.connect:
			return nil

		case err := <-connectErr:
			if err != nil {
				c.markClosed(conn)
				cancel()

				return fmt.Errorf("connect to server: %w", err)
			}

			// The dial returned; keep waiting for the connect event. Setting
			// connectErr to nil makes this case block forever, and marks the
			// channel as spent by a transport that is up.
			connectErr = nil

		case <-ctx.Done():
			if connectErr == nil {
				c.disconnectOrphan(conn)
			} else {
				dialErr := c.abandonDial(conn, connectErr)
				if dialErr != nil {
					return fmt.Errorf("connect to server: %w", dialErr)
				}
			}

			return fmt.Errorf("connect to server: %w", ctx.Err())
		}
	}
}

// abandonDial cleans up the dial a caller walks away from when its wait is
// over, and returns the transport error if the dial already failed with one.
// A connection that was, or still gets, established has no owner: it is torn
// down instead of left running until the client's context is done.
func (c *Client) abandonDial(conn *socketConn, connectErr <-chan error) error {
	select {
	case err := <-connectErr:
		if err != nil {
			c.markClosed(conn)
			conn.cancel()

			return err
		}

		c.disconnectOrphan(conn)

	default:
		// The dial is still in flight (it uses the client's context, so this
		// wait does not abort it). Drain connectErr in the background so a
		// late failure is not silently discarded and a late success is not
		// leaked.
		go func() {
			err := <-connectErr
			if err != nil {
				c.socketioLogger.Warnf("dial failed after the wait for it ended: %s", err)
				c.markClosed(conn)
				conn.cancel()

				return
			}

			c.disconnectOrphan(conn)
		}()
	}

	return nil
}

// markClosed records that conn is out of service, so its handlers drop
// whatever still arrives on it. It reports whether it was the first to do so.
func (c *Client) markClosed(conn *socketConn) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn.closed {
		return false
	}

	conn.closed = true

	return true
}

// teardown takes a connected conn out of service and closes it. The context
// is cancelled first, which aborts the long poll that would otherwise hold
// the close until the server ends it.
func (c *Client) teardown(conn *socketConn) error {
	if !c.markClosed(conn) {
		return nil
	}

	if conn.cancel != nil {
		conn.cancel()
	}

	err := conn.socket.Close()
	if err != nil {
		return fmt.Errorf("close socket.io client: %w", err)
	}

	return nil
}

// disconnectOrphan closes a connection no caller holds a handle to. The close
// is best-effort and asynchronous, because it waits for the transport and the
// message loop to wind down, and the caller it would block is on its way out.
func (c *Client) disconnectOrphan(conn *socketConn) {
	go func() {
		err := c.teardown(conn)
		if err != nil {
			c.socketioLogger.Warnf("disconnect orphaned connection: %s", err)
		}
	}()
}

// installConn makes conn the connection the client's commands go out on,
// unless the client was disconnected in the meantime, which takes no new
// connection.
func (c *Client) installConn(conn *socketConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.conn = conn
}

// currentConn returns the connection the client's commands go out on.
func (c *Client) currentConn() *socketConn {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn
}
