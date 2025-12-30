package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// DNS ...
type DNS struct {
	Base
	DNSDetails
}

// Type ...
func (d DNS) Type() string {
	return d.DNSDetails.Type()
}

// String ...
func (d DNS) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(d.Base, false), formatMonitor(d.DNSDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
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

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return data, nil
}

// DNSDetails ...
type DNSDetails struct {
	Hostname       string         `json:"hostname"`
	ResolverServer string         `json:"dns_resolve_server"`
	ResolveType    DNSResolveType `json:"dns_resolve_type"`
	Port           int            `json:"port"`
}

// Type ...
func (DNSDetails) Type() string {
	return "dns"
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
