package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TCPPort represents a tcpport monitor.
type TCPPort struct {
	Base
	TCPPortDetails
}

// Type returns the monitor type.
func (t TCPPort) Type() string {
	return t.TCPPortDetails.Type()
}

// String returns a string representation of the monitor.
func (t TCPPort) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(t.Base, false), formatMonitor(t.TCPPortDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (t *TCPPort) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := TCPPortDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*t = TCPPort{
		Base:           base,
		TCPPortDetails: details,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (t TCPPort) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = t.ID
	raw["type"] = "port"
	raw["name"] = t.Name
	raw["description"] = t.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = t.PathName
	raw["parent"] = t.Parent
	raw["interval"] = t.Interval
	raw["retryInterval"] = t.RetryInterval
	raw["resendInterval"] = t.ResendInterval
	raw["maxretries"] = t.MaxRetries
	raw["upsideDown"] = t.UpsideDown
	raw["active"] = t.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range t.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current TCP Port-specific field values.
	raw["hostname"] = t.Hostname
	raw["port"] = t.Port
	raw["smtpSecurity"] = t.SMTPSecurity
	raw["expiryNotification"] = t.ExpiryNotification
	raw["expectedTlsAlert"] = t.ExpectedTLSAlert

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

// TCPPortDetails contains tcpport-specific monitor configuration.
type TCPPortDetails struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	// SMTPSecurity selects the TLS handshake mode used by the port
	// monitor. Valid values are "nostarttls" (plain TCP, no TLS handshake
	// or certificate validation), "secure" (immediate TLS handshake) and
	// "starttls" (negotiate TLS via STARTTLS). Only "secure" and
	// "starttls" perform a TLS handshake and certificate validation.
	// Leave nil/unset to fall back to a plain TCP check.
	SMTPSecurity *string `json:"smtpSecurity"`
	// ExpiryNotification enables certificate expiry notifications. It is
	// only honoured by the upstream server when SMTPSecurity is set to
	// "secure" or "starttls".
	ExpiryNotification bool `json:"expiryNotification"`
	// ExpectedTLSAlert is the TLS alert name that the server is expected
	// to return during the handshake (e.g. "certificate_required" for
	// mTLS verification). Leave nil/unset for a successful TLS handshake.
	ExpectedTLSAlert *string `json:"expectedTlsAlert"`
}

// Type returns the monitor type.
func (TCPPortDetails) Type() string {
	return "port"
}
