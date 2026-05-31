package notification

import (
	"fmt"
)

// HaloPSA represents a HaloPSA notification.
type HaloPSA struct {
	Base
	HaloPSADetails
}

// HaloPSADetails contains HaloPSA-specific notification configuration.
type HaloPSADetails struct {
	WebhookURL string `json:"halowebhookurl"`
	Username   string `json:"haloUsername"`
	Password   string `json:"haloPassword"`
}

// Type returns the notification type.
func (h HaloPSA) Type() string {
	return h.HaloPSADetails.Type()
}

// Type returns the notification type.
func (HaloPSADetails) Type() string {
	return "HaloPSA"
}

// String returns a string representation of the notification.
func (h HaloPSA) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(h.Base, false), formatNotification(h.HaloPSADetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (h *HaloPSA) UnmarshalJSON(data []byte) error {
	detail := HaloPSADetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*h = HaloPSA{
		Base:           base,
		HaloPSADetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (h HaloPSA) MarshalJSON() ([]byte, error) {
	return marshalJSON(h.Base, &h.HaloPSADetails)
}
