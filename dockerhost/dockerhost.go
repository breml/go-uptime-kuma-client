package dockerhost

import (
	"encoding/json"
)

// DockerHost represents a Docker host configuration in Uptime Kuma.
// Docker hosts can be used to connect to Docker daemons for container monitoring.
type DockerHost struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"userId"`
	DockerDaemon string `json:"dockerDaemon"` // Connection string (e.g., "unix:///var/run/docker.sock", "tcp://host:2375")
	DockerType   string `json:"dockerType"`   // "socket" or "tcp"
	Name         string `json:"name"`         // Human-readable name
}

func (d DockerHost) GetID() int64 {
	return d.ID
}

func (d DockerHost) String() string {
	return formatDockerHost(d)
}

// Config represents the configuration for creating or updating a Docker host.
// It includes an optional ID field for updates.
type Config struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name"`
	DockerDaemon string `json:"dockerDaemon"`
	DockerType   string `json:"dockerType"`
}

// TestResult contains the result of testing a Docker host connection.
type TestResult struct {
	OK      bool   `json:"ok"`
	Msg     string `json:"msg"`
	Version string `json:"version"` // Docker version if connection successful
}

// UnmarshalJSON implements custom JSON unmarshaling for TestResult.
// This is needed because Uptime Kuma may send the version as a nested object.
func (t *TestResult) UnmarshalJSON(data []byte) error {
	// First try to unmarshal directly
	type Alias TestResult
	aux := &struct {
		*Alias

		Version any `json:"version"`
	}{
		Alias: (*Alias)(t),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	// Handle version field which might be a string or an object
	switch v := aux.Version.(type) {
	case string:
		t.Version = v

	case map[string]any:
		// If version is an object, try to extract Version field
		if version, ok := v["Version"].(string); ok {
			t.Version = version
		}
	}

	return nil
}
