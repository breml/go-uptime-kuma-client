# PEP 0002: Add Parent Field for Monitors

Add support for the `parent` field in the monitor package to enable hierarchical organization of monitors.

## Background

Uptime Kuma supports organizing monitors into a hierarchical structure where monitors can have a parent monitor (typically a "group" type monitor). This relationship is defined through a `parent` field that references another monitor's ID.

## Implementation Details

### 1. Update Base Monitor Structure

Add `Parent` field to `monitor.Base` struct in @monitor/monitor_base.go :

- **Field**: `Parent *int64` (pointer to allow `null` representation)
- **JSON tag**: `"parent"`
- **Default value**: `nil` (represents no parent/top-level monitor)

### 2. Update JSON Marshaling/Unmarshaling

Modify both `UnmarshalJSON` and `MarshalJSON` methods in `monitor.Base`:

- **UnmarshalJSON**: Extract `parent` field from JSON and handle `null` values
- **MarshalJSON**: Include `parent` field in the marshaled JSON output

### 3. Validation Considerations

Based on Uptime Kuma's implementation (@.scratch/uptime-kuma/server/server.js):

- Prevent circular references: A monitor cannot have a descendant as its parent
- This validation is server-side, but document it for client users

### 4. Special Values

- `nil`: No parent (top-level monitor)
- `-1`: Special value used in UI to indicate "create new parent group" (implementation detail, may not need client support)
- Any other valid monitor ID: References an existing monitor as parent

### 5. Type Considerations

- Any monitor can have a parent, not just specific types
- Group monitors commonly serve as parents but this is not enforced
- Parent monitors are typically of type "group" but this is not required

## Testing Requirements

### Unit Tests

Add tests to @monitor/monitor_base_test.go :

1. **Unmarshal parent field**: Test JSON with `parent: null`, `parent: 123`, etc.
2. **Marshal parent field**: Verify parent field is correctly included in output
3. **Round-trip test**: Unmarshal → Marshal → Unmarshal should preserve parent value

### Integration Tests

Add tests to @monitor_test.go :

1. **Create monitor with parent**: Create child monitor referencing parent ID
2. **Update parent field**: Change a monitor's parent
3. **Remove parent**: Set parent to `null`
4. **Parent-child relationship**: Create group monitor and add child monitors to it

## Reference Implementation

See Uptime Kuma implementation:
- Server-side handling: @.scratch/uptime-kuma/server/server.js (validation logic)
- Frontend defaults: @.scratch/uptime-kuma/src/pages/EditMonitor.vue (default value)
- UI implementation: @.scratch/uptime-kuma/src/pages/EditMonitor.vue (parent selector)
