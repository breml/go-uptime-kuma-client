# PEP 0017: Monitor-Tag Association Management

**PEP Number:** 0017  
**Title:** Add Monitor-Tag Association Management  
**Status:** Ready

## Abstract

Extend the tag functionality to support associating tags with monitors, including managing tag values and colors for individual monitor-tag relationships. This completes the tag management API by allowing users to organize and categorize monitors using tags.

## Motivation

While the client currently supports basic tag CRUD operations, it lacks the ability to associate tags with monitors. Tags in Uptime Kuma serve multiple purposes:

- Organizing monitors into logical groups
- Filtering and searching monitors
- Visual categorization with colors
- Associating metadata with monitors via tag values
- Displaying tags on status pages

Monitor-tag associations are many-to-many relationships where:

- A monitor can have multiple tags
- A tag can be associated with multiple monitors
- Each association can have a custom value (optional)
- The association inherits the color from the tag definition (no per-association color override)

## Specification

### Enhanced Data Structures

Add monitor-tag association structures to the `tag` package:

```go
// MonitorTag represents the association between a monitor and a tag
type MonitorTag struct {
    ID        int64  `json:"id"`         // Association ID
    TagID     int64  `json:"tag_id"`     // Tag ID
    MonitorID int64  `json:"monitor_id"` // Monitor ID
    Value     string `json:"value"`      // Optional tag value for this monitor
    Name      string `json:"name"`       // Tag name (from joined tag table)
    Color     string `json:"color"`      // Tag color (from joined tag table, not customizable per association)
}

// Extended Tag with monitor associations
type TagWithMonitors struct {
    Tag
    Monitors []int64 `json:"monitors"` // List of associated monitor IDs
}

// Monitor with its tags (for convenience)
type MonitorTags struct {
    MonitorID int64        `json:"monitor_id"`
    Tags      []MonitorTag `json:"tags"`
}
```

### Client Methods

Add the following methods to the main `Client` type in `tag.go`:

```go
// Add tag to monitor
func (c *Client) AddMonitorTag(ctx context.Context, tagID int64, monitorID int64, value string) (*tag.MonitorTag, error)

// Update monitor-tag association (update value only)
func (c *Client) UpdateMonitorTag(ctx context.Context, tagID int64, monitorID int64, value string) error

// Remove all associations of a tag from a monitor (all values)
func (c *Client) DeleteMonitorTag(ctx context.Context, tagID int64, monitorID int64) error

// Remove specific tag association from monitor by value
func (c *Client) DeleteMonitorTagWithValue(ctx context.Context, tagID int64, monitorID int64, value string) error

// Get all tags for a specific monitor
func (c *Client) GetMonitorTags(ctx context.Context, monitorID int64) ([]tag.MonitorTag, error)

// Get all monitors for a specific tag
func (c *Client) GetTagMonitors(ctx context.Context, tagID int64) ([]int64, error)
```

### Socket.io Events

Map to the following Uptime Kuma socket.io events:

- `addMonitorTag` - Associate tag with monitor
- `editMonitorTag` - Edit (update) monitor-tag association
- `deleteMonitorTag` - Remove tag from monitor

Note: Getting monitor tags and tag monitors may be derived from existing `getMonitor` and `getTags` responses which include association data.

### Integration with Existing Code

The `monitor.Base` struct should be extended to include tags:

```go
type Base struct {
    // ... existing fields ...
    Tags []tag.MonitorTag `json:"tags"`
}
```

This allows monitors to carry their tag associations when retrieved.

## Backward Compatibility

Not a concern for this early-stage project. Breaking changes are acceptable if they improve the overall API design.

The addition of `Tags` field to `monitor.Base` should be backward compatible as it will be empty for existing code.

## Test Plan

### Unit Tests

- Test adding tag to monitor
- Test editing monitor-tag value
- Test removing tag from monitor
- Test adding multiple tags to single monitor
- Test adding same tag to multiple monitors
- Test handling empty/nil values
- Test duplicate tag addition (should be idempotent or error)

### Integration Tests

Add integration tests in `main_test.go` that:

1. Create multiple tags
2. Create multiple monitors
3. Associate tags with monitors
4. Verify tag appears on monitor with correct name and color (from tag definition)
5. Update tag value for specific monitor
6. Retrieve monitor and verify tags are included
7. Retrieve all monitors with a specific tag
8. Remove tag from monitor
9. Verify tag deletion removes all monitor associations
10. Verify monitor deletion removes all tag associations
11. Test with monitors in groups
12. Test tag display on status pages

Integration tests must run against a real Uptime Kuma instance in Docker (following existing test patterns).

### Edge Cases to Test

- Adding tag to non-existent monitor
- Adding non-existent tag to monitor
- Removing tag that isn't associated
- Very long tag values
- Special characters in tag values
- Concurrent tag additions to same monitor

## Implementation

### Phase 1: Data Structures

- Add `MonitorTag` struct to `tag` package
- Add `TagWithMonitors` and `MonitorTags` convenience types
- Extend `monitor.Base` to include `Tags` field
- Add JSON marshaling/unmarshaling tests

### Phase 2: Validation Rules Research

- Review server implementation to find tag value length limits
- Document validation rules for implementation
- Implement client-side validation matching server rules

### Phase 3: Basic Association Operations

- Implement `AddMonitorTag` method
- Implement `DeleteMonitorTag` method
- Add tag value length validation
- Add unit tests for add/delete operations
- Handle socket.io event responses

### Phase 4: Edit Operations

- Implement `UpdateMonitorTag` method
- Apply same validation for value length
- Add tests for edit operations
- Test value updates
- Test validation error cases

### Phase 5: Query Operations

- Implement `GetMonitorTags` method using cached state
- Implement `GetTagMonitors` method using cached state
- Add tests for query operations
- Leverage `c.state.monitors` cache which includes tag information from `monitorList` event

### Phase 6: Integration with Existing Code

- Update monitor retrieval to include tags
- Update tag retrieval to include monitor count
- Ensure tag deletions cascade properly
- Ensure monitor deletions remove tag associations

### Phase 7: Integration Testing

- Add comprehensive integration tests
- Test against real Uptime Kuma instance
- Test validation rules match server behavior
- Test interaction with monitor groups
- Test display on status pages
- Document usage examples

### Phase 8: Documentation

- Add godoc comments for all public APIs
- Create usage examples for common scenarios
- Update main README with tag association examples
- Document best practices for tag organization

## Reference Implementation

Refer to the following Uptime Kuma source files:

- `.scratch/uptime-kuma/server/server.js` - Tag association socket handlers
- `.scratch/uptime-kuma/server/model/tag.js` - Tag and monitor_tag data models
- `.scratch/uptime-kuma/src/components/TagsManager.vue` - Frontend tag management
- `.scratch/uptime-kuma/db/knex_migrations/*.js` - Database schema for monitor_tag table

## Decisions

1. **Bulk Operations**: Not needed at this stage. Single tag operations are sufficient.
2. **Tag Color Handling**: Colors are always inherited from the tag definition. There is no per-association color override capability in Uptime Kuma. The `MonitorTag` struct includes the `Color` field for convenience (populated from the joined tag table), but it is read-only and cannot be modified at the association level.
3. **Tag Value Length Limits**: The server uses TEXT data type with no explicit length validation. The client will not enforce artificial limits. Very long values are allowed per server implementation.
4. **Helper Methods for Filtering**: No helper methods for filtering monitors by tag are required at this stage.
5. **Client-Side Caching**: Query operations (`GetMonitorTags`, `GetTagMonitors`) use the client's cached state (`c.state.monitors`) instead of making API calls. The server's `monitorList` event already includes complete tag information for all monitors via efficient database JOINs. Using cached state provides O(1) performance and eliminates N+1 query patterns. The cache is automatically kept synchronized via socket.io events (`monitorList`, `updateMonitorIntoList`, `deleteMonitorFromList`).
6. **Copy Tags Between Monitors**: No methods needed to copy tags from one monitor to another.
7. **Duplicate Tag Handling**: The server allows the same tag to be associated with a monitor multiple times with different values. The combination of (tag_id, monitor_id, value) forms a unique association. The client allows this behavior.
8. **Delete Methods**: Two delete methods are provided:
   - `DeleteMonitorTag(tagID, monitorID)`: Removes ALL associations between the tag and monitor (regardless of value). This is implemented by fetching all tag associations for the monitor and deleting each one.
   - `DeleteMonitorTagWithValue(tagID, monitorID, value)`: Removes only the specific association with the given value. This maps directly to the server's `deleteMonitorTag` event which requires the value parameter in its WHERE clause.
