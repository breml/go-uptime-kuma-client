package notification

import (
	"fmt"
)

// SendGrid represents a SendGrid notification provider.
// SendGrid is an email delivery service that provides reliable email sending through an API.
type SendGrid struct {
	Base
	SendGridDetails
}

// SendGridDetails contains the SendGrid-specific configuration.
type SendGridDetails struct {
	APIKey    string `json:"sendgridApiKey"`    // SendGrid API key
	ToEmail   string `json:"sendgridToEmail"`   // Recipient email address
	FromEmail string `json:"sendgridFromEmail"` // Sender email address
	Subject   string `json:"sendgridSubject"`   // Email subject
	CcEmail   string `json:"sendgridCcEmail"`   // CC email addresses
	BccEmail  string `json:"sendgridBccEmail"`  // BCC email addresses
}

// Type returns the notification type identifier for SendGrid.
func (s SendGrid) Type() string {
	return s.SendGridDetails.Type()
}

// Type returns the notification type identifier for SendGrid details.
func (s SendGridDetails) Type() string {
	return "SendGrid"
}

// String returns a human-readable representation of the SendGrid notification.
func (s SendGrid) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SendGridDetails, true))
}

// UnmarshalJSON deserializes a SendGrid notification from JSON.
func (s *SendGrid) UnmarshalJSON(data []byte) error {
	detail := SendGridDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SendGrid{
		Base:            base,
		SendGridDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a SendGrid notification to JSON.
func (s SendGrid) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SendGridDetails)
}
