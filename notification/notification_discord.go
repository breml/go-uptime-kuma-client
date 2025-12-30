package notification

import (
	"fmt"
)

// Discord represents a discord notification.
type Discord struct {
	Base
	DiscordDetails
}

// DiscordDetails contains discord-specific notification configuration.
type DiscordDetails struct {
	WebhookURL    string `json:"discordWebhookUrl"`
	Username      string `json:"discordUsername"`
	ChannelType   string `json:"discordChannelType"`
	ThreadID      string `json:"threadId"`
	PostName      string `json:"postName"`
	PrefixMessage string `json:"discordPrefixMessage"`
	DisableURL    bool   `json:"disableUrl"`
}

// Type returns the notification type.
func (d Discord) Type() string {
	return d.DiscordDetails.Type()
}

// Type returns the notification type.
func (DiscordDetails) Type() string {
	return "discord"
}

// String returns a string representation of the notification.
func (d Discord) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(d.Base, false), formatNotification(d.DiscordDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (d *Discord) UnmarshalJSON(data []byte) error {
	detail := DiscordDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*d = Discord{
		Base:           base,
		DiscordDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (d Discord) MarshalJSON() ([]byte, error) {
	return marshalJSON(d.Base, &d.DiscordDetails)
}
