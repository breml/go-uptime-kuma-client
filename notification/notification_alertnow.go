package notification

import (
	"fmt"
)

// AlertNow represents a alertnow notification.
type AlertNow struct {
	Base
	AlertNowDetails
}

// AlertNowDetails contains alertnow-specific notification configuration.
type AlertNowDetails struct {
	WebhookURL string `json:"alertNowWebhookURL"`
}

// Type returns the notification type.
func (a AlertNow) Type() string {
	return a.AlertNowDetails.Type()
}

// Type returns the notification type.
func (AlertNowDetails) Type() string {
	return "AlertNow"
}

// String returns a string representation of the notification.
func (a AlertNow) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AlertNowDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (a AlertNow) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, &a.AlertNowDetails)
}
