package notification

import (
	"fmt"
)

// TurboSMTP represents a TurboSMTP notification provider.
// TurboSMTP (https://serversmtp.com/turbo-api/) is a transactional email
// relay, the alert is sent as an email through its HTTP API.
type TurboSMTP struct {
	Base
	TurboSMTPDetails
}

// TurboSMTPDetails contains the configuration fields for TurboSMTP
// notifications.
type TurboSMTPDetails struct {
	// ConsumerKey identifies the TurboSMTP API account. It is sent verbatim as
	// the "consumerKey" request header.
	ConsumerKey string `json:"turbosmtpConsumerKey"`
	// ConsumerSecret authenticates the request. It is sent verbatim as the
	// "consumerSecret" request header.
	ConsumerSecret string `json:"turbosmtpConsumerSecret"`
	// Region selects the TurboSMTP API host the alert is sent to. The server
	// uses the EU host for TurboSMTPRegionEU and the US host for anything
	// else, including an unset region.
	//
	// Unlike the other fields Uptime Kuma marks as required, the region carries
	// omitempty. The server reads an absent and an empty region exactly like
	// TurboSMTPRegionUS, so omitting the key keeps the config free of a value
	// the form never stores. An empty region read from the server is therefore
	// written back as an absent key.
	Region TurboSMTPRegion `json:"turbosmtpRegion,omitempty"`
	// FromEmail is the sender address. The server trims the whitespace around
	// it before the mail is sent.
	FromEmail string `json:"turbosmtpFromEmail"`
	// ToEmail is a comma-separated list of recipient email addresses. Unlike
	// CcEmail and BccEmail, it is passed on verbatim, so the whitespace around
	// an individual address is not trimmed.
	ToEmail string `json:"turbosmtpToEmail"`
	// CcEmail is a comma-separated list of CC email addresses. The server trims
	// the whitespace around every address. If unset or empty, the mail carries
	// no CC recipients.
	CcEmail *string `json:"turbosmtpCcEmail,omitempty"`
	// BccEmail is a comma-separated list of BCC email addresses. The server
	// trims the whitespace around every address. If unset or empty, the mail
	// carries no BCC recipients.
	BccEmail *string `json:"turbosmtpBccEmail,omitempty"`
	// Subject is the subject of the email. If unset or empty, the server uses
	// "Notification from Your Uptime Kuma".
	Subject *string `json:"turbosmtpSubject,omitempty"`
}

// TurboSMTPRegion represents the TurboSMTP API region.
type TurboSMTPRegion string

// TurboSMTP regions.
const (
	TurboSMTPRegionUS TurboSMTPRegion = "us"
	TurboSMTPRegionEU TurboSMTPRegion = "eu"
)

// String returns the string representation of the TurboSMTP region.
func (r TurboSMTPRegion) String() string {
	return string(r)
}

// Type returns the notification type identifier for TurboSMTP.
func (t TurboSMTP) Type() string {
	return t.TurboSMTPDetails.Type()
}

// Type returns the notification type identifier for TurboSMTPDetails.
func (TurboSMTPDetails) Type() string {
	return "TurboSMTP"
}

// String returns a string representation of the TurboSMTP notification.
func (t TurboSMTP) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TurboSMTPDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a TurboSMTP notification.
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

// MarshalJSON marshals the TurboSMTP notification into JSON.
func (t TurboSMTP) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TurboSMTPDetails)
}
