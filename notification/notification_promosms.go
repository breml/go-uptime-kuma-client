package notification

import (
	"fmt"
)

// PromoSMS represents a PromoSMS notification provider.
// PromoSMS is an SMS gateway service for sending SMS notifications.
type PromoSMS struct {
	Base
	PromoSMSDetails
}

// PromoSMSDetails contains the configuration fields for PromoSMS notifications.
type PromoSMSDetails struct {
	// Login is the PromoSMS login/username for authentication.
	Login string `json:"promosmsLogin"`
	// Password is the PromoSMS password for authentication.
	Password string `json:"promosmsPassword"`
	// PhoneNumber is the recipient phone number for SMS messages.
	PhoneNumber string `json:"promosmsPhoneNumber"`
	// SenderName is the sender name/ID displayed in the SMS.
	SenderName string `json:"promosmsSenderName"`
	// SMSType is the SMS type identifier (as a string representation of a number).
	SMSType string `json:"promosmsSMSType"`
	// AllowLongSMS indicates whether long SMS messages (up to 4 SMS) are allowed.
	// When true, messages can be up to 639 characters; when false, limited to 159 characters.
	AllowLongSMS bool `json:"promosmsAllowLongSMS"`
}

// Type returns the notification type identifier for PromoSMS.
func (p PromoSMS) Type() string {
	return p.PromoSMSDetails.Type()
}

// Type returns the notification type identifier for PromoSMSDetails.
func (PromoSMSDetails) Type() string {
	return "promosms"
}

// String returns a string representation of the PromoSMS notification.
func (p PromoSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(p.Base, false), formatNotification(p.PromoSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a PromoSMS notification.
func (p *PromoSMS) UnmarshalJSON(data []byte) error {
	detail := PromoSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*p = PromoSMS{
		Base:            base,
		PromoSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the PromoSMS notification into JSON.
func (p PromoSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(p.Base, p.PromoSMSDetails)
}
