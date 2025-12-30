package notification

import (
	"fmt"
)

// Pushbullet represents a pushbullet notification.
type Pushbullet struct {
	Base
	PushbulletDetails
}

// PushbulletDetails contains pushbullet-specific notification configuration.
type PushbulletDetails struct {
	AccessToken string `json:"pushbulletAccessToken"`
}

// Type returns the notification type.
func (p Pushbullet) Type() string {
	return p.PushbulletDetails.Type()
}

// Type returns the notification type.
func (PushbulletDetails) Type() string {
	return "pushbullet"
}

// String returns a string representation of the notification.
func (p Pushbullet) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushbulletDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (p *Pushbullet) UnmarshalJSON(data []byte) error {
	detail := PushbulletDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Pushbullet{
		Base:              base,
		PushbulletDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (p Pushbullet) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PushbulletDetails)
}
