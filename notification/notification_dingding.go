package notification

import (
	"fmt"
)

// DingDing represents a DingDing (钉钉) notification provider.
// DingDing is a Chinese enterprise messaging platform that supports webhooks for notifications.
type DingDing struct {
	Base
	DingDingDetails
}

// DingDingDetails contains the DingDing-specific configuration.
type DingDingDetails struct {
	WebHookURL string `json:"webHookUrl"`
	SecretKey  string `json:"secretKey"`
	Mentioning string `json:"mentioning"`
}

// Type returns the notification type identifier for DingDing.
func (d DingDing) Type() string {
	return d.DingDingDetails.Type()
}

// Type returns the notification type identifier for DingDing details.
func (DingDingDetails) Type() string {
	return "DingDing"
}

// String returns a human-readable representation of the DingDing notification.
func (d DingDing) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(d.Base, false), formatNotification(d.DingDingDetails, true))
}

// UnmarshalJSON deserializes a DingDing notification from JSON.
func (d *DingDing) UnmarshalJSON(data []byte) error {
	detail := DingDingDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*d = DingDing{
		Base:            base,
		DingDingDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a DingDing notification to JSON.
func (d DingDing) MarshalJSON() ([]byte, error) {
	return marshalJSON(d.Base, d.DingDingDetails)
}
