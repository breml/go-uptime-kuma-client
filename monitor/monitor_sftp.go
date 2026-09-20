package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// defaultSFTPTimeout is the connection timeout in seconds that the SFTP check
// falls back to. The client sends it explicitly, because the monitor.timeout
// column is NOT NULL and the server rejects an explicit null.
const defaultSFTPTimeout = 10

// SFTPAuthMethod represents the SSH authentication method of an SFTP monitor.
type SFTPAuthMethod string

const (
	// SFTPAuthMethodPassword authenticates with SFTPDetails.SSHPassword.
	SFTPAuthMethodPassword SFTPAuthMethod = "password"
	// SFTPAuthMethodPrivateKey authenticates with SFTPDetails.SSHPrivateKey.
	SFTPAuthMethodPrivateKey SFTPAuthMethod = "privateKey"
)

// SFTP represents an SFTP monitor, which opens an SSH connection to a host and
// optionally checks that a remote path exists.
type SFTP struct {
	Base
	SFTPDetails
}

// Type returns the monitor type.
func (s SFTP) Type() string {
	return s.SFTPDetails.Type()
}

// String returns a string representation of the SFTP monitor.
func (s SFTP) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(s.Base, false), formatMonitor(s.SFTPDetails, true))
}

// UnmarshalJSON unmarshals an SFTP monitor from JSON data.
func (s *SFTP) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := SFTPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*s = SFTP{
		Base:        base,
		SFTPDetails: details,
	}

	return nil
}

// MarshalJSON marshals an SFTP monitor to JSON data.
func (s SFTP) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = s.ID
	raw["type"] = "sftp"
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

	// Always override with current SFTP-specific field values.
	raw["hostname"] = s.Hostname
	raw["port"] = s.Port
	raw["sshUsername"] = s.SSHUsername
	raw["sshPassword"] = s.SSHPassword
	raw["sshPrivateKey"] = s.SSHPrivateKey
	raw["sshPassphrase"] = s.SSHPassphrase
	raw["sftpPath"] = s.SFTPPath

	// The server reads this column back as `sshAuthMethod || "password"`, so a
	// value is always sent to keep the field round-tripping.
	authMethod := s.SSHAuthMethod
	if authMethod == "" {
		authMethod = SFTPAuthMethodPassword
	}

	raw["sshAuthMethod"] = authMethod

	// The monitor.timeout column is NOT NULL, so an unset Timeout is sent as
	// the value the check would fall back to instead of as null.
	timeout := float64(defaultSFTPTimeout)
	if s.Timeout != nil {
		timeout = *s.Timeout
	}

	raw["timeout"] = timeout

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

// SFTPDetails contains SFTP-specific monitor configuration.
//
// The behaviour described below was verified against Uptime Kuma 2.5.5.
//
// SSHUsername, SSHPassword, SSHPrivateKey and SSHPassphrase are only part of
// the sensitive section of the server's monitor representation, so they read
// back empty for anyone but the owner of the monitor.
type SFTPDetails struct {
	// Hostname is the SFTP server address. It is required, the check fails on
	// every heartbeat while it is empty.
	Hostname string `json:"hostname"`
	// Port is the TCP port of the SFTP server. The column has no default;
	// while it is NULL the check falls back to 22.
	Port *int64 `json:"port"`
	// Timeout is the connection timeout in seconds, applied to both the SSH
	// ready timeout and the socket timeout. The column is NOT NULL, so unlike
	// the other optional fields a nil value is not sent as null: MarshalJSON
	// substitutes 10, the value the check itself falls back to. The column is
	// a floating point column, so fractional values round-trip unchanged.
	Timeout *float64 `json:"timeout"`
	// SSHUsername is the user to authenticate as. It is required, the check
	// fails on every heartbeat while it is empty.
	SSHUsername string `json:"sshUsername"`
	// SSHAuthMethod selects how the check authenticates. An empty value is
	// sent as SFTPAuthMethodPassword, which is what the server reads back for
	// a NULL column.
	SSHAuthMethod SFTPAuthMethod `json:"sshAuthMethod"`
	// SSHPassword is the password used when SSHAuthMethod is not
	// SFTPAuthMethodPrivateKey.
	SSHPassword *string `json:"sshPassword"`
	// SSHPrivateKey is the private key used when SSHAuthMethod is
	// SFTPAuthMethodPrivateKey. The check errors out when the method is
	// SFTPAuthMethodPrivateKey and no key is set.
	SSHPrivateKey *string `json:"sshPrivateKey"`
	// SSHPassphrase is the optional passphrase of SSHPrivateKey. It is only
	// applied while it is non-empty.
	SSHPassphrase *string `json:"sshPassphrase"`
	// SFTPPath is an optional remote path. While it is non-empty the check
	// additionally verifies that the path exists on the server.
	SFTPPath *string `json:"sftpPath"`
}

// Type returns the monitor type.
func (SFTPDetails) Type() string {
	return "sftp"
}
