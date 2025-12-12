package notification

import (
	"fmt"
)

type Discord struct {
	Base
	DiscordDetails
}

type DiscordDetails struct {
	WebhookURL       string `json:"discordWebhookUrl"`
	Username         string `json:"discordUsername"`
	ChannelType      string `json:"discordChannelType"`
	ThreadID         string `json:"threadId"`
	PostName         string `json:"postName"`
	PrefixMessage    string `json:"discordPrefixMessage"`
	DisableURL       bool   `json:"disableUrl"`
}

func (d Discord) Type() string {
	return d.DiscordDetails.Type()
}

func (n DiscordDetails) Type() string {
	return "discord"
}

func (d Discord) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(d.Base, false), formatNotification(d.DiscordDetails, true))
}

func (d *Discord) UnmarshalJSON(data []byte) error {
	detail := DiscordDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*d = Discord{
		Base:            base,
		DiscordDetails: detail,
	}

	return nil
}

func (d Discord) MarshalJSON() ([]byte, error) {
	return marshalJSON(d.Base, d.DiscordDetails)
}
