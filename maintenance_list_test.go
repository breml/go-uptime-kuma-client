package kuma_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/maintenance"
)

// TestMaintenanceListKeepsDecodableEntries pins that a maintenance window the
// client cannot decode costs only that window, not the whole maintenanceList:
// the others reach the cache and GetMaintenances names the one left out.
func TestMaintenanceListKeepsDecodableEntries(t *testing.T) {
	t.Parallel()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	kumaClient := newFakeAckClient(t, fake)

	// Uptime Kuma accepts any date JavaScript's new Date(...) reads, such as an
	// offset without a colon, and stores it as sent.
	fake.send([]byte(`42["maintenanceList",{` +
		`"1":{"id":1,"title":"good","strategy":"manual","active":true,"dateRange":[]},` +
		`"2":{"id":2,"title":"bad","strategy":"single","active":true,` +
		`"dateRange":["2026-06-29T13:39:00+0200","2026-06-29T14:39:00+0200"]}}]`))

	var (
		maintenances []maintenance.Maintenance
		err          error
	)

	require.Eventually(t, func() bool {
		maintenances, err = kumaClient.GetMaintenances(t.Context())

		return len(maintenances) == 1
	}, 5*time.Second, 10*time.Millisecond, "the decodable window must reach the cache")

	require.ErrorIs(t, err, kuma.ErrMaintenanceDecode)
	require.ErrorContains(t, err, "maintenance 2")
	require.Equal(t, int64(1), maintenances[0].ID)

	fake.send([]byte(`42["maintenanceList",{` +
		`"1":{"id":1,"title":"good","strategy":"manual","active":true,"dateRange":[]},` +
		`"3":{"id":3,"title":"also good","strategy":"manual","active":true,"dateRange":[]}}]`))

	require.Eventually(t, func() bool {
		maintenances, err = kumaClient.GetMaintenances(t.Context())

		return len(maintenances) == 2
	}, 5*time.Second, 10*time.Millisecond, "the next list must replace the cache")

	require.NoError(t, err, "a list that decodes completely clears the earlier error")
}

// fakeMaintenanceAnswer is how the fake answers a maintenance command: with
// ack, and with list as the maintenanceList broadcast. Uptime Kuma broadcasts
// before it acknowledges; a positive listDelay acknowledges first and
// broadcasts that much later, the order in which the two can reach a client.
type fakeMaintenanceAnswer struct {
	ack       string
	list      string
	listDelay time.Duration
}

func (s *fakeAckServer) setMaintenanceAnswer(command string, answer fakeMaintenanceAnswer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.maintenanceAnswers == nil {
		s.maintenanceAnswers = map[string]fakeMaintenanceAnswer{}
	}

	s.maintenanceAnswers[command] = answer
}

func (s *fakeAckServer) answerMaintenanceWrite(session *fakeSession, ackID string, payload string) bool {
	var frame []json.RawMessage

	if json.Unmarshal([]byte(payload), &frame) != nil || len(frame) < 1 {
		return false
	}

	var command string

	if json.Unmarshal(frame[0], &command) != nil {
		return false
	}

	s.mu.Lock()
	answer, ok := s.maintenanceAnswers[command]
	s.mu.Unlock()

	if !ok {
		return false
	}

	ack := fmt.Appendf(nil, `43%s[%s]`, ackID, answer.ack)
	list := []byte(`42["maintenanceList",` + answer.list + `]`)

	if answer.listDelay <= 0 {
		s.sendTo(session, list)
		s.sendTo(session, ack)

		return true
	}

	s.sendTo(session, ack)

	go func() {
		time.Sleep(answer.listDelay)
		s.sendTo(session, list)
	}()

	return true
}

const (
	goodMaintenanceEntry = `"1":{"id":1,"title":"good","strategy":"manual","active":true,"dateRange":[]}`
	badMaintenanceDates  = `"strategy":"single","active":true,` +
		`"dateRange":["2026-06-29T13:39:00+0200","2026-06-29T14:39:00+0200"]`
)

// TestCreateMaintenanceReportsUndecodableBroadcast pins that a create whose
// window the server broadcasts in a form the client cannot decode returns the
// decode error with the created ID once the broadcast arrives, instead of
// waiting out the operation budget for a window that never reaches the cache.
func TestCreateMaintenanceReportsUndecodableBroadcast(t *testing.T) {
	t.Parallel()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	kumaClient := newFakeAckClient(t, fake, kuma.WithOperationTimeout(5*time.Second))

	fake.setMaintenanceAnswer("addMaintenance", fakeMaintenanceAnswer{
		ack: fmt.Sprintf(`{"ok":true,"msg":"ok","maintenanceID":%d}`, fakeAckCreatedID),
		list: fmt.Sprintf(`{%s,"%d":{"id":%d,"title":"new",%s}}`,
			goodMaintenanceEntry, fakeAckCreatedID, fakeAckCreatedID, badMaintenanceDates),
	})

	start := time.Now()

	created, err := kumaClient.CreateMaintenance(t.Context(), &maintenance.Maintenance{
		Title:    "new",
		Strategy: "manual",
	})

	require.ErrorIs(t, err, kuma.ErrMaintenanceDecode)
	require.NotErrorIs(t, err, kuma.ErrUpdateEventTimeout)
	require.Less(t, time.Since(start), 2*time.Second, "the create must not wait out the operation budget")
	require.NotNil(t, created)
	require.Equal(t, int64(fakeAckCreatedID), created.ID)
}

// TestDeleteMaintenanceWaitsForUndecodableWindowToGo pins that deleting a
// window the cache only knows as undecodable waits for the broadcast that
// drops it, so GetMaintenances no longer reports it once the delete returned.
func TestDeleteMaintenanceWaitsForUndecodableWindowToGo(t *testing.T) {
	t.Parallel()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	kumaClient := newFakeAckClient(t, fake)

	fake.send([]byte(`42["maintenanceList",{` + goodMaintenanceEntry +
		`,"2":{"id":2,"title":"bad",` + badMaintenanceDates + `}}]`))

	require.Eventually(t, func() bool {
		_, err := kumaClient.GetMaintenances(t.Context())

		return errors.Is(err, kuma.ErrMaintenanceDecode)
	}, 5*time.Second, 10*time.Millisecond, "the undecodable window must be reported")

	fake.setMaintenanceAnswer("deleteMaintenance", fakeMaintenanceAnswer{
		ack:       `{"ok":true,"msg":"ok"}`,
		list:      `{` + goodMaintenanceEntry + `}`,
		listDelay: 200 * time.Millisecond,
	})

	err := kumaClient.DeleteMaintenance(t.Context(), 2)
	require.NoError(t, err)

	maintenances, err := kumaClient.GetMaintenances(t.Context())
	require.NoError(t, err, "the deleted window must not be reported any more")
	require.Len(t, maintenances, 1)
}

// TestGetMaintenanceWrapsDecodeError pins that GetMaintenance on a window the
// client cannot decode wraps ErrMaintenanceDecode, as GetMaintenances does.
func TestGetMaintenanceWrapsDecodeError(t *testing.T) {
	t.Parallel()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	kumaClient := newFakeAckClient(t, fake)

	fake.setMaintenanceAnswer("getMaintenance", fakeMaintenanceAnswer{
		ack:  `{"ok":true,"maintenance":{"id":2,"title":"bad",` + badMaintenanceDates + `}}`,
		list: `{` + goodMaintenanceEntry + `}`,
	})

	_, err := kumaClient.GetMaintenance(t.Context(), 2)
	require.ErrorIs(t, err, kuma.ErrMaintenanceDecode)
}
