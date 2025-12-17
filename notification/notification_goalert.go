package notification

import (
	"fmt"
)

// GoAlert represents a GoAlert notification provider.
// GoAlert is an on-call scheduling and alerting platform.
type GoAlert struct {
	Base
	GoAlertDetails
}

// GoAlertDetails contains the configuration fields for GoAlert notifications.
type GoAlertDetails struct {
	// BaseURL is the GoAlert base URL.
	BaseURL string `json:"goAlertBaseURL"`
	// Token is the GoAlert integration token.
	Token string `json:"goAlertToken"`
}

// Type returns the notification type identifier for GoAlert.
func (g GoAlert) Type() string {
	return g.GoAlertDetails.Type()
}

// Type returns the notification type identifier for GoAlertDetails.
func (n GoAlertDetails) Type() string {
	return "GoAlert"
}

// String returns a string representation of the GoAlert notification.
func (g GoAlert) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GoAlertDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a GoAlert notification.
func (g *GoAlert) UnmarshalJSON(data []byte) error {
	detail := GoAlertDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GoAlert{
		Base:           base,
		GoAlertDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the GoAlert notification into JSON.
func (g GoAlert) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GoAlertDetails)
}
