package notification

import (
	"fmt"
)

type Teams struct {
	Base
	TeamsDetails
}

type TeamsDetails struct {
	WebhookURL string `json:"webhookUrl"`
}

func (t Teams) Type() string {
	return t.TeamsDetails.Type()
}

func (t TeamsDetails) Type() string {
	return "teams"
}

func (t Teams) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TeamsDetails, true))
}

func (t *Teams) UnmarshalJSON(data []byte) error {
	detail := TeamsDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Teams{
		Base:         base,
		TeamsDetails: detail,
	}

	return nil
}

func (t Teams) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.TeamsDetails)
}
