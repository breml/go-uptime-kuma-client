package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Redis struct {
	Base
	RedisDetails
}

func (r Redis) Type() string {
	return r.RedisDetails.Type()
}

func (r Redis) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(r.Base, false), formatMonitor(r.RedisDetails, true))
}

func (r *Redis) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := RedisDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*r = Redis{
		Base:         base,
		RedisDetails: details,
	}

	return nil
}

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

	return json.Marshal(raw)
}

type RedisDetails struct {
	ConnectionString string `json:"databaseConnectionString"`
	IgnoreTLS        bool   `json:"ignoreTls"`
}

func (r RedisDetails) Type() string {
	return "redis"
}
