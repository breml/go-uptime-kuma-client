# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Commands

- **Run Tests**: `task test`, use `task test-debug` for verbose output
- **Run Tests for a Single Package**: `go test ./path/to/package`
- **Run Single Test**: `go test -run TestName ./path/to/package`
- **Run Linter**: `task lint`
- **Compile code**: `task build`

## Architecture & Structure

### Core Architecture

This is a Go client library for the Uptime Kuma API using Socket.IO for real-time communication.

**Client Design (client.go:52-60)**:

- `Client` struct wraps a Socket.IO client connection
- Maintains local state cache for monitors and notifications
- Uses signal-based event system for state updates
- Protects state with mutex for concurrent access
- Supports autosetup mode for first-time server initialization

**Entity-Based Package Structure**:

- Each Uptime Kuma entity (Monitor, Notification) has its own package under root
- Root-level files (`monitor.go`, `notification.go`) contain client methods
- Entity packages contain type definitions and implementations
- Example: `monitor/` contains `monitor_base.go`, `monitor_http.go`, `monitor_group.go`

### Key Patterns

**Type Conversion Pattern**:

- Base types (`monitor.Base`, `notification.Base`) store raw JSON for flexible conversion
- Use `.As(target)` method to convert base types to specific types
- Example: `monitor.Base` can be converted to `monitor.HTTP` or `monitor.Group`
- Client methods: `GetMonitorAs()` and `GetNotificationAs()` for type-safe retrieval

**Event-Driven Updates**:

- `syncEmitWithUpdateEvent()` waits for both ack response AND update event
- Uses UUID-based listeners for tracking specific operations
- Ensures state consistency before returning from operations

**NotificationIDList Handling**:

- Server expects `map[string]bool` format for notification IDs
- Client uses `[]int64` for NotificationIDs field
- Conversion happens in monitor CRUD operations (monitor.go:71-75, 98-102)

### Integration Tests

- Integration tests launch real Uptime Kuma instance via Docker (main_test.go:21-92)
- **IMPORTANT**: Assume Uptime Kuma is running when executing integration tests
- Uses dockertest library with automatic container cleanup
- Global `client` variable shared across test suite
- Container expires after 60 seconds
- Test credentials: username="admin", password="admin1"

### Directory Structure

- `.scratch/uptime-kuma/`: Copy of Uptime Kuma source code for reference
- `.scratch/`: Temporary code for experiments (not linted, tested, or committed)

## Code Style & Conventions

- **Formatting**: Use `gofumpt` (stricter than `gofmt`)
- **Linting**: `golangci-lint` enforces style
- **Documentation**: Self-documenting code preferred, minimal inline comments
- **Go Version**: 1.24

## Working with Monitors

Monitor types include `http`, `group`, and others. Each has:

- Base fields in `monitor.Base` (ID, Name, Interval, etc.)
- Type-specific fields in dedicated structs (e.g., `monitor.HTTP`)
- Parent field for hierarchical organization (groups can contain monitors)

## Working with Notifications

Notification types include `Generic`, `Ntfy`, `Slack`, `Teams`, etc.:

- Base fields in `notification.Base` (ID, Name, IsActive, etc.)
- Type-specific configuration stored as JSON string internally
- Configuration details unmarshaled to type-specific structs
