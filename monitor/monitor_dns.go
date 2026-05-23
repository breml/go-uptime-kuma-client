package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// DNS represents a DNS monitor.
type DNS struct {
	Base
	DNSDetails
}

// Type returns the monitor type.
func (d DNS) Type() string {
	return d.DNSDetails.Type()
}

// String returns a string representation of the DNS monitor.
func (d DNS) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(d.Base, false), formatMonitor(d.DNSDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a DNS monitor.
func (d *DNS) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := DNSDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*d = DNS{
		Base:       base,
		DNSDetails: details,
	}

	return nil
}

// MarshalJSON marshals a DNS monitor into a JSON byte slice.
func (d DNS) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = d.ID
	raw["type"] = "dns"
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

	// Always override with current DNS-specific field values.
	raw["hostname"] = d.Hostname
	raw["dns_resolve_server"] = d.ResolverServer
	raw["dns_resolve_type"] = d.ResolveType
	raw["port"] = d.Port

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	if d.Conditions == nil {
		raw["conditions"] = []any{}
	} else {
		raw["conditions"] = d.Conditions
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return data, nil
}

// Condition carries the information that allow setting a condition on the monitor.
type Condition struct {
	Type     string `json:"type"`
	AndOr    string `json:"andOr"`
	Variable string `json:"variable"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// DNSDetails contains DNS-specific monitor configuration.
type DNSDetails struct {
	Hostname string `json:"hostname"`
	// ResolverServer holds one or more DNS resolver servers as a
	// comma-separated list of IP addresses or hostnames (e.g.
	// "1.1.1.1,8.8.8.8"). Use ResolverServers and SetResolverServers for
	// access as a []string.
	ResolverServer string         `json:"dns_resolve_server"`
	ResolveType    DNSResolveType `json:"dns_resolve_type"`
	Port           int            `json:"port"`
	Conditions     []Condition    `json:"conditions"`
}

// Type returns the monitor type.
func (DNSDetails) Type() string {
	return "dns"
}

// ResolverServers returns the configured DNS resolver servers as a slice.
// Whitespace around each entry is trimmed and empty entries are dropped,
// matching the parsing performed by the Uptime Kuma server.
func (d DNSDetails) ResolverServers() []string {
	if d.ResolverServer == "" {
		return nil
	}

	parts := strings.Split(d.ResolverServer, ",")
	servers := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		servers = append(servers, p)
	}

	if len(servers) == 0 {
		return nil
	}

	return servers
}

// SetResolverServers sets the DNS resolver servers from a slice. The values
// are joined with "," into the comma-separated form expected by the Uptime
// Kuma server. Empty and whitespace-only entries are dropped.
func (d *DNSDetails) SetResolverServers(servers []string) {
	cleaned := make([]string, 0, len(servers))
	for _, s := range servers {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		cleaned = append(cleaned, s)
	}

	d.ResolverServer = strings.Join(cleaned, ",")
}

// DNSResolveType represents the DNS record type to resolve.
type DNSResolveType string

// DNS resolve types.
const (
	DNSResolveTypeA     DNSResolveType = "A"
	DNSResolveTypeAAAA  DNSResolveType = "AAAA"
	DNSResolveTypeCAA   DNSResolveType = "CAA"
	DNSResolveTypeCNAME DNSResolveType = "CNAME"
	DNSResolveTypeMX    DNSResolveType = "MX"
	DNSResolveTypeNS    DNSResolveType = "NS"
	DNSResolveTypePTR   DNSResolveType = "PTR"
	DNSResolveTypeSOA   DNSResolveType = "SOA"
	DNSResolveTypeSRV   DNSResolveType = "SRV"
	DNSResolveTypeTXT   DNSResolveType = "TXT"
)
