package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// defaultNTPTimeout is the query timeout in seconds that the NTP check falls
// back to. The client sends it explicitly, because the monitor.timeout column
// is NOT NULL and the server rejects an explicit null.
const defaultNTPTimeout = 10

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
	if n.Hostname == "" {
		return nil, errors.New("marshal: hostname is required")
	}

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
	raw["ntpStratumThreshold"] = n.NTPStratumThreshold
	raw["ntpTimeOffsetThreshold"] = n.NTPTimeOffsetThreshold
	raw["ntpRootDispersionThreshold"] = n.NTPRootDispersionThreshold

	// The monitor.timeout column is NOT NULL, so an unset Timeout is sent as
	// the value the check would fall back to instead of as null.
	timeout := int64(defaultNTPTimeout)
	if n.Timeout != nil {
		timeout = *n.Timeout
	}

	raw["timeout"] = timeout

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
// The behaviour described below was verified against Uptime Kuma 2.5.0.
//
// The Uptime Kuma web UI defaults Base.Interval to 300 seconds for NTP
// monitors, because public NTP servers rate-limit frequent queries. This is a
// UI convention only, the server does not enforce it.
//
// The optional fields below are pointers, but the server has no notion of an
// unset field: a nil pointer is sent as JSON null and stored as SQL NULL, so
// it reads back as nil rather than as the fallback value. The fallbacks are
// applied per check, never at rest, so they do not round-trip. Because the
// check tests these values for JavaScript truthiness, a pointer to 0 is
// indistinguishable from nil, which makes nil the only unambiguous way to
// request the fallback. Timeout is the exception, see its comment.
type NTPDetails struct {
	// Hostname is the NTP server to query. It is required: MarshalJSON
	// rejects an empty value, because the server accepts it and then fails
	// on every heartbeat instead.
	Hostname string `json:"hostname"`
	// Port is the UDP port of the NTP server. The column has no default;
	// while it is NULL the check falls back to 123.
	Port *int64 `json:"port"`
	// Timeout is the query timeout in seconds. The column is NOT NULL, so
	// unlike the other optional fields a nil value is not sent as null:
	// MarshalJSON substitutes 10, the value the check itself falls back to.
	Timeout *int64 `json:"timeout"`
	// NTPStratumThreshold is the stratum at which the monitor is considered
	// down. The check fails when the reported stratum is greater than or
	// equal to this value, so a threshold of 5 already rejects stratum 5.
	// Stratum 16 is down regardless of the threshold. The UI restricts the
	// value to the range 1 to 15. While NULL the check falls back to 5.
	NTPStratumThreshold *int64 `json:"ntpStratumThreshold"`
	// NTPTimeOffsetThreshold is the absolute time offset in milliseconds at
	// which the monitor is considered down. The check fails when the
	// absolute offset is greater than or equal to this value. While NULL the
	// check falls back to 1000.
	NTPTimeOffsetThreshold *int64 `json:"ntpTimeOffsetThreshold"`
	// NTPRootDispersionThreshold is the root dispersion in milliseconds at
	// which the monitor is considered down. The check fails when the root
	// dispersion is greater than or equal to this value. While NULL the
	// check falls back to 500.
	NTPRootDispersionThreshold *int64 `json:"ntpRootDispersionThreshold"`
}

// Type returns the monitor type.
func (NTPDetails) Type() string {
	return "ntp"
}
