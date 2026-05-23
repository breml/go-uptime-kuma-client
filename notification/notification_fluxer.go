package notification

import (
	"fmt"
)

// Fluxer represents a Fluxer notification.
type Fluxer struct {
	Base
	FluxerDetails
}

// FluxerDetails contains fluxer-specific notification configuration.
type FluxerDetails struct {
	WebhookURL         string  `json:"fluxerWebhookUrl"`
	Username           string  `json:"fluxerUsername"`
	PrefixMessage      string  `json:"fluxerPrefixMessage"`
	DisableURL         bool    `json:"disableUrl"`
	UseMessageTemplate *bool   `json:"fluxerUseMessageTemplate,omitempty"`
	MessageFormat      *string `json:"fluxerMessageFormat,omitempty"`
	MessageTemplate    *string `json:"fluxerMessageTemplate,omitempty"`
}

// Type returns the notification type.
func (f Fluxer) Type() string {
	return f.FluxerDetails.Type()
}

// Type returns the notification type.
func (FluxerDetails) Type() string {
	return "fluxer"
}

// String returns a string representation of the notification.
func (f Fluxer) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FluxerDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (f *Fluxer) UnmarshalJSON(data []byte) error {
	detail := FluxerDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = Fluxer{
		Base:          base,
		FluxerDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (f Fluxer) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, &f.FluxerDetails)
}
