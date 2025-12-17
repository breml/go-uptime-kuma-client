package notification

import (
	"fmt"
)

// LineNotify represents a LINE Notify notification provider.
// LINE Notify is a notification service for sending messages to LINE users.
type LineNotify struct {
	Base
	LineNotifyDetails
}

// LineNotifyDetails contains the configuration fields for LINE Notify notifications.
type LineNotifyDetails struct {
	// AccessToken is the LINE Notify access token for authentication.
	AccessToken string `json:"lineNotifyAccessToken"`
}

// Type returns the notification type identifier for LineNotify.
func (l LineNotify) Type() string {
	return l.LineNotifyDetails.Type()
}

// Type returns the notification type identifier for LineNotifyDetails.
func (n LineNotifyDetails) Type() string {
	return "LineNotify"
}

// String returns a string representation of the LineNotify notification.
func (l LineNotify) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(l.Base, false), formatNotification(l.LineNotifyDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a LineNotify notification.
func (l *LineNotify) UnmarshalJSON(data []byte) error {
	detail := LineNotifyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*l = LineNotify{
		Base:              base,
		LineNotifyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the LineNotify notification into JSON.
func (l LineNotify) MarshalJSON() ([]byte, error) {
	return marshalJSON(l.Base, l.LineNotifyDetails)
}
