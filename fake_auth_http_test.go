package kuma_test

import (
	"fmt"
	"net/http"
)

// This file holds the engine.io sessions of fakeAckServer: one per socket.io
// connection the client makes, so that a client that connects again is served
// a session of its own.

// fakeSession is one engine.io session of the fake: the long polls of one
// socket.io connection are served from its channel until it is replaced.
type fakeSession struct {
	sid      string
	messages chan []byte
}

// sendTo queues a frame for one session.
func (*fakeAckServer) sendTo(session *fakeSession, msg []byte) {
	session.messages <- msg
}

// newSession starts the engine.io session a handshake asks for and makes it
// the current one. Earlier sessions stay served, the way a real server keeps
// every connection a client holds open: which one the client reads from is the
// client's business.
func (s *fakeAckServer) newSession() *fakeSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sidCounter++

	session := &fakeSession{
		sid: fmt.Sprintf("ack-test-sid-%d", s.sidCounter),
	}

	if s.current == nil {
		// The first session adopts the channel the test created the fake
		// with, so that frames a test queued before connecting are seen.
		session.messages = s.messages
	} else {
		session.messages = make(chan []byte, cap(s.messages))
	}

	if s.sessions == nil {
		s.sessions = map[string]*fakeSession{}
	}

	s.sessions[session.sid] = session
	s.current = session
	s.messages = session.messages

	return session
}

func (s *fakeAckServer) session(sid string) *fakeSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.sessions[sid]
}

// unknownSession answers a poll for a session the fake does not have, the way
// engine.io does.
func unknownSession(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"code":1,"message":"Session ID unknown"}`))
}
