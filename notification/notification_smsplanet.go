package notification

import (
	"fmt"
)

// SMSPlanet represents an SMS Planet notification provider.
// SMS Planet is an SMS service provider for sending text message notifications.
type SMSPlanet struct {
	Base
	SMSPlanetDetails
}

// SMSPlanetDetails contains the configuration fields for SMS Planet notifications.
type SMSPlanetDetails struct {
	// APIToken is the SMS Planet API token for authentication.
	APIToken string `json:"smsplanetApiToken"`
	// PhoneNumbers is the recipient phone numbers.
	PhoneNumbers string `json:"smsplanetPhoneNumbers"`
	// SenderName is the sender name or identifier.
	SenderName string `json:"smsplanetSenderName"`
}

// Type returns the notification type identifier for SMSPlanet.
func (s SMSPlanet) Type() string {
	return s.SMSPlanetDetails.Type()
}

// Type returns the notification type identifier for SMSPlanetDetails.
func (SMSPlanetDetails) Type() string {
	return "SMSPlanet"
}

// String returns a string representation of the SMS Planet notification.
func (s SMSPlanet) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSPlanetDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMS Planet notification.
func (s *SMSPlanet) UnmarshalJSON(data []byte) error {
	detail := SMSPlanetDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSPlanet{
		Base:             base,
		SMSPlanetDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMS Planet notification into JSON.
func (s SMSPlanet) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSPlanetDetails)
}
