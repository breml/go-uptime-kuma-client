package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MySQL represents a MySQL/MariaDB monitor for testing database connectivity and queries.
type MySQL struct {
	Base
	MySQLDetails
}

// Type returns the monitor type.
func (m MySQL) Type() string {
	return m.MySQLDetails.Type()
}

// String returns a string representation of the MySQL monitor.
func (m MySQL) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.MySQLDetails, true))
}

// UnmarshalJSON unmarshals a MySQL monitor from JSON data.
func (m *MySQL) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := MySQLDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*m = MySQL{
		Base:         base,
		MySQLDetails: details,
	}

	return nil
}

// MarshalJSON marshals a MySQL monitor to JSON data.
func (m MySQL) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = m.ID
	raw["type"] = "mysql"
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

	// Always override with current MySQL-specific field values.
	raw["databaseConnectionString"] = m.DatabaseConnectionString
	raw["databaseQuery"] = m.DatabaseQuery

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

// MySQLDetails contains MySQL-specific monitor configuration.
type MySQLDetails struct {
	// DatabaseConnectionString is the MySQL connection string.
	DatabaseConnectionString string `json:"databaseConnectionString"`
	// DatabaseQuery is an optional SQL query to execute (default: SELECT 1).
	DatabaseQuery *string `json:"databaseQuery"`
}

// Type returns the monitor type.
func (MySQLDetails) Type() string {
	return "mysql"
}
