package notification

import (
	"fmt"
)

// NextcloudTalk represents a Nextcloud Talk notification provider.
// Nextcloud Talk is a communication platform integrated with Nextcloud.
type NextcloudTalk struct {
	Base
	NextcloudTalkDetails
}

// NextcloudTalkDetails contains the configuration fields for Nextcloud Talk notifications.
type NextcloudTalkDetails struct {
	// Host is the Nextcloud instance host URL.
	Host string `json:"host"`
	// ConversationToken is the conversation/room token for the target chat.
	ConversationToken string `json:"conversationToken"`
	// BotSecret is the bot secret for authentication and message signing.
	BotSecret string `json:"botSecret"`
	// SendSilentUp determines if UP notifications should be sent silently.
	SendSilentUp bool `json:"sendSilentUp"`
	// SendSilentDown determines if DOWN notifications should be sent silently.
	SendSilentDown bool `json:"sendSilentDown"`
}

// Type returns the notification type identifier for NextcloudTalk.
func (n NextcloudTalk) Type() string {
	return n.NextcloudTalkDetails.Type()
}

// Type returns the notification type identifier for NextcloudTalkDetails.
func (NextcloudTalkDetails) Type() string {
	return "NextcloudTalk"
}

// String returns a string representation of the NextcloudTalk notification.
func (n NextcloudTalk) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NextcloudTalkDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a NextcloudTalk notification.
func (n *NextcloudTalk) UnmarshalJSON(data []byte) error {
	detail := NextcloudTalkDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*n = NextcloudTalk{
		Base:                 base,
		NextcloudTalkDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the NextcloudTalk notification into JSON.
func (n NextcloudTalk) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, n.NextcloudTalkDetails)
}
