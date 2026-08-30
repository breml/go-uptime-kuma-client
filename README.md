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

- **Monitor Management**: HTTP, TCP, Ping, DNS, Redis, PostgreSQL, gRPC, Real Browser,
  and more
- **Notification Providers**: Ntfy, Slack, Teams, Generic, and other notification
  types
- **Tag Management**: Organize monitors with tags
- **Proxy Configuration**: Route monitor requests through HTTP/HTTPS/SOCKS proxies
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Create and manage public status pages
- **Real-time Updates**: Socket.IO-based event system for state synchronization
- **Authentication**: Username and password, two-factor authentication without a
  human at the keyboard, session token reuse, and servers with authentication
  disabled

## Usage

Here's a complete example showing how to connect to Uptime Kuma, create a notification
provider, and set up an HTTP monitor with notifications:

```go
package main

import (
 "context"
 "log"

 kuma "github.com/breml/go-uptime-kuma-client"
 "github.com/breml/go-uptime-kuma-client/monitor"
 "github.com/breml/go-uptime-kuma-client/notification"
)

func main() {
 ctx := context.Background()

 // Create client and connect to Uptime Kuma
 client, err := kuma.New(ctx, "http://localhost:3001", "admin", "password")
 if err != nil {
  log.Fatal(err)
 }
 defer client.Disconnect()

 // Create a notification provider
 ntfyNotification := notification.Ntfy{
  Base: notification.Base{
   Name:     "My Ntfy Alert",
   IsActive: true,
  },
  NtfyDetails: notification.NtfyDetails{
   ServerURL:            "https://ntfy.sh",
   Topic:                "uptime-alerts",
   Priority:             5,
   AuthenticationMethod: "none",
  },
 }

 notificationID, err := client.CreateNotification(ctx, ntfyNotification)
 if err != nil {
  log.Fatal(err)
 }
 log.Printf("Created notification with ID: %d", notificationID)

 // Create an HTTP monitor that uses the notification
 httpMonitor := &monitor.HTTP{
  Base: monitor.Base{
   Name:            "Example Website",
   Interval:        60,
   NotificationIDs: []int64{notificationID},
  },
  HTTPDetails: monitor.HTTPDetails{
   URL:                 "https://example.com",
   Method:              "GET",
   AcceptedStatusCodes: []string{"200-299"},
  },
 }

 monitorID, err := client.CreateMonitor(ctx, httpMonitor)
 if err != nil {
  log.Fatal(err)
 }
 log.Printf("Created monitor with ID: %d", monitorID)
}
```

## Authentication

A username and password is the usual way in. An account with two-factor
authentication enabled additionally needs the shared secret, which the client
turns into the one-time code the server asks for:

```go
client, err := kuma.New(ctx, url, "username", "password",
 kuma.WithTOTPSecret(secret))
```

`WithTOTPCode` takes a callback instead, for a secret the client cannot read
itself.

A successful login answers with a session token. Persisting it lets a later
client connect with neither the password nor a one-time code:

```go
token := client.SessionToken()

// ... later, in another process
client, err := kuma.New(ctx, url, "", "", kuma.WithSessionToken(token))
```

This is also how several clients start at once against an account with
two-factor authentication: the server refuses a one-time code it has already
accepted, but takes the same token any number of times.

The token is a bearer credential that does not expire, and for an account with
two-factor authentication it is a stronger one than the password. Store it the
way the password would be stored; `WithSessionToken` documents what invalidates
it.

A server with authentication disabled logs the client in by itself, so it takes
no credentials at all:

```go
client, err := kuma.New(ctx, url, "", "")
```

As of Uptime Kuma 2.5.0, these are the only ways in. API keys are not among
them: the server accepts those for HTTP basic auth on `/metrics` only, never for
the Socket.IO API this client speaks. There is likewise no OIDC, SSO, LDAP or
client-certificate login.

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

This is a work in progress, more Uptime Kuma features might be added in the future.

## License

MIT License - see [LICENSE](LICENSE) for details
