package notification

import (
	"fmt"
)

// Twilio represents a twilio notification.
type Twilio struct {
	Base
	TwilioDetails
}

// TwilioDetails contains twilio-specific notification configuration.
type TwilioDetails struct {
	AccountSID string `json:"twilioAccountSID"`
	APIKey     string `json:"twilioApiKey"`
	AuthToken  string `json:"twilioAuthToken"`
	ToNumber   string `json:"twilioToNumber"`
	FromNumber string `json:"twilioFromNumber"`
}

// Type returns the notification type.
func (t Twilio) Type() string {
	return t.TwilioDetails.Type()
}

// Type returns the notification type.
func (TwilioDetails) Type() string {
	return "twilio"
}

// String returns a string representation of the notification.
func (t Twilio) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TwilioDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (t *Twilio) UnmarshalJSON(data []byte) error {
	detail := TwilioDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Twilio{
		Base:          base,
		TwilioDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (t Twilio) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TwilioDetails)
}
