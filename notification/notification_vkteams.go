package notification

import (
	"fmt"
)

// VKTeams represents a VKTeams notification.
type VKTeams struct {
	Base
	VKTeamsDetails
}

// VKTeamsDetails contains VKTeams-specific notification configuration.
type VKTeamsDetails struct {
	BotToken       string `json:"vkteamsBotToken"`
	ChatID         string `json:"vkteamsChatId"`
	BaseURL        string `json:"vkteamsBaseUrl"`
	UseTemplate    bool   `json:"vkteamsUseTemplate"`
	Template       string `json:"vkteamsTemplate"`
	TemplateFormat string `json:"vkteamsTemplateFormat"`
}

// Type returns the notification type.
func (v VKTeams) Type() string {
	return v.VKTeamsDetails.Type()
}

// Type returns the notification type.
func (VKTeamsDetails) Type() string {
	return "VKTeams"
}

// String returns a string representation of the notification.
func (v VKTeams) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(v.Base, false), formatNotification(v.VKTeamsDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (v *VKTeams) UnmarshalJSON(data []byte) error {
	detail := VKTeamsDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*v = VKTeams{
		Base:           base,
		VKTeamsDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (v VKTeams) MarshalJSON() ([]byte, error) {
	return marshalJSON(v.Base, &v.VKTeamsDetails)
}
