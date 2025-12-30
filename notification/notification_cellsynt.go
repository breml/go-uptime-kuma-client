package notification

import (
	"fmt"
)

// Cellsynt represents a Cellsynt SMS notification provider.
// Cellsynt is an SMS service provider for sending text message notifications.
type Cellsynt struct {
	Base
	CellsyntDetails
}

// CellsyntDetails contains the configuration fields for Cellsynt SMS notifications.
type CellsyntDetails struct {
	// Login is the Cellsynt account username.
	Login string `json:"cellsyntLogin"`
	// Password is the Cellsynt account password.
	Password string `json:"cellsyntPassword"`
	// Destination is the recipient phone number.
	Destination string `json:"cellsyntDestination"`
	// Originator is the sender name or phone number.
	Originator string `json:"cellsyntOriginator"`
	// OriginatorType specifies the type of originator (Numeric or Alphanumeric).
	OriginatorType string `json:"cellsyntOriginatortype"`
	// AllowLongSMS allows sending SMS messages longer than 160 characters.
	AllowLongSMS bool `json:"cellsyntAllowLongSMS"`
}

// Type returns the notification type identifier for Cellsynt.
func (c Cellsynt) Type() string {
	return c.CellsyntDetails.Type()
}

// Type returns the notification type identifier for CellsyntDetails.
func (CellsyntDetails) Type() string {
	return "Cellsynt"
}

// String returns a string representation of the Cellsynt notification.
func (c Cellsynt) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(c.Base, false), formatNotification(c.CellsyntDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Cellsynt notification.
func (c *Cellsynt) UnmarshalJSON(data []byte) error {
	detail := CellsyntDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*c = Cellsynt{
		Base:            base,
		CellsyntDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Cellsynt notification into JSON.
func (c Cellsynt) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Base, c.CellsyntDetails)
}
