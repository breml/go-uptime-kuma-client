package notification

import (
	"fmt"
)

// Apprise represents an Apprise notification provider.
// Apprise is a versatile notification library that supports many notification services.
type Apprise struct {
	Base
	AppriseDetails
}

// AppriseDetails contains the Apprise-specific configuration.
type AppriseDetails struct {
	AppriseURL string `json:"appriseURL"`
	Title      string `json:"title"`
}

// Type returns the notification type identifier for Apprise.
func (a Apprise) Type() string {
	return a.AppriseDetails.Type()
}

// Type returns the notification type identifier for Apprise details.
func (d AppriseDetails) Type() string {
	return "apprise"
}

// String returns a human-readable representation of the Apprise notification.
func (a Apprise) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AppriseDetails, true))
}

// UnmarshalJSON deserializes an Apprise notification from JSON.
func (a *Apprise) UnmarshalJSON(data []byte) error {
	detail := AppriseDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*a = Apprise{
		Base:           base,
		AppriseDetails: detail,
	}

	return nil
}

// MarshalJSON serializes an Apprise notification to JSON.
func (a Apprise) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, a.AppriseDetails)
}
