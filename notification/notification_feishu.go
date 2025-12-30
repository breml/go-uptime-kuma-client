package notification

import (
	"fmt"
)

// Feishu represents a Feishu notification provider.
// Feishu is a Chinese enterprise collaboration platform that supports webhooks for notifications.
type Feishu struct {
	Base
	FeishuDetails
}

// FeishuDetails contains the Feishu-specific configuration.
type FeishuDetails struct {
	WebHookURL string `json:"feishuWebHookUrl"`
}

// Type returns the notification type identifier for Feishu.
func (f Feishu) Type() string {
	return f.FeishuDetails.Type()
}

// Type returns the notification type identifier for Feishu details.
func (FeishuDetails) Type() string {
	return "Feishu"
}

// String returns a human-readable representation of the Feishu notification.
func (f Feishu) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FeishuDetails, true))
}

// UnmarshalJSON deserializes a Feishu notification from JSON.
func (f *Feishu) UnmarshalJSON(data []byte) error {
	detail := FeishuDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = Feishu{
		Base:          base,
		FeishuDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a Feishu notification to JSON.
func (f Feishu) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, f.FeishuDetails)
}
