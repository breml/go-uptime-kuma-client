package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// OracleDB represents an OracleDB monitor for testing Oracle Database connectivity and queries.
type OracleDB struct {
	Base
	OracleDBDetails
}

// Type returns the monitor type.
func (o OracleDB) Type() string {
	return o.OracleDBDetails.Type()
}

// String returns a string representation of the OracleDB monitor.
func (o OracleDB) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(o.Base, false), formatMonitor(o.OracleDBDetails, true))
}

// UnmarshalJSON unmarshals an OracleDB monitor from JSON data.
func (o *OracleDB) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := OracleDBDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*o = OracleDB{
		Base:            base,
		OracleDBDetails: details,
	}

	return nil
}

// MarshalJSON marshals an OracleDB monitor to JSON data.
func (o OracleDB) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = o.ID
	raw["type"] = "oracledb"
	raw["name"] = o.Name
	raw["description"] = o.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = o.PathName
	raw["parent"] = o.Parent
	raw["interval"] = o.Interval
	raw["retryInterval"] = o.RetryInterval
	raw["resendInterval"] = o.ResendInterval
	raw["maxretries"] = o.MaxRetries
	raw["upsideDown"] = o.UpsideDown
	raw["active"] = o.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range o.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current OracleDB-specific field values.
	raw["databaseConnectionString"] = o.DatabaseConnectionString
	raw["databaseQuery"] = o.DatabaseQuery
	raw["basic_auth_user"] = o.Username
	raw["basic_auth_pass"] = o.Password

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = conditionsForWire(o.Conditions)

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return data, nil
}

// OracleDBDetails contains OracleDB-specific monitor configuration.
type OracleDBDetails struct {
	// DatabaseConnectionString is the Oracle EZCONNECT string, e.g. host:port/service.
	DatabaseConnectionString string `json:"databaseConnectionString"`
	// DatabaseQuery is an optional SQL query to execute (default: SELECT 1 FROM DUAL).
	DatabaseQuery *string `json:"databaseQuery"`
	// Username is the Oracle Database user.
	Username string `json:"basic_auth_user"`
	// Password is the Oracle Database password.
	Password string `json:"basic_auth_pass"`
	// Conditions is an optional list of assertion clauses evaluated against the
	// query result. When set, the query is expected to return a single value
	// that is matched against the conditions.
	Conditions []Condition `json:"conditions,omitempty"`
}

// Type returns the monitor type.
func (OracleDBDetails) Type() string {
	return "oracledb"
}
