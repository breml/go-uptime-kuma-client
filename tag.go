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
		return nil, fmt.Errorf("get tags: %v", err)
	}

	// Convert to []tag.Tag
	var tags []tag.Tag
	err = convertToStruct(response.Tags, &tags)
	if err != nil {
		return nil, fmt.Errorf("get tags: %v", err)
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
		return 0, fmt.Errorf("create tag: %v", err)
	}

	response, err := c.syncEmit(ctx, "addTag", tagData)
	if err != nil {
		return 0, fmt.Errorf("create tag: %v", err)
	}

	// The response contains the created tag with ID.
	var createdTag tag.Tag
	err = convertToStruct(response.Tag, &createdTag)
	if err != nil {
		return 0, fmt.Errorf("create tag: %v", err)
	}

	return createdTag.ID, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(ctx context.Context, t tag.Tag) error {
	tagData, err := structToMap(t)
	if err != nil {
		return fmt.Errorf("update tag: %v", err)
	}

	_, err = c.syncEmit(ctx, "editTag", tagData)
	if err != nil {
		return fmt.Errorf("update tag %d: %v", t.ID, err)
	}

	return nil
}

// DeleteTag deletes a tag by ID.
func (c *Client) DeleteTag(ctx context.Context, tagID int64) error {
	_, err := c.syncEmit(ctx, "deleteTag", tagID)
	if err != nil {
		return fmt.Errorf("delete tag %d: %v", tagID, err)
	}

	return nil
}
