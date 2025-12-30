package notification

import (
	"fmt"
)

// Evolution represents an Evolution API notification provider.
// Evolution API is a WhatsApp Business API solution for sending notifications via WhatsApp.
type Evolution struct {
	Base
	EvolutionDetails
}

// EvolutionDetails contains the configuration fields for Evolution API notifications.
type EvolutionDetails struct {
	// APIURL is the Evolution API URL endpoint.
	APIURL string `json:"evolutionApiUrl"`
	// InstanceName is the Evolution API instance name.
	InstanceName string `json:"evolutionInstanceName"`
	// AuthToken is the Evolution API authentication token.
	AuthToken string `json:"evolutionAuthToken"`
	// Recipient is the recipient phone number for WhatsApp messages.
	Recipient string `json:"evolutionRecipient"`
}

// Type returns the notification type identifier for Evolution.
func (e Evolution) Type() string {
	return e.EvolutionDetails.Type()
}

// Type returns the notification type identifier for EvolutionDetails.
func (EvolutionDetails) Type() string {
	return "EvolutionApi"
}

// String returns a string representation of the Evolution notification.
func (e Evolution) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(e.Base, false), formatNotification(e.EvolutionDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an Evolution notification.
func (e *Evolution) UnmarshalJSON(data []byte) error {
	detail := EvolutionDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*e = Evolution{
		Base:             base,
		EvolutionDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Evolution notification into JSON.
func (e Evolution) MarshalJSON() ([]byte, error) {
	return marshalJSON(e.Base, e.EvolutionDetails)
}
