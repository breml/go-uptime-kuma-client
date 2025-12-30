package notification

import (
	"fmt"
)

// Notifery represents a Notifery notification provider.
// Notifery is a simple API-based notification service.
type Notifery struct {
	Base
	NotiferyDetails
}

// NotiferyDetails contains the configuration fields for Notifery notifications.
type NotiferyDetails struct {
	// APIKey is the API key for authentication with Notifery service.
	APIKey string `json:"notiferyApiKey"`
	// Title is the title of the notification (default: "Uptime Kuma Alert").
	Title string `json:"notiferyTitle"`
	// Group is the notification group for organizing notifications.
	Group string `json:"notiferyGroup"`
}

// Type returns the notification type identifier for Notifery.
func (n Notifery) Type() string {
	return n.NotiferyDetails.Type()
}

// Type returns the notification type identifier for NotiferyDetails.
func (NotiferyDetails) Type() string {
	return "notifery"
}

// String returns a string representation of the Notifery notification.
func (n Notifery) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NotiferyDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Notifery notification.
func (n *Notifery) UnmarshalJSON(data []byte) error {
	detail := NotiferyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*n = Notifery{
		Base:            base,
		NotiferyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Notifery notification into JSON.
func (n Notifery) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, n.NotiferyDetails)
}
