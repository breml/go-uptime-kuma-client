package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestClient_MonitorCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test HTTP monitor
	httpMonitor := monitor.HTTP{
		Base: monitor.Base{
			Name:           "Test HTTP Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 "https://httpbin.org/status/200",
			Timeout:             48,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
			AuthMethod:          monitor.AuthMethodNone,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var httpMonitorRetrieved monitor.HTTP
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, httpMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test HTTP Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs

		err = client.GetMonitorAs(ctx, monitorID, &httpMonitorRetrieved)
		require.NoError(t, err)
		httpMonitor.ID = monitorID
		httpMonitor.PathName = httpMonitor.Name
		require.EqualExportedValues(t, httpMonitor, httpMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		httpMonitorRetrieved.Name = "Updated HTTP Monitor"
		httpMonitorRetrieved.URL = "https://httpbin.org/status/201"
		err := client.UpdateMonitor(ctx, httpMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated HTTP Monitor", updatedMonitor.Name)

		var updatedHTTP monitor.HTTP
		err = client.GetMonitorAs(ctx, monitorID, &updatedHTTP)
		require.NoError(t, err)
		require.Equal(t, "https://httpbin.org/status/201", updatedHTTP.URL)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorGroupCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test Group monitor
	groupMonitor := monitor.Group{
		Base: monitor.Base{
			Name:           "Test Group Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     0,
			UpsideDown:     false,
			IsActive:       false,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var groupMonitorRetrieved monitor.Group
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, groupMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test Group Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &groupMonitorRetrieved)
		require.NoError(t, err)
		groupMonitor.ID = monitorID
		groupMonitor.PathName = groupMonitor.Name
		require.EqualExportedValues(t, groupMonitor, groupMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		groupMonitorRetrieved.Name = "Updated Group Monitor"
		err := client.UpdateMonitor(ctx, groupMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated Group Monitor", updatedMonitor.Name)

		var updatedGroup monitor.Group
		err = client.GetMonitorAs(ctx, monitorID, &updatedGroup)
		require.NoError(t, err)
		require.Equal(t, "Updated Group Monitor", updatedGroup.Name)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorParent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	groupMonitor := monitor.Group{
		Base: monitor.Base{
			Name:           "Parent Group",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     0,
			UpsideDown:     false,
			IsActive:       true,
		},
	}

	var parentID int64
	t.Run("create_parent_group", func(t *testing.T) {
		parentID, err = client.CreateMonitor(ctx, groupMonitor)
		require.NoError(t, err)
		require.Greater(t, parentID, int64(0))
	})

	var childID int64
	t.Run("create_monitor_with_parent", func(t *testing.T) {
		httpMonitor := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Child Monitor",
				Parent:         &parentID,
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     3,
				UpsideDown:     false,
				IsActive:       true,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                 "https://httpbin.org/status/200",
				Timeout:             48,
				Method:              "GET",
				ExpiryNotification:  false,
				IgnoreTLS:           false,
				MaxRedirects:        10,
				AcceptedStatusCodes: []string{"200-299"},
				AuthMethod:          monitor.AuthMethodNone,
			},
		}

		childID, err = client.CreateMonitor(ctx, httpMonitor)
		require.NoError(t, err)
		require.Greater(t, childID, int64(0))

		var childMonitorRetrieved monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &childMonitorRetrieved)
		require.NoError(t, err)
		require.NotNil(t, childMonitorRetrieved.Parent)
		require.Equal(t, parentID, *childMonitorRetrieved.Parent)
	})

	var newParentID int64
	t.Run("update_parent", func(t *testing.T) {
		newGroupMonitor := monitor.Group{
			Base: monitor.Base{
				Name:           "New Parent Group",
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     0,
				UpsideDown:     false,
				IsActive:       true,
			},
		}

		newParentID, err = client.CreateMonitor(ctx, newGroupMonitor)
		require.NoError(t, err)
		require.Greater(t, newParentID, int64(0))

		var httpMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &httpMonitor)
		require.NoError(t, err)

		httpMonitor.Parent = &newParentID

		err = client.UpdateMonitor(ctx, httpMonitor)
		require.NoError(t, err)

		var updatedMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &updatedMonitor)
		require.NoError(t, err)
		require.NotNil(t, updatedMonitor.Parent)
		require.Equal(t, newParentID, *updatedMonitor.Parent)
	})

	t.Run("remove_parent", func(t *testing.T) {
		var httpMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &httpMonitor)
		require.NoError(t, err)

		httpMonitor.Parent = nil

		err = client.UpdateMonitor(ctx, httpMonitor)
		require.NoError(t, err)

		var updatedMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &updatedMonitor)
		require.NoError(t, err)
		require.Nil(t, updatedMonitor.Parent)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteMonitor(ctx, childID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, newParentID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, parentID)
		require.NoError(t, err)
	})
}
