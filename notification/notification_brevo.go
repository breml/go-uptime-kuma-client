package notification

import (
	"fmt"
)

// Brevo represents a Brevo (formerly Sendinblue) notification provider.
// Brevo is an email service provider for sending notifications via email.
type Brevo struct {
	Base
	BrevoDetails
}

// BrevoDetails contains the configuration fields for Brevo notifications.
type BrevoDetails struct {
	// APIKey is the Brevo API key for authentication.
	APIKey string `json:"brevoApiKey"`
	// ToEmail is the recipient email address.
	ToEmail string `json:"brevoToEmail"`
	// FromEmail is the sender email address.
	FromEmail string `json:"brevoFromEmail"`
	// FromName is the sender name.
	FromName string `json:"brevoFromName"`
	// Subject is the email subject line.
	Subject string `json:"brevoSubject"`
	// CCEmail is a comma-separated list of CC email addresses.
	CCEmail string `json:"brevoCcEmail"`
	// BCCEmail is a comma-separated list of BCC email addresses.
	BCCEmail string `json:"brevoBccEmail"`
}

// Type returns the notification type identifier for Brevo.
func (b Brevo) Type() string {
	return b.BrevoDetails.Type()
}

// Type returns the notification type identifier for BrevoDetails.
func (BrevoDetails) Type() string {
	return "brevo"
}

// String returns a string representation of the Brevo notification.
func (b Brevo) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(b.Base, false), formatNotification(b.BrevoDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Brevo notification.
func (b *Brevo) UnmarshalJSON(data []byte) error {
	detail := BrevoDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*b = Brevo{
		Base:         base,
		BrevoDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Brevo notification into JSON.
func (b Brevo) MarshalJSON() ([]byte, error) {
	return marshalJSON(b.Base, b.BrevoDetails)
}
