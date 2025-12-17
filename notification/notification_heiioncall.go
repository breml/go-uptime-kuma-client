package notification

import (
	"fmt"
)

// HeiiOnCall represents a Heii On-Call notification provider.
// Heii On-Call is an on-call scheduling and incident management platform.
type HeiiOnCall struct {
	Base
	HeiiOnCallDetails
}

// HeiiOnCallDetails contains the configuration fields for Heii On-Call notifications.
type HeiiOnCallDetails struct {
	// APIKey is the Heii On-Call API key for authentication.
	APIKey string `json:"heiiOnCallApiKey"`
	// TriggerID is the Heii On-Call trigger ID.
	TriggerID string `json:"heiiOnCallTriggerId"`
}

// Type returns the notification type identifier for HeiiOnCall.
func (h HeiiOnCall) Type() string {
	return h.HeiiOnCallDetails.Type()
}

// Type returns the notification type identifier for HeiiOnCallDetails.
func (n HeiiOnCallDetails) Type() string {
	return "HeiiOnCall"
}

// String returns a string representation of the HeiiOnCall notification.
func (h HeiiOnCall) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(h.Base, false), formatNotification(h.HeiiOnCallDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a HeiiOnCall notification.
func (h *HeiiOnCall) UnmarshalJSON(data []byte) error {
	detail := HeiiOnCallDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*h = HeiiOnCall{
		Base:              base,
		HeiiOnCallDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the HeiiOnCall notification into JSON.
func (h HeiiOnCall) MarshalJSON() ([]byte, error) {
	return marshalJSON(h.Base, h.HeiiOnCallDetails)
}
