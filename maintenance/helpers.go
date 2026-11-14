package maintenance

import "time"

// NewSingleMaintenance creates a one-time maintenance window. The dates are
// sent as RFC 3339 with the offset of their location, a zero time.Time as
// "0001-01-01T00:00:00Z". A server backed by MariaDB rejects such dates, and
// Uptime Kuma reads them differently from this client unless the timezone is
// UTC, see DateTime. To send wall-clock time in the window's timezone, as the
// Uptime Kuma UI does, set DateRange to NewFloatingDateTime values.
func NewSingleMaintenance(
	title string,
	description string,
	startDate time.Time,
	endDate time.Time,
	timezone string,
) *Maintenance {
	start := NewDateTime(startDate)
	end := NewDateTime(endDate)

	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "single",
		Active:         true,
		DateRange:      []*DateTime{&start, &end},
		TimezoneOption: timezone,
	}
}

// NewRecurringIntervalMaintenance creates a maintenance window that repeats every N days.
func NewRecurringIntervalMaintenance(
	title string, description string,
	intervalDay int,
	timeRange []TimeOfDay,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "recurring-interval",
		Active:         true,
		IntervalDay:    intervalDay,
		DateRange:      []*DateTime{nil, nil},
		TimeRange:      timeRange,
		TimezoneOption: timezone,
	}
}

// NewRecurringWeekdayMaintenance creates a maintenance window that repeats on specific days of the week.
// Weekdays: 1=Monday, 2=Tuesday, ..., 7=Sunday.
func NewRecurringWeekdayMaintenance(
	title string, description string,
	weekdays []int,
	timeRange []TimeOfDay,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "recurring-weekday",
		Active:         true,
		DateRange:      []*DateTime{nil, nil},
		TimeRange:      timeRange,
		Weekdays:       weekdays,
		TimezoneOption: timezone,
	}
}

// NewRecurringDayOfMonthMaintenance creates a maintenance window that repeats on specific days of the month.
// daysOfMonth can contain integers 1-31 or special strings "lastDay1"-"lastDay4".
func NewRecurringDayOfMonthMaintenance(
	title string, description string,
	daysOfMonth []any,
	timeRange []TimeOfDay,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "recurring-day-of-month",
		Active:         true,
		DateRange:      []*DateTime{nil, nil},
		TimeRange:      timeRange,
		DaysOfMonth:    daysOfMonth,
		TimezoneOption: timezone,
	}
}

// NewCronMaintenance creates a maintenance window using a cron expression.
// cronExpr: standard cron expression (e.g., "0 2 * * *" for daily at 2 AM).
// durationMinutes: how long the maintenance window should last.
func NewCronMaintenance(
	title string,
	description string,
	cronExpr string,
	durationMinutes int,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:           title,
		Description:     description,
		Strategy:        "cron",
		Active:          true,
		DateRange:       []*DateTime{nil, nil},
		Cron:            cronExpr,
		DurationMinutes: durationMinutes,
		TimezoneOption:  timezone,
	}
}

// NewManualMaintenance creates a manually-triggered maintenance window.
// Manual maintenance windows do not have a schedule and must be activated manually.
func NewManualMaintenance(title string, description string) *Maintenance {
	return &Maintenance{
		Title:       title,
		Description: description,
		Strategy:    "manual",
		Active:      true,
		DateRange:   []*DateTime{nil, nil},
	}
}
