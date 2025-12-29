package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/tag"
)

// GetTags retrieves all tags.
func (c *Client) GetTags(ctx context.Context) ([]tag.Tag, error) {
	response, err := c.syncEmit(ctx, "getTags")
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}

	// Convert to []tag.Tag
	var tags []tag.Tag
	err = convertToStruct(response.Tags, &tags)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}

	return tags, nil
}

// GetTag retrieves a specific tag by ID.
func (c *Client) GetTag(ctx context.Context, tagID int64) (tag.Tag, error) {
	tags, err := c.GetTags(ctx)
	if err != nil {
		return tag.Tag{}, err
	}

	for _, t := range tags {
		if t.ID == tagID {
			return t, nil
		}
	}

	return tag.Tag{}, fmt.Errorf("get tag %d: %w", tagID, ErrNotFound)
}

// CreateTag creates a new tag.
func (c *Client) CreateTag(ctx context.Context, t tag.Tag) (int64, error) {
	tagData, err := structToMap(t)
	if err != nil {
		return 0, fmt.Errorf("create tag: %w", err)
	}

	response, err := c.syncEmit(ctx, "addTag", tagData)
	if err != nil {
		return 0, fmt.Errorf("create tag: %w", err)
	}

	// The response contains the created tag with ID.
	var createdTag tag.Tag
	err = convertToStruct(response.Tag, &createdTag)
	if err != nil {
		return 0, fmt.Errorf("create tag: %w", err)
	}

	return createdTag.ID, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(ctx context.Context, t tag.Tag) error {
	tagData, err := structToMap(t)
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}

	_, err = c.syncEmit(ctx, "editTag", tagData)
	if err != nil {
		return fmt.Errorf("update tag %d: %w", t.ID, err)
	}

	return nil
}

// DeleteTag deletes a tag by ID.
// This also removes all monitor-tag associations for this tag via cascade delete.
func (c *Client) DeleteTag(ctx context.Context, tagID int64) error {
	_, err := c.syncEmit(ctx, "deleteTag", tagID)
	if err != nil {
		return fmt.Errorf("delete tag %d: %w", tagID, err)
	}

	// Manually remove this tag from all monitors in the cache
	// The server removes all monitor-tag associations when a tag is deleted
	c.mu.Lock()
	for i := range c.state.monitors {
		filteredTags := make([]tag.MonitorTag, 0, len(c.state.monitors[i].Tags))
		for _, t := range c.state.monitors[i].Tags {
			if t.TagID != tagID {
				filteredTags = append(filteredTags, t)
			}
		}

		c.state.monitors[i].Tags = filteredTags
	}

	c.mu.Unlock()

	return nil
}

// refreshMonitorInCache fetches a monitor and updates it in the cache.
// This is used by tag operations since the server doesn't emit update events for them.
// If the monitor is not found in the cache, it will be added.
func (c *Client) refreshMonitorInCache(ctx context.Context, monitorID int64) error {
	mon, err := c.GetMonitor(ctx, monitorID)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Try to find and update existing monitor in cache
	for i, cachedMon := range c.state.monitors {
		if cachedMon.ID == monitorID {
			c.state.monitors[i] = mon
			return nil
		}
	}

	// Monitor not found in cache, add it
	c.state.monitors = append(c.state.monitors, mon)

	return nil
}

// AddMonitorTag adds a tag to a monitor with an optional value.
func (c *Client) AddMonitorTag(
	ctx context.Context,
	tagID int64,
	monitorID int64,
	value string,
) (*tag.MonitorTag, error) {
	response, err := c.syncEmit(ctx, "addMonitorTag", tagID, monitorID, value)
	if err != nil {
		return nil, fmt.Errorf("add monitor tag (tag %d, monitor %d): %w", tagID, monitorID, err)
	}

	if !response.OK {
		return nil, fmt.Errorf("add monitor tag (tag %d, monitor %d): %s", tagID, monitorID, response.Msg)
	}

	// Refresh the monitor in cache to get the updated tags
	err = c.refreshMonitorInCache(ctx, monitorID)
	if err != nil {
		return nil, fmt.Errorf(
			"add monitor tag (tag %d, monitor %d): failed to refresh cache: %w",
			tagID,
			monitorID,
			err,
		)
	}

	// Find the tag we just added from the cache
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, mon := range c.state.monitors {
		if mon.ID == monitorID {
			for _, t := range mon.Tags {
				if t.TagID == tagID && t.Value == value {
					return &t, nil
				}
			}

			break
		}
	}

	return nil, fmt.Errorf(
		"add monitor tag (tag %d, monitor %d): tag added but not found in monitor tags",
		tagID,
		monitorID,
	)
}

// UpdateMonitorTag updates the value of a monitor-tag association.
func (c *Client) UpdateMonitorTag(ctx context.Context, tagID int64, monitorID int64, value string) error {
	response, err := c.syncEmit(ctx, "editMonitorTag", tagID, monitorID, value)
	if err != nil {
		return fmt.Errorf("update monitor tag (tag %d, monitor %d): %w", tagID, monitorID, err)
	}

	if !response.OK {
		return fmt.Errorf("update monitor tag (tag %d, monitor %d): %s", tagID, monitorID, response.Msg)
	}

	// Refresh the monitor in cache to get the updated tags
	err = c.refreshMonitorInCache(ctx, monitorID)
	if err != nil {
		return fmt.Errorf("update monitor tag (tag %d, monitor %d): failed to refresh cache: %w", tagID, monitorID, err)
	}

	return nil
}

// DeleteMonitorTagWithValue removes a specific tag association from a monitor by value.
func (c *Client) DeleteMonitorTagWithValue(ctx context.Context, tagID int64, monitorID int64, value string) error {
	response, err := c.syncEmit(ctx, "deleteMonitorTag", tagID, monitorID, value)
	if err != nil {
		return fmt.Errorf(
			"delete monitor tag with value (tag %d, monitor %d, value %q): %w",
			tagID,
			monitorID,
			value,
			err,
		)
	}

	if !response.OK {
		return fmt.Errorf(
			"delete monitor tag with value (tag %d, monitor %d, value %q): %s",
			tagID,
			monitorID,
			value,
			response.Msg,
		)
	}

	// Refresh the monitor in cache to get the updated tags
	err = c.refreshMonitorInCache(ctx, monitorID)
	if err != nil {
		return fmt.Errorf(
			"delete monitor tag with value (tag %d, monitor %d, value %q): failed to refresh cache: %w",
			tagID,
			monitorID,
			value,
			err,
		)
	}

	return nil
}

// DeleteMonitorTag removes all associations of a tag from a monitor (all values).
func (c *Client) DeleteMonitorTag(ctx context.Context, tagID int64, monitorID int64) error {
	// Get all tags for the monitor
	tags, err := c.GetMonitorTags(ctx, monitorID)
	if err != nil {
		return fmt.Errorf("delete monitor tag (tag %d, monitor %d): failed to fetch tags: %w", tagID, monitorID, err)
	}

	// Delete all associations with the given tagID
	for _, t := range tags {
		if t.TagID == tagID {
			err := c.DeleteMonitorTagWithValue(ctx, tagID, monitorID, t.Value)
			if err != nil {
				return fmt.Errorf("delete monitor tag (tag %d, monitor %d): %w", tagID, monitorID, err)
			}
		}
	}

	return nil
}

// GetMonitorTags retrieves all tags for a specific monitor.
// This method uses the client's local state cache, which is kept synchronized
// via socket.io events. The cache is automatically updated by monitorList,
// updateMonitorIntoList, and deleteMonitorFromList events.
func (c *Client) GetMonitorTags(ctx context.Context, monitorID int64) ([]tag.MonitorTag, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, mon := range c.state.monitors {
		if mon.ID == monitorID {
			return mon.Tags, nil
		}
	}

	return nil, fmt.Errorf("get monitor tags (monitor %d): %w", monitorID, ErrNotFound)
}

// GetTagMonitors retrieves all monitor IDs associated with a specific tag.
// This method uses the client's local state cache, which is kept synchronized
// via socket.io events. The cache is automatically updated by monitorList,
// updateMonitorIntoList, and deleteMonitorFromList events.
func (c *Client) GetTagMonitors(ctx context.Context, tagID int64) ([]int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var monitorIDs []int64
	for _, mon := range c.state.monitors {
		for _, t := range mon.Tags {
			if t.TagID == tagID {
				monitorIDs = append(monitorIDs, mon.ID)
				break
			}
		}
	}

	return monitorIDs, nil
}
