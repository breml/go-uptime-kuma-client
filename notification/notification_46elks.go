package notification

import (
	"fmt"
)

// FortySixElks represents a fortysixelks notification.
type FortySixElks struct {
	Base
	FortySixElksDetails
}

// FortySixElksDetails contains fortysixelks-specific notification configuration.
type FortySixElksDetails struct {
	Username   string `json:"elksUsername"`
	AuthToken  string `json:"elksAuthToken"`
	FromNumber string `json:"elksFromNumber"`
	ToNumber   string `json:"elksToNumber"`
}

// Type returns the notification type.
func (f FortySixElks) Type() string {
	return f.FortySixElksDetails.Type()
}

// Type returns the notification type.
func (FortySixElksDetails) Type() string {
	return "46elks"
}

// String returns a string representation of the notification.
func (f FortySixElks) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FortySixElksDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (f *FortySixElks) UnmarshalJSON(data []byte) error {
	detail := FortySixElksDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = FortySixElks{
		Base:                base,
		FortySixElksDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (f FortySixElks) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, &f.FortySixElksDetails)
}
