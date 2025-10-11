package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/tag"
)

func TestClient_TagCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		require.Greater(t, tagID, int64(0))

		// Test GetTags after creation
		tags, err := client.GetTags(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(tags))

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
		require.Equal(t, initialCount, len(tags))

		// Verify tag is not found
		_, err = client.GetTag(ctx, tagID)
		require.Error(t, err)
		require.ErrorIs(t, err, kuma.ErrNotFound)
	})
}
