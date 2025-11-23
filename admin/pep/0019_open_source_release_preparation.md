# PEP 0019: Open Source Release Preparation

**PEP Number:** 0019
**Title:** Prepare Repository for Open Source Release
**Status:** Implemented

## Abstract

Prepare the go-uptime-kuma-client repository for public open source release by adding essential documentation, licensing, release infrastructure, and examples following patterns established in other breml repositories.

## Motivation

The go-uptime-kuma-client is functionally complete with working code, comprehensive tests, and CI/CD pipeline. However, it lacks the documentation and community health files necessary for a successful open source release.

A public release enables:

- Community contributions and feedback
- Wider adoption of the Uptime Kuma API client
- Package availability via go get and pkg.go.dev
- Transparency in development and feature roadmap
- Collaboration on bug fixes and feature enhancements

Analysis of reference repositories (breml/rootcerts, breml/bidichk, breml/tfreveal, breml/errchkjson) reveals a consistent pattern of essential files and infrastructure that this repository currently lacks.

## Current State

### What Exists

**Core Files:**

- `go.mod` and `go.sum` - Properly configured Go module
- `doc.go` - Minimal package documentation
- `.gitignore` - Present but empty (0 bytes)
- `Taskfile.yml` - Build automation
- `.golangci.yml` - Linter configuration
- Source code organized in entity-based packages

**CI/CD:**

- `.github/workflows/ci.yml` - Build, test, and lint jobs

**Development Files:**

- `CLAUDE.md` - Internal development documentation
- `AGENTS.md` - Development guide
- `admin/pep/` - 18 Project Enhancement Proposals

### What is Missing

**Critical Files:**

- No README.md in repository root
- No LICENSE file
- Empty .gitignore file

**High Priority:**

- No release workflow (.github/workflows/release.yml)
- No dependabot configuration
- No Go example functions for pkg.go.dev
- No git tags or releases
- Limited package documentation

## Specification

### Phase 1: Critical Pre-Release Items

#### 1. README.md

Create comprehensive README with:

```markdown
# go-uptime-kuma-client

[![Test Status](https://github.com/breml/go-uptime-kuma-client/actions/workflows/ci.yml/badge.svg)](https://github.com/breml/go-uptime-kuma-client/actions)
[![Go Report Card](https://goreportcard.com/report/github.com/breml/go-uptime-kuma-client)](https://goreportcard.com/report/github.com/breml/go-uptime-kuma-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go client library for the Uptime Kuma API using Socket.IO for real-time communication.

## Installation

```bash
go get github.com/breml/go-uptime-kuma-client
```

## Features

- **Monitor Management**: HTTP, TCP, Ping, DNS, Redis, PostgreSQL, gRPC, and more
- **Notification Providers**: Generic, Ntfy, Slack, Teams, and other notification types
- **Tag Management**: Organize monitors with tags
- **Proxy Configuration**: Route monitor requests through HTTP/HTTPS/SOCKS proxies
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Create and manage public status pages
- **Real-time Updates**: Socket.IO-based event system for state synchronization

## Usage

[Basic usage example showing client connection, monitor creation, notification setup]

## Documentation

Full documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/breml/go-uptime-kuma-client)

## API Coverage

This client supports the following Uptime Kuma features:

- ✅ Monitors (all types)
- ✅ Notifications
- ✅ Tags
- ✅ Proxies
- ✅ Maintenance Windows
- ✅ Status Pages

## License

MIT License - see [LICENSE](LICENSE) for details
```

**Content Requirements:**
- Project title and description
- Three badges: Test Status, Go Report Card, License
- Installation section with go get command
- Features list highlighting main capabilities
- Basic usage example
- Link to pkg.go.dev documentation
- API coverage status
- License section

#### 2. LICENSE File

Create MIT License file:

```plaintext
MIT License

Copyright (c) 2025 Lucas Bremgartner

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

[Standard MIT License text]
```

**Rationale:** MIT License chosen for consistency with other breml repositories (bidichk, tfreveal, errchkjson).

#### 3. Populate .gitignore

Add standard Go patterns:

```gitignore
# Build artifacts
*.exe
*.dll
*.so
*.dylib

# Test coverage
*.out
coverage.html

# IDE files
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
```

### Phase 2: High Priority Items

#### 4. Add .github/dependabot.yml

Enable automated dependency updates:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

**Rationale:** All reference repositories use Dependabot for weekly dependency monitoring.

#### 5. Add Go Example Functions

Create example_test.go files demonstrating common usage:

**File: `example_test.go`**

```go
package kuma_test

import (
    "context"
    "log"

    kuma "github.com/breml/go-uptime-kuma-client"
    "github.com/breml/go-uptime-kuma-client/monitor"
)

// Example demonstrates basic client connection and usage
func Example() {
    client := kuma.New("http://localhost:3001")
    defer client.Disconnect()

    ctx := context.Background()
    if err := client.Connect(ctx, "admin", "password"); err != nil {
        log.Fatal(err)
    }

    // Client is now connected and ready to use
}

// Example_httpMonitor shows how to create an HTTP monitor
func Example_httpMonitor() {
    client := kuma.New("http://localhost:3001")
    defer client.Disconnect()

    ctx := context.Background()
    client.Connect(ctx, "admin", "password")

    httpMonitor := &monitor.HTTP{
        Base: monitor.Base{
            Name:     "Example Website",
            Interval: 60,
        },
        URL:    "https://example.com",
        Method: "GET",
    }

    id, err := client.CreateMonitor(ctx, httpMonitor)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created monitor with ID: %d", id)
}

// Example_notification shows how to create a notification provider
func Example_notification() {
    client := kuma.New("http://localhost:3001")
    defer client.Disconnect()

    ctx := context.Background()
    client.Connect(ctx, "admin", "password")

    // Example showing notification creation
    // [Implementation details]
}
```

**Additional Examples:**

- `example_monitor_test.go` - Various monitor types (TCP, DNS, Ping)
- `example_tag_test.go` - Tag creation and association
- `example_proxy_test.go` - Proxy configuration
- `example_statuspage_test.go` - Status page management

**Purpose:** Example functions automatically appear in pkg.go.dev documentation, providing immediate value to users.

#### 6. Enhance doc.go

Expand package-level documentation:

```go
// Package kuma provides a Go client library for the Uptime Kuma API.
//
// This library enables programmatic interaction with Uptime Kuma instances
// using Socket.IO for real-time communication. It supports comprehensive
// management of monitors, notifications, tags, proxies, maintenance windows,
// and status pages.
//
// # Architecture
//
// The library uses an entity-based package structure:
//   - monitor/   - Monitor types (HTTP, TCP, DNS, gRPC, etc.)
//   - notification/ - Notification provider types
//   - tag/       - Tag management
//   - proxy/     - Proxy configuration
//   - maintenance/ - Maintenance windows
//   - statuspage/ - Public status pages
//
// The Client type maintains a local state cache synchronized via Socket.IO
// events, ensuring consistency with the Uptime Kuma server.
//
// # Basic Usage
//
//     client := kuma.New("http://localhost:3001")
//     defer client.Disconnect()
//
//     ctx := context.Background()
//     err := client.Connect(ctx, "username", "password")
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     // Create a monitor
//     monitor := &monitor.HTTP{
//         Base: monitor.Base{Name: "Example", Interval: 60},
//         URL: "https://example.com",
//     }
//     id, err := client.CreateMonitor(ctx, monitor)
//
// # Supported Monitor Types
//
// HTTP, TCP, Ping, DNS, gRPC, Redis, PostgreSQL, Real Browser, and more.
// Each monitor type has its own struct with type-specific fields.
//
// # Type Conversion
//
// The library provides base types (monitor.Base, notification.Base) that
// can be converted to specific types using the .As(target) method:
//
//     base, _ := client.GetMonitor(ctx, id)
//     var httpMon monitor.HTTP
//     base.As(&httpMon)
//
// For detailed documentation and examples, see the individual package documentation
// and https://pkg.go.dev/github.com/breml/go-uptime-kuma-client
package kuma
```

**Improvements:**

- Comprehensive overview of capabilities
- Architecture explanation
- Basic usage example
- Explanation of type conversion pattern
- Link to detailed documentation

#### 7. Create Initial Release

**Steps:**

1. Ensure all previous items are completed
2. Run full test suite to verify everything works
3. Create git tag: `git tag v0.1.0`
4. Push tag: `git push origin v0.1.0`
5. GitHub Actions will automatically create release
6. Edit release notes to summarize features

**Release Notes Template:**

```markdown
# v0.1.0 - Initial Public Release

## Features

This is the initial public release of go-uptime-kuma-client, a Go client library for the Uptime Kuma API.

### Supported Entities

- **Monitors**: HTTP, TCP, Ping, DNS, Redis, PostgreSQL, gRPC, Real Browser, and more
- **Notifications**: Generic, Ntfy, Slack, Teams, and other providers
- **Tags**: Monitor organization with tags
- **Proxies**: HTTP/HTTPS/SOCKS proxy support
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Public status page management

### Architecture

- Socket.IO-based real-time communication
- Entity-based package structure
- Local state caching with event-driven updates
- Type-safe monitor and notification handling

## Installation

```bash
go get github.com/breml/go-uptime-kuma-client@v0.1.0
```

## Documentation

Full documentation: https://pkg.go.dev/github.com/breml/go-uptime-kuma-client@v0.1.0

## Backward Compatibility

Not applicable - this is the initial public release. Future releases will follow semantic versioning.

## Test Plan

### Pre-Release Verification

Before creating the v0.1.0 tag, verify:

1. **All tests pass**: `task test` completes successfully
2. **Linter passes**: `task lint` reports no issues
3. **Build succeeds**: `task build` completes without errors
4. **Examples compile**: All example functions compile and are valid Go code
5. **Documentation renders**: Check that pkg.go.dev will render correctly (test with `go doc`)

### Post-Release Verification

After creating the release:

1. Verify GitHub Release is created with correct notes
2. Confirm pkg.go.dev updates (may take a few minutes)
3. Test `go get github.com/breml/go-uptime-kuma-client@v0.1.0` from a clean environment
4. Verify README badges display correctly
5. Check that all links in README work

## Implementation

### Step 1: Create README.md

- Write comprehensive README following reference repo pattern
- Include badges (test status, go report card, license)
- Add installation instructions
- Provide feature overview
- Include basic usage example
- Link to pkg.go.dev

### Step 2: Add LICENSE

- Create MIT License file
- Copyright year: 2025
- Copyright holder: Lucas Bremgartner

### Step 3: Populate .gitignore

- Add build artifacts patterns
- Add test coverage patterns
- Add IDE and OS file patterns

### Step 4: Add Dependabot Configuration

- Create .github/dependabot.yml
- Configure gomod ecosystem monitoring
- Set weekly update schedule

### Step 5: Add Example Functions

- Create example_test.go with basic examples
- Add Example() for client connection
- Add Example_httpMonitor() for monitor creation
- Add Example_notification() for notification setup
- Consider additional examples for tags, proxies, status pages

### Step 6: Enhance Package Documentation

- Update doc.go with comprehensive overview
- Explain architecture and package structure
- Include usage examples
- Document type conversion pattern
- Link to detailed documentation

### Step 7: Create Initial Release

- Verify all tests pass
- Create git tag v0.1.0
- Push tag to trigger release workflow
- Edit GitHub Release notes
- Verify pkg.go.dev updates

## Reference

Analysis based on the following reference repositories:

- <https://github.com/breml/rootcerts>
- <https://github.com/breml/bidichk>
- <https://github.com/breml/tfreveal>
- <https://github.com/breml/errchkjson>

These repositories demonstrate consistent patterns:

- README with badges and examples
- MIT License
- CI/CD workflows
- Dependabot configuration
- Release automation
- No CONTRIBUTING, CODE_OF_CONDUCT, or SECURITY files
- No CHANGELOG.md (use GitHub Releases)
- No issue/PR templates

## Decisions

1. **License Choice**: Use MIT License for consistency with other breml repositories (bidichk, tfreveal, errchkjson).

2. **GoReleaser**: Skip .goreleaser.yml as this is a library, not a CLI tool. Reference repos include it but it's not essential for a pure library.

3. **Community Health Files**: Skip CONTRIBUTING.md, CODE_OF_CONDUCT.md, and SECURITY.md as they're not present in any reference repositories. Keep the release minimal and focused.

4. **Changelog**: Use GitHub Releases for changelog functionality instead of maintaining a separate CHANGELOG.md file (consistent with reference repos).

5. **Issue/PR Templates**: Not included in reference repos; skip for initial release. Can be added later if community contributions increase.

6. **TODO.md**: Already excluded by global .gitignore. No action needed.

7. **Version Number**: Start at v0.1.0 to indicate early stage and allow for API refinement before 1.0.

8. **Badge Services**: Use standard badges (GitHub Actions, Go Report Card, License badge) consistent with reference repos.

9. **Documentation Strategy**: Prioritize pkg.go.dev documentation through godoc comments and example functions over separate documentation site.

10. **Examples Location**: Use example test functions (Example*, ExampleFoo*) in test files rather than separate examples/ directory. This ensures examples are tested and appear in pkg.go.dev automatically.
