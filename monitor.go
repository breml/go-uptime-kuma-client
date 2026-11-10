package kuma

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

// GetMonitors retrieves all monitors for the authenticated user.
//
// The server answers getMonitorList with a monitorList broadcast that replaces
// the cached list. The client can see the ack before it has applied that
// broadcast, so GetMonitors waits for the broadcast itself. It holds the
// exclusive side of the monitorList lock meanwhile, so that no other
// GetMonitors and no single-monitor write of this client interleaves with it,
// see sharedListUpdateEvents.
//
// An error wrapping ErrUpdateEventTimeout means the list never arrived; no
// monitors are returned with it, the cache being the stale list.
//
// monitorList broadcasts it did not ask for, from Resync, from another session
// of the same user or from addProxy with applyExisting, are not ordered with it.
//
// The list identifies monitors, it does not describe them. An Uptime Kuma 3.x
// server sends it in a compact form that carries only ID, Name, Type, IsActive,
// Parent, PathName, Interval and Tags; every other field of the returned
// monitor.Base values, and of anything decoded from them with As, is the zero
// value there. (The compact form also carries the children, the path and the
// weight of a monitor, which Base has no field for.) A 2.x server sends the
// full monitors, and so does the updateMonitorIntoList a write is answered
// with on either line, so the list can hold both shapes until the next
// broadcast replaces it. GetMonitor and GetMonitorAs return the full monitor
// on both lines; read the details there.
func (c *Client) GetMonitors(ctx context.Context) ([]monitor.Base, error) {
	_, err := c.syncEmitWithUpdateEvent(ctx, "getMonitorList", "monitorList")
	if err != nil {
		return nil, fmt.Errorf("get monitors: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	monitors := make([]monitor.Base, len(c.state.monitors))
	copy(monitors, c.state.monitors)

	return monitors, nil
}

// GetMonitor retrieves a specific monitor by ID.
func (c *Client) GetMonitor(ctx context.Context, monitorID int64) (monitor.Base, error) {
	response, err := c.syncEmit(ctx, "getMonitor", monitorID)
	if err != nil {
		return monitor.Base{}, fmt.Errorf("get monitor %d: %w", monitorID, err)
	}

	if response.Monitor == nil {
		return monitor.Base{}, fmt.Errorf("get monitor %d: monitor not found in response", monitorID)
	}

	// Convert the monitor data to a monitor.Base.
	var mon monitor.Base
	err = convertToStruct(response.Monitor, &mon)
	if err != nil {
		return monitor.Base{}, fmt.Errorf("get monitor %d: %w", monitorID, err)
	}

	return mon, nil
}

// GetMonitorAs retrieves a specific monitor by ID and converts it to the target type.
func (c *Client) GetMonitorAs(ctx context.Context, monitorID int64, target any) error {
	mon, err := c.GetMonitor(ctx, monitorID)
	if err != nil {
		return err
	}

	err = mon.As(target)
	if err != nil {
		return fmt.Errorf("get monitor %d as %T: %w", monitorID, target, err)
	}

	return nil
}

// CreateMonitor creates a new monitor.
//
// An error wrapping ErrUpdateEventTimeout means the monitor was created and only
// the update event is missing; the returned ID identifies it, as does the ID an
// UpdateEventTimeoutError carries, so retrying the call would create a
// duplicate.
func (c *Client) CreateMonitor(ctx context.Context, mon monitor.Monitor) (int64, error) {
	monitorData, err := structToMap(mon)
	if err != nil {
		return 0, fmt.Errorf("create monitor: %w", err)
	}

	// The server expects notificationIDList as a map[string]bool.
	notificationIDList := map[string]bool{}
	for _, id := range mon.GetNotificationIDs() {
		notificationIDList[strconv.FormatInt(id, 10)] = true
	}

	monitorData["notificationIDList"] = notificationIDList
	c.adaptMonitorPayload(monitorData)

	response, err := c.syncEmitWithConfirmedUpdateEvent(
		ctx, "add", "updateMonitorIntoList",
		func(response ackResponse) bool { return c.hasMonitor(response.MonitorID) },
		monitorData,
	)
	if err != nil {
		if errors.Is(err, ErrUpdateEventTimeout) {
			return response.MonitorID, fmt.Errorf("create monitor: %w", withCreatedID(err, response.MonitorID))
		}

		return 0, fmt.Errorf("create monitor: %w", err)
	}

	// Get monitorID from the response.
	if response.MonitorID == 0 {
		return 0, errors.New("create monitor: no monitor ID in response")
	}

	return response.MonitorID, nil
}

// UpdateMonitor updates an existing monitor.
//
// An error wrapping ErrUpdateEventTimeout means the monitor was updated and only
// the update event is missing.
//
// Unlike CreateMonitor and DeleteMonitor, this and the two pause/resume calls
// wait for the broadcast by name alone, because they change no property the
// cache can be checked against: the monitor is in the list before and after
// them. With several writes in flight the broadcast one of them observes may
// therefore be a neighbour's, and the cache can hold the previous values until
// the next broadcast arrives. GetMonitor is unaffected, it asks the server.
func (c *Client) UpdateMonitor(ctx context.Context, mon monitor.Monitor) error {
	monitorData, err := structToMap(mon)
	if err != nil {
		return fmt.Errorf("update monitor: %w", err)
	}

	// The server expects notificationIDList as a map[string]bool.
	notificationIDList := map[string]bool{}
	for _, id := range mon.GetNotificationIDs() {
		notificationIDList[strconv.FormatInt(id, 10)] = true
	}

	monitorData["notificationIDList"] = notificationIDList
	c.adaptMonitorPayload(monitorData)

	_, err = c.syncEmitWithUpdateEvent(ctx, "editMonitor", "updateMonitorIntoList", monitorData)
	if err != nil {
		return fmt.Errorf("update monitor %d: %w", mon.GetID(), err)
	}

	return nil
}

// adaptMonitorPayload rewrites a create or edit payload into the shape the
// server line expects.
//
// Uptime Kuma 3.x keeps the name of a pm2 or system-service monitor in the
// monitor's config, see AuthSchemeHTTP: it takes the name only as
// config.systemServiceName and has dropped the system_service_name column. A
// 2.x server has no config column and takes system_service_name. Both servers
// import an add payload into the frozen monitor table as is, so a payload
// carrying both fields fails on either line; the field is moved instead.
func (c *Client) adaptMonitorPayload(monitorData map[string]any) {
	if c.scheme != AuthSchemeHTTP {
		return
	}

	name, ok := monitorData["system_service_name"]
	if !ok {
		return
	}

	delete(monitorData, "system_service_name")

	if name == nil {
		return
	}

	config, ok := monitorData["config"].(map[string]any)
	if !ok {
		config = map[string]any{}
	}

	config["systemServiceName"] = name
	monitorData["config"] = config
}

// DeleteMonitor deletes a monitor by ID.
//
// An error wrapping ErrUpdateEventTimeout means the monitor was deleted and only
// the update event is missing.
func (c *Client) DeleteMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithConfirmedUpdateEvent(
		ctx, "deleteMonitor", "deleteMonitorFromList",
		func(ackResponse) bool { return !c.hasMonitor(monitorID) },
		monitorID,
	)
	if err != nil {
		return fmt.Errorf("delete monitor %d: %w", monitorID, err)
	}

	return nil
}

// PauseMonitor pauses a monitor by ID.
//
// An error wrapping ErrUpdateEventTimeout means the monitor was paused and only
// the update event is missing.
func (c *Client) PauseMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "pauseMonitor", "updateMonitorIntoList", monitorID)
	if err != nil {
		return fmt.Errorf("pause monitor %d: %w", monitorID, err)
	}

	return nil
}

// ResumeMonitor resumes a monitor by ID.
//
// An error wrapping ErrUpdateEventTimeout means the monitor was resumed and only
// the update event is missing.
func (c *Client) ResumeMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "resumeMonitor", "updateMonitorIntoList", monitorID)
	if err != nil {
		return fmt.Errorf("resume monitor %d: %w", monitorID, err)
	}

	return nil
}

// hasMonitor reports whether the state cache holds a monitor with the given ID.
func (c *Client) hasMonitor(id int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return cachedByID(c.state.monitors, id)
}
