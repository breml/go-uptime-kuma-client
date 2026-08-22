package notification

import (
	"fmt"
)

// Ooredoo represents an Ooredoo notification provider.
// Ooredoo Maldives is a bulk SMS gateway for sending text message notifications.
type Ooredoo struct {
	Base
	OoredooDetails
}

// OoredooDetails contains the configuration fields for Ooredoo notifications.
type OoredooDetails struct {
	// Username is the Ooredoo bulk SMS account username.
	Username string `json:"ooredooUsername"`
	// AccessKey is the Ooredoo account access key. It is Base64 encoded by the
	// server before it is passed on to the gateway.
	AccessKey string `json:"ooredooAccessKey"`
	// BearerToken authenticates the request against the gateway and is sent as
	// an "Authorization: Bearer" header.
	BearerToken string `json:"ooredooBearerToken"`
	// ToNumber holds one or more recipient phone numbers, separated by comma,
	// semicolon or whitespace. Since whitespace separates recipients, an
	// individual number must not contain spaces. The server strips every "+"
	// character and prefixes bare 7 digit numbers with the Maldives country
	// code 960. The numbers are not validated: the server only rejects the
	// notification when it is sent and no recipient remains at all, any other
	// value is passed on to the gateway as is.
	ToNumber string `json:"ooredooToNumber"`
	// ServerURL is the Ooredoo API endpoint. If unset or empty, the server
	// falls back to https://o-papi1-lb01.ooredoo.mv/bulk_sms/v2.
	ServerURL *string `json:"ooredooServerUrl,omitempty"`
}

// Type returns the notification type identifier for Ooredoo.
func (o Ooredoo) Type() string {
	return o.OoredooDetails.Type()
}

// Type returns the notification type identifier for OoredooDetails.
func (OoredooDetails) Type() string {
	return "Ooredoo"
}

// String returns a string representation of the Ooredoo notification.
func (o Ooredoo) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OoredooDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an Ooredoo notification.
func (o *Ooredoo) UnmarshalJSON(data []byte) error {
	detail := OoredooDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = Ooredoo{
		Base:           base,
		OoredooDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Ooredoo notification into JSON.
func (o Ooredoo) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, &o.OoredooDetails)
}
