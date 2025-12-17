package notification

import (
	"fmt"
)

// Pumble represents a Pumble notification provider.
// Pumble is a team communication platform that supports webhook integrations for notifications.
type Pumble struct {
	Base
	PumbleDetails
}

// PumbleDetails contains the Pumble-specific configuration.
type PumbleDetails struct {
	WebhookURL string `json:"webhookURL"`
}

// Type returns the notification type identifier for Pumble.
func (p Pumble) Type() string {
	return p.PumbleDetails.Type()
}

// Type returns the notification type identifier for Pumble details.
func (p PumbleDetails) Type() string {
	return "Pumble"
}

// String returns a human-readable representation of the Pumble notification.
func (p Pumble) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PumbleDetails, true))
}

// UnmarshalJSON deserializes a Pumble notification from JSON.
func (p *Pumble) UnmarshalJSON(data []byte) error {
	detail := PumbleDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Pumble{
		Base:          base,
		PumbleDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a Pumble notification to JSON.
func (p Pumble) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PumbleDetails)
}
