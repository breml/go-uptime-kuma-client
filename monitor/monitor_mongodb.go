package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MongoDB represents a MongoDB monitor for testing database connectivity and commands.
type MongoDB struct {
	Base
	MongoDBDetails
}

// Type returns the monitor type.
func (m MongoDB) Type() string {
	return m.MongoDBDetails.Type()
}

// String returns a string representation of the MongoDB monitor.
func (m MongoDB) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.MongoDBDetails, true))
}

// UnmarshalJSON unmarshals a MongoDB monitor from JSON data.
func (m *MongoDB) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := MongoDBDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*m = MongoDB{
		Base:           base,
		MongoDBDetails: details,
	}

	return nil
}

// MarshalJSON marshals a MongoDB monitor to JSON data.
func (m MongoDB) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = m.ID
	raw["type"] = "mongodb"
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

	// Always override with current MongoDB-specific field values.
	raw["databaseConnectionString"] = m.DatabaseConnectionString
	raw["databaseQuery"] = m.DatabaseQuery
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

// MongoDBDetails contains MongoDB-specific monitor configuration.
type MongoDBDetails struct {
	// DatabaseConnectionString is the MongoDB connection string.
	DatabaseConnectionString string `json:"databaseConnectionString"`
	// DatabaseQuery is an optional MongoDB command as JSON (default: {ping: 1}).
	DatabaseQuery *string `json:"databaseQuery"`
	// JSONPath is an optional JSONata expression for result validation.
	JSONPath *string `json:"jsonPath"`
	// ExpectedValue is the expected value when using jsonPath.
	ExpectedValue *string `json:"expectedValue"`
}

// Type returns the monitor type.
func (MongoDBDetails) Type() string {
	return "mongodb"
}
