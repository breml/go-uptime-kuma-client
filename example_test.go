package kuma_test

import (
	"context"
	"fmt"
	"log"

	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/notification"
)

func Example() {
	ctx := context.Background()

	// Using pre-initialized client from main_test.go connecting to Uptime
	// Kuma running in Docker container.
	//
	// client, err := kuma.New(ctx, "http://localhost:3001", "admin", "password")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() { _ = client.Disconnect() }()

	// Create a notification provider.
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

	retrievedNotification := notification.Ntfy{}
	err = client.GetNotificationAs(ctx, notificationID, &retrievedNotification)
	if err != nil {
		log.Fatal(err)
	}

	// Create an HTTP monitor that uses the notification.
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

	retrievedMonitor := monitor.HTTP{}
	err = client.GetMonitorAs(ctx, monitorID, &retrievedMonitor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Notification Topic:", retrievedNotification.Topic)
	fmt.Println("HTTP Monitor URL:", retrievedMonitor.URL)

	// Output:
	// Notification Topic: uptime-alerts
	// HTTP Monitor URL: https://example.com
}
