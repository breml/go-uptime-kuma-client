package notification

import (
	"fmt"
)

type AlertNow struct {
	Base
	AlertNowDetails
}

type AlertNowDetails struct {
	WebhookURL string `json:"alertNowWebhookURL"`
}

func (a AlertNow) Type() string {
	return a.AlertNowDetails.Type()
}

func (n AlertNowDetails) Type() string {
	return "AlertNow"
}

func (a AlertNow) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AlertNowDetails, true))
}

func (a *AlertNow) UnmarshalJSON(data []byte) error {
	detail := AlertNowDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*a = AlertNow{
		Base:            base,
		AlertNowDetails: detail,
	}

	return nil
}

func (a AlertNow) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, a.AlertNowDetails)
}
