package notification

import (
	"fmt"
)

// Nostr represents a Nostr notification provider.
// Nostr is a decentralized social protocol for sending messages and notifications.
type Nostr struct {
	Base
	NostrDetails
}

// NostrDetails contains the configuration fields for Nostr notifications.
type NostrDetails struct {
	// Sender is the sender's private key in Nostr format (nsec encoded).
	Sender string `json:"sender"`
	// Recipients is a newline-delimited list of recipient public keys (npub encoded).
	Recipients string `json:"recipients"`
	// Relays is a newline-delimited list of Nostr relay URLs.
	Relays string `json:"relays"`
}

// Type returns the notification type identifier for Nostr.
func (n Nostr) Type() string {
	return n.NostrDetails.Type()
}

// Type returns the notification type identifier for NostrDetails.
func (n NostrDetails) Type() string {
	return "nostr"
}

// String returns a string representation of the Nostr notification.
func (n Nostr) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), formatNotification(n.NostrDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Nostr notification.
func (n *Nostr) UnmarshalJSON(data []byte) error {
	detail := NostrDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*n = Nostr{
		Base:         base,
		NostrDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Nostr notification into JSON.
func (n Nostr) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Base, n.NostrDetails)
}
