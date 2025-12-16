package notification

import (
	"fmt"
)

// ClickSendSMS represents a ClickSend SMS notification provider.
// ClickSend SMS is an SMS service provider for sending text message notifications.
type ClickSendSMS struct {
	Base
	ClickSendSMSDetails
}

// ClickSendSMSDetails contains the configuration fields for ClickSend SMS notifications.
type ClickSendSMSDetails struct {
	// Login is the ClickSend account username.
	Login string `json:"clicksendsmsLogin"`
	// Password is the ClickSend API key.
	Password string `json:"clicksendsmsPassword"`
	// ToNumber is the recipient phone number.
	ToNumber string `json:"clicksendsmsToNumber"`
	// SenderName is the sender name or phone number.
	SenderName string `json:"clicksendsmsSenderName"`
}

// Type returns the notification type identifier for ClickSendSMS.
func (c ClickSendSMS) Type() string {
	return c.ClickSendSMSDetails.Type()
}

// Type returns the notification type identifier for ClickSendSMSDetails.
func (n ClickSendSMSDetails) Type() string {
	return "clicksendsms"
}

// String returns a string representation of the ClickSendSMS notification.
func (c ClickSendSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(c.Base, false), formatNotification(c.ClickSendSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a ClickSendSMS notification.
func (c *ClickSendSMS) UnmarshalJSON(data []byte) error {
	detail := ClickSendSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*c = ClickSendSMS{
		Base:                base,
		ClickSendSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the ClickSendSMS notification into JSON.
func (c ClickSendSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Base, c.ClickSendSMSDetails)
}
