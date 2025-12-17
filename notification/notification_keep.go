package notification

import (
	"fmt"
)

// Keep represents a Keep (Google) notification provider.
// Keep is a Google service for capturing and organizing notes and alerts.
type Keep struct {
	Base
	KeepDetails
}

// KeepDetails contains the configuration fields for Keep notifications.
type KeepDetails struct {
	// WebhookURL is the Keep webhook URL.
	WebhookURL string `json:"webhookURL"`
	// APIKey is the Keep API key for authentication.
	APIKey string `json:"webhookAPIKey"`
}

// Type returns the notification type identifier for Keep.
func (k Keep) Type() string {
	return k.KeepDetails.Type()
}

// Type returns the notification type identifier for KeepDetails.
func (n KeepDetails) Type() string {
	return "Keep"
}

// String returns a string representation of the Keep notification.
func (k Keep) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(k.Base, false), formatNotification(k.KeepDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Keep notification.
func (k *Keep) UnmarshalJSON(data []byte) error {
	detail := KeepDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*k = Keep{
		Base:        base,
		KeepDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Keep notification into JSON.
func (k Keep) MarshalJSON() ([]byte, error) {
	return marshalJSON(k.Base, k.KeepDetails)
}
