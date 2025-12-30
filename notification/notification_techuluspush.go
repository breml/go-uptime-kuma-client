package notification

import (
	"fmt"
)

// TechulusPush represents a Techulus Push notification provider.
// TechulusPush is a push notification service for sending alerts via Techulus API.
type TechulusPush struct {
	Base
	TechulusPushDetails
}

// TechulusPushDetails contains the configuration fields for TechulusPush notifications.
type TechulusPushDetails struct {
	// APIKey is the Techulus Push API key for authentication.
	APIKey string `json:"pushAPIKey"`
	// Title is the notification title.
	Title string `json:"pushTitle"`
	// Sound is the notification sound identifier.
	Sound string `json:"pushSound"`
	// Channel is the push notification channel.
	Channel string `json:"pushChannel"`
	// TimeSensitive indicates if the notification is time-sensitive.
	TimeSensitive bool `json:"pushTimeSensitive"`
}

// Type returns the notification type identifier for TechulusPush.
func (t TechulusPush) Type() string {
	return t.TechulusPushDetails.Type()
}

// Type returns the notification type identifier for TechulusPushDetails.
func (TechulusPushDetails) Type() string {
	return "PushByTechulus"
}

// String returns a string representation of the TechulusPush notification.
func (t TechulusPush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TechulusPushDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a TechulusPush notification.
func (t *TechulusPush) UnmarshalJSON(data []byte) error {
	detail := TechulusPushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = TechulusPush{
		Base:                base,
		TechulusPushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the TechulusPush notification into JSON.
func (t TechulusPush) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.TechulusPushDetails)
}
