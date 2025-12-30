package notification

import (
	"fmt"
)

// PushPlus represents a PushPlus notification provider.
// PushPlus is a push notification service that requires a send key for authentication.
type PushPlus struct {
	Base
	PushPlusDetails
}

// PushPlusDetails contains the PushPlus-specific configuration.
type PushPlusDetails struct {
	SendKey string `json:"pushPlusSendKey"` // PushPlus send key for authentication
}

// Type returns the notification type identifier for PushPlus.
func (p PushPlus) Type() string {
	return p.PushPlusDetails.Type()
}

// Type returns the notification type identifier for PushPlus details.
func (PushPlusDetails) Type() string {
	return "PushPlus"
}

// String returns a human-readable representation of the PushPlus notification.
func (p PushPlus) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PushPlusDetails, true))
}

// UnmarshalJSON deserializes a PushPlus notification from JSON.
func (p *PushPlus) UnmarshalJSON(data []byte) error {
	detail := PushPlusDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = PushPlus{
		Base:            base,
		PushPlusDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a PushPlus notification to JSON.
func (p PushPlus) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PushPlusDetails)
}
