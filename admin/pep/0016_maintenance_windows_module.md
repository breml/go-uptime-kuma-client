# PEP 0016: Maintenance Windows Module

**PEP Number:** 0016  
**Title:** Add Maintenance Windows Module  
**Status:** Ready

## Abstract

Add a new `maintenance` package to provide complete management of Uptime Kuma maintenance windows. This includes creating, editing, deleting maintenance windows, scheduling recurring maintenance, and associating monitors and status pages with maintenance periods.

## Motivation

Maintenance windows are essential for planned downtime management in monitoring systems. They prevent false alerts during scheduled maintenance activities and provide visibility into planned service interruptions on status pages.

Maintenance windows enable:

- Scheduled monitoring pauses for planned maintenance
- Recurring maintenance schedules (daily, weekly, monthly)
- Association of multiple monitors with a single maintenance window
- Status page integration for public maintenance notifications
- Prevention of alert fatigue during planned downtime

## Specification

### New Package Structure

Create a new package `maintenance/` with the following structure:

```none
maintenance/
├── maintenance.go          # Core maintenance window entity
├── maintenance_test.go     # Unit tests
├── schedule.go             # Scheduling and recurrence logic
├── helpers.go              # Helper functions
```

### Core Data Structures

#### Maintenance Window Entity

```go
type Maintenance struct {
    ID                int64             `json:"id"`
    Title             string            `json:"title"`
    Description       string            `json:"description"`
    Strategy          string            `json:"strategy"`          // "single", "recurring-interval", "recurring-weekday", "recurring-day-of-month", "cron", "manual"
    Active            bool              `json:"active"`
    IntervalDay       int               `json:"intervalDay"`       // For recurring-interval strategy
    DateRange         []*time.Time      `json:"dateRange"`         // Start and end time [start, end], can be nil
    TimeRange         []TimeOfDay       `json:"timeRange"`         // Time of day range [start, end]
    Weekdays          []int             `json:"weekdays"`          // 1-7 (Monday=1, Tuesday=2, ..., Sunday=7) for recurring-weekday
    DaysOfMonth       []interface{}     `json:"daysOfMonth"`       // 1-31 or "lastDay1"-"lastDay4" for recurring-day-of-month
    Cron              string            `json:"cron"`              // Cron expression for cron strategy
    Duration          int               `json:"duration"`          // Duration in seconds
    DurationMinutes   int               `json:"durationMinutes"`   // Duration in minutes (calculated: duration / 60)
    Timezone          string            `json:"timezone"`          // Resolved IANA timezone (e.g., "America/New_York")
    TimezoneOption    string            `json:"timezoneOption"`    // Original timezone option: "UTC", "SAME_AS_SERVER", or IANA timezone
    TimezoneOffset    string            `json:"timezoneOffset"`    // e.g., "+02:00"
    Status            string            `json:"status"`            // "inactive", "scheduled", "under-maintenance", "ended", "unknown" (read-only, computed by server)
    TimeslotList      []Timeslot        `json:"timeslotList"`      // Scheduled maintenance windows (read-only, computed by server)
}

type TimeOfDay struct {
    Hours   int `json:"hours"`   // 0-23
    Minutes int `json:"minutes"` // 0-59
    Seconds int `json:"seconds"` // 0-59
}

type Timeslot struct {
    StartDate time.Time `json:"startDate"`
    EndDate   time.Time `json:"endDate"`
}

type MaintenanceMonitor struct {
    MaintenanceID int64 `json:"maintenance_id"`
    MonitorID     int64 `json:"monitor_id"`
}

type MaintenanceStatusPage struct {
    MaintenanceID  int64  `json:"maintenance_id"`
    StatusPageID   int64  `json:"status_page_id"`
}
```

### Client Methods

Add the following methods to the main `Client` type in `maintenance.go`:

```go
// Maintenance Window Management
func (c *Client) GetMaintenance(ctx context.Context, id int64) (*maintenance.Maintenance, error)
func (c *Client) GetMaintenances(ctx context.Context) ([]maintenance.Maintenance, error)
func (c *Client) CreateMaintenance(ctx context.Context, m *maintenance.Maintenance) (*maintenance.Maintenance, error)
func (c *Client) UpdateMaintenance(ctx context.Context, m *maintenance.Maintenance) error
func (c *Client) DeleteMaintenance(ctx context.Context, id int64) error
func (c *Client) PauseMaintenance(ctx context.Context, id int64) error
func (c *Client) ResumeMaintenance(ctx context.Context, id int64) error

// Monitor Association
// Note: SetMonitorMaintenance replaces all monitors, it doesn't append
// Internally converts []int64 to []MonitorID{ID: int64} format expected by server
func (c *Client) SetMonitorMaintenance(ctx context.Context, maintenanceID int64, monitorIDs []int64) error
func (c *Client) GetMonitorMaintenance(ctx context.Context, maintenanceID int64) ([]int64, error)

// Status Page Association
// Note: SetMaintenanceStatusPage replaces all status pages, it doesn't append
// Internally converts []int64 to []StatusPageID{ID: int64} format expected by server
func (c *Client) SetMaintenanceStatusPage(ctx context.Context, maintenanceID int64, statusPageIDs []int64) error
func (c *Client) GetMaintenanceStatusPage(ctx context.Context, maintenanceID int64) ([]int64, error)
```

### Socket.io Events

Map to the following Uptime Kuma socket.io events:

- `getMaintenance` - Get specific maintenance window
- `getMaintenanceList` - Get all maintenance windows
- `addMaintenance` - Create new maintenance window
- `editMaintenance` - Update maintenance window
- `deleteMaintenance` - Delete maintenance window
- `pauseMaintenance` - Pause/deactivate maintenance window
- `resumeMaintenance` - Resume/activate maintenance window
- `addMonitorMaintenance` - Set monitors for maintenance (replaces existing list)
- `getMonitorMaintenance` - Get monitors for maintenance
- `addMaintenanceStatusPage` - Set status pages for maintenance (replaces existing list)
- `getMaintenanceStatusPage` - Get status pages for maintenance

### State Management

Add maintenance list to the client state:

```go
type state struct {
    notifications []notification.Base
    monitors      []monitor.Base
    statusPages   []statuspage.StatusPage
    maintenances  []maintenance.Maintenance  // New field
}
```

Listen for `maintenanceList` event during client initialization to populate the cache.

**Note:** The `maintenanceList` event from Uptime Kuma sends a map `{[id: string]: maintenance}` rather than an array. The client should convert this to a slice of `Maintenance` objects for internal storage.

### Data Format Details

#### DateRange Format

- Array with exactly 2 elements: `[startDate, endDate]`
- Each element can be ISO 8601 datetime string or `null`
- Example: `["2025-12-25T10:30:00.000Z", "2025-12-25T18:30:00.000Z"]`
- For recurring strategies, start date can be `null` to indicate "start now"

#### TimeRange Format

- Array with exactly 2 elements: `[startTime, endTime]`
- Each element is a `TimeOfDay` object: `{hours: int, minutes: int, seconds: int}`
- Used only for recurring strategies and cron
- Example: `[{hours: 9, minutes: 0, seconds: 0}, {hours: 17, minutes: 0, seconds: 0}]`

#### Weekdays Format

- Array of integers: 1-7 where Monday=1, Tuesday=2, ..., Sunday=7
- Used only for "recurring-weekday" strategy
- Example: `[1, 3, 5]` for Monday, Wednesday, Friday

#### DaysOfMonth Format

- Array of integers (1-31) or special strings
- Integers: day of month (1-31)
- Special strings: "lastDay1", "lastDay2", "lastDay3", "lastDay4" (last day, 2nd to last, etc.)
- Used only for "recurring-day-of-month" strategy
- Example: `[1, 15, "lastDay1"]` for 1st, 15th, and last day of month

#### Monitor/StatusPage Association Format

- When sending to server: array of objects with `id` field: `[{id: 1}, {id: 2}]`
- When receiving from server: array of objects with `id` field (monitors) or `id` and `title` (status pages)
- Client methods should accept/return simple `[]int64` and handle conversion

## Backward Compatibility

Not a concern for this early-stage project. Breaking changes are acceptable if they improve the overall API design.

## Test Plan

### Unit Tests

- Test CRUD operations for maintenance windows
- Test all strategy types (single, recurring-interval, recurring-weekday, recurring-day-of-month, cron, manual)
- Test date/time range validation
- Test timezone handling
- Test weekday and day-of-month validation
- Test cron expression validation
- Test monitor association operations
- Test status page association operations
- Test pause/resume functionality

### Integration Tests

Add integration tests in `maintenance_test.go` that:

1. Create a single maintenance window
2. Create a recurring maintenance window (weekly)
3. Set monitors for maintenance window
4. Verify monitor list can be retrieved
5. Update monitor list (add/remove monitors by replacing full list)
6. Set status pages for maintenance window
7. Verify status page list can be retrieved
8. Update status page list (add/remove by replacing full list)
9. Verify monitors are paused during active maintenance
10. Edit maintenance window schedule
11. Pause maintenance window
12. Resume maintenance window
13. Delete maintenance window and verify associations are cleaned up
14. Test different strategies (interval, weekday, day-of-month, cron)
15. Test timezone handling with different timezones

Integration tests must run against a real Uptime Kuma instance in Docker (following existing test patterns).

### Edge Cases to Test

- Maintenance window spanning multiple days
- Overlapping maintenance windows (server allows this)
- Maintenance windows with invalid cron expressions (verify server rejects)
- Maintenance windows with invalid date ranges (verify server rejects)
- Association with non-existent monitors/status pages (verify server behavior)
- Deletion of maintenance with active associations (should cascade delete)
- Maintenance crossing midnight (endTime < startTime)
- Different timezone options: "UTC", "SAME_AS_SERVER", and specific IANA timezones
- Empty monitor/status page lists (clearing all associations)

## Implementation

### Phase 1: Core Maintenance CRUD

- Create `maintenance` package structure
- Implement base `Maintenance` struct with all strategy types
- Add `GetMaintenance`, `GetMaintenanceList`, `AddMaintenance`, `UpdateMaintenance`, `DeleteMaintenance` methods
- Add unit tests for CRUD operations

### Phase 2: Schedule Strategies

- Implement schedule validation for each strategy type
- Add helper functions for cron validation
- Add helper functions for date/time range validation
- Test all scheduling strategies

### Phase 3: Pause/Resume Functionality

- Implement `PauseMaintenance` method
- Implement `ResumeMaintenance` method
- Add tests for pause/resume operations
- Verify behavior when maintenance is paused and resumed

### Phase 4: Monitor Associations

- Implement `SetMonitorMaintenance` method (replaces entire list)
- Implement `GetMonitorMaintenance` method
- Add tests for setting/retrieving monitor associations
- Test replacing monitor lists (add/remove by replacement)
- Verify monitors are correctly paused during maintenance

### Phase 5: Status Page Associations

- Implement `SetMaintenanceStatusPage` method (replaces entire list)
- Implement `GetMaintenanceStatusPage` method
- Add tests for setting/retrieving status page associations
- Test replacing status page lists (add/remove by replacement)
- Verify status pages display maintenance information

### Phase 6: Integration Testing

- Add comprehensive integration tests covering all scenarios
- Test against real Uptime Kuma instance
- Test interaction with monitors and status pages
- Document usage examples

### Phase 7: Documentation

- Add godoc comments for all public APIs
- Create usage examples for each strategy type
- Update main README with maintenance window examples
- Document best practices for recurring maintenance

## Reference Implementation

Refer to the following Uptime Kuma source files:

- `.scratch/uptime-kuma/server/socket-handlers/maintenance-socket-handler.js` - Server-side socket handlers
- `.scratch/uptime-kuma/server/model/maintenance.js` - Maintenance data model
- `.scratch/uptime-kuma/src/pages/EditMaintenance.vue` - Frontend implementation
- `.scratch/uptime-kuma/server/jobs.js` - Maintenance scheduling logic

## Decisions

1. **Helper Functions**: Provide helper functions for common maintenance schedules (e.g., "every weekend", "every night") as they add value for users.
2. **Timezone Handling**: Store timezone information as provided by the user. The server handles timezone conversions and scheduling logic.
3. **Cron Validation**: Rely on server-side validation for cron expressions. The client should not duplicate this validation logic.
4. **Active Status Check**: No dedicated methods needed to check if a maintenance window is currently active. Use the `Status` field from the `Maintenance` struct.
5. **Deleting Active Maintenance**: Perform the delete operation and let the server decide if it's possible. Don't make assumptions client-side about whether an active maintenance can be deleted.
6. **Duration Calculation**: When creating/updating maintenance with recurring strategies, the client should send `TimeRange`. The server automatically calculates `Duration` (in seconds) from the difference between start and end times. For cron strategy, send `DurationMinutes` and the server converts to `Duration`.
7. **Read-Only Fields**: `Status`, `TimeslotList`, `DurationMinutes` (when using recurring strategies), `Timezone` (when using "SAME_AS_SERVER"), and `TimezoneOffset` are computed by the server and should be treated as read-only in responses.
