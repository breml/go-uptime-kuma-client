package notification

import (
	"fmt"
)

// ZohoCliq represents a Zoho Cliq notification provider.
// Zoho Cliq is a team collaboration platform for sending notifications via webhooks.
type ZohoCliq struct {
	Base
	ZohoCliqDetails
}

// ZohoCliqDetails contains the configuration fields for Zoho Cliq notifications.
type ZohoCliqDetails struct {
	// WebhookURL is the Zoho Cliq webhook URL where notifications will be sent.
	WebhookURL string `json:"webhookUrl"`
}

// Type returns the notification type identifier for ZohoCliq.
func (z ZohoCliq) Type() string {
	return z.ZohoCliqDetails.Type()
}

// Type returns the notification type identifier for ZohoCliqDetails.
func (ZohoCliqDetails) Type() string {
	return "ZohoCliq"
}

// String returns a string representation of the ZohoCliq notification.
func (z ZohoCliq) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(z.Base, false), formatNotification(z.ZohoCliqDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a ZohoCliq notification.
func (z *ZohoCliq) UnmarshalJSON(data []byte) error {
	detail := ZohoCliqDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*z = ZohoCliq{
		Base:            base,
		ZohoCliqDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the ZohoCliq notification into JSON.
func (z ZohoCliq) MarshalJSON() ([]byte, error) {
	return marshalJSON(z.Base, z.ZohoCliqDetails)
}
