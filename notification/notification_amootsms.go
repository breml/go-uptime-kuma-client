package notification

import (
	"fmt"
)

// AmootSMS represents an Amoot SMS notification provider.
// Amoot SMS (https://portal.amootsms.com/) is an Iranian SMS gateway.
type AmootSMS struct {
	Base
	AmootSMSDetails
}

// AmootSMSDetails contains the configuration fields for Amoot SMS notifications.
//
// The send mode is selected by the combination of UsePattern and UseOwnLine:
//
//   - UsePattern unset or false: plain SMS via /rest/SendSimple, LineNumber is
//     required.
//   - UsePattern true, UseOwnLine unset or false: pattern SMS via
//     /rest/SendWithPattern, PatternCodeID is required and LineNumber is
//     ignored.
//   - UsePattern true, UseOwnLine true: pattern SMS from an own line via
//     /rest/SendWithPatternOWN, PatternCodeID and LineNumber are required.
//
// The server checks these requirements only when a notification is sent, an
// inconsistent configuration is accepted on create and update. A zero
// PatternCodeID or an empty LineNumber counts as missing.
type AmootSMSDetails struct {
	// APIToken authenticates against the gateway. It is a secret, the Uptime
	// Kuma form masks it.
	APIToken string `json:"amootApiToken"`
	// Mobiles holds one or more comma separated recipient numbers. The server
	// trims every entry and drops empty ones.
	Mobiles string `json:"amootMobiles"`
	// UsePattern sends the message as the values of a pre-registered pattern
	// instead of as plain text.
	UsePattern *bool `json:"amootUsePattern,omitempty"`
	// PatternCodeID is the id of the pattern registered with Amoot SMS, used
	// when UsePattern is set.
	PatternCodeID *int64 `json:"amootPatternCodeId,omitempty"`
	// UseOwnLine sends pattern messages from LineNumber, used when UsePattern
	// is set.
	UseOwnLine *bool `json:"amootUseOwnLine,omitempty"`
	// LineNumber is the sender line. The Uptime Kuma form defaults it to
	// "public".
	LineNumber *string `json:"amootLineNumber,omitempty"`
}

// Type returns the notification type identifier for Amoot SMS.
func (a AmootSMS) Type() string {
	return a.AmootSMSDetails.Type()
}

// Type returns the notification type identifier for AmootSMSDetails.
func (AmootSMSDetails) Type() string {
	return "amootsms"
}

// String returns a string representation of the Amoot SMS notification.
func (a AmootSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AmootSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an Amoot SMS notification.
func (a *AmootSMS) UnmarshalJSON(data []byte) error {
	detail := AmootSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*a = AmootSMS{
		Base:            base,
		AmootSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Amoot SMS notification into JSON.
func (a AmootSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, a.AmootSMSDetails)
}
