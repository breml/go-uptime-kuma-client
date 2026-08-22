package notification

import (
	"fmt"
)

// Plivo represents a plivo notification.
type Plivo struct {
	Base
	PlivoDetails
}

// PlivoDetails contains plivo-specific notification configuration.
type PlivoDetails struct {
	AuthID     string `json:"plivoAuthID"`
	AuthToken  string `json:"plivoAuthToken"`
	FromNumber string `json:"plivoFromNumber"`
	ToNumber   string `json:"plivoToNumber"`
	// MessageType selects whether the alert is delivered as an SMS or as a
	// voice call. If unset, the server delivers the alert as an SMS.
	MessageType PlivoMessageType `json:"plivoMessageType,omitempty"`
	// AnswerURL is fetched by Plivo with an HTTP GET to obtain the Plivo XML
	// driving the call. It is only used and only required if MessageType is
	// PlivoMessageTypeCall, where it must be an absolute URL. The server sets
	// the alert text as the "message" query parameter, replacing any "message"
	// parameter already present.
	AnswerURL *string `json:"plivoAnswerUrl,omitempty"`
}

// PlivoMessageType represents the delivery type of a plivo notification.
type PlivoMessageType string

// Plivo message types.
const (
	PlivoMessageTypeSMS  PlivoMessageType = "sms"
	PlivoMessageTypeCall PlivoMessageType = "call"
)

// String returns the string representation of the plivo message type.
func (t PlivoMessageType) String() string {
	return string(t)
}

// Type returns the notification type.
func (p Plivo) Type() string {
	return p.PlivoDetails.Type()
}

// Type returns the notification type.
func (PlivoDetails) Type() string {
	return "plivo"
}

// String returns a string representation of the notification.
func (p Plivo) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PlivoDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (p *Plivo) UnmarshalJSON(data []byte) error {
	detail := PlivoDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = Plivo{
		Base:         base,
		PlivoDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (p Plivo) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, &p.PlivoDetails)
}
