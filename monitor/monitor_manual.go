package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Manual represents a manual monitor that allows manual control of the status.
type Manual struct {
	Base
	ManualDetails
}

// Type returns the monitor type string.
func (m Manual) Type() string {
	return m.ManualDetails.Type()
}

// String returns a string representation of the manual monitor.
func (m Manual) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.ManualDetails, true))
}

// UnmarshalJSON unmarshals a manual monitor from JSON data.
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

// MarshalJSON marshals a manual monitor to JSON.
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

	// Always override with current Manual-specific field values.
	raw["manual_status"] = m.ManualStatus

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// ManualDetails contains manual monitor specific fields.
type ManualDetails struct {
	// ManualStatus is the manually set status (0=DOWN, 1=UP, 2=PENDING, null=not set).
	ManualStatus *int `json:"manual_status"`
}

// Type returns the monitor type string.
func (m ManualDetails) Type() string {
	return "manual"
}
