package notification

import (
	"fmt"
)

// Signal ...
type Signal struct {
	Base
	SignalDetails
}

// SignalDetails ...
type SignalDetails struct {
	URL        string `json:"signalURL"`
	Number     string `json:"signalNumber"`
	Recipients string `json:"signalRecipients"`
}

// Type ...
func (s Signal) Type() string {
	return s.SignalDetails.Type()
}

// Type ...
func (n SignalDetails) Type() string {
	return "signal"
}

// String ...
func (s Signal) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SignalDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (s Signal) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, &s.SignalDetails)
}
