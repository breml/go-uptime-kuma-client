package notification

import (
	"fmt"
)

// EgoSMS represents an EgoSMS notification.
type EgoSMS struct {
	Base
	EgoSMSDetails
}

// EgoSMSDetails contains EgoSMS-specific notification configuration.
type EgoSMSDetails struct {
	Username    string  `json:"egosmsUsername"`
	Password    string  `json:"egosmsPassword"`
	Sender      *string `json:"egosmsSender,omitempty"`
	PhoneNumber string  `json:"egosmsPhoneNumber"`
}

// Type returns the notification type.
func (e EgoSMS) Type() string {
	return e.EgoSMSDetails.Type()
}

// Type returns the notification type.
func (EgoSMSDetails) Type() string {
	return "egosms"
}

// String returns a string representation of the notification.
func (e EgoSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(e.Base, false), formatNotification(e.EgoSMSDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (e *EgoSMS) UnmarshalJSON(data []byte) error {
	detail := EgoSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*e = EgoSMS{
		Base:          base,
		EgoSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (e EgoSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(e.Base, &e.EgoSMSDetails)
}
