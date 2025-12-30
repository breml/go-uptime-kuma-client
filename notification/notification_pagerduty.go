package notification

import (
	"fmt"
)

// PagerDuty represents a pagerduty notification.
type PagerDuty struct {
	Base
	PagerDutyDetails
}

// PagerDutyDetails contains pagerduty-specific notification configuration.
type PagerDutyDetails struct {
	IntegrationURL string `json:"pagerdutyIntegrationUrl"`
	IntegrationKey string `json:"pagerdutyIntegrationKey"`
	Priority       string `json:"pagerdutyPriority"`
	AutoResolve    string `json:"pagerdutyAutoResolve"`
}

// Type returns the notification type.
func (p PagerDuty) Type() string {
	return p.PagerDutyDetails.Type()
}

// Type returns the notification type.
func (PagerDutyDetails) Type() string {
	return "PagerDuty"
}

// String returns a string representation of the notification.
func (p PagerDuty) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PagerDutyDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (p *PagerDuty) UnmarshalJSON(data []byte) error {
	detail := PagerDutyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = PagerDuty{
		Base:             base,
		PagerDutyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (p PagerDuty) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PagerDutyDetails)
}
