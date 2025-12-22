package notification

import (
	"fmt"
)

// WPush represents a WPush notification provider.
// WPush is a push notification service for sending notifications via the WPush API.
type WPush struct {
	Base
	WPushDetails
}

// WPushDetails contains the configuration fields for WPush notifications.
type WPushDetails struct {
	// APIKey is the WPush API key for authentication.
	APIKey string `json:"wpushAPIkey"`
	// Channel is the WPush push channel identifier.
	Channel string `json:"wpushChannel"`
}

// Type returns the notification type identifier for WPush.
func (w WPush) Type() string {
	return w.WPushDetails.Type()
}

// Type returns the notification type identifier for WPushDetails.
func (n WPushDetails) Type() string {
	return "WPush"
}

// String returns a string representation of the WPush notification.
func (w WPush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WPushDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a WPush notification.
func (w *WPush) UnmarshalJSON(data []byte) error {
	detail := WPushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = WPush{
		Base:         base,
		WPushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the WPush notification into JSON.
func (w WPush) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WPushDetails)
}
