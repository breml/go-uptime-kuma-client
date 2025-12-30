package notification

import (
	"fmt"
)

// Threema represents a Threema notification provider.
// Threema is a privacy-focused messaging service for sending notifications.
type Threema struct {
	Base
	ThreemaDetails
}

// ThreemaDetails contains the configuration fields for Threema notifications.
type ThreemaDetails struct {
	// SenderIdentity is the Threema Gateway ID for the sender.
	SenderIdentity string `json:"threemaSenderIdentity"`
	// Secret is the Threema API secret for authentication.
	Secret string `json:"threemaSecret"`
	// Recipient is the recipient identifier (ID, phone number, or email address).
	Recipient string `json:"threemaRecipient"`
	// RecipientType specifies the type of recipient (identity, phone, or email).
	RecipientType string `json:"threemaRecipientType"`
}

// Type returns the notification type identifier for Threema.
func (t Threema) Type() string {
	return t.ThreemaDetails.Type()
}

// Type returns the notification type identifier for ThreemaDetails.
func (ThreemaDetails) Type() string {
	return "threema"
}

// String returns a string representation of the Threema notification.
func (t Threema) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.ThreemaDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Threema notification.
func (t *Threema) UnmarshalJSON(data []byte) error {
	detail := ThreemaDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Threema{
		Base:           base,
		ThreemaDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Threema notification into JSON.
func (t Threema) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.ThreemaDetails)
}
