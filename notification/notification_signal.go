package notification

import (
	"fmt"
)

// Signal represents a signal notification.
type Signal struct {
	Base
	SignalDetails
}

// SignalDetails contains signal-specific notification configuration.
type SignalDetails struct {
	URL        string `json:"signalURL"`
	Number     string `json:"signalNumber"`
	Recipients string `json:"signalRecipients"`
}

// Type returns the notification type.
func (s Signal) Type() string {
	return s.SignalDetails.Type()
}

// Type returns the notification type.
func (SignalDetails) Type() string {
	return "signal"
}

// String returns a string representation of the notification.
func (s Signal) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SignalDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (s *Signal) UnmarshalJSON(data []byte) error {
	detail := SignalDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Signal{
		Base:          base,
		SignalDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (s Signal) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, &s.SignalDetails)
}
