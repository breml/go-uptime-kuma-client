package notification

import (
	"fmt"
)

// OneChat represents a OneChat notification provider.
// OneChat is a Thai messaging platform for sending notifications through a bot API.
type OneChat struct {
	Base
	OneChatDetails
}

// OneChatDetails contains the configuration fields for OneChat notifications.
type OneChatDetails struct {
	// AccessToken is the access token for authentication with OneChat API.
	AccessToken string `json:"accessToken"`
	// ReceiverID is the recipient ID (user or group ID) in OneChat.
	ReceiverID string `json:"recieverId"`
	// BotID is the bot ID for sending messages through OneChat.
	BotID string `json:"botId"`
}

// Type returns the notification type identifier for OneChat.
func (o OneChat) Type() string {
	return o.OneChatDetails.Type()
}

// Type returns the notification type identifier for OneChatDetails.
func (o OneChatDetails) Type() string {
	return "OneChat"
}

// String returns a string representation of the OneChat notification.
func (o OneChat) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OneChatDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a OneChat notification.
func (o *OneChat) UnmarshalJSON(data []byte) error {
	detail := OneChatDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = OneChat{
		Base:           base,
		OneChatDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the OneChat notification into JSON.
func (o OneChat) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OneChatDetails)
}
