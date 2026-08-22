package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// NTP represents a Network Time Protocol monitor.
type NTP struct {
	Base
	NTPDetails
}

// Type returns the monitor type.
func (n NTP) Type() string {
	return n.NTPDetails.Type()
}

// String returns a string representation of the monitor.
func (n NTP) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(n.Base, false), formatMonitor(n.NTPDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (n *NTP) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := NTPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*n = NTP{
		Base:       base,
		NTPDetails: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (n NTP) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = n.ID
	raw["type"] = "ntp"
	raw["name"] = n.Name
	raw["description"] = n.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = n.PathName
	raw["parent"] = n.Parent
	raw["interval"] = n.Interval
	raw["retryInterval"] = n.RetryInterval
	raw["resendInterval"] = n.ResendInterval
	raw["maxretries"] = n.MaxRetries
	raw["upsideDown"] = n.UpsideDown
	raw["active"] = n.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range n.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current NTP-specific field values.
	raw["hostname"] = n.Hostname
	raw["port"] = n.Port
	raw["timeout"] = n.Timeout
	raw["ntpStratumThreshold"] = n.NTPStratumThreshold
	raw["ntpTimeOffsetThreshold"] = n.NTPTimeOffsetThreshold
	raw["ntpRootDispersionThreshold"] = n.NTPRootDispersionThreshold

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

// NTPDetails contains NTP-specific monitor configuration.
//
// The Uptime Kuma web UI defaults Base.Interval to 300 seconds for NTP
// monitors, because public NTP servers rate-limit frequent queries. This is a
// UI convention only, the server does not enforce it.
type NTPDetails struct {
	// Hostname is the NTP server to query. It is required.
	Hostname string `json:"hostname"`
	// Port is the UDP port of the NTP server. Leave nil/unset to fall back
	// to the server default of 123.
	Port *int64 `json:"port"`
	// Timeout is the query timeout in seconds. Leave nil/unset to fall back
	// to the server default of 10.
	Timeout *int64 `json:"timeout"`
	// NTPStratumThreshold is the maximum acceptable stratum. The monitor
	// goes down when the reported stratum is greater than or equal to the
	// threshold (stratum 16 is always down). The UI restricts the value to
	// the range 1 to 15. Leave nil/unset to fall back to the server default
	// of 5.
	NTPStratumThreshold *int64 `json:"ntpStratumThreshold"`
	// NTPTimeOffsetThreshold is the maximum acceptable absolute time offset
	// in milliseconds. The monitor goes down when the absolute offset is
	// greater than or equal to the threshold. Leave nil/unset to fall back
	// to the server default of 1000.
	NTPTimeOffsetThreshold *int64 `json:"ntpTimeOffsetThreshold"`
	// NTPRootDispersionThreshold is the maximum acceptable root dispersion
	// in milliseconds. The monitor goes down when the root dispersion is
	// greater than or equal to the threshold. Leave nil/unset to fall back
	// to the server default of 500.
	NTPRootDispersionThreshold *int64 `json:"ntpRootDispersionThreshold"`
}

// Type returns the monitor type.
func (NTPDetails) Type() string {
	return "ntp"
}
