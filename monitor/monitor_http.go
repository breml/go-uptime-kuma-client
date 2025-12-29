package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type HTTP struct {
	Base
	HTTPDetails
}

func (h HTTP) Type() string {
	return h.HTTPDetails.Type()
}

func (h HTTP) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(h.Base, false), formatMonitor(h.HTTPDetails, true))
}

func (h *HTTP) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := HTTPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*h = HTTP{
		Base:        base,
		HTTPDetails: details,
	}

	return nil
}

func (h HTTP) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = h.ID
	raw["type"] = "http"
	raw["name"] = h.Name
	raw["description"] = h.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = h.PathName
	raw["parent"] = h.Parent
	raw["interval"] = h.Interval
	raw["retryInterval"] = h.RetryInterval
	raw["resendInterval"] = h.ResendInterval
	raw["maxretries"] = h.MaxRetries
	raw["upsideDown"] = h.UpsideDown
	raw["active"] = h.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range h.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current HTTP-specific field values.
	raw["url"] = h.URL
	raw["timeout"] = h.Timeout
	raw["expiryNotification"] = h.ExpiryNotification
	raw["ignoreTls"] = h.IgnoreTLS
	raw["maxredirects"] = h.MaxRedirects
	raw["accepted_statuscodes"] = h.AcceptedStatusCodes
	raw["proxyId"] = h.ProxyID
	raw["method"] = h.Method
	raw["httpBodyEncoding"] = h.HTTPBodyEncoding
	raw["body"] = h.Body
	raw["headers"] = h.Headers
	raw["authMethod"] = h.AuthMethod
	raw["basic_auth_user"] = h.BasicAuthUser
	raw["basic_auth_pass"] = h.BasicAuthPass
	raw["authDomain"] = h.AuthDomain
	raw["authWorkstation"] = h.AuthWorkstation
	raw["tlsCert"] = h.TLSCert
	raw["tlsKey"] = h.TLSKey
	raw["tlsCa"] = h.TLSCa
	raw["oauth_auth_method"] = h.OAuthAuthMethod
	raw["oauth_token_url"] = h.OAuthTokenURL
	raw["oauth_client_id"] = h.OAuthClientID
	raw["oauth_client_secret"] = h.OAuthClientSecret
	raw["oauth_scopes"] = h.OAuthScopes
	raw["cacheBust"] = h.CacheBust

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return data, nil
}

type HTTPDetails struct {
	URL                 string     `json:"url"`
	Timeout             int64      `json:"timeout"`
	ExpiryNotification  bool       `json:"expiryNotification"`
	IgnoreTLS           bool       `json:"ignoreTls"`
	MaxRedirects        int        `json:"maxredirects"`
	AcceptedStatusCodes []string   `json:"accepted_statuscodes"`
	Method              string     `json:"method"`
	HTTPBodyEncoding    string     `json:"httpBodyEncoding"`
	Body                string     `json:"body"`
	Headers             string     `json:"headers"`
	AuthMethod          AuthMethod `json:"authMethod"`
	BasicAuthUser       string     `json:"basic_auth_user"`
	BasicAuthPass       string     `json:"basic_auth_pass"`
	AuthDomain          string     `json:"authDomain"`
	AuthWorkstation     string     `json:"authWorkstation"`
	TLSCert             string     `json:"tlsCert"`
	TLSKey              string     `json:"tlsKey"`
	TLSCa               string     `json:"tlsCa"`
	OAuthAuthMethod     string     `json:"oauth_auth_method"`
	OAuthTokenURL       string     `json:"oauth_token_url"`
	OAuthClientID       string     `json:"oauth_client_id"`
	OAuthClientSecret   string     `json:"oauth_client_secret"`
	OAuthScopes         string     `json:"oauth_scopes"`
	CacheBust           bool       `json:"cacheBust"`
}

func (h HTTPDetails) Type() string {
	return "http"
}

// AuthMethod represents the authentication method for monitors.
type AuthMethod string

// Auth methods.
const (
	AuthMethodNone     AuthMethod = ""
	AuthMethodBasic    AuthMethod = "basic"
	AuthMethodNTLM     AuthMethod = "ntlm"
	AuthMethodMTLS     AuthMethod = "mtls"
	AuthMethodOAtuh2CC AuthMethod = "oauth2-cc"
)
