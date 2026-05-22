package notification

import (
	"fmt"
)

// Whatsapp360messenger represents a 360messenger WhatsApp notification provider.
// 360messenger is a WhatsApp messaging service for sending notifications.
type Whatsapp360messenger struct {
	Base
	Whatsapp360messengerDetails
}

// Whatsapp360messengerDetails contains the configuration fields for 360messenger notifications.
type Whatsapp360messengerDetails struct {
	// AuthToken is the Bearer authentication token.
	AuthToken string `json:"Whatsapp360messengerAuthToken"`
	// Recipient is a comma- or semicolon-separated list of phone numbers.
	Recipient string `json:"Whatsapp360messengerRecipient"`
	// GroupIDs is the list of WhatsApp group IDs (multi-select).
	GroupIDs []string `json:"Whatsapp360messengerGroupIds,omitzero"`
	// GroupID is the legacy single group ID kept for back-compat in upstream payloads.
	GroupID *string `json:"Whatsapp360messengerGroupId,omitempty"`
	// UseTemplate enables the use of a custom message template.
	UseTemplate *bool `json:"Whatsapp360messengerUseTemplate,omitempty"`
	// Template is the custom message template used when UseTemplate is enabled.
	Template *string `json:"Whatsapp360messengerTemplate,omitempty"`
}

// Type returns the notification type identifier for Whatsapp360messenger.
func (w Whatsapp360messenger) Type() string {
	return w.Whatsapp360messengerDetails.Type()
}

// Type returns the notification type identifier for Whatsapp360messengerDetails.
func (Whatsapp360messengerDetails) Type() string {
	return "Whatsapp360messenger"
}

// String returns a string representation of the Whatsapp360messenger notification.
func (w Whatsapp360messenger) String() string {
	return fmt.Sprintf(
		"%s, %s",
		formatNotification(w.Base, false),
		formatNotification(w.Whatsapp360messengerDetails, true),
	)
}

// UnmarshalJSON unmarshals JSON data into a Whatsapp360messenger notification.
func (w *Whatsapp360messenger) UnmarshalJSON(data []byte) error {
	detail := Whatsapp360messengerDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = Whatsapp360messenger{
		Base:                        base,
		Whatsapp360messengerDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Whatsapp360messenger notification into JSON.
func (w Whatsapp360messenger) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.Whatsapp360messengerDetails)
}
