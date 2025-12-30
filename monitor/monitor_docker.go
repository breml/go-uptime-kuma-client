package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Docker represents a Docker container monitor.
type Docker struct {
	Base
	DockerDetails
}

// Type returns the monitor type.
func (d Docker) Type() string {
	return d.DockerDetails.Type()
}

// String returns a string representation of the Docker monitor.
func (d Docker) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(d.Base, false), formatMonitor(d.DockerDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a Docker monitor.
func (d *Docker) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := DockerDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*d = Docker{
		Base:          base,
		DockerDetails: details,
	}

	return nil
}

// MarshalJSON marshals a Docker monitor into a JSON byte slice.
func (d Docker) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = d.ID
	raw["type"] = "docker"
	raw["name"] = d.Name
	raw["description"] = d.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = d.PathName
	raw["parent"] = d.Parent
	raw["interval"] = d.Interval
	raw["retryInterval"] = d.RetryInterval
	raw["resendInterval"] = d.ResendInterval
	raw["maxretries"] = d.MaxRetries
	raw["upsideDown"] = d.UpsideDown
	raw["active"] = d.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range d.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current Docker-specific field values.
	raw["docker_host"] = d.DockerHost
	raw["docker_container"] = d.DockerContainer

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

// DockerDetails contains Docker-specific monitor configuration.
type DockerDetails struct {
	DockerHost      int64  `json:"docker_host"`
	DockerContainer string `json:"docker_container"`
}

// Type returns the monitor type.
func (DockerDetails) Type() string {
	return "docker"
}
