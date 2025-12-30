package notification

import (
	"fmt"
)

// SpugPush represents a SpugPush notification provider.
// SpugPush is a push notification service for sending alerts via spug.cc.
type SpugPush struct {
	Base
	SpugPushDetails
}

// SpugPushDetails contains the configuration fields for SpugPush notifications.
type SpugPushDetails struct {
	// TemplateKey is the SpugPush template key for the notification.
	TemplateKey string `json:"templateKey"`
}

// Type returns the notification type identifier for SpugPush.
func (s SpugPush) Type() string {
	return s.SpugPushDetails.Type()
}

// Type returns the notification type identifier for SpugPushDetails.
func (SpugPushDetails) Type() string {
	return "SpugPush"
}

// String returns a string representation of the SpugPush notification.
func (s SpugPush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SpugPushDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a SpugPush notification.
func (s *SpugPush) UnmarshalJSON(data []byte) error {
	detail := SpugPushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SpugPush{
		Base:            base,
		SpugPushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SpugPush notification into JSON.
func (s SpugPush) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SpugPushDetails)
}
