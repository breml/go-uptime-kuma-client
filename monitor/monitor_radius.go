package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Radius represents a Radius monitor for testing Radius server authentication.
type Radius struct {
	Base
	RadiusDetails
}

// Type returns the monitor type.
func (r Radius) Type() string {
	return r.RadiusDetails.Type()
}

// String returns a string representation of the Radius monitor.
func (r Radius) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(r.Base, false), formatMonitor(r.RadiusDetails, true))
}

// UnmarshalJSON unmarshals a Radius monitor from JSON data.
func (r *Radius) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := RadiusDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*r = Radius{
		Base:          base,
		RadiusDetails: details,
	}

	return nil
}

// MarshalJSON marshals a Radius monitor to JSON data.
func (r Radius) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = r.ID
	raw["type"] = "radius"
	raw["name"] = r.Name
	raw["description"] = r.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = r.PathName
	raw["parent"] = r.Parent
	raw["interval"] = r.Interval
	raw["retryInterval"] = r.RetryInterval
	raw["resendInterval"] = r.ResendInterval
	raw["maxretries"] = r.MaxRetries
	raw["upsideDown"] = r.UpsideDown
	raw["active"] = r.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range r.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current Radius-specific field values.
	raw["hostname"] = r.Hostname
	raw["port"] = r.Port
	raw["radiusUsername"] = r.Username
	raw["radiusPassword"] = r.Password
	raw["radiusSecret"] = r.Secret
	raw["radiusCalledStationId"] = r.CalledStationID
	raw["radiusCallingStationId"] = r.CallingStationID

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

// RadiusDetails contains Radius-specific monitor configuration.
type RadiusDetails struct {
	// Hostname is the Radius server address.
	Hostname string `json:"hostname"`
	// Port is the Radius server port (default: 1812).
	Port *int64 `json:"port"`
	// Username is the username for Radius authentication.
	Username string `json:"radiusUsername"`
	// Password is the password for Radius authentication.
	Password string `json:"radiusPassword"`
	// Secret is the shared secret for Radius server.
	Secret string `json:"radiusSecret"`
	// CalledStationID is the optional Called-Station-ID attribute.
	CalledStationID *string `json:"radiusCalledStationId"`
	// CallingStationID is the optional Calling-Station-ID attribute.
	CallingStationID *string `json:"radiusCallingStationId"`
}

// Type returns the monitor type.
func (r RadiusDetails) Type() string {
	return "radius"
}
