package notification

import (
	"fmt"
)

// FreeMobile represents a Free Mobile SMS notification provider.
// Free Mobile is an SMS service provider for sending text message notifications.
type FreeMobile struct {
	Base
	FreeMobileDetails
}

// FreeMobileDetails contains the configuration fields for Free Mobile SMS notifications.
type FreeMobileDetails struct {
	// User is the Free Mobile user ID.
	User string `json:"freemobileUser"`
	// Pass is the Free Mobile API key.
	Pass string `json:"freemobilePass"`
}

// Type returns the notification type identifier for FreeMobile.
func (f FreeMobile) Type() string {
	return f.FreeMobileDetails.Type()
}

// Type returns the notification type identifier for FreeMobileDetails.
func (FreeMobileDetails) Type() string {
	return "FreeMobile"
}

// String returns a string representation of the FreeMobile notification.
func (f FreeMobile) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FreeMobileDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a FreeMobile notification.
func (f *FreeMobile) UnmarshalJSON(data []byte) error {
	detail := FreeMobileDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = FreeMobile{
		Base:              base,
		FreeMobileDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the FreeMobile notification into JSON.
func (f FreeMobile) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, f.FreeMobileDetails)
}
