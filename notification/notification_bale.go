package notification

import (
	"fmt"
)

// Bale represents a Bale messenger notification.
type Bale struct {
	Base
	BaleDetails
}

// BaleDetails contains bale-specific notification configuration.
type BaleDetails struct {
	BotToken string `json:"baleBotToken"`
	ChatID   string `json:"baleChatID"`
}

// Type returns the notification type.
func (b Bale) Type() string {
	return b.BaleDetails.Type()
}

// Type returns the notification type.
func (BaleDetails) Type() string {
	return "bale"
}

// String returns a string representation of the notification.
func (b Bale) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(b.Base, false), formatNotification(b.BaleDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (b *Bale) UnmarshalJSON(data []byte) error {
	detail := BaleDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*b = Bale{
		Base:        base,
		BaleDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (b Bale) MarshalJSON() ([]byte, error) {
	return marshalJSON(b.Base, &b.BaleDetails)
}
