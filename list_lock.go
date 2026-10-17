package kuma

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

// wholeListUpdateEvents names the update events whose payload is a complete
// list that replaces the cached one.
//
// Such a broadcast is only as fresh as the moment the server read the list, and
// Uptime Kuma reads it after the write it belongs to and emits it afterwards,
// over a connection that preserves that order. Two writes in flight at once
// therefore interleave into a read-without, emit-after order, and the older
// snapshot lands last and drops the newer write from the cache again. Writes to
// these lists are serialized so that no two of them can interleave that way.
//
// GetMonitors takes the monitorList lock too: its getMonitorList request is
// answered with the same whole-list broadcast, which goes stale the same way.
//
// updateMonitorIntoList and deleteMonitorFromList are absent on purpose: they
// carry a single monitor and are merged into the cache rather than replacing
// it, so they cannot go stale this way and stay concurrent with each other,
// see sharedListUpdateEvents.
func wholeListUpdateEvents() []string {
	return []string{
		"notificationList",
		"proxyList",
		"dockerHostList",
		"maintenanceList",
		"monitorList",
		"statusPageList",
	}
}

// sharedListUpdateEvents maps the update events that merge a single entry into
// a cached list to the whole-list event of that list.
//
// A merge cannot go stale itself, but a whole-list broadcast read before it and
// delivered after it replaces the cache and drops it again: a monitor created
// while a GetMonitors is in flight disappears from the cache once the older
// list arrives. Commands on these events therefore take the shared side of the
// whole-list lock, which keeps them concurrent with each other and excludes
// them from a whole-list command in flight.
func sharedListUpdateEvents() map[string]string {
	return map[string]string{
		"updateMonitorIntoList": "monitorList",
		"deleteMonitorFromList": "monitorList",
	}
}

// exclusiveListWeight is the weight a whole-list command acquires, and the size
// of every list lock. Each single-entry command acquires a weight of one, so
// the two exclude each other while single-entry commands share the lock.
const exclusiveListWeight = 1 << 30

// lockConn takes the shared side of the connection lock and returns the
// release. Every command holds it from the moment it is emitted until its ack,
// and its update event where it waits for one, is in, so that the connection
// is never replaced under a command in flight, see lockConnExclusive.
//
// Like lockList it is a weighted semaphore rather than a sync.RWMutex: a
// caller waiting its turn has to be able to give up with its context. Waiters
// are served in order, so a resync waiting for the commands in flight holds
// back those that arrive after it and cannot be starved by them.
func (c *Client) lockConn(ctx context.Context) (func(), error) {
	err := c.connLock.Acquire(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("waiting for the resync in progress: %w", contextErr(ctx))
	}

	return func() { c.connLock.Release(1) }, nil
}

// lockConnExclusive takes the exclusive side of the connection lock, which
// excludes every command, and returns the release. It is what a resync that
// reconnects holds while it replaces the connection and until the lists the
// new one answers with have rebuilt the cache, so that no command is emitted
// on a connection about to be closed and none observes the cache half built.
//
// Lock order: a list lock, see lockList, is taken before the connection lock
// and the connection lock is never held while another command is emitted.
func (c *Client) lockConnExclusive(ctx context.Context) (func(), error) {
	err := c.connLock.Acquire(ctx, exclusiveListWeight)
	if err != nil {
		return nil, fmt.Errorf("waiting for the commands in flight: %w", contextErr(ctx))
	}

	return func() { c.connLock.Release(exclusiveListWeight) }, nil
}

// newListLocks builds the lock per list-replacing update event. A weighted
// semaphore rather than a sync.RWMutex, because a caller waiting its turn has
// to be able to give up, see lockList.
func newListLocks() map[string]*semaphore.Weighted {
	locks := make(map[string]*semaphore.Weighted, len(wholeListUpdateEvents()))
	for _, event := range wholeListUpdateEvents() {
		locks[event] = semaphore.NewWeighted(exclusiveListWeight)
	}

	return locks
}

// lockList takes the lock updateEvent belongs to and returns the release: the
// exclusive side for a whole-list event, the shared side for a single-entry
// event of sharedListUpdateEvents. It returns a nil release for any other
// event, which needs no ordering, and an error for a caller whose context was
// done before its turn came.
//
// Waiters are served in order, so a whole-list command waiting for the
// single-entry commands in flight holds back those that arrive after it, and a
// steady stream of creates cannot starve a GetMonitors.
//
// The wait is on ctx and not on the operation budget: the budget bounds a
// command's own round trip and is taken once the lock is won, see
// syncEmitWithConfirmedUpdateEvent, so what bounds the queue is whatever bounds
// the caller. Giving up here is the only way out of a queue behind commands
// that are each waiting out a server that stopped answering.
func (c *Client) lockList(ctx context.Context, updateEvent string) (func(), error) {
	listEvent, weight := updateEvent, int64(exclusiveListWeight)
	if shared, ok := sharedListUpdateEvents()[updateEvent]; ok {
		listEvent, weight = shared, 1
	}

	lock, ok := c.listLocks[listEvent]
	if !ok {
		return nil, nil //nolint:nilnil // no lock to take is not a failure.
	}

	err := lock.Acquire(ctx, weight)
	if err != nil {
		return nil, fmt.Errorf("waiting for the other commands on %s: %w", listEvent, contextErr(ctx))
	}

	return func() { lock.Release(weight) }, nil
}
