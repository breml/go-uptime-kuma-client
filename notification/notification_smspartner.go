package notification

import (
	"fmt"
)

// SMSPartner represents an SMSPartner notification provider.
// SMSPartner is an SMS service provider for sending text message notifications.
type SMSPartner struct {
	Base
	SMSPartnerDetails
}

// SMSPartnerDetails contains the configuration fields for SMSPartner notifications.
type SMSPartnerDetails struct {
	// APIKey is the SMSPartner API key for authentication.
	APIKey string `json:"smspartnerApikey"`
	// PhoneNumber is the recipient phone number.
	PhoneNumber string `json:"smspartnerPhoneNumber"`
	// SenderName is the sender name or identifier.
	SenderName string `json:"smspartnerSenderName"`
}

// Type returns the notification type identifier for SMSPartner.
func (s SMSPartner) Type() string {
	return s.SMSPartnerDetails.Type()
}

// Type returns the notification type identifier for SMSPartnerDetails.
func (n SMSPartnerDetails) Type() string {
	return "SMSPartner"
}

// String returns a string representation of the SMSPartner notification.
func (s SMSPartner) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSPartnerDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSPartner notification.
func (s *SMSPartner) UnmarshalJSON(data []byte) error {
	detail := SMSPartnerDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSPartner{
		Base:              base,
		SMSPartnerDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSPartner notification into JSON.
func (s SMSPartner) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSPartnerDetails)
}
