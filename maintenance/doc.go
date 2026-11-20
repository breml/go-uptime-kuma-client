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
//
// Dates:
//
// DateRange and the dates of TimeslotList are DateTime values, either floating,
// wall-clock time without UTC offset meant in the window's Timezone, such as
// the "2026-06-29T13:39" the Uptime Kuma UI writes, or with an offset, such as
// the RFC 3339 NewSingleMaintenance writes. An absent date is the zero
// DateTime, check it with IsZero.
//
// Prefer floating dates, see DateTime for why. Take the digits in the window's
// timezone:
//
//	loc, err := time.LoadLocation("Europe/Zurich")
//	if err != nil {
//	    return err
//	}
//	m := maintenance.NewSingleMaintenance("Upgrade", "", startTime, endTime, "Europe/Zurich")
//	m.DateRange = []maintenance.DateTime{
//	    maintenance.NewFloatingDateTime(startTime.In(loc)),
//	    maintenance.NewFloatingDateTime(endTime.In(loc)),
//	}
//
// For a floating date, DateTime.In with the window's timezone gives the
// instant it means. For a date with an offset, In converts the instant, which
// is not how Uptime Kuma's status computation reads it, see DateTime.
package maintenance
