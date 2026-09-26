package notification

import (
	"fmt"
)

// NotifyApp represents a Notify! push notification provider.
// Notify! (https://getnotifyapp.com/) is a push notification app. The server
// posts the alert as JSON to
// "https://push.getnotifyapp.com/notify-json/<deviceId>?token=<token>"; the
// title and the time sensitive flag for DOWN heartbeats are derived by the
// server and not configurable.
type NotifyApp struct {
	Base
	NotifyAppDetails
}

// NotifyAppDetails contains the configuration fields for Notify! notifications.
type NotifyAppDetails struct {
	// DeviceID is the Notify! device or group id the push is sent to.
	DeviceID string `json:"notifyAppDeviceId"`
	// Token authenticates the push. It is a secret, the Uptime Kuma form
	// masks it.
	Token string `json:"notifyAppToken"`
	// IconURL is sent as "iconUrl" only when non-empty.
	IconURL *string `json:"notifyAppIconUrl,omitempty"`
}

// Type returns the notification type identifier for Notify!.
func (n NotifyApp) Type() string {
	return n.NotifyAppDetails.Type()
}

// Type returns the notification type identifier for NotifyAppDetails.
func (NotifyAppDetails) Type() string {
	return "notifyapp"
}

// String returns a string representation of the Notify! notification.
func (n NotifyApp) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NotifyAppDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Notify! notification.
func (n *NotifyApp) UnmarshalJSON(data []byte) error {
	detail := NotifyAppDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*n = NotifyApp{
		Base:             base,
		NotifyAppDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Notify! notification into JSON.
func (n NotifyApp) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, n.NotifyAppDetails)
}
