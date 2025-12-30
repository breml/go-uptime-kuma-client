package notification

import (
	"fmt"
)

// GrafanaOncall represents a grafanaoncall notification.
type GrafanaOncall struct {
	Base
	GrafanaOncallDetails
}

// GrafanaOncallDetails contains grafanaoncall-specific notification configuration.
type GrafanaOncallDetails struct {
	GrafanaOncallURL string `json:"GrafanaOncallURL"`
}

// Type returns the notification type.
func (g GrafanaOncall) Type() string {
	return g.GrafanaOncallDetails.Type()
}

// Type returns the notification type.
func (GrafanaOncallDetails) Type() string {
	return "GrafanaOncall"
}

// String returns a string representation of the notification.
func (g GrafanaOncall) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GrafanaOncallDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (g GrafanaOncall) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, &g.GrafanaOncallDetails)
}
