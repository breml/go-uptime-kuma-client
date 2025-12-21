package notification

import (
	"fmt"
)

// Stackfield represents a Stackfield notification provider.
// Stackfield is a team collaboration platform for sending alerts via webhooks.
type Stackfield struct {
	Base
	StackfieldDetails
}

// StackfieldDetails contains the configuration fields for Stackfield notifications.
type StackfieldDetails struct {
	// WebhookURL is the Stackfield webhook URL for sending notifications.
	WebhookURL string `json:"stackfieldwebhookURL"`
}

// Type returns the notification type identifier for Stackfield.
func (s Stackfield) Type() string {
	return s.StackfieldDetails.Type()
}

// Type returns the notification type identifier for StackfieldDetails.
func (n StackfieldDetails) Type() string {
	return "stackfield"
}

// String returns a string representation of the Stackfield notification.
func (s Stackfield) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.StackfieldDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Stackfield notification.
func (s *Stackfield) UnmarshalJSON(data []byte) error {
	detail := StackfieldDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Stackfield{
		Base:              base,
		StackfieldDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Stackfield notification into JSON.
func (s Stackfield) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.StackfieldDetails)
}
