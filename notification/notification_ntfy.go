package notification

import (
	"fmt"
)

// Ntfy represents a ntfy notification.
type Ntfy struct {
	Base
	NtfyDetails
}

// NtfyDetails contains ntfy-specific notification configuration.
type NtfyDetails struct {
	AccessToken          string `json:"ntfyaccesstoken"`
	AuthenticationMethod string `json:"ntfyAuthenticationMethod"`
	Icon                 string `json:"ntfyIcon"`
	Password             string `json:"ntfypassword"`
	Priority             int64  `json:"ntfyPriority"`
	ServerURL            string `json:"ntfyserverurl"`
	Topic                string `json:"ntfytopic"`
	Username             string `json:"ntfyusername"`
}

// Type returns the notification type.
func (n Ntfy) Type() string {
	return n.NtfyDetails.Type()
}

// Type returns the notification type.
func (NtfyDetails) Type() string {
	return "ntfy"
}

// String returns a string representation of the notification.
func (n Ntfy) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NtfyDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (n *Ntfy) UnmarshalJSON(data []byte) error {
	detail := NtfyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*n = Ntfy{
		Base:        base,
		NtfyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (n Ntfy) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, &n.NtfyDetails)
}
