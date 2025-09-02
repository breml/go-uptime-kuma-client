package notification

import (
	"fmt"
)

type Ntfy struct {
	Base
	NtfyDetails
}

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

func (n Ntfy) Type() string {
	return n.NtfyDetails.Type()
}

func (n NtfyDetails) Type() string {
	return "ntfy"
}

func (n Ntfy) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NtfyDetails, true))
}

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

func (n Ntfy) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, n.NtfyDetails)
}
