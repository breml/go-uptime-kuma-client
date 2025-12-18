package notification

import (
	"fmt"
)

// SevenIO represents a Seven.io SMS notification provider.
// Seven.io is an SMS service provider for sending text message notifications.
type SevenIO struct {
	Base
	SevenIODetails
}

// SevenIODetails contains the configuration fields for Seven.io SMS notifications.
type SevenIODetails struct {
	// APIKey is the Seven.io API key for authentication.
	APIKey string `json:"sevenioApiKey"`
	// Sender is the sender name or phone number.
	Sender string `json:"sevenioSender"`
	// To is the recipient phone number.
	To string `json:"sevenioTo"`
}

// Type returns the notification type identifier for SevenIO.
func (s SevenIO) Type() string {
	return s.SevenIODetails.Type()
}

// Type returns the notification type identifier for SevenIODetails.
func (n SevenIODetails) Type() string {
	return "sevenio"
}

// String returns a string representation of the SevenIO notification.
func (s SevenIO) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SevenIODetails, true))
}

// UnmarshalJSON unmarshals JSON data into a SevenIO notification.
func (s *SevenIO) UnmarshalJSON(data []byte) error {
	detail := SevenIODetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SevenIO{
		Base:           base,
		SevenIODetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SevenIO notification into JSON.
func (s SevenIO) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SevenIODetails)
}
