package notification

import (
	"fmt"
)

type Gotify struct {
	Base
	GotifyDetails
}

type GotifyDetails struct {
	ServerURL        string `json:"gotifyserverurl"`
	ApplicationToken string `json:"gotifyapplicationToken"`
	Priority         int    `json:"gotifyPriority"`
}

func (g Gotify) Type() string {
	return g.GotifyDetails.Type()
}

func (n GotifyDetails) Type() string {
	return "gotify"
}

func (g Gotify) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GotifyDetails, true))
}

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

func (g Gotify) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GotifyDetails)
}
