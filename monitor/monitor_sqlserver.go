package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SQLServer represents a SQL Server monitor for testing SQL Server connectivity and queries.
type SQLServer struct {
	Base
	SQLServerDetails
}

// Type returns the monitor type.
func (s SQLServer) Type() string {
	return s.SQLServerDetails.Type()
}

// String returns a string representation of the SQL Server monitor.
func (s SQLServer) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SQLServerDetails, true))
}

// UnmarshalJSON unmarshals a SQL Server monitor from JSON data.
func (s *SQLServer) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := SQLServerDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*s = SQLServer{
		Base:             base,
		SQLServerDetails: details,
	}

	return nil
}

// MarshalJSON marshals a SQL Server monitor to JSON data.
func (s SQLServer) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "sqlserver"
	raw["name"] = s.Name
	raw["description"] = s.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = s.PathName
	raw["parent"] = s.Parent
	raw["interval"] = s.Interval
	raw["retryInterval"] = s.RetryInterval
	raw["resendInterval"] = s.ResendInterval
	raw["maxretries"] = s.MaxRetries
	raw["upsideDown"] = s.UpsideDown
	raw["active"] = s.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range s.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current SQL Server-specific field values.
	raw["databaseConnectionString"] = s.DatabaseConnectionString
	raw["databaseQuery"] = s.DatabaseQuery

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// SQLServerDetails contains SQL Server-specific monitor configuration.
type SQLServerDetails struct {
	// DatabaseConnectionString is the SQL Server connection string.
	DatabaseConnectionString string `json:"databaseConnectionString"`
	// DatabaseQuery is an optional SQL query to execute (default: basic connection test).
	DatabaseQuery *string `json:"databaseQuery"`
}

// Type returns the monitor type.
func (s SQLServerDetails) Type() string {
	return "sqlserver"
}
