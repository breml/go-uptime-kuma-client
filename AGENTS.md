# Go Client for Uptime Kuma API

## Build & Commands

- **Run Tests**: `task test`, use `task test-debug` for verbose output
- **Run Linter**: `task lint`
- **Compile code**: `task build`

## Architecture & Structure

- Each entity of Uptime Kuma (Monitor, Notification, etc.) has its own package
- The client it self is in the project's root directory
- The client contains integration tests, which launch a real Uptime Kuma instance run as docker container (see `main_test.go`), expect Uptime Kuma to be running when executing the integration tests.
- **.scratch/uptime-kuma/**: Code of Uptime Kuma itself, copied here for reference
- **.scratch/**: Temporary code for testing ideas, not linted, not tested, not checked into git
- **admin/pep/**: Project Enhancement Proposals (PEP) for changes to the codebase

## Code Style & Conventions

- **Formatting**: `gofumpt` for code formatting
- **Linting**: `golangci-lint` for static code analysis
- **Documentation**: Self-documenting code, avoid inline comments

## Tools & Dependencies

- **Go**: Version 1.24
- **Task**: Task runner for build commands (Taskfile.yml)
