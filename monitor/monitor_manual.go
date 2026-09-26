package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Manual represents a manual monitor, a passive monitor whose status is set
// by the user instead of being determined by a check.
type Manual struct {
	Base
	ManualDetails
}

// Type returns the monitor type.
func (m Manual) Type() string {
	return m.ManualDetails.Type()
}

// String returns a string representation of the monitor.
func (m Manual) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(m.Base, false), formatMonitor(m.ManualDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (m *Manual) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal manual monitor base: %w", err)
	}

	details := ManualDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal manual monitor details: %w", err)
	}

	*m = Manual{
		Base:          base,
		ManualDetails: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
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

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal manual monitor: %w", err)
	}

	return data, nil
}

// ManualDetails contains manual-specific monitor configuration.
type ManualDetails struct {
	// ManualStatus is the status every heartbeat reports: 0 (down), 1 (up)
	// or 2 (pending). While nil the monitor reports pending with the message
	// "Manual monitoring - No status set".
	//
	// The field is write-only: the server stores it but never includes it
	// in the monitor list, so it always reads back as nil. An edit that does
	// not set it therefore resets the status to pending.
	// Note: the upstream API uses snake_case for this field, unlike most other Uptime Kuma fields.
	ManualStatus *int64 `json:"manual_status"`
}

// Type returns the monitor type.
func (ManualDetails) Type() string {
	return "manual"
}
