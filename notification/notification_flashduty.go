package notification

import (
	"fmt"
)

// FlashDuty represents a FlashDuty notification provider.
// FlashDuty is an incident management and alerting platform for sending notifications.
type FlashDuty struct {
	Base
	FlashDutyDetails
}

// FlashDutyDetails contains the configuration fields for FlashDuty notifications.
type FlashDutyDetails struct {
	// IntegrationKey is the FlashDuty integration key or full webhook URL.
	IntegrationKey string `json:"flashdutyIntegrationKey"`
	// Severity is the alert severity level (Info, Warning, Critical, Ok).
	Severity string `json:"flashdutySeverity"`
}

// Type returns the notification type identifier for FlashDuty.
func (f FlashDuty) Type() string {
	return f.FlashDutyDetails.Type()
}

// Type returns the notification type identifier for FlashDutyDetails.
func (FlashDutyDetails) Type() string {
	return "FlashDuty"
}

// String returns a string representation of the FlashDuty notification.
func (f FlashDuty) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FlashDutyDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a FlashDuty notification.
func (f *FlashDuty) UnmarshalJSON(data []byte) error {
	detail := FlashDutyDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = FlashDuty{
		Base:             base,
		FlashDutyDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the FlashDuty notification into JSON.
func (f FlashDuty) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, f.FlashDutyDetails)
}
