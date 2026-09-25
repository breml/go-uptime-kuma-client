package notification

import (
	"fmt"
)

// OpenWa represents an OpenWA notification provider.
// OpenWA (https://www.open-wa.org/) is a self-hosted WhatsApp gateway. The
// server posts the alert to "<apiUrl>/api/sessions/<session>/messages/send-text",
// once per chat ID.
type OpenWa struct {
	Base
	OpenWaDetails
}

// OpenWaDetails contains the configuration fields for OpenWA notifications.
type OpenWaDetails struct {
	// APIURL is the base URL of the OpenWA gateway (e.g., "http://localhost:2785/").
	// The server strips trailing slashes before it appends the endpoint. The URL
	// is stored verbatim, a trailing slash survives the round trip.
	APIURL string `json:"openwaApiUrl"`
	// APIKey is the OpenWA API key, sent as "X-Api-Key" header.
	APIKey string `json:"openwaApiKey"`
	// Session is the name of the OpenWA session, URL-encoded into the path.
	Session string `json:"openwaSession"`
	// ChatID is a comma-separated list of WhatsApp chat IDs
	// (e.g., "00117612345678@c.us,123456789012345678@g.us,1234567890@lid").
	// The server trims each entry and drops empty ones, sending fails if no
	// valid chat ID remains.
	ChatID string `json:"openwaChatId"`
	// UseCustomMessage enables the use of a custom message template.
	UseCustomMessage *bool `json:"openwaUseCustomMessage,omitempty"`
	// CustomMessage is the custom message template. The server only renders it
	// if UseCustomMessage is enabled and the trimmed template is not empty,
	// otherwise the default message is sent.
	CustomMessage *string `json:"openwaCustomMessage,omitempty"`
}

// Type returns the notification type identifier for OpenWa.
func (o OpenWa) Type() string {
	return o.OpenWaDetails.Type()
}

// Type returns the notification type identifier for OpenWaDetails.
//
// Note: this is the provider name of Uptime Kuma
// (server/notification-providers/openwa.js) and deliberately does not follow the
// casing of the Go type. Do not "fix" it, the server dispatches on this exact
// string.
func (OpenWaDetails) Type() string {
	return "openwa"
}

// String returns a string representation of the OpenWa notification.
func (o OpenWa) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OpenWaDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an OpenWa notification.
func (o *OpenWa) UnmarshalJSON(data []byte) error {
	detail := OpenWaDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = OpenWa{
		Base:          base,
		OpenWaDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the OpenWa notification into JSON.
func (o OpenWa) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OpenWaDetails)
}
