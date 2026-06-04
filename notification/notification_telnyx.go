package notification

import (
	"fmt"
)

// Telnyx represents a telnyx notification.
type Telnyx struct {
	Base
	TelnyxDetails
}

// TelnyxDetails contains telnyx-specific notification configuration.
type TelnyxDetails struct {
	APIKey             string  `json:"telnyxApiKey"`
	MessagingProfileID *string `json:"telnyxMessagingProfileId,omitempty"`
	PhoneNumber        string  `json:"telnyxPhoneNumber"`
	ToNumber           string  `json:"telnyxToNumber"`
}

// Type returns the notification type.
func (t Telnyx) Type() string {
	return t.TelnyxDetails.Type()
}

// Type returns the notification type.
func (TelnyxDetails) Type() string {
	return "telnyx"
}

// String returns a string representation of the notification.
func (t Telnyx) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TelnyxDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (t *Telnyx) UnmarshalJSON(data []byte) error {
	detail := TelnyxDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Telnyx{
		Base:          base,
		TelnyxDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (t Telnyx) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TelnyxDetails)
}
