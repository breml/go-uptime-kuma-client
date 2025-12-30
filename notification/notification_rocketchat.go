package notification

import (
	"fmt"
)

// RocketChat ...
type RocketChat struct {
	Base
	RocketChatDetails
}

// RocketChatDetails ...
type RocketChatDetails struct {
	WebhookURL string `json:"rocketwebhookURL"`
	Channel    string `json:"rocketchannel"`
	Username   string `json:"rocketusername"`
	IconEmoji  string `json:"rocketiconemo"`
	Button     string `json:"rocketbutton"`
}

// Type ...
func (r RocketChat) Type() string {
	return r.RocketChatDetails.Type()
}

// Type ...
func (RocketChatDetails) Type() string {
	return "rocket.chat"
}

// String ...
func (r RocketChat) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(r.Base, false), formatNotification(r.RocketChatDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (r RocketChat) MarshalJSON() ([]byte, error) {
	return marshalJSON(r.Base, &r.RocketChatDetails)
}
