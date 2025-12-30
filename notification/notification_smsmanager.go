package notification

import (
	"fmt"
)

// SMSManager represents an SMS Manager notification provider.
// SMS Manager is an SMS service provider for sending text message notifications.
type SMSManager struct {
	Base
	SMSManagerDetails
}

// SMSManagerDetails contains the configuration fields for SMS Manager notifications.
type SMSManagerDetails struct {
	// APIKey is the SMS Manager API key for authentication.
	APIKey string `json:"smsmanagerApiKey"`
	// Numbers is the recipient phone number.
	Numbers string `json:"numbers"`
	// MessageType is the message gateway type.
	MessageType string `json:"messageType"`
}

// Type returns the notification type identifier for SMSManager.
func (s SMSManager) Type() string {
	return s.SMSManagerDetails.Type()
}

// Type returns the notification type identifier for SMSManagerDetails.
func (SMSManagerDetails) Type() string {
	return "SMSManager"
}

// String returns a string representation of the SMSManager notification.
func (s SMSManager) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSManagerDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSManager notification.
func (s *SMSManager) UnmarshalJSON(data []byte) error {
	detail := SMSManagerDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSManager{
		Base:              base,
		SMSManagerDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSManager notification into JSON.
func (s SMSManager) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSManagerDetails)
}
