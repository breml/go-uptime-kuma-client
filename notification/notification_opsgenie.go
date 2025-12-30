package notification

import (
	"fmt"
)

// Opsgenie represents a opsgenie notification.
type Opsgenie struct {
	Base
	OpsgenieDetails
}

// OpsgenieDetails contains opsgenie-specific notification configuration.
type OpsgenieDetails struct {
	APIKey   string `json:"opsgenieApiKey"`
	Region   string `json:"opsgenieRegion"`
	Priority int    `json:"opsgeniePriority"`
}

// Type returns the notification type.
func (o Opsgenie) Type() string {
	return o.OpsgenieDetails.Type()
}

// Type returns the notification type.
func (OpsgenieDetails) Type() string {
	return "Opsgenie"
}

// String returns a string representation of the notification.
func (o Opsgenie) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OpsgenieDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (o *Opsgenie) UnmarshalJSON(data []byte) error {
	detail := OpsgenieDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = Opsgenie{
		Base:            base,
		OpsgenieDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (o Opsgenie) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, &o.OpsgenieDetails)
}
