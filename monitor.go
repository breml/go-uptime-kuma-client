package kuma

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

// GetMonitors retrieves all monitors for the authenticated user.
func (c *Client) GetMonitors(ctx context.Context) ([]monitor.Base, error) {
	_, err := c.syncEmit(ctx, "getMonitorList")
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

	response, err := c.syncEmitWithUpdateEvent(ctx, "add", "updateMonitorIntoList", monitorData)
	if err != nil {
		return 0, fmt.Errorf("create monitor: %w", err)
	}

	// Get monitorID from the response.
	if response.MonitorID == 0 {
		return 0, errors.New("create monitor: no monitor ID in response")
	}

	return response.MonitorID, nil
}

// UpdateMonitor updates an existing monitor.
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

	_, err = c.syncEmitWithUpdateEvent(ctx, "editMonitor", "updateMonitorIntoList", monitorData)
	if err != nil {
		return fmt.Errorf("update monitor %d: %w", mon.GetID(), err)
	}

	return nil
}

// DeleteMonitor deletes a monitor by ID.
func (c *Client) DeleteMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteMonitor", "deleteMonitorFromList", monitorID)
	if err != nil {
		return fmt.Errorf("delete monitor %d: %w", monitorID, err)
	}

	return nil
}

// PauseMonitor pauses a monitor by ID.
func (c *Client) PauseMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "pauseMonitor", "updateMonitorIntoList", monitorID)
	if err != nil {
		return fmt.Errorf("pause monitor %d: %w", monitorID, err)
	}

	return nil
}

// ResumeMonitor resumes a monitor by ID.
func (c *Client) ResumeMonitor(ctx context.Context, monitorID int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "resumeMonitor", "updateMonitorIntoList", monitorID)
	if err != nil {
		return fmt.Errorf("resume monitor %d: %w", monitorID, err)
	}

	return nil
}
