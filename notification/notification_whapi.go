package notification

import (
	"fmt"
)

// Whapi represents a Whapi (WhatsApp API) notification provider.
// Whapi is a WhatsApp API service for sending notifications via WhatsApp.
type Whapi struct {
	Base
	WhapiDetails
}

// WhapiDetails contains the configuration fields for Whapi notifications.
type WhapiDetails struct {
	// APIURL is the Whapi API endpoint URL.
	APIURL string `json:"whapiApiUrl"`
	// AuthToken is the Whapi API authorization token.
	AuthToken string `json:"whapiAuthToken"`
	// Recipient is the recipient phone number.
	Recipient string `json:"whapiRecipient"`
}

// Type returns the notification type identifier for Whapi.
func (w Whapi) Type() string {
	return w.WhapiDetails.Type()
}

// Type returns the notification type identifier for WhapiDetails.
func (n WhapiDetails) Type() string {
	return "whapi"
}

// String returns a string representation of the Whapi notification.
func (w Whapi) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WhapiDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Whapi notification.
func (w *Whapi) UnmarshalJSON(data []byte) error {
	detail := WhapiDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = Whapi{
		Base:         base,
		WhapiDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Whapi notification into JSON.
func (w Whapi) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WhapiDetails)
}
