package notification

import (
	"fmt"
)

type Alerta struct {
	Base
	AlertaDetails
}

type AlertaDetails struct {
	APIEndpoint  string `json:"alertaApiEndpoint"`
	APIKey       string `json:"alertaApiKey"`
	Environment  string `json:"alertaEnvironment"`
	AlertState   string `json:"alertaAlertState"`
	RecoverState string `json:"alertaRecoverState"`
}

func (a Alerta) Type() string {
	return a.AlertaDetails.Type()
}

func (n AlertaDetails) Type() string {
	return "alerta"
}

func (a Alerta) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AlertaDetails, true))
}

func (a *Alerta) UnmarshalJSON(data []byte) error {
	detail := AlertaDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*a = Alerta{
		Base:          base,
		AlertaDetails: detail,
	}

	return nil
}

func (a Alerta) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, a.AlertaDetails)
}
