package notification

import (
	"fmt"
)

// BearSMS represents a BearSMS notification provider.
// BearSMS (https://app.bearsms.com/) is an SMS gateway, the alert is sent as
// an HTTP GET with the credentials and the message in the query string.
type BearSMS struct {
	Base
	BearSMSDetails
}

// BearSMSDetails contains the configuration fields for BearSMS notifications.
type BearSMSDetails struct {
	// Username is the BearSMS account username.
	Username string `json:"bearsmsUsername"`
	// HashKey is the secret belonging to the account.
	HashKey string `json:"bearsmsHashKey"`
	// SenderID is the sender name or number shown to the recipient. Uptime
	// Kuma's form caps it at 11 characters and notes that the sender has to be
	// approved in the BearSMS account, neither of which this client or the
	// server enforces.
	//
	// Upstream appends it to the request only when it is non-empty, so an
	// absent key and an empty string behave alike and the pointer buys no
	// semantics, only round trip fidelity: it keeps an unset sender
	// distinguishable from an explicitly empty one, so an edit does not rewrite
	// a stored empty sender as an absent one.
	SenderID *string `json:"bearsmsSenderId,omitempty"`
	// PhoneNumber is the recipient number with country code and without a "+"
	// or "00" prefix, for example 9725XXXXXXXX.
	PhoneNumber string `json:"bearsmsPhoneNumber"`
}

// Type returns the notification type identifier for BearSMS.
func (b BearSMS) Type() string {
	return b.BearSMSDetails.Type()
}

// Type returns the notification type identifier for BearSMSDetails.
func (BearSMSDetails) Type() string {
	return "bearsms"
}

// String returns a string representation of the BearSMS notification.
func (b BearSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(b.Base, false), formatNotification(b.BearSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a BearSMS notification.
func (b *BearSMS) UnmarshalJSON(data []byte) error {
	detail := BearSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*b = BearSMS{
		Base:           base,
		BearSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the BearSMS notification into JSON.
func (b BearSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(b.Base, b.BearSMSDetails)
}
