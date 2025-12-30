package notification

import (
	"fmt"
)

// Pushbullet ...
type Pushbullet struct {
	Base
	PushbulletDetails
}

// PushbulletDetails ...
type PushbulletDetails struct {
	AccessToken string `json:"pushbulletAccessToken"`
}

// Type ...
func (p Pushbullet) Type() string {
	return p.PushbulletDetails.Type()
}

// Type ...
func (PushbulletDetails) Type() string {
	return "pushbullet"
}

// String ...
func (p Pushbullet) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushbulletDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (p Pushbullet) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PushbulletDetails)
}
