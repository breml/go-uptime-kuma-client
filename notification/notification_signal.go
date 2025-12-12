package notification

import (
	"fmt"
)

type Signal struct {
	Base
	SignalDetails
}

type SignalDetails struct {
	URL        string `json:"signalURL"`
	Number     string `json:"signalNumber"`
	Recipients string `json:"signalRecipients"`
}

func (s Signal) Type() string {
	return s.SignalDetails.Type()
}

func (n SignalDetails) Type() string {
	return "signal"
}

func (s Signal) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SignalDetails, true))
}

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

func (s Signal) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SignalDetails)
}
