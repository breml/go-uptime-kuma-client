package notification

import (
	"fmt"
)

// Pushy represents a Pushy notification provider.
// Pushy is a push notification service that requires an API key and device token.
type Pushy struct {
	Base
	PushyDetails
}

// PushyDetails contains the Pushy-specific configuration.
type PushyDetails struct {
	APIKey string `json:"pushyAPIKey"` // Pushy API secret key
	Token  string `json:"pushyToken"`  // Device token
}

// Type returns the notification type identifier for Pushy.
func (p Pushy) Type() string {
	return p.PushyDetails.Type()
}

// Type returns the notification type identifier for Pushy details.
func (p PushyDetails) Type() string {
	return "pushy"
}

// String returns a human-readable representation of the Pushy notification.
func (p Pushy) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushyDetails, true))
}

// UnmarshalJSON deserializes a Pushy notification from JSON.
func (p *Pushy) UnmarshalJSON(data []byte) error {
	detail := PushyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Pushy{
		Base:         base,
		PushyDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a Pushy notification to JSON.
func (p Pushy) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PushyDetails)
}
