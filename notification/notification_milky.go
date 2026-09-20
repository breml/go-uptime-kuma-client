package notification

import (
	"fmt"
)

// Milky represents a Milky (QQ) notification provider.
// Milky (https://milky.ntqqrev.org/) is a QQ bot protocol implementation, the
// alert is posted as a single text message segment to a QQ group or to a
// private chat, prefixed with "UptimeKuma Alert: ".
type Milky struct {
	Base
	MilkyDetails
}

// MilkyDetails contains the configuration fields for Milky notifications.
//
// The wire field names are shared with the OneBot provider, so both configs
// look identical on the wire and are only told apart by their type. A config
// stored for one of the two therefore unmarshals into the other without an
// error, the type is only rewritten on the next marshal.
type MilkyDetails struct {
	// HTTPAddr is the address of the Milky service (e.g., "http://localhost:3000").
	// The server prepends "http://" if the value does not start with "http".
	//
	// It then appends "/api/" if the value does not end with "/", before it
	// appends the endpoint. A trailing slash therefore does not merely avoid a
	// double slash, it skips the "/api/" segment altogether and the alert is
	// posted to an endpoint that does not exist. The address is stored
	// verbatim, a trailing slash survives the round trip.
	HTTPAddr string `json:"httpAddr"`
	// AccessToken is the access token sent as bearer token to the Milky service.
	AccessToken string `json:"accessToken"`
	// MsgType selects the endpoint the alert is sent to and whether the payload
	// carries a group or a user ID. The server treats anything other than
	// MilkyMessageTypeGroup as a private message.
	//
	// Unlike the other fields Uptime Kuma marks as required, the message type
	// carries omitempty. The form stores no default, so a notification created
	// in the UI has no msgType key at all, and the server reads an absent and
	// an empty message type exactly like MilkyMessageTypePrivate. An empty
	// message type read from the server is therefore written back as an absent
	// key.
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
