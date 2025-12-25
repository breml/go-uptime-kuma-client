package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Manual represents a Manual monitor that allows manual control of the monitor status.
type Manual struct {
	Base
	ManualDetails
}

// Type returns the monitor type.
func (m Manual) Type() string {
	return m.ManualDetails.Type()
}

// String returns a string representation of the Manual monitor.
func (m Manual) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.ManualDetails, true))
}

// UnmarshalJSON unmarshals a Manual monitor from JSON data.
func (m *Manual) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := ManualDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*m = Manual{
		Base:          base,
		ManualDetails: details,
	}

	return nil
}

// MarshalJSON marshals a Manual monitor to JSON data.
func (m Manual) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = m.ID
	raw["type"] = "manual"
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

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// ManualDetails contains Manual monitor-specific configuration.
// Manual monitors don't have additional fields beyond the Base monitor configuration.
type ManualDetails struct {
	// Manual monitors have no additional fields.
}

// Type returns the monitor type.
func (m ManualDetails) Type() string {
	return "manual"
}
