package notification

import (
	"fmt"
)

// SIGNL4 represents a SIGNL4 notification provider.
// SIGNL4 is an alert and incident management platform for sending notifications.
type SIGNL4 struct {
	Base
	SIGNL4Details
}

// SIGNL4Details contains the configuration fields for SIGNL4 notifications.
type SIGNL4Details struct {
	// WebhookURL is the SIGNL4 webhook URL for receiving notifications.
	WebhookURL string `json:"webhookURL"`
}

// Type returns the notification type identifier for SIGNL4.
func (s SIGNL4) Type() string {
	return s.SIGNL4Details.Type()
}

// Type returns the notification type identifier for SIGNL4Details.
func (SIGNL4Details) Type() string {
	return "SIGNL4"
}

// String returns a string representation of the SIGNL4 notification.
func (s SIGNL4) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SIGNL4Details, true))
}

// UnmarshalJSON unmarshals JSON data into a SIGNL4 notification.
func (s *SIGNL4) UnmarshalJSON(data []byte) error {
	detail := SIGNL4Details{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SIGNL4{
		Base:          base,
		SIGNL4Details: detail,
	}

	return nil
}

// MarshalJSON marshals the SIGNL4 notification into JSON.
func (s SIGNL4) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SIGNL4Details)
}
