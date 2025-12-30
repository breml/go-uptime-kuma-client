package notification

import (
	"fmt"
)

// Opsgenie ...
type Opsgenie struct {
	Base
	OpsgenieDetails
}

// OpsgenieDetails ...
type OpsgenieDetails struct {
	APIKey   string `json:"opsgenieApiKey"`
	Region   string `json:"opsgenieRegion"`
	Priority int    `json:"opsgeniePriority"`
}

// Type ...
func (o Opsgenie) Type() string {
	return o.OpsgenieDetails.Type()
}

// Type ...
func (OpsgenieDetails) Type() string {
	return "Opsgenie"
}

// String ...
func (o Opsgenie) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OpsgenieDetails, true))
}

// UnmarshalJSON ...
func (o *Opsgenie) UnmarshalJSON(data []byte) error {
	detail := OpsgenieDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = Opsgenie{
		Base:            base,
		OpsgenieDetails: detail,
	}

	return nil
}

// MarshalJSON ...
func (o Opsgenie) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, &o.OpsgenieDetails)
}
