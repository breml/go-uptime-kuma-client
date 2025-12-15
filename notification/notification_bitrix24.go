package notification

import (
	"fmt"
)

// Bitrix24 represents a Bitrix24 notification provider.
// Bitrix24 is a communication and collaboration platform that supports webhooks for sending notifications.
type Bitrix24 struct {
	Base
	Bitrix24Details
}

// Bitrix24Details contains the configuration fields for Bitrix24 notifications.
type Bitrix24Details struct {
	// WebhookURL is the Bitrix24 webhook URL.
	WebhookURL string `json:"bitrix24WebhookURL"`
	// NotificationUserID is the Bitrix24 user ID to receive the notification.
	NotificationUserID string `json:"bitrix24UserID"`
}

// Type returns the notification type identifier for Bitrix24.
func (b Bitrix24) Type() string {
	return b.Bitrix24Details.Type()
}

// Type returns the notification type identifier for Bitrix24 details.
func (n Bitrix24Details) Type() string {
	return "Bitrix24"
}

// String returns a human-readable representation of the Bitrix24 notification.
func (b Bitrix24) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(b.Base, false), formatNotification(b.Bitrix24Details, true))
}

// UnmarshalJSON unmarshals JSON data into a Bitrix24 notification.
func (b *Bitrix24) UnmarshalJSON(data []byte) error {
	detail := Bitrix24Details{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*b = Bitrix24{
		Base:            base,
		Bitrix24Details: detail,
	}

	return nil
}

// MarshalJSON marshals the Bitrix24 notification into JSON.
func (b Bitrix24) MarshalJSON() ([]byte, error) {
	return marshalJSON(b.Base, b.Bitrix24Details)
}
