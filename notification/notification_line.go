package notification

import (
	"fmt"
)

// Line represents a LINE Messaging API notification provider.
// LINE is a messaging platform for sending notifications to LINE users.
type Line struct {
	Base
	LineDetails
}

// LineDetails contains the configuration fields for LINE notifications.
type LineDetails struct {
	// ChannelAccessToken is the LINE channel access token for authentication.
	ChannelAccessToken string `json:"lineChannelAccessToken"`
	// UserID is the LINE user ID to send messages to.
	UserID string `json:"lineUserID"`
}

// Type returns the notification type identifier for Line.
func (l Line) Type() string {
	return l.LineDetails.Type()
}

// Type returns the notification type identifier for LineDetails.
func (n LineDetails) Type() string {
	return "line"
}

// String returns a string representation of the Line notification.
func (l Line) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(l.Base, false), formatNotification(l.LineDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Line notification.
func (l *Line) UnmarshalJSON(data []byte) error {
	detail := LineDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*l = Line{
		Base:        base,
		LineDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Line notification into JSON.
func (l Line) MarshalJSON() ([]byte, error) {
	return marshalJSON(l.Base, l.LineDetails)
}
