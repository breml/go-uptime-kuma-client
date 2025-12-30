package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TailscalePing represents a Tailscale Ping monitor for checking Tailscale node availability.
type TailscalePing struct {
	Base
	TailscalePingDetails
}

// Type returns the monitor type.
func (t TailscalePing) Type() string {
	return t.TailscalePingDetails.Type()
}

// String returns a string representation of the Tailscale Ping monitor.
func (t TailscalePing) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(t.Base, false), formatMonitor(t.TailscalePingDetails, true))
}

// UnmarshalJSON unmarshals a Tailscale Ping monitor from JSON data.
func (t *TailscalePing) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := TailscalePingDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*t = TailscalePing{
		Base:                 base,
		TailscalePingDetails: details,
	}

	return nil
}

// MarshalJSON marshals a Tailscale Ping monitor to JSON data.
func (t TailscalePing) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = t.ID
	raw["type"] = "tailscale-ping"
	raw["name"] = t.Name
	raw["description"] = t.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = t.PathName
	raw["parent"] = t.Parent
	raw["interval"] = t.Interval
	raw["retryInterval"] = t.RetryInterval
	raw["resendInterval"] = t.ResendInterval
	raw["maxretries"] = t.MaxRetries
	raw["upsideDown"] = t.UpsideDown
	raw["active"] = t.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range t.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current Tailscale Ping-specific field values.
	raw["hostname"] = t.Hostname

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

// TailscalePingDetails contains Tailscale Ping-specific monitor configuration.
type TailscalePingDetails struct {
	// Hostname is the Tailscale hostname or IP address to ping.
	Hostname string `json:"hostname"`
}

// Type returns the monitor type.
func (TailscalePingDetails) Type() string {
	return "tailscale-ping"
}
