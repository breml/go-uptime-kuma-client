package notification

import (
	"fmt"
)

type PagerDuty struct {
	Base
	PagerDutyDetails
}

type PagerDutyDetails struct {
	IntegrationURL string `json:"pagerdutyIntegrationUrl"`
	IntegrationKey string `json:"pagerdutyIntegrationKey"`
	Priority       string `json:"pagerdutyPriority"`
	AutoResolve    string `json:"pagerdutyAutoResolve"`
}

func (p PagerDuty) Type() string {
	return p.PagerDutyDetails.Type()
}

func (n PagerDutyDetails) Type() string {
	return "PagerDuty"
}

func (p PagerDuty) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PagerDutyDetails, true))
}

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

func (p PagerDuty) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PagerDutyDetails)
}
