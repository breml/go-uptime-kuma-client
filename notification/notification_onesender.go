package notification

import (
	"fmt"
)

// OneSender represents a OneSender WhatsApp notification provider.
// OneSender is a WhatsApp messaging service that enables notifications through WhatsApp channels.
type OneSender struct {
	Base
	OneSenderDetails
}

// OneSenderDetails contains the configuration fields for OneSender notifications.
type OneSenderDetails struct {
	// URL is the OneSender API endpoint URL for sending messages.
	URL string `json:"onesenderURL"`
	// Token is the API token for authentication with OneSender service.
	Token string `json:"onesenderToken"`
	// Receiver is the recipient identifier (phone number or group ID).
	Receiver string `json:"onesenderReceiver"`
	// TypeReceiver is the type of receiver: "private" for individual or "group" for group.
	TypeReceiver string `json:"onesenderTypeReceiver"`
}

// Type returns the notification type identifier for OneSender.
func (o OneSender) Type() string {
	return o.OneSenderDetails.Type()
}

// Type returns the notification type identifier for OneSenderDetails.
func (o OneSenderDetails) Type() string {
	return "onesender"
}

// String returns a string representation of the OneSender notification.
func (o OneSender) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OneSenderDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a OneSender notification.
func (o *OneSender) UnmarshalJSON(data []byte) error {
	detail := OneSenderDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = OneSender{
		Base:             base,
		OneSenderDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the OneSender notification into JSON.
func (o OneSender) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OneSenderDetails)
}
