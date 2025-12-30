package notification

import (
	"fmt"
)

// OneBot represents a OneBot (QQ) notification provider.
// OneBot is a protocol for QQ robot integrations, allowing notifications to be sent to QQ groups or private users.
type OneBot struct {
	Base
	OneBotDetails
}

// OneBotDetails contains the configuration fields for OneBot notifications.
type OneBotDetails struct {
	// HTTPAddr is the HTTP address of the OneBot service (e.g., "http://localhost:5700").
	HTTPAddr string `json:"httpAddr"`
	// AccessToken is the access token for authentication with the OneBot service.
	AccessToken string `json:"accessToken"`
	// MsgType is the message type: "group" for group messages or "private" for private messages.
	MsgType string `json:"msgType"`
	// ReceiverID is the QQ group ID (when MsgType is "group") or user ID (when MsgType is "private").
	ReceiverID string `json:"recieverId"`
}

// Type returns the notification type identifier for OneBot.
func (o OneBot) Type() string {
	return o.OneBotDetails.Type()
}

// Type returns the notification type identifier for OneBotDetails.
func (OneBotDetails) Type() string {
	return "OneBot"
}

// String returns a string representation of the OneBot notification.
func (o OneBot) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OneBotDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a OneBot notification.
func (o *OneBot) UnmarshalJSON(data []byte) error {
	detail := OneBotDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = OneBot{
		Base:          base,
		OneBotDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the OneBot notification into JSON.
func (o OneBot) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OneBotDetails)
}
