package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Push ...
type Push struct {
	Base
	PushDetails
}

// Type ...
func (p Push) Type() string {
	return p.PushDetails.Type()
}

// String ...
func (p Push) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(p.Base, false), formatMonitor(p.PushDetails, true))
}

// UnmarshalJSON ...
func (p *Push) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := PushDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*p = Push{
		Base:        base,
		PushDetails: details,
	}

	return nil
}

// MarshalJSON ...
func (p Push) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = p.ID
	raw["type"] = "push"
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

	// Always override with current Push-specific field values.
	raw["pushToken"] = p.PushToken

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

// PushDetails ...
type PushDetails struct {
	PushToken string `json:"pushToken"`
}

// Type ...
func (PushDetails) Type() string {
	return "push"
}
