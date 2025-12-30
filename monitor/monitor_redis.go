package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Redis represents a redis monitor.
type Redis struct {
	Base
	RedisDetails
}

// Type returns the monitor type.
func (r Redis) Type() string {
	return r.RedisDetails.Type()
}

// String returns a string representation of the monitor.
func (r Redis) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(r.Base, false), formatMonitor(r.RedisDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (r *Redis) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := RedisDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*r = Redis{
		Base:         base,
		RedisDetails: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (r Redis) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = r.ID
	raw["type"] = "redis"
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

	// Always override with current Redis-specific field values.
	raw["databaseConnectionString"] = r.ConnectionString
	raw["ignoreTls"] = r.IgnoreTLS

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

// RedisDetails contains redis-specific monitor configuration.
type RedisDetails struct {
	ConnectionString string `json:"databaseConnectionString"`
	IgnoreTLS        bool   `json:"ignoreTls"`
}

// Type returns the monitor type.
func (RedisDetails) Type() string {
	return "redis"
}
