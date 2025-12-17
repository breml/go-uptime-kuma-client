package notification

import (
	"fmt"
)

// PagerTree represents a PagerTree notification provider.
// PagerTree is an incident management platform for sending and managing alerts.
type PagerTree struct {
	Base
	PagerTreeDetails
}

// PagerTreeDetails contains the configuration fields for PagerTree notifications.
type PagerTreeDetails struct {
	// IntegrationUrl is the PagerTree integration endpoint URL.
	IntegrationUrl string `json:"pagertreeIntegrationUrl"`
	// Urgency is the urgency level of the alert (e.g., "high", "medium", "low").
	Urgency string `json:"pagertreeUrgency"`
	// AutoResolve determines if alerts should be auto-resolved when monitor recovers.
	// Valid values: "resolve" to auto-resolve, empty string or other to not auto-resolve.
	AutoResolve string `json:"pagertreeAutoResolve"`
}

// Type returns the notification type identifier for PagerTree.
func (p PagerTree) Type() string {
	return p.PagerTreeDetails.Type()
}

// Type returns the notification type identifier for PagerTreeDetails.
func (p PagerTreeDetails) Type() string {
	return "PagerTree"
}

// String returns a string representation of the PagerTree notification.
func (p PagerTree) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PagerTreeDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a PagerTree notification.
func (p *PagerTree) UnmarshalJSON(data []byte) error {
	detail := PagerTreeDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = PagerTree{
		Base:             base,
		PagerTreeDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the PagerTree notification into JSON.
func (p PagerTree) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PagerTreeDetails)
}
