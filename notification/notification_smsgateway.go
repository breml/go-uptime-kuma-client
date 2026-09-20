package notification

import (
	"fmt"
)

// SMSGateway represents an SMS Gateway notification provider.
// SMS Gateway (https://github.com/mattboston/sms-gateway) is a self-hosted
// gateway that sends text messages through an attached USB GSM modem.
type SMSGateway struct {
	Base
	SMSGatewayDetails
}

// SMSGatewayDetails contains the configuration fields for SMS Gateway notifications.
type SMSGatewayDetails struct {
	// URL is the base URL of the SMS Gateway server, without an API path
	// (e.g. http://localhost:8080). The server strips trailing slashes and
	// appends /api/v1/sms/send itself.
	URL string `json:"smsgatewayUrl"`
	// APIKey authenticates the request against the gateway and is sent as an
	// "X-API-Key" header.
	APIKey string `json:"smsgatewayApiKey"`
	// Recipients holds one or more recipient phone numbers in international
	// format (e.g. "+15551234567, +15559876543"), separated by comma. The
	// server trims the whitespace around every number and drops empty entries.
	// A value that leaves no recipient at all is not an error: the server sends
	// nothing and still reports the notification as sent successfully.
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
