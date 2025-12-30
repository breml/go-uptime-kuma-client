package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SMTP represents an SMTP monitor for testing SMTP server connectivity.
type SMTP struct {
	Base
	SMTPDetails
}

// Type returns the monitor type.
func (s SMTP) Type() string {
	return s.SMTPDetails.Type()
}

// String returns a string representation of the SMTP monitor.
func (s SMTP) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SMTPDetails, true))
}

// UnmarshalJSON unmarshals an SMTP monitor from JSON data.
func (s *SMTP) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := SMTPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*s = SMTP{
		Base:        base,
		SMTPDetails: details,
	}

	return nil
}

// MarshalJSON marshals an SMTP monitor to JSON data.
func (s SMTP) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "smtp"
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

	// Always override with current SMTP-specific field values.
	raw["hostname"] = s.Hostname
	raw["port"] = s.Port
	raw["smtpSecurity"] = s.SMTPSecurity

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

// SMTPDetails contains SMTP-specific monitor configuration.
type SMTPDetails struct {
	Hostname     string  `json:"hostname"`
	Port         *int64  `json:"port"`
	SMTPSecurity *string `json:"smtpSecurity"`
}

// Type returns the monitor type.
func (SMTPDetails) Type() string {
	return "smtp"
}
