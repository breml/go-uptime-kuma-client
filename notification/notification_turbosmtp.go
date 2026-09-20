package notification

import (
	"fmt"
)

// TurboSMTP represents a turboSMTP notification.
type TurboSMTP struct {
	Base
	TurboSMTPDetails
}

// TurboSMTPDetails contains turboSMTP-specific notification configuration.
type TurboSMTPDetails struct {
	ConsumerKey    string `json:"turbosmtpConsumerKey"`
	ConsumerSecret string `json:"turbosmtpConsumerSecret"`
	// Region selects the turboSMTP API host the alert is sent to. The server
	// uses the EU host for TurboSMTPRegionEU and the US host for anything
	// else, including an unset region.
	Region    TurboSMTPRegion `json:"turbosmtpRegion,omitempty"`
	FromEmail string          `json:"turbosmtpFromEmail"`
	// ToEmail is a comma-separated list of recipient email addresses.
	ToEmail string `json:"turbosmtpToEmail"`
	// CcEmail is a comma-separated list of CC email addresses.
	CcEmail *string `json:"turbosmtpCcEmail,omitempty"`
	// BccEmail is a comma-separated list of BCC email addresses.
	BccEmail *string `json:"turbosmtpBccEmail,omitempty"`
	// Subject is the subject of the email. If unset, the server uses
	// "Notification from Your Uptime Kuma".
	Subject *string `json:"turbosmtpSubject,omitempty"`
}

// TurboSMTPRegion represents the turboSMTP API region.
type TurboSMTPRegion string

// TurboSMTP regions.
const (
	TurboSMTPRegionUS TurboSMTPRegion = "us"
	TurboSMTPRegionEU TurboSMTPRegion = "eu"
)

// String returns the string representation of the turboSMTP region.
func (r TurboSMTPRegion) String() string {
	return string(r)
}

// Type returns the notification type.
func (t TurboSMTP) Type() string {
	return t.TurboSMTPDetails.Type()
}

// Type returns the notification type.
func (TurboSMTPDetails) Type() string {
	return "TurboSMTP"
}

// String returns a string representation of the notification.
func (t TurboSMTP) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TurboSMTPDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (t *TurboSMTP) UnmarshalJSON(data []byte) error {
	detail := TurboSMTPDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = TurboSMTP{
		Base:             base,
		TurboSMTPDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (t TurboSMTP) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TurboSMTPDetails)
}
