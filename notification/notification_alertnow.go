package notification

import (
	"fmt"
)

// AlertNow ...
type AlertNow struct {
	Base
	AlertNowDetails
}

// AlertNowDetails ...
type AlertNowDetails struct {
	WebhookURL string `json:"alertNowWebhookURL"`
}

// Type ...
func (a AlertNow) Type() string {
	return a.AlertNowDetails.Type()
}

// Type ...
func (AlertNowDetails) Type() string {
	return "AlertNow"
}

// String ...
func (a AlertNow) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AlertNowDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (a AlertNow) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, &a.AlertNowDetails)
}
