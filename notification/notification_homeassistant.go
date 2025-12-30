package notification

import (
	"fmt"
)

// HomeAssistant ...
type HomeAssistant struct {
	Base
	HomeAssistantDetails
}

// HomeAssistantDetails ...
type HomeAssistantDetails struct {
	HomeAssistantURL     string `json:"homeAssistantUrl"`
	LongLivedAccessToken string `json:"longLivedAccessToken"`
	NotificationService  string `json:"notificationService"`
}

// Type ...
func (h HomeAssistant) Type() string {
	return h.HomeAssistantDetails.Type()
}

// Type ...
func (HomeAssistantDetails) Type() string {
	return "HomeAssistant"
}

// String ...
func (h HomeAssistant) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(h.Base, false), formatNotification(h.HomeAssistantDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (h HomeAssistant) MarshalJSON() ([]byte, error) {
	return marshalJSON(h.Base, &h.HomeAssistantDetails)
}
