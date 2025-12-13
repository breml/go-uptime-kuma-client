package notification

import (
	"fmt"
)

type Pushbullet struct {
	Base
	PushbulletDetails
}

type PushbulletDetails struct {
	AccessToken string `json:"pushbulletAccessToken"`
}

func (p Pushbullet) Type() string {
	return p.PushbulletDetails.Type()
}

func (n PushbulletDetails) Type() string {
	return "pushbullet"
}

func (p Pushbullet) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushbulletDetails, true))
}

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

func (p Pushbullet) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PushbulletDetails)
}
