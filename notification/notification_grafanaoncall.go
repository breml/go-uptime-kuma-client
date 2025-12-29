package notification

import (
	"fmt"
)

// GrafanaOncall ...
type GrafanaOncall struct {
	Base
	GrafanaOncallDetails
}

// GrafanaOncallDetails ...
type GrafanaOncallDetails struct {
	GrafanaOncallURL string `json:"GrafanaOncallURL"`
}

// Type ...
func (g GrafanaOncall) Type() string {
	return g.GrafanaOncallDetails.Type()
}

// Type ...
func (n GrafanaOncallDetails) Type() string {
	return "GrafanaOncall"
}

// String ...
func (g GrafanaOncall) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GrafanaOncallDetails, true))
}

// UnmarshalJSON ...
func (g *GrafanaOncall) UnmarshalJSON(data []byte) error {
	detail := GrafanaOncallDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GrafanaOncall{
		Base:                 base,
		GrafanaOncallDetails: detail,
	}

	return nil
}

// MarshalJSON ...
func (g GrafanaOncall) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, &g.GrafanaOncallDetails)
}
