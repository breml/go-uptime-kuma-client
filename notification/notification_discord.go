package notification

import (
	"fmt"
)

// Discord ...
type Discord struct {
	Base
	DiscordDetails
}

// DiscordDetails ...
type DiscordDetails struct {
	WebhookURL    string `json:"discordWebhookUrl"`
	Username      string `json:"discordUsername"`
	ChannelType   string `json:"discordChannelType"`
	ThreadID      string `json:"threadId"`
	PostName      string `json:"postName"`
	PrefixMessage string `json:"discordPrefixMessage"`
	DisableURL    bool   `json:"disableUrl"`
}

// Type ...
func (d Discord) Type() string {
	return d.DiscordDetails.Type()
}

// Type ...
func (DiscordDetails) Type() string {
	return "discord"
}

// String ...
func (d Discord) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(d.Base, false), formatNotification(d.DiscordDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (d Discord) MarshalJSON() ([]byte, error) {
	return marshalJSON(d.Base, &d.DiscordDetails)
}
