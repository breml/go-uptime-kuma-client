package notification

import (
	"fmt"
)

// Milky represents a Milky (QQ) notification provider.
// Milky is a QQ bot protocol implementation, allowing notifications to be sent to QQ groups or private users.
type Milky struct {
	Base
	MilkyDetails
}

// MilkyDetails contains the configuration fields for Milky notifications.
//
// The wire field names are shared with the OneBot provider, so both configs
// look identical on the wire and are only told apart by their type.
type MilkyDetails struct {
	// HTTPAddr is the address of the Milky service (e.g., "http://localhost:3000").
	// The server prepends "http://" if the value does not start with "http" and
	// appends "/api/" if it does not end with "/".
	HTTPAddr string `json:"httpAddr"`
	// AccessToken is the access token sent as bearer token to the Milky service.
	AccessToken string `json:"accessToken"`
	// MsgType selects the endpoint the alert is sent to and whether the payload
	// carries a group or a user ID. The server treats anything other than
	// MilkyMessageTypeGroup as a private message.
	MsgType MilkyMessageType `json:"msgType,omitempty"`
	// ReceiverID is the QQ group ID (when MsgType is MilkyMessageTypeGroup) or
	// user ID (when MsgType is MilkyMessageTypePrivate).
	ReceiverID string `json:"recieverId"`
}

// MilkyMessageType represents the message type of a Milky notification.
type MilkyMessageType string

// Milky message types.
const (
	MilkyMessageTypeGroup   MilkyMessageType = "group"
	MilkyMessageTypePrivate MilkyMessageType = "private"
)

// String returns the string representation of the Milky message type.
func (t MilkyMessageType) String() string {
	return string(t)
}

// Type returns the notification type identifier for Milky.
func (m Milky) Type() string {
	return m.MilkyDetails.Type()
}

// Type returns the notification type identifier for MilkyDetails.
func (MilkyDetails) Type() string {
	return "Milky"
}

// String returns a string representation of the Milky notification.
func (m Milky) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MilkyDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Milky notification.
func (m *Milky) UnmarshalJSON(data []byte) error {
	detail := MilkyDetails{}

	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*m = Milky{
		Base:         base,
		MilkyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Milky notification into JSON.
func (m Milky) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, m.MilkyDetails)
}
