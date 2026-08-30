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
//	monitor := &monitor.HTTP{
//	    Base: monitor.Base{Name: "Example", Interval: 60},
//	    HTTPDetails: monitor.HTTPDetails{
//	        URL:    "https://example.com",
//	        Method: "GET",
//	    },
//	}
//	id, err := client.CreateMonitor(ctx, monitor)
//
// # Authentication
//
// A username and password is the usual way in. An account with two-factor
// authentication enabled additionally needs the shared secret, which the
// client turns into the one-time code the server asks for:
//
//	client, err := kuma.New(ctx, url, "username", "password",
//	    kuma.WithTOTPSecret(secret))
//
// A successful login answers with a session token. Persisting it lets a later
// client connect with neither the password nor a one-time code, which is also
// how several clients start at once against an account with two-factor
// authentication: the server refuses a one-time code it has already accepted,
// but takes the same token any number of times.
//
//	token := client.SessionToken()
//	// ... later, in another process
//	client, err := kuma.New(ctx, url, "", "", kuma.WithSessionToken(token))
//
// The token is a bearer credential that does not expire. Store it the way the
// password would be stored, see [WithSessionToken] for what invalidates it and
// for the fallback a client with a password gets when it is refused.
//
// A server with authentication disabled logs the client in by itself, so it
// takes no credentials at all:
//
//	client, err := kuma.New(ctx, url, "", "")
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
package kuma
