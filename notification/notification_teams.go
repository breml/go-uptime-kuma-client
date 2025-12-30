package notification

import (
	"fmt"
)

// Teams ...
type Teams struct {
	Base
	TeamsDetails
}

// TeamsDetails ...
type TeamsDetails struct {
	WebhookURL string `json:"webhookUrl"`
}

// Type ...
func (t Teams) Type() string {
	return t.TeamsDetails.Type()
}

// Type ...
func (TeamsDetails) Type() string {
	return "teams"
}

// String ...
func (t Teams) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TeamsDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (t Teams) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TeamsDetails)
}
