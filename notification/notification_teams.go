package notification

import (
	"fmt"
)

// Teams represents a teams notification.
type Teams struct {
	Base
	TeamsDetails
}

// TeamsDetails contains teams-specific notification configuration.
type TeamsDetails struct {
	WebhookURL string `json:"webhookUrl"`
}

// Type returns the notification type.
func (t Teams) Type() string {
	return t.TeamsDetails.Type()
}

// Type returns the notification type.
func (TeamsDetails) Type() string {
	return "teams"
}

// String returns a string representation of the notification.
func (t Teams) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TeamsDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (t Teams) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TeamsDetails)
}
