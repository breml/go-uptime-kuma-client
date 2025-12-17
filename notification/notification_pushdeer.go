package notification

import (
	"fmt"
)

// PushDeer represents a PushDeer notification provider.
// PushDeer is a push notification service that supports custom push key and optional server URL.
type PushDeer struct {
	Base
	PushDeerDetails
}

// PushDeerDetails contains the PushDeer-specific configuration.
type PushDeerDetails struct {
	Key    string `json:"pushdeerKey"`    // PushDeer push key
	Server string `json:"pushdeerServer"` // Custom server URL (optional, defaults to https://api2.pushdeer.com)
}

// Type returns the notification type identifier for PushDeer.
func (p PushDeer) Type() string {
	return p.PushDeerDetails.Type()
}

// Type returns the notification type identifier for PushDeer details.
func (p PushDeerDetails) Type() string {
	return "PushDeer"
}

// String returns a human-readable representation of the PushDeer notification.
func (p PushDeer) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushDeerDetails, true))
}

// UnmarshalJSON deserializes a PushDeer notification from JSON.
func (p *PushDeer) UnmarshalJSON(data []byte) error {
	detail := PushDeerDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = PushDeer{
		Base:            base,
		PushDeerDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a PushDeer notification to JSON.
func (p PushDeer) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PushDeerDetails)
}
