package notification

import (
	"fmt"
)

// Flowtriq represents a Flowtriq notification provider.
// Flowtriq is a DDoS detection and incident platform, which receives the
// notifications via a webhook.
type Flowtriq struct {
	Base
	FlowtriqDetails
}

// FlowtriqDetails contains the configuration fields for Flowtriq notifications.
type FlowtriqDetails struct {
	// WebhookURL is the Flowtriq webhook endpoint the notification is posted
	// to, for example https://app.flowtriq.com/api/webhooks/...
	WebhookURL string `json:"flowtriqWebhookUrl"`
	// APIKey authenticates the request against Flowtriq. If set, the server
	// sends it as the "X-API-Key" header, otherwise the header is omitted.
	APIKey *string `json:"flowtriqApiKey,omitempty"`
}

// Type returns the notification type identifier for Flowtriq.
func (f Flowtriq) Type() string {
	return f.FlowtriqDetails.Type()
}

// Type returns the notification type identifier for FlowtriqDetails.
func (FlowtriqDetails) Type() string {
	return "Flowtriq"
}

// String returns a string representation of the Flowtriq notification.
func (f Flowtriq) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FlowtriqDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Flowtriq notification.
func (f *Flowtriq) UnmarshalJSON(data []byte) error {
	detail := FlowtriqDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = Flowtriq{
		Base:            base,
		FlowtriqDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Flowtriq notification into JSON.
func (f Flowtriq) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, &f.FlowtriqDetails)
}
