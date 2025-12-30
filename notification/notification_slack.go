package notification

import (
	"fmt"
)

// Slack ...
type Slack struct {
	Base
	SlackDetails
}

// SlackDetails ...
type SlackDetails struct {
	WebhookURL    string `json:"slackwebhookURL"`
	Username      string `json:"slackusername"`
	IconEmoji     string `json:"slackiconemo"`
	Channel       string `json:"slackchannel"`
	RichMessage   bool   `json:"slackrichmessage"`
	ChannelNotify bool   `json:"slackchannelnotify"`
}

// Type ...
func (s Slack) Type() string {
	return s.SlackDetails.Type()
}

// Type ...
func (SlackDetails) Type() string {
	return "slack"
}

// String ...
func (s Slack) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SlackDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (s Slack) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, &s.SlackDetails)
}
