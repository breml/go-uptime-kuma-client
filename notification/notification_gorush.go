package notification

import (
	"fmt"
)

// Gorush represents a Gorush push notification provider.
// Gorush is a push notification server for sending notifications to mobile devices.
type Gorush struct {
	Base
	GorushDetails
}

// GorushDetails contains the configuration fields for Gorush notifications.
type GorushDetails struct {
	// ServerURL is the Gorush server URL.
	ServerURL string `json:"gorushServerURL"`
	// DeviceToken is the device token for push notifications.
	DeviceToken string `json:"gorushDeviceToken"`
	// Platform is the target platform (ios, android, huawei).
	Platform string `json:"gorushPlatform"`
	// Title is the notification title (optional).
	Title string `json:"gorushTitle"`
	// Priority is the notification priority level.
	Priority string `json:"gorushPriority"`
	// Retry is the number of retries for failed deliveries.
	Retry int `json:"gorushRetry"`
	// Topic is the APNs topic for iOS notifications.
	Topic string `json:"gorushTopic"`
}

// Type returns the notification type identifier for Gorush.
func (g Gorush) Type() string {
	return g.GorushDetails.Type()
}

// Type returns the notification type identifier for GorushDetails.
func (GorushDetails) Type() string {
	return "gorush"
}

// String returns a string representation of the Gorush notification.
func (g Gorush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GorushDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Gorush notification.
func (g *Gorush) UnmarshalJSON(data []byte) error {
	detail := GorushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = Gorush{
		Base:          base,
		GorushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Gorush notification into JSON.
func (g Gorush) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GorushDetails)
}
