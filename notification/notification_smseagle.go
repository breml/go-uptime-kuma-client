package notification

import (
	"fmt"
)

// SMSEagle represents an SMSEagle notification provider.
// SMSEagle is an SMS and voice call service (https://www.smseagle.eu) that supports
// multiple API versions (v1 and v2) and message types including SMS, MMS, ring calls,
// and text-to-speech calls. More info: https://www.smseagle.eu/api/
type SMSEagle struct {
	Base
	SMSEagleDetails
}

// SMSEagleDetails contains the configuration fields for SMSEagle notifications.
// Based on the SMSEagle API documentation (https://www.smseagle.eu/docs/apiv2/ and https://www.smseagle.eu/apiv1/)
type SMSEagleDetails struct {
	// URL is the SMSEagle device URL (e.g., https://192.168.1.100 or https://smseagle.example.com).
	URL string `json:"smseagleUrl"`
	// Token is the API access token for authentication.
	Token string `json:"smseagleToken"`
	// RecipientType specifies the recipient type. Valid values: "smseagle-to", "smseagle-contact", "smseagle-group".
	// "smseagle-to" uses phone number recipients, "smseagle-contact" uses contact IDs, "smseagle-group" uses group IDs.
	RecipientType string `json:"smseagleRecipientType"`
	// Recipient is the recipient identifier used with RecipientType (phone number for "smseagle-to").
	Recipient string `json:"smseagleRecipient"`
	// RecipientTo is the recipient phone number(s) for API v2 (comma-separated for multiple).
	RecipientTo string `json:"smseagleRecipientTo"`
	// RecipientContact is the contact recipient ID(s) for API v2 (comma-separated for multiple).
	RecipientContact string `json:"smseagleRecipientContact"`
	// RecipientGroup is the group recipient ID(s) for API v2 (comma-separated for multiple).
	RecipientGroup string `json:"smseagleRecipientGroup"`
	// MsgType specifies the message type. Valid values: "smseagle-sms", "smseagle-ring", "smseagle-tts", "smseagle-tts-advanced".
	// "smseagle-sms" sends text messages, others send voice calls.
	MsgType string `json:"smseagleMsgType"`
	// Priority is the message priority level (0-9, where 9 is highest). Default is 0.
	// Only used for SMS messages in API v1; for API v2 this indicates priority in the message queue.
	Priority int `json:"smseaglePriority"`
	// Encoding indicates message encoding. false = standard (default), true = unicode for non-ASCII characters.
	Encoding bool `json:"smseagleEncoding"`
	// Duration is the duration for voice calls in seconds (default 10). Applicable for ring, tts, and tts-advanced message types.
	Duration int `json:"smseagleDuration"`
	// TtsModel is the TTS voice ID/model for text-to-speech calls (default 1). Only used with "smseagle-tts-advanced".
	TtsModel int `json:"smseagleTtsModel"`
	// APIType specifies the API version. Valid values: "smseagle-apiv1", "smseagle-apiv2".
	// API v1: https://www.smseagle.eu/apiv1/ (simple HTTP GET-based API)
	// API v2: https://www.smseagle.eu/docs/apiv2/ (modern RESTful API with more features)
	APIType string `json:"smseagleApiType"`
}

// Type returns the notification type identifier for SMSEagle.
func (s SMSEagle) Type() string {
	return s.SMSEagleDetails.Type()
}

// Type returns the notification type identifier for SMSEagleDetails.
func (SMSEagleDetails) Type() string {
	return "SMSEagle"
}

// String returns a string representation of the SMSEagle notification.
func (s SMSEagle) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMSEagleDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an SMSEagle notification.
func (s *SMSEagle) UnmarshalJSON(data []byte) error {
	detail := SMSEagleDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMSEagle{
		Base:            base,
		SMSEagleDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the SMSEagle notification into JSON.
func (s SMSEagle) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SMSEagleDetails)
}
