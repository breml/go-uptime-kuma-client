package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type TCPPort struct {
	Base
	TCPPortDetails
}

func (t TCPPort) Type() string {
	return t.TCPPortDetails.Type()
}

func (t TCPPort) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(t.Base, false), formatMonitor(t.TCPPortDetails, true))
}

func (t *TCPPort) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := TCPPortDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*t = TCPPort{
		Base:           base,
		TCPPortDetails: details,
	}

	return nil
}

func (t TCPPort) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = t.ID
	raw["type"] = "port"
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

	// Always override with current TCP Port-specific field values.
	raw["hostname"] = t.Hostname
	raw["port"] = t.Port

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

type TCPPortDetails struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func (t TCPPortDetails) Type() string {
	return "port"
}
