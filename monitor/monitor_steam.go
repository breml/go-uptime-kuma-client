package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Steam represents a Steam Game Server monitor that checks the status of a game server via Steam API.
type Steam struct {
	Base
	SteamDetails
}

// Type returns the monitor type string.
func (s Steam) Type() string {
	return s.SteamDetails.Type()
}

// String returns a string representation of the Steam monitor.
func (s Steam) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SteamDetails, true))
}

// UnmarshalJSON unmarshals a Steam monitor from JSON data.
func (s *Steam) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := SteamDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*s = Steam{
		Base:         base,
		SteamDetails: details,
	}

	return nil
}

// MarshalJSON marshals a Steam monitor to JSON.
func (s Steam) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "steam"
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

	// Always override with current Steam-specific field values.
	raw["hostname"] = s.Hostname
	raw["port"] = s.Port
	raw["timeout"] = s.Timeout

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// SteamDetails contains Steam monitor specific fields.
type SteamDetails struct {
	// Hostname is the Steam game server IP address.
	Hostname string `json:"hostname"`
	// Port is the Steam game server port.
	Port int `json:"port"`
	// Timeout is the request timeout in seconds.
	Timeout *int64 `json:"timeout"`
}

// Type returns the monitor type string.
func (s SteamDetails) Type() string {
	return "steam"
}
