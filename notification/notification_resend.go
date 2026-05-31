package notification

import (
	"fmt"
)

// Resend represents a Resend notification provider.
// Resend is a transactional email API for sending notifications via email.
type Resend struct {
	Base
	ResendDetails
}

// ResendDetails contains the configuration fields for Resend notifications.
type ResendDetails struct {
	// APIKey is the Resend API key for authentication.
	APIKey string `json:"resendApiKey"`
	// FromEmail is the sender email address.
	FromEmail string `json:"resendFromEmail"`
	// FromName is the sender name.
	FromName string `json:"resendFromName"`
	// ToEmail is a comma-separated list of recipient email addresses.
	ToEmail string `json:"resendToEmail"`
	// Subject is the email subject line.
	Subject string `json:"resendSubject"`
}

// Type returns the notification type identifier for Resend.
func (r Resend) Type() string {
	return r.ResendDetails.Type()
}

// Type returns the notification type identifier for ResendDetails.
func (ResendDetails) Type() string {
	return "Resend"
}

// String returns a string representation of the Resend notification.
func (r Resend) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(r.Base, false), formatNotification(r.ResendDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Resend notification.
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

// MarshalJSON marshals the Resend notification into JSON.
func (r Resend) MarshalJSON() ([]byte, error) {
	return marshalJSON(r.Base, r.ResendDetails)
}
