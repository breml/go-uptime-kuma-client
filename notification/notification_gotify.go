package notification

import (
	"fmt"
)

// Gotify ...
type Gotify struct {
	Base
	GotifyDetails
}

// GotifyDetails ...
type GotifyDetails struct {
	ServerURL        string `json:"gotifyserverurl"`
	ApplicationToken string `json:"gotifyapplicationToken"`
	Priority         int    `json:"gotifyPriority"`
}

// Type ...
func (g Gotify) Type() string {
	return g.GotifyDetails.Type()
}

// Type ...
func (GotifyDetails) Type() string {
	return "gotify"
}

// String ...
func (g Gotify) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GotifyDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (g Gotify) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, &g.GotifyDetails)
}
