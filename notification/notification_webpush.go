package notification

import (
	"fmt"
)

// Webpush represents a Web Push notification.
// VAPID keys are managed server-side via the webpushPublicVapidKey /
// webpushPrivateVapidKey settings and are not part of this notification payload.
type Webpush struct {
	Base
	WebpushDetails
}

// WebpushDetails contains the Web Push-specific notification configuration.
type WebpushDetails struct {
	Subscription WebpushSubscription `json:"subscription"`
}

// WebpushSubscription is a W3C PushSubscription object identifying the push endpoint.
type WebpushSubscription struct {
	Endpoint string                  `json:"endpoint"`
	Keys     WebpushSubscriptionKeys `json:"keys"`
}

// WebpushSubscriptionKeys holds the encryption keys for a Web Push subscription.
type WebpushSubscriptionKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

// Type returns the notification type.
func (w Webpush) Type() string {
	return w.WebpushDetails.Type()
}

// Type returns the notification type.
func (WebpushDetails) Type() string {
	return "Webpush"
}

// String returns a string representation of the notification.
func (w Webpush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WebpushDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (w *Webpush) UnmarshalJSON(data []byte) error {
	detail := WebpushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = Webpush{
		Base:           base,
		WebpushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (w Webpush) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WebpushDetails)
}
