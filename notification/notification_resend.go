package notification

import (
	"fmt"
)

// Resend represents a Resend notification.
type Resend struct {
	Base
	ResendDetails
}

// ResendDetails contains Resend-specific notification configuration.
type ResendDetails struct {
	APIKey    string `json:"resendApiKey"`
	FromEmail string `json:"resendFromEmail"`
	FromName  string `json:"resendFromName,omitempty"`
	// ToEmail is a comma-separated list of recipient email addresses.
	ToEmail string `json:"resendToEmail"`
	Subject string `json:"resendSubject,omitempty"`
}

// Type returns the notification type.
func (r Resend) Type() string {
	return r.ResendDetails.Type()
}

// Type returns the notification type.
func (ResendDetails) Type() string {
	return "Resend"
}

// String returns a string representation of the notification.
func (r Resend) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(r.Base, false), formatNotification(r.ResendDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (r *Resend) UnmarshalJSON(data []byte) error {
	detail := ResendDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*r = Resend{
		Base:          base,
		ResendDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (r Resend) MarshalJSON() ([]byte, error) {
	return marshalJSON(r.Base, &r.ResendDetails)
}
