package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Group ...
type Group struct {
	Base
	GroupDetails
}

// Type ...
func (g Group) Type() string {
	return g.GroupDetails.Type()
}

// String ...
func (g Group) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(g.Base, false), formatMonitor(g.GroupDetails, true))
}

// UnmarshalJSON ...
func (g *Group) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := GroupDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*g = Group{
		Base:         base,
		GroupDetails: details,
	}

	return nil
}

// MarshalJSON ...
func (g Group) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = g.ID
	raw["type"] = "group"
	raw["name"] = g.Name
	raw["description"] = g.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = g.PathName
	raw["parent"] = g.Parent
	raw["interval"] = g.Interval
	raw["retryInterval"] = g.RetryInterval
	raw["resendInterval"] = g.ResendInterval
	raw["maxretries"] = g.MaxRetries
	raw["upsideDown"] = g.UpsideDown
	raw["active"] = g.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range g.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

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

// GroupDetails ...
type GroupDetails struct {
	// Groups don't have additional fields beyond Base.
}

// Type ...
func (g GroupDetails) Type() string {
	return "group"
}
