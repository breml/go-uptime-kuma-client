package notification

import (
	"fmt"
)

// SerwerSMS represents a SerwerSMS notification provider.
// SerwerSMS is an SMS service provider for sending text message notifications.
type SerwerSMS struct {
	Base
	SerwerSMSDetails
}

// SerwerSMSDetails contains the configuration fields for SerwerSMS notifications.
type SerwerSMSDetails struct {
	// Username is the SerwerSMS account username.
	Username string `json:"serwersmsUsername"`
	// Password is the SerwerSMS account password.
	Password string `json:"serwersmsPassword"`
	// PhoneNumber is the recipient phone number.
	PhoneNumber string `json:"serwersmsPhoneNumber"`
	// SenderName is the sender name or identifier.
	SenderName string `json:"serwersmsSenderName"`
}

// Type returns the notification type identifier for SerwerSMS.
func (s SerwerSMS) Type() string {
	return s.SerwerSMSDetails.Type()
}

// Type returns the notification type identifier for SerwerSMSDetails.
func (SerwerSMSDetails) Type() string {
	return "serwersms"
}

// String returns a string representation of the SerwerSMS notification.
func (s SerwerSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SerwerSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a SerwerSMS notification.
func (s *SerwerSMS) UnmarshalJSON(data []byte) error {
	detail := SerwerSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SerwerSMS{
		Base:             base,
		SerwerSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SerwerSMS notification into JSON.
func (s SerwerSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SerwerSMSDetails)
}
