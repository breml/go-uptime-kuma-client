package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/maintenance"
)

// GetMaintenances retrieves all maintenance windows from the client cache.
func (c *Client) GetMaintenances(ctx context.Context) ([]maintenance.Maintenance, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state.maintenances == nil {
		return []maintenance.Maintenance{}, nil
	}

	// Return a copy to prevent external modifications
	maintenances := make([]maintenance.Maintenance, len(c.state.maintenances))
	copy(maintenances, c.state.maintenances)

	return maintenances, nil
}

// GetMaintenance retrieves a specific maintenance window by ID.
func (c *Client) GetMaintenance(ctx context.Context, id int64) (*maintenance.Maintenance, error) {
	response, err := c.syncEmit(ctx, "getMaintenance", id)
	if err != nil {
		return nil, fmt.Errorf("get maintenance %d: %w", id, err)
	}

	if response.Maintenance == nil {
		return nil, fmt.Errorf("get maintenance %d: maintenance not found in response", id)
	}

	var m maintenance.Maintenance
	err = convertToStruct(response.Maintenance, &m)
	if err != nil {
		return nil, fmt.Errorf("get maintenance %d: %w", id, err)
	}

	return &m, nil
}

// CreateMaintenance creates a new maintenance window.
func (c *Client) CreateMaintenance(ctx context.Context, m *maintenance.Maintenance) (*maintenance.Maintenance, error) {
	maintenanceData, err := structToMap(m)
	if err != nil {
		return nil, fmt.Errorf("create maintenance: %w", err)
	}

	response, err := c.syncEmitWithUpdateEvent(ctx, "addMaintenance", "maintenanceList", maintenanceData)
	if err != nil {
		return nil, fmt.Errorf("create maintenance: %w", err)
	}

	// The server returns the maintenance ID in the response
	m.ID = response.MaintenanceID

	// Fetch the complete maintenance object from the server
	return c.GetMaintenance(ctx, m.ID)
}

// UpdateMaintenance updates an existing maintenance window.
func (c *Client) UpdateMaintenance(ctx context.Context, m *maintenance.Maintenance) error {
	maintenanceData, err := structToMap(m)
	if err != nil {
		return fmt.Errorf("update maintenance: %w", err)
	}

	_, err = c.syncEmitWithUpdateEvent(ctx, "editMaintenance", "maintenanceList", maintenanceData)
	if err != nil {
		return fmt.Errorf("update maintenance: %w", err)
	}

	return nil
}

// DeleteMaintenance deletes a maintenance window by ID.
func (c *Client) DeleteMaintenance(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "deleteMaintenance", "maintenanceList", id)
	if err != nil {
		return fmt.Errorf("delete maintenance: %w", err)
	}

	return nil
}

// PauseMaintenance pauses (deactivates) a maintenance window.
func (c *Client) PauseMaintenance(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "pauseMaintenance", "maintenanceList", id)
	if err != nil {
		return fmt.Errorf("pause maintenance: %w", err)
	}

	return nil
}

// ResumeMaintenance resumes (activates) a maintenance window.
func (c *Client) ResumeMaintenance(ctx context.Context, id int64) error {
	_, err := c.syncEmitWithUpdateEvent(ctx, "resumeMaintenance", "maintenanceList", id)
	if err != nil {
		return fmt.Errorf("resume maintenance: %w", err)
	}

	return nil
}

// SetMonitorMaintenance sets the monitors associated with a maintenance window.
// This replaces all existing monitor associations.
func (c *Client) SetMonitorMaintenance(ctx context.Context, maintenanceID int64, monitorIDs []int64) error {
	// Convert []int64 to []map[string]any format expected by server
	monitors := make([]map[string]any, len(monitorIDs))
	for i, id := range monitorIDs {
		monitors[i] = map[string]any{"id": id}
	}

	_, err := c.syncEmit(ctx, "addMonitorMaintenance", maintenanceID, monitors)
	if err != nil {
		return fmt.Errorf("set monitor maintenance: %w", err)
	}

	return nil
}

// GetMonitorMaintenance retrieves the monitors associated with a maintenance window.
func (c *Client) GetMonitorMaintenance(ctx context.Context, maintenanceID int64) ([]int64, error) {
	response, err := c.syncEmit(ctx, "getMonitorMaintenance", maintenanceID)
	if err != nil {
		return nil, fmt.Errorf("get monitor maintenance: %w", err)
	}

	if response.Monitors == nil {
		return []int64{}, nil
	}

	monitorIDs := make([]int64, 0, len(response.Monitors))
	for _, monitorItem := range response.Monitors {
		monitorMap, ok := monitorItem.(map[string]interface{})
		if !ok {
			continue
		}

		if id, idOk := monitorMap["id"].(float64); idOk {
			monitorIDs = append(monitorIDs, int64(id))
		}
	}

	return monitorIDs, nil
}

// SetMaintenanceStatusPage sets the status pages associated with a maintenance window.
// This replaces all existing status page associations.
func (c *Client) SetMaintenanceStatusPage(ctx context.Context, maintenanceID int64, statusPageIDs []int64) error {
	// Convert []int64 to []map[string]any format expected by server
	statusPages := make([]map[string]any, len(statusPageIDs))
	for i, id := range statusPageIDs {
		statusPages[i] = map[string]any{"id": id}
	}

	_, err := c.syncEmit(ctx, "addMaintenanceStatusPage", maintenanceID, statusPages)
	if err != nil {
		return fmt.Errorf("set maintenance status page: %w", err)
	}

	return nil
}

// GetMaintenanceStatusPage retrieves the status pages associated with a maintenance window.
func (c *Client) GetMaintenanceStatusPage(ctx context.Context, maintenanceID int64) ([]int64, error) {
	response, err := c.syncEmit(ctx, "getMaintenanceStatusPage", maintenanceID)
	if err != nil {
		return nil, fmt.Errorf("get maintenance status page: %w", err)
	}

	if response.StatusPages == nil {
		return []int64{}, nil
	}

	statusPageIDs := make([]int64, 0, len(response.StatusPages))
	for _, statusPageItem := range response.StatusPages {
		statusPageMap, ok := statusPageItem.(map[string]interface{})
		if !ok {
			continue
		}

		if id, idOk := statusPageMap["id"].(float64); idOk {
			statusPageIDs = append(statusPageIDs, int64(id))
		}
	}

	return statusPageIDs, nil
}
