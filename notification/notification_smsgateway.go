package notification

import (
	"fmt"
)

// SMSGateway represents an SMS Gateway notification provider.
// SMS Gateway for Android is a self-hosted gateway that sends text messages
// through an Android phone.
type SMSGateway struct {
	Base
	SMSGatewayDetails
}

// SMSGatewayDetails contains the configuration fields for SMS Gateway notifications.
type SMSGatewayDetails struct {
	// URL is the base URL of the SMS Gateway server.
	URL string `json:"smsgatewayUrl"`
	// APIKey is the SMS Gateway API key for authentication.
	APIKey string `json:"smsgatewayApiKey"`
	// Recipients is the comma separated list of recipient phone numbers.
	Recipients string `json:"smsgatewayTo"`
}

// Type returns the notification type identifier for SMSGateway.
func (s SMSGateway) Type() string {
	return s.SMSGatewayDetails.Type()
}

// Type returns the notification type identifier for SMSGatewayDetails.
func (SMSGatewayDetails) Type() string {
	return "SMSGateway"
}

// String returns a string representation of the SMSGateway notification.
func (s SMSGateway) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSGatewayDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSGateway notification.
func (s *SMSGateway) UnmarshalJSON(data []byte) error {
	detail := SMSGatewayDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSGateway{
		Base:              base,
		SMSGatewayDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSGateway notification into JSON.
func (s SMSGateway) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSGatewayDetails)
}
