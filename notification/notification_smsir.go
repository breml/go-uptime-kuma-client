package notification

import (
	"fmt"
)

// SMSIR represents an SMS.ir notification provider.
// SMS.ir is an Iranian SMS gateway for sending text message notifications.
type SMSIR struct {
	Base
	SMSIRDetails
}

// SMSIRDetails contains the configuration fields for SMS.ir notifications.
type SMSIRDetails struct {
	// APIKey is the SMS.ir API key for authentication.
	APIKey string `json:"smsirApiKey"`
	// Number is the comma separated list of recipient phone numbers.
	Number string `json:"smsirNumber"`
	// Template is the pre-approved template ID used for sending the message.
	Template string `json:"smsirTemplate"`
}

// Type returns the notification type identifier for SMSIR.
func (s SMSIR) Type() string {
	return s.SMSIRDetails.Type()
}

// Type returns the notification type identifier for SMSIRDetails.
func (SMSIRDetails) Type() string {
	return "smsir"
}

// String returns a string representation of the SMSIR notification.
func (s SMSIR) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSIRDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSIR notification.
func (s *SMSIR) UnmarshalJSON(data []byte) error {
	detail := SMSIRDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSIR{
		Base:         base,
		SMSIRDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSIR notification into JSON.
func (s SMSIR) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSIRDetails)
}
