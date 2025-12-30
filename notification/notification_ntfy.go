package notification

import (
	"fmt"
)

// Ntfy ...
type Ntfy struct {
	Base
	NtfyDetails
}

// NtfyDetails ...
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

// Type ...
func (n Ntfy) Type() string {
	return n.NtfyDetails.Type()
}

// Type ...
func (NtfyDetails) Type() string {
	return "ntfy"
}

// String ...
func (n Ntfy) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NtfyDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (n Ntfy) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, &n.NtfyDetails)
}
