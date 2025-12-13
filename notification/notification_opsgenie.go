package notification

import (
	"fmt"
)

type Opsgenie struct {
	Base
	OpsgenieDetails
}

type OpsgenieDetails struct {
	ApiKey   string `json:"opsgenieApiKey"`
	Region   string `json:"opsgenieRegion"`
	Priority int    `json:"opsgeniePriority"`
}

func (o Opsgenie) Type() string {
	return o.OpsgenieDetails.Type()
}

func (n OpsgenieDetails) Type() string {
	return "Opsgenie"
}

func (o Opsgenie) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OpsgenieDetails, true))
}

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

func (o Opsgenie) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OpsgenieDetails)
}
