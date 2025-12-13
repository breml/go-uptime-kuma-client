package notification

import (
	"fmt"
)

type HomeAssistant struct {
	Base
	HomeAssistantDetails
}

type HomeAssistantDetails struct {
	HomeAssistantURL     string `json:"homeAssistantUrl"`
	LongLivedAccessToken string `json:"longLivedAccessToken"`
	NotificationService  string `json:"notificationService"`
}

func (h HomeAssistant) Type() string {
	return h.HomeAssistantDetails.Type()
}

func (n HomeAssistantDetails) Type() string {
	return "HomeAssistant"
}

func (h HomeAssistant) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(h.Base, false), formatNotification(h.HomeAssistantDetails, true))
}

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

func (h HomeAssistant) MarshalJSON() ([]byte, error) {
	return marshalJSON(h.Base, h.HomeAssistantDetails)
}
