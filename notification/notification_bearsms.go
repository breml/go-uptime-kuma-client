package notification

import (
	"fmt"
)

// BearSMS represents a BearSMS notification provider.
// BearSMS is an Israeli SMS gateway; the alert is sent as an HTTP GET with the
// credentials and the message in the query string.
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
	// SenderID is the sender name or number shown to the recipient, at most 11
	// characters.
	//
	// Upstream appends it to the request only when it is non-empty, so an
	// absent key and an empty string behave alike. The pointer keeps an unset
	// sender distinguishable from an explicitly empty one, so an edit does not
	// send back a value the server never stored.
	SenderID *string `json:"bearsmsSenderId,omitempty"`
	// PhoneNumber is the recipient number in international format, without a
	// leading plus, for example 9725XXXXXXXX.
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
