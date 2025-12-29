package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Ping struct {
	Base
	PingDetails
}

func (p Ping) Type() string {
	return p.PingDetails.Type()
}

func (p Ping) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(p.Base, false), formatMonitor(p.PingDetails, true))
}

func (p *Ping) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := PingDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*p = Ping{
		Base:        base,
		PingDetails: details,
	}

	return nil
}

func (p Ping) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = p.ID
	raw["type"] = "ping"
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

	// Always override with current Ping-specific field values.
	raw["hostname"] = p.Hostname
	raw["packetSize"] = p.PacketSize

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

type PingDetails struct {
	Hostname   string `json:"hostname"`
	PacketSize int    `json:"packetSize"`
}

func (p PingDetails) Type() string {
	return "ping"
}
