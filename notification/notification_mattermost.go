package notification

import (
	"fmt"
)

// Mattermost represents a mattermost notification.
type Mattermost struct {
	Base
	MattermostDetails
}

// MattermostDetails contains mattermost-specific notification configuration.
type MattermostDetails struct {
	WebhookURL string `json:"mattermostWebhookUrl"`
	Username   string `json:"mattermostusername"`
	Channel    string `json:"mattermostchannel"`
	IconEmoji  string `json:"mattermosticonemo"`
	IconURL    string `json:"mattermosticonurl"`
}

// Type returns the notification type.
func (m Mattermost) Type() string {
	return m.MattermostDetails.Type()
}

// Type returns the notification type.
func (MattermostDetails) Type() string {
	return "mattermost"
}

// String returns a string representation of the notification.
func (m Mattermost) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MattermostDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (m *Mattermost) UnmarshalJSON(data []byte) error {
	detail := MattermostDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*m = Mattermost{
		Base:              base,
		MattermostDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (m Mattermost) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, &m.MattermostDetails)
}
