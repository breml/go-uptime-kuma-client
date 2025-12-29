package notification

import (
	"fmt"
)

// Twilio ...
type Twilio struct {
	Base
	TwilioDetails
}

// TwilioDetails ...
type TwilioDetails struct {
	AccountSID string `json:"twilioAccountSID"`
	APIKey     string `json:"twilioApiKey"`
	AuthToken  string `json:"twilioAuthToken"`
	ToNumber   string `json:"twilioToNumber"`
	FromNumber string `json:"twilioFromNumber"`
}

// Type ...
func (t Twilio) Type() string {
	return t.TwilioDetails.Type()
}

// Type ...
func (n TwilioDetails) Type() string {
	return "twilio"
}

// String ...
func (t Twilio) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TwilioDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (t Twilio) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TwilioDetails)
}
