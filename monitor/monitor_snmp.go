package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SNMP represents a Simple Network Management Protocol monitor.
type SNMP struct {
	Base
	SNMPDetails
}

// Type ...
func (s SNMP) Type() string {
	return s.SNMPDetails.Type()
}

// String ...
func (s SNMP) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SNMPDetails, true))
}

// UnmarshalJSON ...
func (s *SNMP) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := SNMPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*s = SNMP{
		Base:        base,
		SNMPDetails: details,
	}

	return nil
}

// MarshalJSON ...
func (s SNMP) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "snmp"
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

	// Always override with current SNMP-specific field values.
	raw["hostname"] = s.Hostname
	raw["port"] = s.Port
	raw["snmpVersion"] = s.SNMPVersion
	raw["snmpOid"] = s.SNMPOID
	raw["radiusPassword"] = s.SNMPCommunity
	raw["jsonPath"] = s.JSONPath
	raw["jsonPathOperator"] = s.JSONPathOperator
	raw["expectedValue"] = s.ExpectedValue

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

// SNMPDetails contains SNMP-specific monitor configuration.
type SNMPDetails struct {
	Hostname         string  `json:"hostname"`
	Port             *int64  `json:"port"`
	SNMPVersion      string  `json:"snmpVersion"`
	SNMPOID          string  `json:"snmpOid"`
	SNMPCommunity    string  `json:"radiusPassword"`
	JSONPath         *string `json:"jsonPath"`
	JSONPathOperator *string `json:"jsonPathOperator"`
	ExpectedValue    *string `json:"expectedValue"`
}

// Type ...
func (s SNMPDetails) Type() string {
	return "snmp"
}
