package notification

import (
	"fmt"
)

type Slack struct {
	Base
	SlackDetails
}

type SlackDetails struct {
	WebhookURL    string `json:"slackwebhookURL"`
	Username      string `json:"slackusername"`
	IconEmoji     string `json:"slackiconemo"`
	Channel       string `json:"slackchannel"`
	RichMessage   bool   `json:"slackrichmessage"`
	ChannelNotify bool   `json:"slackchannelnotify"`
}

func (s Slack) Type() string {
	return s.SlackDetails.Type()
}

func (n SlackDetails) Type() string {
	return "slack"
}

func (s Slack) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SlackDetails, true))
}

func (s *Slack) UnmarshalJSON(data []byte) error {
	detail := SlackDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Slack{
		Base:         base,
		SlackDetails: detail,
	}

	return nil
}

func (s Slack) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SlackDetails)
}
