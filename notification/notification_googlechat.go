package notification

import (
	"fmt"
)

// GoogleChat represents a Google Chat notification provider.
// Google Chat is a messaging platform that supports webhooks for notifications.
type GoogleChat struct {
	Base
	GoogleChatDetails
}

// GoogleChatDetails contains the Google Chat-specific configuration.
type GoogleChatDetails struct {
	WebhookURL  string `json:"googleChatWebhookURL"`
	UseTemplate bool   `json:"googleChatUseTemplate"`
	Template    string `json:"googleChatTemplate"`
}

// Type returns the notification type identifier for Google Chat.
func (g GoogleChat) Type() string {
	return g.GoogleChatDetails.Type()
}

// Type returns the notification type identifier for Google Chat details.
func (GoogleChatDetails) Type() string {
	return "GoogleChat"
}

// String returns a human-readable representation of the Google Chat notification.
func (g GoogleChat) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GoogleChatDetails, true))
}

// UnmarshalJSON deserializes a Google Chat notification from JSON.
func (g *GoogleChat) UnmarshalJSON(data []byte) error {
	detail := GoogleChatDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GoogleChat{
		Base:              base,
		GoogleChatDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a Google Chat notification to JSON.
func (g GoogleChat) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GoogleChatDetails)
}
