package notification

import (
	"fmt"
)

// CallMeBot represents a CallMeBot notification provider.
// CallMeBot is a service for sending notifications via HTTP GET requests.
type CallMeBot struct {
	Base
	CallMeBotDetails
}

// CallMeBotDetails contains the configuration fields for CallMeBot notifications.
type CallMeBotDetails struct {
	// Endpoint is the CallMeBot endpoint URL where notifications will be sent.
	Endpoint string `json:"callMeBotEndpoint"`
}

// Type returns the notification type identifier for CallMeBot.
func (c CallMeBot) Type() string {
	return c.CallMeBotDetails.Type()
}

// Type returns the notification type identifier for CallMeBotDetails.
func (CallMeBotDetails) Type() string {
	return "CallMeBot"
}

// String returns a string representation of the CallMeBot notification.
func (c CallMeBot) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(c.Base, false), formatNotification(c.CallMeBotDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a CallMeBot notification.
func (c *CallMeBot) UnmarshalJSON(data []byte) error {
	detail := CallMeBotDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*c = CallMeBot{
		Base:             base,
		CallMeBotDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the CallMeBot notification into JSON.
func (c CallMeBot) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Base, c.CallMeBotDetails)
}
