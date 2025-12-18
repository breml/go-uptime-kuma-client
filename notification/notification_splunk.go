package notification

import (
	"fmt"
)

// Splunk represents a Splunk On-Call (formerly VictorOps) notification provider.
// Splunk On-Call is an incident management and alerting platform for sending notifications.
type Splunk struct {
	Base
	SplunkDetails
}

// SplunkDetails contains the configuration fields for Splunk On-Call notifications.
type SplunkDetails struct {
	// RestURL is the Splunk On-Call REST API URL endpoint.
	RestURL string `json:"splunkRestURL"`
	// Severity is the alert severity level for triggered events.
	Severity string `json:"splunkSeverity"`
	// AutoResolve is the action to take when a monitor recovers.
	AutoResolve string `json:"splunkAutoResolve"`
	// IntegrationKey is the Splunk On-Call routing key for alert routing.
	IntegrationKey string `json:"pagerdutyIntegrationKey"`
}

// Type returns the notification type identifier for Splunk.
func (s Splunk) Type() string {
	return s.SplunkDetails.Type()
}

// Type returns the notification type identifier for SplunkDetails.
func (n SplunkDetails) Type() string {
	return "Splunk"
}

// String returns a string representation of the Splunk notification.
func (s Splunk) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SplunkDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a Splunk notification.
func (s *Splunk) UnmarshalJSON(data []byte) error {
	detail := SplunkDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = Splunk{
		Base:          base,
		SplunkDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the Splunk notification into JSON.
func (s Splunk) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.SplunkDetails)
}
