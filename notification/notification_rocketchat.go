package notification

import (
	"fmt"
)

// RocketChat represents a rocketchat notification.
type RocketChat struct {
	Base
	RocketChatDetails
}

// RocketChatDetails contains rocketchat-specific notification configuration.
type RocketChatDetails struct {
	WebhookURL string `json:"rocketwebhookURL"`
	Channel    string `json:"rocketchannel"`
	Username   string `json:"rocketusername"`
	IconEmoji  string `json:"rocketiconemo"`
	Button     string `json:"rocketbutton"`
}

// Type returns the notification type.
func (r RocketChat) Type() string {
	return r.RocketChatDetails.Type()
}

// Type returns the notification type.
func (RocketChatDetails) Type() string {
	return "rocket.chat"
}

// String returns a string representation of the notification.
func (r RocketChat) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(r.Base, false), formatNotification(r.RocketChatDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (r RocketChat) MarshalJSON() ([]byte, error) {
	return marshalJSON(r.Base, &r.RocketChatDetails)
}
