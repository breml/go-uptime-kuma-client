package notification

import (
	"fmt"
)

// Pushover ...
type Pushover struct {
	Base
	PushoverDetails
}

// PushoverDetails ...
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

// Type ...
func (p Pushover) Type() string {
	return p.PushoverDetails.Type()
}

// Type ...
func (PushoverDetails) Type() string {
	return "pushover"
}

// String ...
func (p Pushover) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushoverDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (p Pushover) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PushoverDetails)
}
