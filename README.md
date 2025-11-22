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

- **Monitor Management**: HTTP, TCP, Ping, DNS, Redis, PostgreSQL, gRPC, Real Browser, and more
- **Notification Providers**: Ntfy, Slack, Teams, Generic, and other notification types
- **Tag Management**: Organize monitors with tags
- **Proxy Configuration**: Route monitor requests through HTTP/HTTPS/SOCKS proxies
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Create and manage public status pages
- **Real-time Updates**: Socket.IO-based event system for state synchronization

## Usage

Here's a complete example showing how to connect to Uptime Kuma, create a notification provider, and set up an HTTP monitor with notifications:

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
