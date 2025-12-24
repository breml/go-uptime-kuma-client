package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GameDig represents a GameDig monitor that checks game server status using the GameDig protocol.
type GameDig struct {
	Base
	GameDigDetails
}

// Type returns the monitor type string.
func (g GameDig) Type() string {
	return g.GameDigDetails.Type()
}

// String returns a string representation of the GameDig monitor.
func (g GameDig) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(g.Base, false), formatMonitor(g.GameDigDetails, true))
}

// UnmarshalJSON unmarshals a GameDig monitor from JSON data.
func (g *GameDig) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := GameDigDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*g = GameDig{
		Base:           base,
		GameDigDetails: details,
	}

	return nil
}

// MarshalJSON marshals a GameDig monitor to JSON.
func (g GameDig) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = g.ID
	raw["type"] = "gamedig"
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

	// Always override with current GameDig-specific field values.
	raw["hostname"] = g.Hostname
	raw["port"] = g.Port
	raw["game"] = g.Game
	raw["gamedigGivenPortOnly"] = g.GameDigGivenPortOnly

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

// GameDigDetails contains GameDig monitor specific fields.
type GameDigDetails struct {
	// Hostname is the game server address.
	Hostname string `json:"hostname"`
	// Port is the game server port.
	Port int `json:"port"`
	// Game is the game type identifier (e.g., minecraft, csgo, etc.).
	Game string `json:"game"`
	// GameDigGivenPortOnly determines if only the given port should be used without auto-detection.
	GameDigGivenPortOnly bool `json:"gamedigGivenPortOnly"`
}

// Type returns the monitor type string.
func (g GameDigDetails) Type() string {
	return "gamedig"
}
