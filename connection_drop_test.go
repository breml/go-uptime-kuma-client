package kuma_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/notification"
)

// dropOperationTimeout is the operation timeout the drop tests give the
// client. A command that is not failed by the drop itself runs into it, which
// is what tells a detected drop from an undetected one.
const dropOperationTimeout = 5 * time.Second

// dropDetectionBudget is how long a command may take to fail once the
// connection it goes out on is gone.
const dropDetectionBudget = time.Second

// dropSilentPingInterval and dropSilentPingTimeout are the ping interval and
// timeout, in milliseconds, the dropping polling server states: a client takes
// a server that went silent for gone after both together, see
// silentDropDetection.
const (
	dropSilentPingInterval = 50
	dropSilentPingTimeout  = 200
)

// silentDropDetection is how long a client takes to notice a server that went
// silent, with a margin.
const silentDropDetection = 2 * (dropSilentPingInterval + dropSilentPingTimeout) * time.Millisecond

// dropMode is how a fake server loses the connection.
type dropMode int

const (
	// dropSessionClosed answers every further request of the session with the
	// 400 an engine.io server answers a session it no longer holds with, the
	// way a restarted server does.
	dropSessionClosed dropMode = iota

	// dropSilent stops answering without closing anything, the way a server
	// behind a network partition looks.
	dropSilent

	// dropWebSocketClosed closes the WebSocket connection, the way a server
	// that shuts down does.
	dropWebSocketClosed
)

// dropper is a fake server that can lose the connection on demand.
type dropper interface {
	http.Handler
	drop()
	// release ends the requests a silent server holds, so the test server can
	// shut down.
	release()
}

// droppingPollingServer is fakeSocketIOServer that loses the session when
// drop is called, by mode.
type droppingPollingServer struct {
	*fakeSocketIOServer

	mode dropMode
	// holdCommands holds the request of every command sent from then on
	// open, without answering it, until the server is released, see
	// heldCommands.
	holdCommands atomic.Bool
	heldCommands atomic.Int64

	dropped     chan struct{}
	dropOnce    sync.Once
	released    chan struct{}
	releaseOnce sync.Once
}

func newDroppingPollingServer(mode dropMode) *droppingPollingServer {
	return &droppingPollingServer{
		fakeSocketIOServer: &fakeSocketIOServer{
			messages: make(chan []byte, 20),
			eventsAfterLogin: []string{
				`42["monitorList",{}]`,
				`42["notificationList",[]]`,
				`42["statusPageList",{}]`,
			},
			// Every command after the login stays unanswered, so one can be
			// in flight when the connection drops.
			onCommand:    func(string, string) {},
			pingInterval: dropSilentPingInterval,
			pingTimeout:  dropSilentPingTimeout,
		},
		mode:     mode,
		dropped:  make(chan struct{}),
		released: make(chan struct{}),
	}
}

func (s *droppingPollingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	if sid == "" || r.URL.Path == "/api/auth/ok" {
		s.fakeSocketIOServer.ServeHTTP(w, r)

		return
	}

	select {
	case <-s.dropped:
		s.answerDropped(w, r)

		return

	default:
	}

	if r.Method != http.MethodGet {
		if s.holdCommands.Load() {
			s.holdCommand(w, r)

			return
		}

		s.fakeSocketIOServer.ServeHTTP(w, r)

		return
	}

	// Until the drop, an idle poll ends in a ping, the way an engine.io server
	// keeps a session alive, so the client has no reason to take the server
	// for gone before it is.
	ping := time.NewTimer(dropSilentPingInterval * time.Millisecond)
	defer ping.Stop()

	select {
	case msg := <-s.messages:
		_, _ = w.Write(msg)

	case <-ping.C:
		_, _ = w.Write([]byte(fakePingPacket))

	case <-s.dropped:
		s.answerDropped(w, r)

	case <-r.Context().Done():
	}
}

// holdCommand holds a POST that carries a socket.io event open until the server
// is released, and answers any other, such as a pong.
func (s *droppingPollingServer) holdCommand(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusInternalServerError)

		return
	}

	if !carriesEvent(body) {
		w.WriteHeader(http.StatusOK)

		return
	}

	s.heldCommands.Add(1)

	select {
	case <-s.released:
	case <-r.Context().Done():
	}
}

// carriesEvent reports whether a polling payload, packets separated by the
// record separator, carries a socket.io event.
func carriesEvent(body []byte) bool {
	for packet := range strings.SplitSeq(string(body), "\x1e") {
		if strings.HasPrefix(packet, "42") {
			return true
		}
	}

	return false
}

func (s *droppingPollingServer) drop() {
	s.dropOnce.Do(func() { close(s.dropped) })
}

func (s *droppingPollingServer) release() {
	s.releaseOnce.Do(func() { close(s.released) })
}

// answerDropped answers a request of the session after the drop.
func (s *droppingPollingServer) answerDropped(w http.ResponseWriter, r *http.Request) {
	if s.mode == dropSilent {
		select {
		case <-s.released:
		case <-r.Context().Done():
		}

		return
	}

	unknownSession(w)
}

// droppingWebSocketServer upgrades to WebSocket, logs the client in, sends the
// lists New waits for and then leaves every command unanswered until drop
// closes the WebSocket connection.
type droppingWebSocketServer struct {
	dropped  chan struct{}
	dropOnce sync.Once
}

func newDroppingWebSocketServer() *droppingWebSocketServer {
	return &droppingWebSocketServer{dropped: make(chan struct{})}
}

func (s *droppingWebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")

	switch {
	case r.URL.Path == "/api/auth/ok":
		http.NotFound(w, r)

	case r.URL.Query().Get("transport") == "websocket":
		s.serveWebSocket(w, r)

	case r.Method == http.MethodGet && sid == "":
		data, err := json.Marshal(map[string]any{
			"sid":          "ws-drop-sid",
			"upgrades":     []string{"websocket"},
			"pingInterval": 25000,
			"pingTimeout":  20000,
			"maxPayload":   1000000,
		})
		if err != nil {
			http.Error(w, "marshal handshake", http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(append([]byte("0"), data...))

	case r.Method == http.MethodPost:
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)

	default:
		<-r.Context().Done()
	}
}

func (s *droppingWebSocketServer) drop() {
	s.dropOnce.Do(func() { close(s.dropped) })
}

func (*droppingWebSocketServer) release() {}

func (s *droppingWebSocketServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsHandshake(w, r)
	if err != nil {
		return
	}

	defer conn.Close()

	// The probe, its answer, the upgrade packet and the socket.io CONNECT.
	for _, step := range []string{"read", "3probe", "read", "read"} {
		if step == "read" {
			_, err = wsReadText(conn)
		} else {
			err = wsSendText(conn, step)
		}

		if err != nil {
			return
		}
	}

	for _, msg := range []string{`40{"sid":"ws-drop-socket"}`, `42["loginRequired"]`} {
		if wsSendText(conn, msg) != nil {
			return
		}
	}

	login, err := wsReadText(conn)
	if err != nil || !strings.HasPrefix(login, "42") {
		return
	}

	ackID := strings.TrimLeft(login[2:], "0123456789")
	ackID = login[2 : len(login)-len(ackID)]

	for _, msg := range []string{
		fmt.Sprintf(`43%s[{"ok":true,"msg":"Logged in successfully.","token":"fake-session-token"}]`, ackID),
		`42["monitorList",{}]`,
		`42["notificationList",[]]`,
		`42["statusPageList",{}]`,
	} {
		if wsSendText(conn, msg) != nil {
			return
		}
	}

	go func() {
		for {
			_, readErr := wsReadText(conn)
			if readErr != nil {
				return
			}
		}
	}()

	select {
	case <-s.dropped:
	case <-r.Context().Done():
	}
}

// newDropServer returns a fake server that loses the connection the way mode
// says, and how long after drop the client can know.
func newDropServer(mode dropMode) (dropper, time.Duration) {
	switch mode {
	case dropWebSocketClosed:
		return newDroppingWebSocketServer(), 0

	case dropSilent:
		return newDroppingPollingServer(mode), silentDropDetection

	case dropSessionClosed:
		return newDroppingPollingServer(mode), 0

	default:
		panic(fmt.Sprintf("unknown drop mode %d", mode))
	}
}

// dropCommand is a command of the client the drop tests send.
type dropCommand func(ctx context.Context, client *kuma.Client) error

// dropCommands are the commands the drop tests send, one for each way a
// command waits: for its ack alone, and for its ack and an update event.
func dropCommands() map[string]dropCommand {
	return map[string]dropCommand{
		"ack": func(ctx context.Context, client *kuma.Client) error {
			_, err := client.GetTags(ctx)

			return err
		},
		"update event": func(ctx context.Context, client *kuma.Client) error {
			return client.PauseMonitor(ctx, 1)
		},
	}
}

// TestConnectionDrop covers the commands of a client whose connection is gone:
// the server shut down, restarted and forgot the session, or became
// unreachable. Each command has to fail once the client can know, rather than
// run into the operation timeout - a command in flight when the connection
// drops as well as one sent after it.
func TestConnectionDrop(t *testing.T) {
	modes := map[string]dropMode{
		"session closed by the server": dropSessionClosed,
		"server gone silent":           dropSilent,
		"websocket closed by server":   dropWebSocketClosed,
	}

	for name, mode := range modes {
		for commandName, command := range dropCommands() {
			t.Run(name+"/"+commandName, func(t *testing.T) {
				t.Run("command in flight", func(t *testing.T) {
					server, detection := newDropServer(mode)
					dropClient := newDropClient(t, server)

					result := startCommand(t, dropClient, command)

					// Long enough for the command to be on the wire before
					// the drop.
					time.Sleep(100 * time.Millisecond)

					dropped := time.Now()

					server.drop()

					requireDropDetected(t, result, dropped, detection)
				})

				t.Run("command after the drop", func(t *testing.T) {
					server, detection := newDropServer(mode)
					dropClient := newDropClient(t, server)

					server.drop()
					time.Sleep(detection + 100*time.Millisecond)

					requireDropDetected(t, startCommand(t, dropClient, command), time.Now(), 0)
				})
			})
		}
	}
}

// TestCommandSendHeld covers a command whose send the server never answers,
// while the connection otherwise stays up: over long polling the send is an
// HTTP request bounded by nothing but the connection. The command has to give
// up with the caller's context, and fail once the connection drops.
func TestCommandSendHeld(t *testing.T) {
	const deadline = 300 * time.Millisecond

	for commandName, command := range dropCommands() {
		t.Run(commandName, func(t *testing.T) {
			t.Run("context done", func(t *testing.T) {
				server := newDroppingPollingServer(dropSessionClosed)
				dropClient := newDropClient(t, server)

				server.holdCommands.Store(true)

				ctx, cancel := context.WithTimeout(t.Context(), deadline)
				t.Cleanup(cancel)

				result := startCommandCtx(ctx, dropClient, command)

				hang := time.NewTimer(deadline + dropDetectionBudget)
				defer hang.Stop()

				select {
				case err := <-result:
					require.ErrorIs(t, err, context.DeadlineExceeded)

				case <-hang.C:
					require.FailNow(t, "a held send must not outlast the caller's context")
				}

				require.Positive(t, server.heldCommands.Load(), "the command must have been on the wire")
			})

			t.Run("connection lost", func(t *testing.T) {
				server := newDroppingPollingServer(dropSessionClosed)
				dropClient := newDropClient(t, server)

				server.holdCommands.Store(true)

				result := startCommand(t, dropClient, command)

				require.Eventually(t, func() bool {
					return server.heldCommands.Load() > 0
				}, time.Second, 10*time.Millisecond, "the command must have been on the wire")

				dropped := time.Now()

				server.drop()

				requireDropDetected(t, result, dropped, 0)
			})
		})
	}
}

// startCommand sends command on client and returns where its error arrives.
// The command runs on a goroutine of its own, so one that never returns fails
// the test instead of hanging it.
func startCommand(t *testing.T, client *kuma.Client, command dropCommand) <-chan error {
	t.Helper()

	return startCommandCtx(t.Context(), client, command)
}

// startCommandCtx is startCommand with ctx.
func startCommandCtx(ctx context.Context, client *kuma.Client, command dropCommand) <-chan error {
	result := make(chan error, 1)

	go func() {
		result <- command(ctx, client)
	}()

	return result
}

// newDropClient connects a client to server with dropOperationTimeout.
func newDropClient(t *testing.T, server dropper) *kuma.Client {
	t.Helper()

	httpServer := httptest.NewServer(server)
	t.Cleanup(func() {
		server.release()
		httpServer.CloseClientConnections()
		httpServer.Close()
	})

	dropClient, err := kuma.New(
		t.Context(),
		httpServer.URL,
		"admin", "admin1",
		kuma.WithConnectTimeout(3*time.Second),
		kuma.WithReadyGracePeriod(0),
		kuma.WithOperationTimeout(dropOperationTimeout),
		kuma.WithLogLevel(kuma.LogLevel(os.Getenv("SOCKETIO_LOG_LEVEL"))),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = dropClient.Disconnect() })

	return dropClient
}

// requireDropDetected asserts that the command whose error result delivers
// failed because of the drop, within dropDetectionBudget of when the client
// could know: detection after since.
func requireDropDetected(t *testing.T, result <-chan error, since time.Time, detection time.Duration) {
	t.Helper()

	// Past the operation timeout, so a command that waits for it is told
	// apart from one that hangs.
	hang := time.NewTimer(dropOperationTimeout + time.Second)
	defer hang.Stop()

	var err error

	select {
	case err = <-result:

	case <-hang.C:
		require.FailNow(t, "the command never returned, not even on the operation timeout")
	}

	elapsed := time.Since(since)

	t.Logf("command returned after %s: %v", elapsed.Round(time.Millisecond), err)

	require.ErrorIs(t, err, kuma.ErrConnectionLost)
	require.NotErrorIs(t, err, kuma.ErrOperationTimeout,
		"the drop must fail the command, not the operation timeout")
	require.Less(t, elapsed, detection+dropDetectionBudget,
		"the command must fail once the client can know the connection is gone")
}

// TestResyncAfterConnectionLost covers the way back from a lost connection:
// the server restarted and forgot the session, the next command reports that,
// and Resync connects again, authenticated the way the server line asks for.
func TestResyncAfterConnectionLost(t *testing.T) {
	tests := map[string]struct {
		newFake func(t *testing.T) (*fakeAckServer, string)
		// username and password are what New is given.
		username, password string
		// wantToken is whether the new connection has to be logged in with
		// the session token, which only a 2.x server asks for.
		wantToken bool
	}{
		"3.x, session cookie": {
			newFake:  newHTTPAuthFake,
			username: "admin", password: "admin1",
		},
		"2.x, session token": {
			newFake:  newAuthFake,
			username: "admin", password: "admin1",
			wantToken: true,
		},
		"2.x, authentication disabled": {
			newFake: func(t *testing.T) (*fakeAckServer, string) {
				t.Helper()

				fake, url := newAuthFake(t)
				fake.setAutoLogin()

				return fake, url
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fake, url := test.newFake(t)
			resyncClient := newLostClient(t, url, test.username, test.password)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			t.Cleanup(cancel)

			fake.forgetSessions()
			requireConnectionLost(ctx, t, resyncClient)

			tokensBefore := len(fake.receivedTokens())

			require.NoError(t, resyncClient.Resync(ctx))

			if test.wantToken {
				tokens := fake.receivedTokens()[tokensBefore:]
				require.Equal(t, []string{fakeSessionToken}, tokens,
					"the new connection has to be logged in with the session token")
			}

			_, err := resyncClient.GetNotification(ctx, fakeSeededNotificationID)
			require.NoError(t, err, "the resync has to rebuild the cache on the new connection")

			msg, err := resyncClient.TestNotification(ctx, notification.Generic{
				Base:           notification.Base{Name: "after the drop"},
				GenericDetails: notification.GenericDetails{},
				TypeName:       "generic",
			})
			require.NoError(t, err, "commands have to go out on the new connection")
			require.Equal(t, "ok", msg)
		})
	}
}

// TestResyncAfterConnectionLostRejectedToken covers a 2.x server that no
// longer accepts the session token on the new connection: the resync reports
// that, and the connection stays lost.
func TestResyncAfterConnectionLostRejectedToken(t *testing.T) {
	fake, url := newAuthFake(t)
	resyncClient := newLostClient(t, url, "admin", "admin1")

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	fake.forgetSessions()
	requireConnectionLost(ctx, t, resyncClient)

	fake.setRejectToken()

	require.ErrorIs(t, resyncClient.Resync(ctx), kuma.ErrInvalidSessionToken)

	_, err := resyncClient.GetTags(ctx)
	require.ErrorIs(t, err, kuma.ErrConnectionLost, "a refused resync must not install the new connection")
}

// newLostClient connects a client to the fake server at url.
func newLostClient(t *testing.T, url string, username string, password string) *kuma.Client {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	lostClient, err := kuma.New(ctx, url, username, password, kuma.WithConnectTimeout(5*time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { _ = lostClient.Disconnect() })

	return lostClient
}

// requireConnectionLost waits for client to report its connection as lost.
func requireConnectionLost(ctx context.Context, t *testing.T, client *kuma.Client) {
	t.Helper()

	require.Eventually(t, func() bool {
		_, err := client.GetTags(ctx)

		return errors.Is(err, kuma.ErrConnectionLost)
	}, 2*time.Second, 50*time.Millisecond, "the forgotten session must be reported as a lost connection")
}
