package notification

import (
	"fmt"
)

// Squadcast represents a Squadcast notification provider.
// Squadcast is an incident management platform for sending alerts via webhooks.
type Squadcast struct {
	Base
	SquadcastDetails
}

// SquadcastDetails contains the configuration fields for Squadcast notifications.
type SquadcastDetails struct {
	// WebhookURL is the Squadcast webhook URL for sending notifications.
	WebhookURL string `json:"squadcastWebhookURL"`
}

// Type returns the notification type identifier for Squadcast.
func (s Squadcast) Type() string {
	return s.SquadcastDetails.Type()
}

// Type returns the notification type identifier for SquadcastDetails.
func (n SquadcastDetails) Type() string {
	return "squadcast"
}

// String returns a string representation of the Squadcast notification.
func (s Squadcast) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SquadcastDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Squadcast notification.
func (s *Squadcast) UnmarshalJSON(data []byte) error {
	detail := SquadcastDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Squadcast{
		Base:             base,
		SquadcastDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Squadcast notification into JSON.
func (s Squadcast) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SquadcastDetails)
}
