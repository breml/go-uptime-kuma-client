package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MQTT represents an MQTT monitor that checks the status of MQTT brokers and topics.
type MQTT struct {
	Base
	MQTTDetails
}

// Type returns the monitor type string.
func (m MQTT) Type() string {
	return m.MQTTDetails.Type()
}

// String returns a string representation of the MQTT monitor.
func (m MQTT) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.MQTTDetails, true))
}

// UnmarshalJSON unmarshals an MQTT monitor from JSON data.
func (m *MQTT) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := MQTTDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*m = MQTT{
		Base:        base,
		MQTTDetails: details,
	}

	return nil
}

// MarshalJSON marshals an MQTT monitor to JSON.
func (m MQTT) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = m.ID
	raw["type"] = "mqtt"
	raw["name"] = m.Name
	raw["description"] = m.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = m.PathName
	raw["parent"] = m.Parent
	raw["interval"] = m.Interval
	raw["retryInterval"] = m.RetryInterval
	raw["resendInterval"] = m.ResendInterval
	raw["maxretries"] = m.MaxRetries
	raw["upsideDown"] = m.UpsideDown
	raw["active"] = m.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range m.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current MQTT-specific field values.
	raw["hostname"] = m.Hostname
	raw["port"] = m.Port
	raw["mqttTopic"] = m.MQTTTopic
	raw["mqttUsername"] = m.MQTTUsername
	raw["mqttPassword"] = m.MQTTPassword
	raw["mqttWebsocketPath"] = m.MQTTWebsocketPath
	raw["mqttCheckType"] = m.MQTTCheckType
	raw["mqttSuccessMessage"] = m.MQTTSuccessMessage
	raw["jsonPath"] = m.JSONPath
	raw["expectedValue"] = m.ExpectedValue

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return data, nil
}

// MQTTDetails contains MQTT monitor specific fields.
type MQTTDetails struct {
	// Hostname is the MQTT broker address (supports mqtt://, ws://, wss:// protocols).
	Hostname string `json:"hostname"`
	// Port is the MQTT broker port.
	Port *int64 `json:"port"`
	// MQTTTopic is the topic to subscribe to.
	MQTTTopic string `json:"mqttTopic"`
	// MQTTUsername is the optional username for MQTT authentication.
	MQTTUsername *string `json:"mqttUsername"`
	// MQTTPassword is the optional password for MQTT authentication.
	MQTTPassword *string `json:"mqttPassword"`
	// MQTTWebsocketPath is the optional WebSocket path for WebSocket connections.
	MQTTWebsocketPath *string `json:"mqttWebsocketPath"`
	// MQTTCheckType is the check type: keyword or json-query.
	MQTTCheckType MQTTCheckType `json:"mqttCheckType"`
	// MQTTSuccessMessage is the expected message for keyword check.
	MQTTSuccessMessage *string `json:"mqttSuccessMessage"`
	// JSONPath is the JSON path for json-query check.
	JSONPath *string `json:"jsonPath"`
	// ExpectedValue is the expected value for json-query check.
	ExpectedValue *string `json:"expectedValue"`
}

// Type returns the monitor type string.
func (MQTTDetails) Type() string {
	return "mqtt"
}

// MQTTCheckType represents the MQTT check type.
type MQTTCheckType string

const (
	// MQTTCheckTypeKeyword checks for a keyword in the MQTT message.
	MQTTCheckTypeKeyword MQTTCheckType = "keyword"
	// MQTTCheckTypeJSONQuery checks the MQTT message using JSON query.
	MQTTCheckTypeJSONQuery MQTTCheckType = "json-query"
)
