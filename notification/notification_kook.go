package notification

import (
	"fmt"
)

// Kook represents a KOOK (formerly KaiHeiLa) notification provider.
// KOOK is a communication platform popular in China for messaging and community management.
type Kook struct {
	Base
	KookDetails
}

// KookDetails contains the configuration fields for KOOK notifications.
type KookDetails struct {
	// BotToken is the KOOK bot token for authentication.
	BotToken string `json:"kookBotToken"`
	// GuildID is the KOOK guild/server ID where messages will be sent.
	GuildID string `json:"kookGuildID"`
}

// Type returns the notification type identifier for Kook.
func (k Kook) Type() string {
	return k.KookDetails.Type()
}

// Type returns the notification type identifier for KookDetails.
func (KookDetails) Type() string {
	return "Kook"
}

// String returns a string representation of the Kook notification.
func (k Kook) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(k.Base, false), formatNotification(k.KookDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Kook notification.
func (k *Kook) UnmarshalJSON(data []byte) error {
	detail := KookDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*k = Kook{
		Base:        base,
		KookDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Kook notification into JSON.
func (k Kook) MarshalJSON() ([]byte, error) {
	return marshalJSON(k.Base, k.KookDetails)
}
