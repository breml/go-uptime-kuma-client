package notification

import (
	"fmt"
)

// Gotify represents a gotify notification.
type Gotify struct {
	Base
	GotifyDetails
}

// GotifyDetails contains gotify-specific notification configuration.
type GotifyDetails struct {
	ServerURL        string `json:"gotifyserverurl"`
	ApplicationToken string `json:"gotifyapplicationToken"`
	Priority         int    `json:"gotifyPriority"`
}

// Type returns the notification type.
func (g Gotify) Type() string {
	return g.GotifyDetails.Type()
}

// Type returns the notification type.
func (GotifyDetails) Type() string {
	return "gotify"
}

// String returns a string representation of the notification.
func (g Gotify) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GotifyDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (g *Gotify) UnmarshalJSON(data []byte) error {
	detail := GotifyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = Gotify{
		Base:          base,
		GotifyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (g Gotify) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, &g.GotifyDetails)
}
