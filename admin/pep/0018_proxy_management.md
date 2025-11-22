# PEP 0018: Proxy Management Module

**PEP Number:** 0018  
**Title:** Add Proxy Management Module  
**Status:** Ready

## Abstract

Add a new `proxy` package to provide management of HTTP/HTTPS/SOCKS proxies used by monitors. This allows monitors to perform checks through proxy servers, essential for monitoring internal services or services behind firewalls.

## Motivation

Uptime Kuma supports routing monitor requests through proxy servers, which is essential for:

- Monitoring services behind corporate firewalls
- Testing services from different geographic locations
- Monitoring internal services without exposing them publicly
- Load distribution across multiple proxy servers
- Privacy and anonymity for monitoring requests

Proxies in Uptime Kuma can be:

- HTTP/HTTPS proxies
- SOCKS4/SOCKS5 proxies
- Authenticated proxies (username/password)
- Associated with specific monitors

## Specification

### New Package Structure

Create a new package `proxy/` with the following structure:

```none
proxy/
├── proxy.go           # Core proxy entity
├── proxy_test.go      # Unit tests
└── helpers.go         # Helper functions
```

### Core Data Structures

#### Proxy Entity

```go
type Proxy struct {
    ID       int64  `json:"id"`
    UserID   int64  `json:"userId"`
    Protocol string `json:"protocol"` // "http", "https", "socks", "socks5", "socks5h", "socks4"
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Auth     bool   `json:"auth"`
    Username string `json:"username"`
    Password string `json:"password"`
    Active   bool   `json:"active"`
    Default  bool   `json:"default"`  // Only one proxy can be default (server enforces)

    // Derived/computed fields
    CreatedDate time.Time `json:"createdDate"`
}

// ProxyConfig is used when creating/updating proxies
type ProxyConfig struct {
    Protocol     string `json:"protocol"`
    Host         string `json:"host"`
    Port         int    `json:"port"`
    Auth         bool   `json:"auth"`
    Username     string `json:"username,omitempty"`
    Password     string `json:"password,omitempty"`
    Active       bool   `json:"active"`
    Default      bool   `json:"default"`
    ApplyExisting bool  `json:"applyExisting"` // Apply this proxy to all existing monitors
}
```

### Client Methods

Add the following methods to the main `Client` type in `proxy.go`:

```go
// Get all proxies
func (c *Client) GetProxyList(ctx context.Context) []proxy.Proxy

// Get specific proxy (query from cached list)
func (c *Client) GetProxy(ctx context.Context, id int64) (*proxy.Proxy, error)

// Create new proxy
func (c *Client) CreateProxy(ctx context.Context, config proxy.Config) (int64, error)

// Update existing proxy
func (c *Client) UpdateProxy(ctx context.Context, config proxy.Config) error

// Delete proxy (automatically removes proxy from monitors)
func (c *Client) DeleteProxy(ctx context.Context, id int64) error
```

### Socket.io Events

Map to the following Uptime Kuma socket.io events:

- `proxyList` - Event sent on connection containing all proxies
- `addProxy(proxy, proxyID, callback)` - Create new proxy (proxyID = null) or update existing (proxyID set)
- `deleteProxy(proxyID, callback)` - Delete proxy

**Note:** The `addProxy` event handles both creation and updates. When `proxyID` is null, it creates a new proxy. When `proxyID` is provided, it updates the existing proxy with that ID.

### State Management

Add proxy list to the client state:

```go
type state struct {
    notifications []notification.Base
    monitors      []monitor.Base
    statusPages   []statuspage.StatusPage
    maintenance   []maintenance.Maintenance
    proxies       []proxy.Proxy  // New field
}
```

Listen for `proxyList` event during client initialization to populate the cache.

### Integration with Monitors

The `monitor.Base` struct should support proxy configuration:

```go
type Base struct {
    // ... existing fields ...
    ProxyID *int64 `json:"proxyId"` // Pointer to allow null (no proxy)
}
```

## Backward Compatibility

Not a concern for this early-stage project. Breaking changes are acceptable if they improve the overall API design.

Adding `ProxyID` to `monitor.Base` should be backward compatible as it will be nil for monitors without proxies.

## Test Plan

### Unit Tests

- Test proxy data structure parsing
- Test protocol validation ("http", "https", "socks", "socks5", "socks5h", "socks4")
- Test port validation (1-65535)
- Test authentication flag with username/password
- Test default proxy flag (only one can be default)
- Test host validation (IP address or hostname)
- Test ApplyExisting flag
- Test empty/nil value handling

### Integration Tests

Add integration tests in `main_test.go` that:

1. Create HTTP proxy without authentication
2. Create SOCKS5 proxy with authentication
3. Retrieve proxy list and verify created proxies
4. Update existing proxy configuration (using SaveProxy with proxyID)
5. Create monitor with proxy configured
6. Verify monitor uses proxy for checks
7. Set proxy as default (verify old default is unset)
8. Create proxy with ApplyExisting=true, verify all monitors updated
9. Delete proxy
10. Verify monitors using deleted proxy have proxy_id set to null
11. Test multiple proxies
12. Test active/inactive proxy states
13. Verify only one proxy can be default at a time

Integration tests must run against a real Uptime Kuma instance in Docker (following existing test patterns).

Note: Testing actual proxy functionality may require setting up a test proxy server.

### Edge Cases to Test

- Creating proxy with invalid protocol (not in supported list)
- Creating proxy with invalid port
- Creating proxy with auth enabled but no credentials
- Deleting proxy that's in use by monitors (server sets proxy_id to null)
- Setting multiple proxies as default (server ensures only one is default)
- Very long hostnames
- IPv6 addresses
- Updating non-existent proxy

## Implementation

### Phase 1: Core Proxy Structures

- Create `proxy` package structure
- Implement `Proxy` and `ProxyConfig` structs
- Add JSON marshaling/unmarshaling tests
- Add validation functions

### Phase 2: Save/Delete Operations

- Implement `GetProxyList` method
- Implement `SaveProxy` method (handles both create and update via proxyID parameter)
- Implement `DeleteProxy` method
- Add protocol validation (match server's supported list)
- Add unit tests for save/delete operations

### Phase 3: Proxy State Management

- Add `proxies` field to client state
- Listen for `proxyList` event
- Update cache on proxy changes
- Add cache synchronization tests

### Phase 4: Monitor Integration

- Add `ProxyID` field to `monitor.Base`
- Update monitor creation/editing to support proxy
- Test monitors with and without proxies
- Test ApplyExisting functionality
- Document proxy usage in monitors

### Phase 5: Default Proxy Handling

- Implement default proxy logic (server enforces only one default)
- Test switching default between proxies
- Verify old default is cleared when new one is set
- Add tests for default proxy behavior

### Phase 6: Integration Testing

- Add comprehensive integration tests
- Test against real Uptime Kuma instance
- Test with actual proxy server if possible
- Document usage examples

### Phase 7: Documentation

- Add godoc comments for all public APIs
- Create usage examples for common scenarios
- Update main README with proxy examples
- Document supported proxy types and protocols

## Reference Implementation

Refer to the following Uptime Kuma source files:

- `.scratch/uptime-kuma/server/socket-handlers/proxy-socket-handler.js` - Server-side socket handlers
- `.scratch/uptime-kuma/server/model/proxy.js` - Proxy data model
- `.scratch/uptime-kuma/src/components/ProxyDialog.vue` - Frontend proxy configuration
- `.scratch/uptime-kuma/server/proxy.js` - Proxy connection logic

## Decisions

1. **Editing Proxies**: The server uses the same `addProxy` event for both creating and updating. When `proxyID` is provided, it updates the existing proxy. The client will expose this via `SaveProxy(config, proxyID)` method.

2. **Testing Proxy Connectivity**: The server does not provide a dedicated proxy testing endpoint. Since this is a client package, we rely on the server to validate proxy configuration. If the server accepts the proxy, it's considered valid from the client's perspective.

3. **Proxy Deletion Impact**: When a proxy is deleted, the server automatically sets `proxy_id = null` on all monitors using that proxy (see `proxy.js` line 81). The client does not need to handle this - the server manages it correctly.

4. **Helper Methods**: No helper methods for common proxy configurations are needed.

5. **Proxy Connectivity Validation**: No client-side connectivity validation. The server handles validation when proxies are used by monitors.

6. **Proxy Credentials Security**: No special handling needed in the client. Rely on the server to handle credentials securely. The client passes credentials as provided by the user.

7. **PAC File Support**: The server does not support proxy auto-configuration (PAC) files. Therefore, the client should not provide PAC file support. The client only provides functionality available in the Uptime Kuma server.

8. **Default Proxy Exclusivity**: Only one proxy can be set as default at a time. The server enforces this by clearing the default flag on all other proxies when a new default is set (see `proxy.js` line 44).
