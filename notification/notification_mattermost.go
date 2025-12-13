package notification

import (
	"fmt"
)

type Mattermost struct {
	Base
	MattermostDetails
}

type MattermostDetails struct {
	WebhookURL string `json:"mattermostWebhookUrl"`
	Username   string `json:"mattermostusername"`
	Channel    string `json:"mattermostchannel"`
	IconEmoji  string `json:"mattermosticonemo"`
	IconURL    string `json:"mattermosticonurl"`
}

func (m Mattermost) Type() string {
	return m.MattermostDetails.Type()
}

func (n MattermostDetails) Type() string {
	return "mattermost"
}

func (m Mattermost) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MattermostDetails, true))
}

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

func (m Mattermost) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, m.MattermostDetails)
}
