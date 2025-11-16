// Package maintenance provides types and utilities for managing Uptime Kuma maintenance windows.
//
// Maintenance windows are used to schedule planned downtime and prevent false alerts
// during scheduled maintenance activities. They can be configured with various
// scheduling strategies including one-time windows, recurring schedules, and cron-based timing.
//
// Supported Strategies:
//   - single: One-time maintenance window
//   - recurring-interval: Repeats every N days
//   - recurring-weekday: Repeats on specific days of the week
//   - recurring-day-of-month: Repeats on specific days of the month
//   - cron: Custom cron expression
//   - manual: Manually triggered maintenance
//
// Example usage:
//
//	// Create a one-time maintenance window
//	m := maintenance.NewSingleMaintenance(
//	    "Server Upgrade",
//	    "Upgrading to new version",
//	    startTime,
//	    endTime,
//	    "UTC",
//	)
//
//	// Create a recurring weekly maintenance (every Monday at 2 AM for 2 hours)
//	m := maintenance.NewRecurringWeekdayMaintenance(
//	    "Weekly Backup",
//	    "Weekly database backup",
//	    []int{1}, // Monday
//	    []maintenance.TimeOfDay{
//	        {Hours: 2, Minutes: 0, Seconds: 0},
//	        {Hours: 4, Minutes: 0, Seconds: 0},
//	    },
//	    "America/New_York",
//	)
package maintenance
