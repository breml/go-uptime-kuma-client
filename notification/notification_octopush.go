package notification

import (
	"fmt"
)

// Octopush represents an Octopush SMS notification provider.
// Octopush is an SMS gateway service that supports both V1 and V2 API versions.
type Octopush struct {
	Base
	OctopushDetails
}

// OctopushDetails contains the configuration fields for Octopush notifications.
type OctopushDetails struct {
	// Version is the Octopush API version: "1" or "2" (default: "2").
	Version string `json:"octopushVersion"`
	// APIKey is the API key for Octopush V2 API authentication.
	APIKey string `json:"octopushAPIKey"`
	// Login is the login/username for Octopush V2 API authentication.
	Login string `json:"octopushLogin"`
	// PhoneNumber is the recipient phone number for Octopush V2 API.
	PhoneNumber string `json:"octopushPhoneNumber"`
	// SMSType is the SMS type for Octopush V2 API (e.g., "sms_premium", "sms_low_cost").
	SMSType string `json:"octopushSMSType"`
	// SenderName is the sender name for Octopush V2 API.
	SenderName string `json:"octopushSenderName"`
	// DMLogin is the login/username for Octopush V1 (Direct Mail) API.
	DMLogin string `json:"octopushDMLogin"`
	// DMAPIKey is the API key for Octopush V1 (Direct Mail) API.
	DMAPIKey string `json:"octopushDMAPIKey"`
	// DMPhoneNumber is the recipient phone number for Octopush V1 API.
	DMPhoneNumber string `json:"octopushDMPhoneNumber"`
	// DMSMSType is the SMS type for Octopush V1 API.
	DMSMSType string `json:"octopushDMSMSType"`
	// DMSenderName is the sender name for Octopush V1 API.
	DMSenderName string `json:"octopushDMSenderName"`
}

// Type returns the notification type identifier for Octopush.
func (o Octopush) Type() string {
	return o.OctopushDetails.Type()
}

// Type returns the notification type identifier for OctopushDetails.
func (OctopushDetails) Type() string {
	return "octopush"
}

// String returns a string representation of the Octopush notification.
func (o Octopush) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(o.Base, false), formatNotification(o.OctopushDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an Octopush notification.
func (o *Octopush) UnmarshalJSON(data []byte) error {
	detail := OctopushDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*o = Octopush{
		Base:            base,
		OctopushDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Octopush notification into JSON.
func (o Octopush) MarshalJSON() ([]byte, error) {
	return marshalJSON(o.Base, o.OctopushDetails)
}
