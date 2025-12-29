package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/tag"
)

func TestClient_TagCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test tag
	testTag := tag.Tag{
		Name:  "Test Tag",
		Color: "#FF0000",
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetTags (should work even if empty)
		tags, err := client.GetTags(ctx)
		require.NoError(t, err)
		initialCount = len(tags)
	})

	var tagID int64
	var tagRetrieved tag.Tag
	t.Run("create", func(t *testing.T) {
		// Test CreateTag
		tagID, err = client.CreateTag(ctx, testTag)
		require.NoError(t, err)
		require.Positive(t, tagID)

		// Test GetTags after creation
		tags, err := client.GetTags(ctx)
		require.NoError(t, err)
		require.Len(t, tags, initialCount+1)

		// Test GetTag
		tagRetrieved, err = client.GetTag(ctx, tagID)
		require.NoError(t, err)
		require.Equal(t, tagID, tagRetrieved.ID)
		require.Equal(t, "Test Tag", tagRetrieved.Name)
		require.Equal(t, "#FF0000", tagRetrieved.Color)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateTag
		tagRetrieved.Name = "Updated Tag"
		tagRetrieved.Color = "#00FF00"
		err := client.UpdateTag(ctx, tagRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedTag, err := client.GetTag(ctx, tagID)
		require.NoError(t, err)
		require.Equal(t, "Updated Tag", updatedTag.Name)
		require.Equal(t, "#00FF00", updatedTag.Color)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteTag
		err := client.DeleteTag(ctx, tagID)
		require.NoError(t, err)

		// Verify deletion
		tags, err := client.GetTags(ctx)
		require.NoError(t, err)
		require.Len(t, tags, initialCount)

		// Verify tag is not found
		_, err = client.GetTag(ctx, tagID)
		require.Error(t, err)
		require.ErrorIs(t, err, kuma.ErrNotFound)
	})
}

func TestClient_MonitorTagAssociations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	// Create test tags
	tag1ID, err := client.CreateTag(ctx, tag.Tag{Name: "Environment", Color: "#FF0000"})
	require.NoError(t, err)
	defer func() { _ = client.DeleteTag(ctx, tag1ID) }()

	tag2ID, err := client.CreateTag(ctx, tag.Tag{Name: "Priority", Color: "#00FF00"})
	require.NoError(t, err)
	defer func() { _ = client.DeleteTag(ctx, tag2ID) }()

	tag3ID, err := client.CreateTag(ctx, tag.Tag{Name: "Team", Color: "#0000FF"})
	require.NoError(t, err)
	defer func() { _ = client.DeleteTag(ctx, tag3ID) }()

	// Create test monitors
	monitorID1, err := client.CreateMonitor(ctx, &monitor.HTTP{
		Base: monitor.Base{
			Name:           "Test Monitor 1",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 "https://example.com",
			Timeout:             48,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
			AuthMethod:          monitor.AuthMethodNone,
		},
	})
	require.NoError(t, err)
	defer func() { _ = client.DeleteMonitor(ctx, monitorID1) }()

	monitorID2, err := client.CreateMonitor(ctx, &monitor.HTTP{
		Base: monitor.Base{
			Name:           "Test Monitor 2",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 "https://example.org",
			Timeout:             48,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
			AuthMethod:          monitor.AuthMethodNone,
		},
	})
	require.NoError(t, err)
	defer func() { _ = client.DeleteMonitor(ctx, monitorID2) }()

	t.Run("add_tag_to_monitor", func(t *testing.T) {
		// Add tag with value to monitor
		monitorTag, err := client.AddMonitorTag(ctx, tag1ID, monitorID1, "production")
		require.NoError(t, err)
		require.NotNil(t, monitorTag)
		require.Equal(t, tag1ID, monitorTag.TagID)
		require.Equal(t, monitorID1, monitorTag.MonitorID)
		require.Equal(t, "production", monitorTag.Value)
		require.Equal(t, "Environment", monitorTag.Name)
		require.Equal(t, "#FF0000", monitorTag.Color)

		// Verify tag appears on monitor
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)
		require.Len(t, tags, 1)
		require.Equal(t, tag1ID, tags[0].TagID)
		require.Equal(t, "production", tags[0].Value)
		require.Equal(t, "Environment", tags[0].Name)
		require.Equal(t, "#FF0000", tags[0].Color)
	})

	t.Run("add_multiple_tags_to_monitor", func(t *testing.T) {
		// Add second tag
		_, err := client.AddMonitorTag(ctx, tag2ID, monitorID1, "high")
		require.NoError(t, err)

		// Add third tag with empty value
		_, err = client.AddMonitorTag(ctx, tag3ID, monitorID1, "")
		require.NoError(t, err)

		// Verify all tags
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)
		require.Len(t, tags, 3)

		// Check each tag
		tagsByID := make(map[int64]tag.MonitorTag)
		for _, mt := range tags {
			tagsByID[mt.TagID] = mt
		}

		require.Equal(t, "production", tagsByID[tag1ID].Value)
		require.Equal(t, "high", tagsByID[tag2ID].Value)
		require.Empty(t, tagsByID[tag3ID].Value)
	})

	t.Run("add_same_tag_with_different_values", func(t *testing.T) {
		// Add the same tag with a different value
		_, err := client.AddMonitorTag(ctx, tag1ID, monitorID1, "staging")
		require.NoError(t, err)

		// Verify both values exist
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)

		// Count environment tags
		envTags := 0
		for _, mt := range tags {
			if mt.TagID == tag1ID {
				envTags++
			}
		}

		require.Equal(t, 2, envTags, "Should have 2 Environment tags with different values")
	})

	t.Run("update_monitor_tag_value", func(t *testing.T) {
		// Update one of the tag values
		err := client.UpdateMonitorTag(ctx, tag2ID, monitorID1, "critical")
		require.NoError(t, err)

		// Verify update
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)

		found := false
		for _, mt := range tags {
			if mt.TagID == tag2ID {
				require.Equal(t, "critical", mt.Value)
				found = true
				break
			}
		}

		require.True(t, found, "Updated tag not found")
	})

	t.Run("add_tag_to_multiple_monitors", func(t *testing.T) {
		// Add tag1 to monitor2
		_, err := client.AddMonitorTag(ctx, tag1ID, monitorID2, "development")
		require.NoError(t, err)

		// Verify tag appears on both monitors
		monitors, err := client.GetTagMonitors(ctx, tag1ID)
		require.NoError(t, err)
		require.Len(t, monitors, 2)
		require.Contains(t, monitors, monitorID1)
		require.Contains(t, monitors, monitorID2)
	})

	t.Run("get_tag_monitors", func(t *testing.T) {
		// Get monitors with tag2
		monitors, err := client.GetTagMonitors(ctx, tag2ID)
		require.NoError(t, err)
		require.Len(t, monitors, 1)
		require.Equal(t, monitorID1, monitors[0])

		// Get monitors with tag3
		monitors, err = client.GetTagMonitors(ctx, tag3ID)
		require.NoError(t, err)
		require.Len(t, monitors, 1)
		require.Equal(t, monitorID1, monitors[0])
	})

	t.Run("delete_specific_tag_value", func(t *testing.T) {
		// Delete one specific Environment tag value (staging)
		err := client.DeleteMonitorTagWithValue(ctx, tag1ID, monitorID1, "staging")
		require.NoError(t, err)

		// Verify only production remains for monitor1
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)

		envTags := 0
		for _, mt := range tags {
			if mt.TagID == tag1ID && mt.MonitorID == monitorID1 {
				envTags++
				require.Equal(t, "production", mt.Value)
			}
		}

		require.Equal(t, 1, envTags, "Should have only 1 Environment tag after deleting staging")
	})

	t.Run("delete_all_tag_associations", func(t *testing.T) {
		// Delete all tag1 associations from monitor1
		err := client.DeleteMonitorTag(ctx, tag1ID, monitorID1)
		require.NoError(t, err)

		// Verify tag1 is removed from monitor1
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)

		for _, mt := range tags {
			require.NotEqual(t, tag1ID, mt.TagID, "tag1 should be removed from monitor1")
		}

		// Verify tag1 still exists on monitor2
		tags, err = client.GetMonitorTags(ctx, monitorID2)
		require.NoError(t, err)

		found := false
		for _, mt := range tags {
			if mt.TagID == tag1ID {
				found = true
				break
			}
		}

		require.True(t, found, "tag1 should still exist on monitor2")
	})

	t.Run("monitor_includes_tags", func(t *testing.T) {
		// Get monitor and verify tags are included
		mon, err := client.GetMonitor(ctx, monitorID2)
		require.NoError(t, err)
		require.NotNil(t, mon.Tags)
		require.NotEmpty(t, mon.Tags)

		// Verify tag details
		found := false
		for _, mt := range mon.Tags {
			if mt.TagID != tag1ID {
				continue
			}
			require.Equal(t, "development", mt.Value)
			require.Equal(t, "Environment", mt.Name)
			require.Equal(t, "#FF0000", mt.Color)
			found = true
			break
		}

		require.True(t, found, "tag1 should be present in monitor2 tags")
	})

	t.Run("delete_tag_removes_associations", func(t *testing.T) {
		// Create a temporary tag
		tempTagID, err := client.CreateTag(ctx, tag.Tag{Name: "Temp", Color: "#AAAAAA"})
		require.NoError(t, err)

		// Add to monitor
		_, err = client.AddMonitorTag(ctx, tempTagID, monitorID1, "temp")
		require.NoError(t, err)

		// Verify it's there
		tags, err := client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)
		found := false
		for _, mt := range tags {
			if mt.TagID == tempTagID {
				found = true
				break
			}
		}

		require.True(t, found, "temp tag should be on monitor")

		// Delete the tag
		err = client.DeleteTag(ctx, tempTagID)
		require.NoError(t, err)

		// Verify associations are removed (cascade delete)
		tags, err = client.GetMonitorTags(ctx, monitorID1)
		require.NoError(t, err)
		for _, mt := range tags {
			require.NotEqual(t, tempTagID, mt.TagID, "temp tag should be removed after tag deletion")
		}
	})

	t.Run("edge_case_nonexistent_tag", func(t *testing.T) {
		// Try to add non-existent tag
		_, err := client.AddMonitorTag(ctx, 99999, monitorID1, "test")
		require.Error(t, err)
	})

	t.Run("edge_case_nonexistent_monitor", func(t *testing.T) {
		// Try to add tag to non-existent monitor
		_, err := client.AddMonitorTag(ctx, tag1ID, 99999, "test")
		require.Error(t, err)
	})
}
