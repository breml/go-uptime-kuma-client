package notification

import (
	"fmt"
)

type RocketChat struct {
	Base
	RocketChatDetails
}

type RocketChatDetails struct {
	WebhookURL string `json:"rocketwebhookURL"`
	Channel    string `json:"rocketchannel"`
	Username   string `json:"rocketusername"`
	IconEmoji  string `json:"rocketiconemo"`
	Button     string `json:"rocketbutton"`
}

func (r RocketChat) Type() string {
	return r.RocketChatDetails.Type()
}

func (n RocketChatDetails) Type() string {
	return "rocket.chat"
}

func (r RocketChat) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(r.Base, false), formatNotification(r.RocketChatDetails, true))
}

func (r *RocketChat) UnmarshalJSON(data []byte) error {
	detail := RocketChatDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*r = RocketChat{
		Base:              base,
		RocketChatDetails: detail,
	}

	return nil
}

func (r RocketChat) MarshalJSON() ([]byte, error) {
	return marshalJSON(r.Base, r.RocketChatDetails)
}
