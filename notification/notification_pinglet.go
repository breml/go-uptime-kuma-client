package notification

import (
	"fmt"
)

// Pinglet represents a Pinglet notification provider.
// Pinglet forwards the raw Uptime Kuma webhook payload to a publish URL, where
// its Uptime Kuma rewriter turns it into a titled message.
type Pinglet struct {
	Base
	PingletDetails
}

// PingletDetails contains the configuration fields for Pinglet notifications.
type PingletDetails struct {
	// PublishURL is the Pinglet topic endpoint the payload is posted to.
	PublishURL string `json:"pingletPublishUrl"`
	// APIKey is the Pinglet API key, sent as a bearer token.
	APIKey string `json:"pingletApiKey"`
}

// Type returns the notification type identifier for Pinglet.
func (p Pinglet) Type() string {
	return p.PingletDetails.Type()
}

// Type returns the notification type identifier for PingletDetails.
func (PingletDetails) Type() string {
	return "pinglet"
}

// String returns a string representation of the Pinglet notification.
func (p Pinglet) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PingletDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Pinglet notification.
func (p *Pinglet) UnmarshalJSON(data []byte) error {
	detail := PingletDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Pinglet{
		Base:           base,
		PingletDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Pinglet notification into JSON.
func (p Pinglet) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PingletDetails)
}
