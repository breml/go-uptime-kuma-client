package notification

import (
	"fmt"
)

// Signalgrid represents a Signalgrid push notification provider.
type Signalgrid struct {
	Base
	SignalgridDetails
}

// SignalgridDetails contains the configuration fields for Signalgrid notifications.
type SignalgridDetails struct {
	// ClientKey is the Signalgrid client key.
	ClientKey string `json:"signalgridClientKey"`
	// Channel is the Signalgrid channel the push is sent to.
	Channel string `json:"signalgridChannel"`
}

// Type returns the notification type identifier for Signalgrid.
func (s Signalgrid) Type() string {
	return s.SignalgridDetails.Type()
}

// Type returns the notification type identifier for SignalgridDetails.
func (SignalgridDetails) Type() string {
	return "signalgrid"
}

// String returns a string representation of the Signalgrid notification.
func (s Signalgrid) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SignalgridDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Signalgrid notification.
func (s *Signalgrid) UnmarshalJSON(data []byte) error {
	detail := SignalgridDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Signalgrid{
		Base:              base,
		SignalgridDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Signalgrid notification into JSON.
func (s Signalgrid) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SignalgridDetails)
}
