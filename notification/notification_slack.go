package notification

import (
	"fmt"
)

// Slack represents a slack notification.
type Slack struct {
	Base
	SlackDetails
}

// SlackDetails contains slack-specific notification configuration.
type SlackDetails struct {
	WebhookURL    string `json:"slackwebhookURL"`
	Username      string `json:"slackusername"`
	IconEmoji     string `json:"slackiconemo"`
	Channel       string `json:"slackchannel"`
	RichMessage   bool   `json:"slackrichmessage"`
	ChannelNotify bool   `json:"slackchannelnotify"`
}

// Type returns the notification type.
func (s Slack) Type() string {
	return s.SlackDetails.Type()
}

// Type returns the notification type.
func (SlackDetails) Type() string {
	return "slack"
}

// String returns a string representation of the notification.
func (s Slack) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SlackDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (s Slack) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, &s.SlackDetails)
}
