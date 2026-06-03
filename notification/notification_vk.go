package notification

import (
	"fmt"
)

// VK represents a VK notification.
type VK struct {
	Base
	VKDetails
}

// VKDetails contains VK-specific notification configuration.
type VKDetails struct {
	AccessToken    string `json:"vkAccessToken"`
	PeerID         string `json:"vkPeerId"`
	APIVersion     string `json:"vkApiVersion"`
	DontParseLinks bool   `json:"vkDontParseLinks"`
}

// Type returns the notification type.
func (v VK) Type() string {
	return v.VKDetails.Type()
}

// Type returns the notification type.
func (VKDetails) Type() string {
	return "VK"
}

// String returns a string representation of the notification.
func (v VK) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(v.Base, false), formatNotification(v.VKDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (v *VK) UnmarshalJSON(data []byte) error {
	detail := VKDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*v = VK{
		Base:      base,
		VKDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (v VK) MarshalJSON() ([]byte, error) {
	return marshalJSON(v.Base, &v.VKDetails)
}
