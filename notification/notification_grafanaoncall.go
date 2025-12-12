package notification

import (
	"fmt"
)

type GrafanaOncall struct {
	Base
	GrafanaOncallDetails
}

type GrafanaOncallDetails struct {
	GrafanaOncallURL string `json:"GrafanaOncallURL"`
}

func (g GrafanaOncall) Type() string {
	return g.GrafanaOncallDetails.Type()
}

func (n GrafanaOncallDetails) Type() string {
	return "GrafanaOncall"
}

func (g GrafanaOncall) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GrafanaOncallDetails, true))
}

func (g *GrafanaOncall) UnmarshalJSON(data []byte) error {
	detail := GrafanaOncallDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GrafanaOncall{
		Base:                   base,
		GrafanaOncallDetails: detail,
	}

	return nil
}

func (g GrafanaOncall) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GrafanaOncallDetails)
}
