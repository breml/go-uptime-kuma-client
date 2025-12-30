package notification

import (
	"fmt"
)

// Alerta represents a alerta notification.
type Alerta struct {
	Base
	AlertaDetails
}

// AlertaDetails contains alerta-specific notification configuration.
type AlertaDetails struct {
	APIEndpoint  string `json:"alertaApiEndpoint"`
	APIKey       string `json:"alertaApiKey"`
	Environment  string `json:"alertaEnvironment"`
	AlertState   string `json:"alertaAlertState"`
	RecoverState string `json:"alertaRecoverState"`
}

// Type returns the notification type.
func (a Alerta) Type() string {
	return a.AlertaDetails.Type()
}

// Type returns the notification type.
func (AlertaDetails) Type() string {
	return "alerta"
}

// String returns a string representation of the notification.
func (a Alerta) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AlertaDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (a Alerta) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, &a.AlertaDetails)
}
