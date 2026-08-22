package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PM2 represents a pm2 monitor.
type PM2 struct {
	Base
	PM2Details
}

// Type returns the monitor type.
func (p PM2) Type() string {
	return p.PM2Details.Type()
}

// String returns a string representation of the monitor.
func (p PM2) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(p.Base, false), formatMonitor(p.PM2Details, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (p *PM2) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal pm2 monitor base: %w", err)
	}

	details := PM2Details{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal pm2 monitor details: %w", err)
	}

	*p = PM2{
		Base:       base,
		PM2Details: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (p PM2) MarshalJSON() ([]byte, error) {
	if strings.TrimSpace(p.ProcessName) == "" {
		return nil, errors.New("marshal: process name is required")
	}

	raw := map[string]any{}
	raw["id"] = p.ID
	raw["type"] = "pm2"
	raw["name"] = p.Name
	raw["description"] = p.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = p.PathName
	raw["parent"] = p.Parent
	raw["interval"] = p.Interval
	raw["retryInterval"] = p.RetryInterval
	raw["resendInterval"] = p.ResendInterval
	raw["maxretries"] = p.MaxRetries
	raw["upsideDown"] = p.UpsideDown
	raw["active"] = p.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range p.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current PM2-specific field values.
	raw["system_service_name"] = p.ProcessName

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal pm2 monitor: %w", err)
	}

	return data, nil
}

// PM2Details contains pm2-specific monitor configuration.
//
// The behaviour described below was verified against Uptime Kuma 2.5.0.
//
// The check runs `pm2 jlist` on the Uptime Kuma host, so the pm2 CLI must be
// installed there and must see the target PM2 daemon. The official container
// image does not ship it, but, unlike the system-service monitor, upstream
// does not otherwise restrict the type in containers.
type PM2Details struct {
	// ProcessName identifies the PM2 process to check. The server matches it
	// against both the process name and the stringified numeric PM2 id
	// reported by `pm2 jlist`, so either form works. Prefer the name, because
	// PM2 reassigns ids after a process is deleted and recreated.
	//
	// The server trims the value and requires it to be non-empty and free of
	// ASCII control characters (U+0000..U+001F, U+007F). MarshalJSON rejects
	// an empty (or whitespace-only) value, the remaining rules are enforced
	// by the server only.
	//
	// Note: the wire field is shared with the system-service monitor and the
	// upstream API uses snake_case for it, unlike most other Uptime Kuma
	// fields.
	ProcessName string `json:"system_service_name"`
}

// Type returns the monitor type.
func (PM2Details) Type() string {
	return "pm2"
}
