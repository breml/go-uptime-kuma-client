package notification

import (
	"fmt"
)

// HomeAssistant represents a homeassistant notification.
type HomeAssistant struct {
	Base
	HomeAssistantDetails
}

// HomeAssistantDetails contains homeassistant-specific notification configuration.
type HomeAssistantDetails struct {
	HomeAssistantURL     string `json:"homeAssistantUrl"`
	LongLivedAccessToken string `json:"longLivedAccessToken"`
	NotificationService  string `json:"notificationService"`
}

// Type returns the notification type.
func (h HomeAssistant) Type() string {
	return h.HomeAssistantDetails.Type()
}

// Type returns the notification type.
func (HomeAssistantDetails) Type() string {
	return "HomeAssistant"
}

// String returns a string representation of the notification.
func (h HomeAssistant) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(h.Base, false), formatNotification(h.HomeAssistantDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (h *HomeAssistant) UnmarshalJSON(data []byte) error {
	detail := HomeAssistantDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*h = HomeAssistant{
		Base:                 base,
		HomeAssistantDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (h HomeAssistant) MarshalJSON() ([]byte, error) {
	return marshalJSON(h.Base, &h.HomeAssistantDetails)
}
