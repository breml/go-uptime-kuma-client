package notification

import (
	"fmt"
)

// GTXMessaging represents a GTX Messaging SMS notification provider.
// GTX Messaging is an SMS service provider for sending text message notifications.
type GTXMessaging struct {
	Base
	GTXMessagingDetails
}

// GTXMessagingDetails contains the configuration fields for GTX Messaging SMS notifications.
type GTXMessagingDetails struct {
	// ApiKey is the GTX Messaging API key.
	ApiKey string `json:"gtxMessagingApiKey"`
	// From is the sender ID.
	From string `json:"gtxMessagingFrom"`
	// To is the recipient phone number.
	To string `json:"gtxMessagingTo"`
}

// Type returns the notification type identifier for GTXMessaging.
func (g GTXMessaging) Type() string {
	return g.GTXMessagingDetails.Type()
}

// Type returns the notification type identifier for GTXMessagingDetails.
func (n GTXMessagingDetails) Type() string {
	return "gtxmessaging"
}

// String returns a string representation of the GTXMessaging notification.
func (g GTXMessaging) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GTXMessagingDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a GTXMessaging notification.
func (g *GTXMessaging) UnmarshalJSON(data []byte) error {
	detail := GTXMessagingDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GTXMessaging{
		Base:                base,
		GTXMessagingDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the GTXMessaging notification into JSON.
func (g GTXMessaging) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GTXMessagingDetails)
}
