# PEP 0015: Status Pages Module

**PEP Number:** 0015  
**Title:** Add Status Pages Module  
**Status:** Ready

## Abstract

Add a new `statuspage` package to provide complete management of Uptime Kuma status pages. This includes creating, editing, deleting status pages, managing incidents, and associating monitors with status pages.

## Motivation

Status pages are a critical feature of Uptime Kuma, allowing users to publicly display the status of their services. Currently, the go-uptime-kuma-client has no support for status page management, limiting its usefulness for automation and infrastructure-as-code workflows.

Status pages enable:

- Public-facing service status dashboards
- Incident management and communication
- Monitor grouping for external visibility
- Customizable branding and appearance

## Specification

### New Package Structure

Create a new package `statuspage/` with the following structure:

```none
statuspage/
├── statuspage.go           # Core status page entity
├── statuspage_test.go      # Unit tests
├── incident.go             # Incident posting/management
├── incident_test.go        # Incident unit tests
└── helpers.go              # Helper functions
```

### Core Data Structures

#### StatusPage Entity

```go
type StatusPage struct {
    ID          int64             `json:"id"`
    Slug        string            `json:"slug"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Theme       string            `json:"theme"`        // "light", "dark", "auto"
    Published   bool              `json:"published"`
    ShowTags    bool              `json:"showTags"`
    DomainNameList []string       `json:"domainNameList"`
    GoogleAnalyticsID string      `json:"googleAnalyticsId"`
    CustomCSS   string            `json:"customCSS"`
    FooterText  string            `json:"footerText"`
    ShowPoweredBy bool            `json:"showPoweredBy"`
    Icon        string            `json:"icon"`
    PublicGroupList []PublicGroup `json:"publicGroupList"`
}

type PublicGroup struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Weight      int     `json:"weight"`
    MonitorList []int64 `json:"monitorList"` // Monitor IDs
}

type Incident struct {
    ID          int64     `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Style       string    `json:"style"`     // "info", "warning", "danger", "primary"
    Created     time.Time `json:"created"`
    LastUpdated time.Time `json:"lastUpdated"`
    Pin         bool      `json:"pin"`
}
```

### Client Methods

Add the following methods to the main `Client` type in `statuspage.go`:

```go
// Status Page Management
func (c *Client) GetStatusPages(ctx context.Context) (map[int64]statuspage.StatusPage, error)
func (c *Client) GetStatusPage(ctx context.Context, slug string) (*statuspage.StatusPage, error)
func (c *Client) AddStatusPage(ctx context.Context, title, slug string) error
func (c *Client) SaveStatusPage(ctx context.Context, sp *statuspage.StatusPage) error
func (c *Client) DeleteStatusPage(ctx context.Context, slug string) error

// Incident Management
func (c *Client) PostIncident(ctx context.Context, slug string, incident *statuspage.Incident) error
func (c *Client) UnpinIncident(ctx context.Context, slug string) error
```

**Note on API Limitations:**
- `GetStatusPage` does not return `PublicGroupList` because the Uptime Kuma server's `getStatusPage` endpoint only returns basic configuration fields.
- `PublicGroupList` must be maintained client-side and provided when calling `SaveStatusPage`.
- The `addStatusPage` socket.io event only accepts `title` and `slug` parameters, not a full StatusPage object.

### Socket.io Events

Map to the following Uptime Kuma socket.io events:

- `getStatusPage(slug)` - Retrieve status page configuration by slug (returns basic config only, not publicGroupList)
- `saveStatusPage(slug, config, imgDataUrl, publicGroupList)` - Save/update status page configuration and public groups
- `addStatusPage(title, slug)` - Create new status page with title and slug
- `deleteStatusPage(slug)` - Delete status page by slug
- `postIncident(slug, incident)` - Post or update incident on status page
- `unpinIncident(slug)` - Unpin currently pinned incident
- `statusPageList` - Event received during initialization containing all status pages (map[int64]StatusPage)

### State Management

Add status page map to the client state:

```go
type state struct {
    notifications []notification.Base
    monitors      []monitor.Base
    statusPages   map[int64]statuspage.StatusPage  // New field
}
```

Listen for `statusPageList` event during client initialization to populate the cache. The server sends status pages as a map indexed by ID.

## Backward Compatibility

Since backward compatibility is not a concern for this early-stage project, breaking changes are acceptable if they improve the overall API design. However, the new status page functionality should follow existing patterns for monitors and notifications to maintain consistency.

## Test Plan

### Unit Tests

- Test CRUD operations for status pages
- Test incident posting and unpinning
- Test validation of required fields (slug, title)
- Test public group and monitor associations
- Test custom domain configuration

### Integration Tests

Add integration tests in `statuspage_test.go` that:

1. Create a status page with monitors
2. Publish and unpublish the status page
3. Post incidents with different styles
4. Pin and unpin incidents
5. Update status page configuration
6. Delete the status page

Integration tests should run against a real Uptime Kuma instance in Docker (following existing test patterns).

### Manual Testing

- Verify status page appears correctly in Uptime Kuma UI
- Verify incident posts display properly
- Verify monitor associations work correctly
- Test with custom domains and SSL

## Implementation

### Phase 1: Core Status Page CRUD

- Create `statuspage` package structure
- Implement base `StatusPage` struct
- Add `GetStatusPage`, `AddStatusPage`, `UpdateStatusPage`, `DeleteStatusPage` methods
- Add unit tests for CRUD operations

### Phase 2: Monitor Associations

- Implement `PublicGroup` struct
- Add monitor association logic
- Test monitor grouping on status pages

### Phase 3: Incident Management

- Implement `Incident` struct
- Add `PostIncident` and `UnpinIncident` methods
- Add incident-specific tests

### Phase 4: Integration Testing

- Add comprehensive integration tests
- Test against real Uptime Kuma instance
- Document usage examples

### Phase 5: Documentation

- Add godoc comments for all public APIs
- Create usage examples in package documentation
- Update main README with status page examples

## Reference Implementation

Refer to the following Uptime Kuma source files:

- `.scratch/uptime-kuma/server/socket-handlers/status-page-socket-handler.js` - Server-side socket handlers
- `.scratch/uptime-kuma/server/model/status_page.js` - Status page data model
- `.scratch/uptime-kuma/src/pages/StatusPage.vue` - Frontend implementation

## Decisions

1. **Customization Options**: Support all customization options from the UI for comprehensive coverage.
2. **Slug Validation**: Rely on server-side validation for slug uniqueness and format. The client may perform basic validation (non-empty string) but should defer to the server for uniqueness checks and detailed format requirements.
3. **Helper Methods**: No helper methods needed for themes and configurations at this stage.
