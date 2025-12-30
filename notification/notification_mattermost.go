package notification

import (
	"fmt"
)

// Mattermost ...
type Mattermost struct {
	Base
	MattermostDetails
}

// MattermostDetails ...
type MattermostDetails struct {
	WebhookURL string `json:"mattermostWebhookUrl"`
	Username   string `json:"mattermostusername"`
	Channel    string `json:"mattermostchannel"`
	IconEmoji  string `json:"mattermosticonemo"`
	IconURL    string `json:"mattermosticonurl"`
}

// Type ...
func (m Mattermost) Type() string {
	return m.MattermostDetails.Type()
}

// Type ...
func (MattermostDetails) Type() string {
	return "mattermost"
}

// String ...
func (m Mattermost) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MattermostDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (m Mattermost) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, &m.MattermostDetails)
}
