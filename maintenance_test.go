package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/maintenance"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestClient_MaintenanceCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	var maintenanceID int64

	t.Run("create_single_maintenance", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)

		m := maintenance.NewSingleMaintenance(
			"Server Upgrade",
			"Planned server upgrade",
			startTime,
			endTime,
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Greater(t, created.ID, int64(0))
		require.Equal(t, "Server Upgrade", created.Title)
		require.Equal(t, "Planned server upgrade", created.Description)
		require.Equal(t, "single", created.Strategy)
		require.True(t, created.Active)

		maintenanceID = created.ID
	})

	var retrievedMaintenance *maintenance.Maintenance
	t.Run("get_maintenance", func(t *testing.T) {
		retrievedMaintenance, err = client.GetMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.NotNil(t, retrievedMaintenance)
		require.Equal(t, maintenanceID, retrievedMaintenance.ID)
		require.Equal(t, "Server Upgrade", retrievedMaintenance.Title)
		require.Equal(t, "single", retrievedMaintenance.Strategy)
	})

	t.Run("get_maintenances", func(t *testing.T) {
		maintenances, err := client.GetMaintenances(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, maintenances)
		t.Logf("Maintenances in cache: %d", len(maintenances))
	})

	t.Run("update_maintenance", func(t *testing.T) {
		retrievedMaintenance.Title = "Updated Server Upgrade"
		retrievedMaintenance.Description = "Updated planned server upgrade"

		err = client.UpdateMaintenance(ctx, retrievedMaintenance)
		require.NoError(t, err)

		updated, err := client.GetMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.Equal(t, "Updated Server Upgrade", updated.Title)
		require.Equal(t, "Updated planned server upgrade", updated.Description)
	})

	t.Run("pause_maintenance", func(t *testing.T) {
		err = client.PauseMaintenance(ctx, maintenanceID)
		require.NoError(t, err)

		paused, err := client.GetMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.False(t, paused.Active)
	})

	t.Run("resume_maintenance", func(t *testing.T) {
		err = client.ResumeMaintenance(ctx, maintenanceID)
		require.NoError(t, err)

		resumed, err := client.GetMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.True(t, resumed.Active)
	})

	t.Run("delete_maintenance", func(t *testing.T) {
		err = client.DeleteMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
	})
}

func TestClient_MaintenanceStrategies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	t.Run("recurring_weekday_maintenance", func(t *testing.T) {
		// Every Monday and Wednesday at 2 AM for 2 hours
		m := maintenance.NewRecurringWeekdayMaintenance(
			"Weekly Backup",
			"Weekly database backup",
			[]int{1, 3}, // Monday and Wednesday
			[]maintenance.TimeOfDay{
				{Hours: 2, Minutes: 0, Seconds: 0},
				{Hours: 4, Minutes: 0, Seconds: 0},
			},
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "recurring-weekday", created.Strategy)
		require.Equal(t, []int{1, 3}, created.Weekdays)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("recurring_interval_maintenance", func(t *testing.T) {
		// Every 3 days at 3 AM for 1 hour
		m := maintenance.NewRecurringIntervalMaintenance(
			"Periodic Cleanup",
			"System cleanup every 3 days",
			3,
			[]maintenance.TimeOfDay{
				{Hours: 3, Minutes: 0, Seconds: 0},
				{Hours: 4, Minutes: 0, Seconds: 0},
			},
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "recurring-interval", created.Strategy)
		require.Equal(t, 3, created.IntervalDay)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("recurring_day_of_month_maintenance", func(t *testing.T) {
		// 1st and 15th of each month at midnight for 2 hours
		m := maintenance.NewRecurringDayOfMonthMaintenance(
			"Monthly Maintenance",
			"Monthly system maintenance",
			[]interface{}{1, 15},
			[]maintenance.TimeOfDay{
				{Hours: 0, Minutes: 0, Seconds: 0},
				{Hours: 2, Minutes: 0, Seconds: 0},
			},
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "recurring-day-of-month", created.Strategy)
		require.NotEmpty(t, created.DaysOfMonth)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("cron_maintenance", func(t *testing.T) {
		// Daily at 2 AM for 30 minutes
		m := maintenance.NewCronMaintenance(
			"Daily Backup",
			"Daily automated backup",
			"0 2 * * *",
			30,
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "cron", created.Strategy)
		require.Equal(t, "0 2 * * *", created.Cron)
		require.Equal(t, 30, created.DurationMinutes)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("manual_maintenance", func(t *testing.T) {
		m := maintenance.NewManualMaintenance(
			"Emergency Maintenance",
			"Manual emergency maintenance window",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "manual", created.Strategy)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})
}

func TestClient_MaintenanceWithMonitors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	var monitor1ID, monitor2ID int64
	var maintenanceID int64

	t.Run("create_monitors", func(t *testing.T) {
		httpMonitor1 := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Maintenance Test Monitor 1",
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

		monitor1ID, err = client.CreateMonitor(ctx, httpMonitor1)
		require.NoError(t, err)
		require.Greater(t, monitor1ID, int64(0))

		httpMonitor2 := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Maintenance Test Monitor 2",
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     3,
				UpsideDown:     false,
				IsActive:       true,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                 "https://httpbin.org/status/201",
				Timeout:             48,
				Method:              "GET",
				ExpiryNotification:  false,
				IgnoreTLS:           false,
				MaxRedirects:        10,
				AcceptedStatusCodes: []string{"200-299"},
				AuthMethod:          monitor.AuthMethodNone,
			},
		}

		monitor2ID, err = client.CreateMonitor(ctx, httpMonitor2)
		require.NoError(t, err)
		require.Greater(t, monitor2ID, int64(0))
	})

	t.Run("create_maintenance", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)

		m := maintenance.NewSingleMaintenance(
			"Maintenance with Monitors",
			"Testing monitor associations",
			startTime,
			endTime,
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)

		maintenanceID = created.ID
	})

	t.Run("set_monitor_maintenance", func(t *testing.T) {
		err = client.SetMonitorMaintenance(ctx, maintenanceID, []int64{monitor1ID, monitor2ID})
		require.NoError(t, err)
	})

	t.Run("get_monitor_maintenance", func(t *testing.T) {
		monitorIDs, err := client.GetMonitorMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.Len(t, monitorIDs, 2)
		require.Contains(t, monitorIDs, monitor1ID)
		require.Contains(t, monitorIDs, monitor2ID)
	})

	t.Run("update_monitor_maintenance", func(t *testing.T) {
		// Remove monitor2, keep only monitor1
		err = client.SetMonitorMaintenance(ctx, maintenanceID, []int64{monitor1ID})
		require.NoError(t, err)

		monitorIDs, err := client.GetMonitorMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.Len(t, monitorIDs, 1)
		require.Contains(t, monitorIDs, monitor1ID)
	})

	t.Run("clear_monitor_maintenance", func(t *testing.T) {
		// Clear all monitor associations
		err = client.SetMonitorMaintenance(ctx, maintenanceID, []int64{})
		require.NoError(t, err)

		monitorIDs, err := client.GetMonitorMaintenance(ctx, maintenanceID)
		require.NoError(t, err)
		require.Empty(t, monitorIDs)
	})

	t.Run("cleanup", func(t *testing.T) {
		err = client.DeleteMaintenance(ctx, maintenanceID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, monitor1ID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, monitor2ID)
		require.NoError(t, err)
	})
}

func TestClient_MaintenanceWithStatusPages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	var maintenanceID int64

	slug := "test-maintenance-status-page"
	title := "Maintenance Status Page Test"

	t.Run("create_status_page", func(t *testing.T) {
		err = client.AddStatusPage(ctx, title, slug)
		require.NoError(t, err)
	})

	// Get the status page to retrieve its ID
	var statusPageID int64
	t.Run("get_status_page_id", func(t *testing.T) {
		sp, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.NotNil(t, sp)
		statusPageID = sp.ID
		require.Greater(t, statusPageID, int64(0))
	})

	t.Run("create_maintenance", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)

		m := maintenance.NewSingleMaintenance(
			"Maintenance with Status Page",
			"Testing status page associations",
			startTime,
			endTime,
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)

		maintenanceID = created.ID
	})

	t.Run("set_maintenance_status_page", func(t *testing.T) {
		err = client.SetMaintenanceStatusPage(ctx, maintenanceID, []int64{statusPageID})
		require.NoError(t, err)
	})

	t.Run("get_maintenance_status_page", func(t *testing.T) {
		statusPageIDs, err := client.GetMaintenanceStatusPage(ctx, maintenanceID)
		require.NoError(t, err)
		require.Len(t, statusPageIDs, 1)
		require.Contains(t, statusPageIDs, statusPageID)
	})

	t.Run("clear_maintenance_status_page", func(t *testing.T) {
		// Clear all status page associations
		err = client.SetMaintenanceStatusPage(ctx, maintenanceID, []int64{})
		require.NoError(t, err)

		statusPageIDs, err := client.GetMaintenanceStatusPage(ctx, maintenanceID)
		require.NoError(t, err)
		require.Empty(t, statusPageIDs)
	})

	t.Run("cleanup", func(t *testing.T) {
		err = client.DeleteMaintenance(ctx, maintenanceID)
		require.NoError(t, err)

		err = client.DeleteStatusPage(ctx, slug)
		require.NoError(t, err)
	})
}

func TestClient_MaintenanceTimezones(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("maintenance_with_utc_timezone", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		m := maintenance.NewSingleMaintenance(
			"UTC Timezone Test",
			"Testing UTC timezone",
			startTime,
			endTime,
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "UTC", created.TimezoneOption)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("maintenance_with_specific_timezone", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		m := maintenance.NewSingleMaintenance(
			"America/New_York Timezone Test",
			"Testing specific timezone",
			startTime,
			endTime,
			"America/New_York",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "America/New_York", created.TimezoneOption)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})
}

func TestClient_MaintenanceEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("maintenance_with_last_day_of_month", func(t *testing.T) {
		// Last day and 2nd to last day of each month
		m := maintenance.NewRecurringDayOfMonthMaintenance(
			"End of Month Maintenance",
			"Maintenance on last days of month",
			[]interface{}{"lastDay1", "lastDay2"},
			[]maintenance.TimeOfDay{
				{Hours: 23, Minutes: 0, Seconds: 0},
				{Hours: 23, Minutes: 59, Seconds: 59},
			},
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, "recurring-day-of-month", created.Strategy)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})

	t.Run("maintenance_crossing_midnight", func(t *testing.T) {
		// 11 PM to 1 AM (crosses midnight)
		m := maintenance.NewRecurringWeekdayMaintenance(
			"Midnight Crossing Maintenance",
			"Maintenance that crosses midnight",
			[]int{7}, // Sunday
			[]maintenance.TimeOfDay{
				{Hours: 23, Minutes: 0, Seconds: 0},
				{Hours: 1, Minutes: 0, Seconds: 0},
			},
			"UTC",
		)

		created, err := client.CreateMaintenance(ctx, m)
		require.NoError(t, err)
		require.NotNil(t, created)

		err = client.DeleteMaintenance(ctx, created.ID)
		require.NoError(t, err)
	})
}
