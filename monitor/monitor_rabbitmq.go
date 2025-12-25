package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// RabbitMQ represents a RabbitMQ monitor for checking RabbitMQ cluster health.
type RabbitMQ struct {
	Base
	RabbitMQDetails
}

// Type returns the monitor type.
func (r RabbitMQ) Type() string {
	return r.RabbitMQDetails.Type()
}

// String returns a string representation of the RabbitMQ monitor.
func (r RabbitMQ) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(r.Base, false), formatMonitor(r.RabbitMQDetails, true))
}

// UnmarshalJSON unmarshals a RabbitMQ monitor from JSON data.
func (r *RabbitMQ) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := RabbitMQDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*r = RabbitMQ{
		Base:            base,
		RabbitMQDetails: details,
	}

	return nil
}

// MarshalJSON marshals a RabbitMQ monitor to JSON data.
func (r RabbitMQ) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = r.ID
	raw["type"] = "rabbitmq"
	raw["name"] = r.Name
	raw["description"] = r.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = r.PathName
	raw["parent"] = r.Parent
	raw["interval"] = r.Interval
	raw["retryInterval"] = r.RetryInterval
	raw["resendInterval"] = r.ResendInterval
	raw["maxretries"] = r.MaxRetries
	raw["upsideDown"] = r.UpsideDown
	raw["active"] = r.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range r.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}
	raw["notificationIDList"] = ids

	// Always override with current RabbitMQ-specific field values.
	raw["rabbitmqNodes"] = r.Nodes
	raw["rabbitmqUsername"] = r.Username
	raw["rabbitmqPassword"] = r.Password
	raw["timeout"] = r.Timeout

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// RabbitMQDetails contains RabbitMQ-specific monitor configuration.
type RabbitMQDetails struct {
	// Nodes is a JSON-encoded array of RabbitMQ node URLs to check.
	Nodes string `json:"rabbitmqNodes"`
	// Username is an optional username for HTTP Basic authentication.
	Username *string `json:"rabbitmqUsername"`
	// Password is an optional password for HTTP Basic authentication.
	Password *string `json:"rabbitmqPassword"`
	// Timeout is an optional request timeout in seconds.
	Timeout *int64 `json:"timeout"`
}

// Type returns the monitor type.
func (r RabbitMQDetails) Type() string {
	return "rabbitmq"
}
