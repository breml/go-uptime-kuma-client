package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SIPOptions represents a sip-options monitor.
type SIPOptions struct {
	Base
	SIPOptionsDetails
}

// Type returns the monitor type.
func (s SIPOptions) Type() string {
	return s.SIPOptionsDetails.Type()
}

// String returns a string representation of the monitor.
func (s SIPOptions) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SIPOptionsDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (s *SIPOptions) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := SIPOptionsDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*s = SIPOptions{
		Base:              base,
		SIPOptionsDetails: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (s SIPOptions) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "sip-options"
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

	// Always override with current SIPOptions-specific field values.
	raw["hostname"] = s.Hostname
	raw["port"] = s.Port

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

// SIPOptionsDetails contains sip-options-specific monitor configuration.
type SIPOptionsDetails struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// Type returns the monitor type.
func (SIPOptionsDetails) Type() string {
	return "sip-options"
}
