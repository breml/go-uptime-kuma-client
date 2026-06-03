package notification

import (
	"fmt"
)

// Teltonika represents a Teltonika notification provider.
// Teltonika uses a Teltonika RUTxxx series router as an SMS gateway over its HTTP API
// (https://developers.teltonika-networks.com/reference/rut241/7.19.4/v1.11.1/messages).
// This provider is only compatible with Teltonika RutOS >= 7.14.0 devices.
type Teltonika struct {
	Base
	TeltonikaDetails
}

// TeltonikaDetails contains the configuration fields for Teltonika notifications.
type TeltonikaDetails struct {
	// URL is the Teltonika router base URL (e.g., https://192.168.100.1 or http://teltonika.example.com:8080).
	URL string `json:"teltonikaUrl"`
	// Username is the router user used to authenticate against the Teltonika HTTP API.
	Username string `json:"teltonikaUsername"`
	// Password is the password for the router user.
	Password string `json:"teltonikaPassword"`
	// Modem is the modem identifier on the router (e.g., "1-1").
	Modem string `json:"teltonikaModem"`
	// PhoneNumber is the recipient phone number (e.g., "+336xxxxxxxx").
	// Multiple recipients can be provided as a comma-separated list.
	PhoneNumber string `json:"teltonikaPhoneNumber"`
	// UnsafeTLS disables TLS certificate validation when communicating with the router.
	// This is useful for routers using a self-signed certificate.
	UnsafeTLS bool `json:"teltonikaUnsafeTls"`
}

// Type returns the notification type identifier for Teltonika.
func (t Teltonika) Type() string {
	return t.TeltonikaDetails.Type()
}

// Type returns the notification type identifier for TeltonikaDetails.
func (TeltonikaDetails) Type() string {
	return "Teltonika"
}

// String returns a string representation of the Teltonika notification.
func (t Teltonika) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TeltonikaDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Teltonika notification.
func (t *Teltonika) UnmarshalJSON(data []byte) error {
	detail := TeltonikaDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Teltonika{
		Base:             base,
		TeltonikaDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Teltonika notification into JSON.
func (t Teltonika) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.TeltonikaDetails)
}
