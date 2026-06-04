package notification

import (
	"fmt"
)

// Max represents a MAX messenger notification.
type Max struct {
	Base
	MaxDetails
}

// MaxDetails contains MAX messenger-specific notification configuration.
type MaxDetails struct {
	APIURL         string `json:"maxApiUrl"`
	BotToken       string `json:"maxBotToken"`
	ChatID         string `json:"maxChatID"`
	UseTemplate    bool   `json:"maxUseTemplate"`
	Template       string `json:"maxTemplate"`
	TemplateFormat string `json:"maxTemplateFormat"`
}

// Type returns the notification type.
func (m Max) Type() string {
	return m.MaxDetails.Type()
}

// Type returns the notification type.
func (MaxDetails) Type() string {
	return "max"
}

// String returns a string representation of the notification.
func (m Max) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MaxDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (m *Max) UnmarshalJSON(data []byte) error {
	detail := MaxDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*m = Max{
		Base:       base,
		MaxDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (m Max) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, &m.MaxDetails)
}
