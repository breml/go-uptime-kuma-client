package notification

import (
	"fmt"
)

// Bark represents a Bark notification provider.
// Bark is an iOS app for sending push notifications to Apple devices.
type Bark struct {
	Base
	BarkDetails
}

// BarkDetails contains the configuration fields for Bark notifications.
type BarkDetails struct {
	// Endpoint is the Bark server endpoint URL.
	Endpoint string `json:"barkEndpoint"`
	// Group is the notification group name.
	Group string `json:"barkGroup"`
	// Sound is the notification sound.
	Sound string `json:"barkSound"`
	// APIVersion is the API version (v1 or v2).
	APIVersion string `json:"apiVersion"`
}

// Type returns the notification type identifier for Bark.
func (b Bark) Type() string {
	return b.BarkDetails.Type()
}

// Type returns the notification type identifier for BarkDetails.
func (n BarkDetails) Type() string {
	return "bark"
}

// String returns a string representation of the Bark notification.
func (b Bark) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(b.Base, false), formatNotification(b.BarkDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Bark notification.
func (b *Bark) UnmarshalJSON(data []byte) error {
	detail := BarkDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*b = Bark{
		Base:        base,
		BarkDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Bark notification into JSON.
func (b Bark) MarshalJSON() ([]byte, error) {
	return marshalJSON(b.Base, b.BarkDetails)
}
