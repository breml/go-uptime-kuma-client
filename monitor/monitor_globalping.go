package monitor

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Globalping represents a globalping monitor.
type Globalping struct {
	Base
	HTTPDetails
	GlobalpingDetails
}

// Type returns the monitor type.
func (g Globalping) Type() string {
	return g.GlobalpingDetails.Type()
}

// String returns a string representation of the monitor.
func (g Globalping) String() string {
	return fmt.Sprintf(
		"%s, %s, %s",
		formatMonitor(g.Base, false),
		formatMonitor(g.HTTPDetails, true),
		formatMonitor(g.GlobalpingDetails, true),
	)
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (g *Globalping) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	httpDetails := HTTPDetails{}
	err = json.Unmarshal(data, &httpDetails)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	globalpingDetails := GlobalpingDetails{}
	err = json.Unmarshal(data, &globalpingDetails)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*g = Globalping{
		Base:              base,
		HTTPDetails:       httpDetails,
		GlobalpingDetails: globalpingDetails,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (g Globalping) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = g.ID
	raw["type"] = "globalping"
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

	raw["notificationIDList"] = notificationIDMap(g.NotificationIDs)

	// Always override with current HTTP-specific field values.
	raw["url"] = g.URL
	raw["timeout"] = g.Timeout
	raw["expiryNotification"] = g.ExpiryNotification
	raw["ignoreTls"] = g.IgnoreTLS
	raw["maxredirects"] = g.MaxRedirects
	raw["accepted_statuscodes"] = g.AcceptedStatusCodes
	raw["proxyId"] = g.ProxyID
	raw["method"] = g.Method
	raw["httpBodyEncoding"] = g.HTTPBodyEncoding
	raw["body"] = g.Body
	raw["headers"] = g.Headers
	raw["authMethod"] = g.AuthMethod
	raw["basic_auth_user"] = g.BasicAuthUser
	raw["basic_auth_pass"] = g.BasicAuthPass
	raw["authDomain"] = g.AuthDomain
	raw["authWorkstation"] = g.AuthWorkstation
	raw["tlsCert"] = g.TLSCert
	raw["tlsKey"] = g.TLSKey
	raw["tlsCa"] = g.TLSCa
	raw["oauth_auth_method"] = g.OAuthAuthMethod
	raw["oauth_token_url"] = g.OAuthTokenURL
	raw["oauth_client_id"] = g.OAuthClientID
	raw["oauth_client_secret"] = g.OAuthClientSecret
	raw["oauth_scopes"] = g.OAuthScopes
	raw["oauth_audience"] = g.OAuthAudience
	raw["cacheBust"] = g.CacheBust

	// Always override with current Globalping-specific field values.
	maps.Copy(raw, g.toMap())

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return data, nil
}

// GlobalpingDetails contains globalping-specific monitor configuration.
type GlobalpingDetails struct {
	Subtype          GlobalpingSubtype  `json:"subtype"`
	Location         string             `json:"location"`
	IPFamily         GlobalpingIPFamily `json:"ipFamily"`
	Protocol         string             `json:"protocol"`
	PingCount        int                `json:"ping_count"`
	Hostname         string             `json:"hostname"`
	Port             int                `json:"port"`
	DNSResolveType   DNSResolveType     `json:"dns_resolve_type"`
	DNSResolveServer string             `json:"dns_resolve_server"`
	Keyword          string             `json:"keyword"`
	InvertKeyword    bool               `json:"invertKeyword"`
	ExpectedValue    string             `json:"expectedValue"`
	JSONPath         string             `json:"jsonPath"`
	JSONPathOperator string             `json:"jsonPathOperator"`
}

// Type returns the monitor type.
func (GlobalpingDetails) Type() string {
	return "globalping"
}

// toMap returns the Globalping-specific fields as a map suitable for merging
// into the JSON payload.
func (d GlobalpingDetails) toMap() map[string]any {
	return map[string]any{
		"subtype":            d.Subtype,
		"location":           d.Location,
		"ipFamily":           d.IPFamily,
		"protocol":           d.Protocol,
		"ping_count":         d.PingCount,
		"hostname":           d.Hostname,
		"port":               d.Port,
		"dns_resolve_type":   d.DNSResolveType,
		"dns_resolve_server": d.DNSResolveServer,
		"keyword":            d.Keyword,
		"invertKeyword":      d.InvertKeyword,
		"expectedValue":      d.ExpectedValue,
		"jsonPath":           d.JSONPath,
		"jsonPathOperator":   d.JSONPathOperator,
	}
}

// GlobalpingSubtype represents the subtype of a Globalping monitor.
type GlobalpingSubtype string

// Globalping subtypes.
const (
	GlobalpingSubtypePing       GlobalpingSubtype = "ping"
	GlobalpingSubtypeTraceroute GlobalpingSubtype = "traceroute"
	GlobalpingSubtypeDNS        GlobalpingSubtype = "dns"
	GlobalpingSubtypeHTTP       GlobalpingSubtype = "http"
)

// String returns the string representation of the Globalping subtype.
func (s GlobalpingSubtype) String() string {
	return string(s)
}

// GlobalpingIPFamily represents the IP family for a Globalping monitor.
type GlobalpingIPFamily string

// Globalping IP families. The empty value selects the IP family automatically.
const (
	GlobalpingIPFamilyAuto GlobalpingIPFamily = ""
	GlobalpingIPFamilyIPv4 GlobalpingIPFamily = "ipv4"
	GlobalpingIPFamilyIPv6 GlobalpingIPFamily = "ipv6"
)

// String returns the string representation of the Globalping IP family.
func (f GlobalpingIPFamily) String() string {
	return string(f)
}
