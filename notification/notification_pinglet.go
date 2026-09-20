package notification

import (
	"fmt"
)

// Pinglet represents a Pinglet notification provider.
// Pinglet (https://app.pinglet.co.uk/) is a topic based push service, the raw
// Uptime Kuma webhook payload is posted to a topic endpoint, where Pinglet's
// Uptime Kuma rewriter turns it into a titled message with level and priority.
type Pinglet struct {
	Base
	PingletDetails
}

// PingletDetails contains the configuration fields for Pinglet notifications.
type PingletDetails struct {
	// PublishURL is the publish URL of the Pinglet topic the payload is posted
	// to (e.g. https://app.pinglet.co.uk/your-namespace/alerts). The topic is
	// created on the first publish, if it does not exist yet. The server strips
	// a single trailing slash before the request and always appends the
	// rewrite=uptimekuma query parameter, so the rewriter is not configurable.
	// The URL is stored verbatim, a trailing slash survives the round trip.
	PublishURL string `json:"pingletPublishUrl"`
	// APIKey authenticates the request against Pinglet and is sent as an
	// "Authorization: Bearer <key>" header.
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
