package maintenance

import "time"

// NewSingleMaintenance creates a one-time maintenance window.
func NewSingleMaintenance(title, description string, startDate, endDate time.Time, timezone string) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "single",
		Active:         true,
		DateRange:      []*time.Time{&startDate, &endDate},
		TimezoneOption: timezone,
	}
}

// NewRecurringIntervalMaintenance creates a maintenance window that repeats every N days.
func NewRecurringIntervalMaintenance(
	title, description string,
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
		DateRange:      []*time.Time{nil, nil},
		TimeRange:      timeRange,
		TimezoneOption: timezone,
	}
}

// NewRecurringWeekdayMaintenance creates a maintenance window that repeats on specific days of the week.
// Weekdays: 1=Monday, 2=Tuesday, ..., 7=Sunday.
func NewRecurringWeekdayMaintenance(
	title, description string,
	weekdays []int,
	timeRange []TimeOfDay,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "recurring-weekday",
		Active:         true,
		DateRange:      []*time.Time{nil, nil},
		TimeRange:      timeRange,
		Weekdays:       weekdays,
		TimezoneOption: timezone,
	}
}

// NewRecurringDayOfMonthMaintenance creates a maintenance window that repeats on specific days of the month.
// daysOfMonth can contain integers 1-31 or special strings "lastDay1"-"lastDay4".
func NewRecurringDayOfMonthMaintenance(
	title, description string,
	daysOfMonth []any,
	timeRange []TimeOfDay,
	timezone string,
) *Maintenance {
	return &Maintenance{
		Title:          title,
		Description:    description,
		Strategy:       "recurring-day-of-month",
		Active:         true,
		DateRange:      []*time.Time{nil, nil},
		TimeRange:      timeRange,
		DaysOfMonth:    daysOfMonth,
		TimezoneOption: timezone,
	}
}

// NewCronMaintenance creates a maintenance window using a cron expression.
// cronExpr: standard cron expression (e.g., "0 2 * * *" for daily at 2 AM).
// durationMinutes: how long the maintenance window should last.
func NewCronMaintenance(title, description, cronExpr string, durationMinutes int, timezone string) *Maintenance {
	return &Maintenance{
		Title:           title,
		Description:     description,
		Strategy:        "cron",
		Active:          true,
		DateRange:       []*time.Time{nil, nil},
		Cron:            cronExpr,
		DurationMinutes: durationMinutes,
		TimezoneOption:  timezone,
	}
}

// NewManualMaintenance creates a manually-triggered maintenance window.
// Manual maintenance windows do not have a schedule and must be activated manually.
func NewManualMaintenance(title, description string) *Maintenance {
	return &Maintenance{
		Title:       title,
		Description: description,
		Strategy:    "manual",
		Active:      true,
		DateRange:   []*time.Time{nil, nil},
	}
}
