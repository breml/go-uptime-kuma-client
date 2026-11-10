package monitor

import (
	"encoding/json"
	"fmt"
)

// monitorConfig holds the keys of an Uptime Kuma 3.x monitor's config that the
// client reads.
type monitorConfig struct {
	SystemServiceName *string `json:"systemServiceName"`
}

// systemServiceNameFromConfig returns the name a pm2 or system-service monitor
// of an Uptime Kuma 3.x server carries as config.systemServiceName, and
// whether the payload had one. A 2.x server sends it as the top-level
// system_service_name instead, which the details structs decode directly.
func systemServiceNameFromConfig(data []byte) (name string, ok bool, err error) {
	var payload struct {
		Config *monitorConfig `json:"config"`
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return "", false, fmt.Errorf("unmarshal config: %w", err)
	}

	if payload.Config == nil || payload.Config.SystemServiceName == nil {
		return "", false, nil
	}

	return *payload.Config.SystemServiceName, true, nil
}
