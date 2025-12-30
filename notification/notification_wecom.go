package notification

import (
	"fmt"
)

// WeCom represents a WeCom (WeChat Work) notification provider.
// WeCom is a Chinese enterprise communication platform that supports webhooks for notifications.
type WeCom struct {
	Base
	WeComDetails
}

// WeComDetails contains the WeCom-specific configuration.
type WeComDetails struct {
	BotKey string `json:"weComBotKey"`
}

// Type returns the notification type identifier for WeCom.
func (w WeCom) Type() string {
	return w.WeComDetails.Type()
}

// Type returns the notification type identifier for WeCom details.
func (WeComDetails) Type() string {
	return "WeCom"
}

// String returns a human-readable representation of the WeCom notification.
func (w WeCom) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WeComDetails, true))
}

// UnmarshalJSON deserializes a WeCom notification from JSON.
func (w *WeCom) UnmarshalJSON(data []byte) error {
	detail := WeComDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = WeCom{
		Base:         base,
		WeComDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a WeCom notification to JSON.
func (w WeCom) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WeComDetails)
}
