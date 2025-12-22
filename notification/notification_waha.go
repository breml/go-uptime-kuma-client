package notification

import (
	"fmt"
)

// WAHA represents a WAHA (WhatsApp HTTP API) notification provider.
// WAHA is a WhatsApp Business API solution for sending notifications via WhatsApp.
type WAHA struct {
	Base
	WAHADetails
}

// WAHADetails contains the configuration fields for WAHA notifications.
type WAHADetails struct {
	// ApiURL is the WAHA API endpoint URL.
	ApiURL string `json:"wahaApiUrl"`
	// Session is the WAHA session name.
	Session string `json:"wahaSession"`
	// ChatID is the recipient chat ID (typically a phone number).
	ChatID string `json:"wahaChatId"`
	// ApiKey is the WAHA API key for authentication (optional).
	ApiKey string `json:"wahaApiKey"`
}

// Type returns the notification type identifier for WAHA.
func (w WAHA) Type() string {
	return w.WAHADetails.Type()
}

// Type returns the notification type identifier for WAHADetails.
func (n WAHADetails) Type() string {
	return "waha"
}

// String returns a string representation of the WAHA notification.
func (w WAHA) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WAHADetails, true))
}

// UnmarshalJSON unmarshals JSON data into a WAHA notification.
func (w *WAHA) UnmarshalJSON(data []byte) error {
	detail := WAHADetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = WAHA{
		Base:        base,
		WAHADetails: detail,
	}

	return nil
}

// MarshalJSON marshals the WAHA notification into JSON.
func (w WAHA) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WAHADetails)
}
