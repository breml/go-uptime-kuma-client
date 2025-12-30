package notification

import (
	"fmt"
)

// Pushover represents a pushover notification.
type Pushover struct {
	Base
	PushoverDetails
}

// PushoverDetails contains pushover-specific notification configuration.
type PushoverDetails struct {
	UserKey  string `json:"pushoveruserkey"`
	AppToken string `json:"pushoverapptoken"`
	Sounds   string `json:"pushoversounds"`
	SoundsUp string `json:"pushoversounds_up"`
	Priority string `json:"pushoverpriority"`
	Title    string `json:"pushovertitle"`
	Device   string `json:"pushoverdevice"`
	TTL      string `json:"pushoverttl"`
}

// Type returns the notification type.
func (p Pushover) Type() string {
	return p.PushoverDetails.Type()
}

// Type returns the notification type.
func (PushoverDetails) Type() string {
	return "pushover"
}

// String returns a string representation of the notification.
func (p Pushover) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushoverDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (p *Pushover) UnmarshalJSON(data []byte) error {
	detail := PushoverDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Pushover{
		Base:            base,
		PushoverDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (p Pushover) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PushoverDetails)
}
