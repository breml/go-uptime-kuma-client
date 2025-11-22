// Package kuma provides a Go client library for the Uptime Kuma API.
//
// This library enables programmatic interaction with Uptime Kuma instances
// using Socket.IO for real-time communication. It supports comprehensive
// management of Uptime Kuma resources like monitors, notifications, etc.
//
// # Architecture
//
// The library uses an entity-based package structure:
//   - monitor/      - Monitor types (HTTP, TCP, DNS, gRPC, etc.)
//   - notification/ - Notification provider types
//   - tag/          - Tag management
//   - proxy/        - Proxy configuration
//   - maintenance/  - Maintenance windows
//   - statuspage/   - Public status pages
//
// The Client type maintains a local state cache synchronized via Socket.IO
// events, ensuring consistency with the Uptime Kuma server.
//
// # Basic Usage
//
//	ctx := context.Background()
//	client, err := kuma.New(ctx, "http://localhost:3001", "username", "password")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Disconnect()
//
//	// Create a monitor
//	monitor := &monitor.HTTP{
//	    Base: monitor.Base{Name: "Example", Interval: 60},
//	    HTTPDetails: monitor.HTTPDetails{
//	        URL:    "https://example.com",
//	        Method: "GET",
//	    },
//	}
//	id, err := client.CreateMonitor(ctx, monitor)
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
//	base, _ := client.GetMonitor(ctx, id)
//	var httpMon monitor.HTTP
//	base.As(&httpMon)
//
package kuma
