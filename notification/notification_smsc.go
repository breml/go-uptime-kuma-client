package notification

import (
	"fmt"
)

// SMSC represents an SMSC SMS notification provider.
// SMSC is an SMS service provider (https://smsc.kz) for sending text message notifications.
type SMSC struct {
	Base
	SMSCDetails
}

// SMSCDetails contains the configuration fields for SMSC SMS notifications.
type SMSCDetails struct {
	// Login is the SMSC account login/username.
	Login string `json:"smscLogin"`
	// Password is the SMSC account password.
	Password string `json:"smscPassword"`
	// ToNumber is the recipient phone number.
	ToNumber string `json:"smscToNumber"`
	// SenderName is the optional sender name or identifier.
	SenderName string `json:"smscSenderName"`
	// Translit indicates whether to transliterate non-ASCII characters.
	Translit string `json:"smscTranslit"`
}

// Type returns the notification type identifier for SMSC.
func (s SMSC) Type() string {
	return s.SMSCDetails.Type()
}

// Type returns the notification type identifier for SMSCDetails.
func (SMSCDetails) Type() string {
	return "smsc"
}

// String returns a string representation of the SMSC notification.
func (s SMSC) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSCDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSC notification.
func (s *SMSC) UnmarshalJSON(data []byte) error {
	detail := SMSCDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSC{
		Base:        base,
		SMSCDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSC notification into JSON.
func (s SMSC) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSCDetails)
}
