package notification

import (
	"fmt"
)

type Twilio struct {
	Base
	TwilioDetails
}

type TwilioDetails struct {
	AccountSID string `json:"twilioAccountSID"`
	APIKey     string `json:"twilioApiKey"`
	AuthToken  string `json:"twilioAuthToken"`
	ToNumber   string `json:"twilioToNumber"`
	FromNumber string `json:"twilioFromNumber"`
}

func (t Twilio) Type() string {
	return t.TwilioDetails.Type()
}

func (n TwilioDetails) Type() string {
	return "twilio"
}

func (t Twilio) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TwilioDetails, true))
}

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

func (t Twilio) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.TwilioDetails)
}
