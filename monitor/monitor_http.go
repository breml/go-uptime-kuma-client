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
		return err
	}

	details := HTTPDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*h = HTTP{
		Base:        base,
		HTTPDetails: details,
	}

	return nil
}

func (h HTTP) MarshalJSON() ([]byte, error) {
	if h.Base.raw == nil || h.Base.internalType == "" {
		return nil, fmt.Errorf("not unmarshaled monitor, unable to marshal")
	}

	raw := map[string]any{}
	err := json.Unmarshal(h.Base.raw, &raw)
	if err != nil {
		return nil, fmt.Errorf("invalid internal state for raw, failed to unmarshal: %w", err)
	}

	// Update base fields
	raw["id"] = h.Base.ID
	raw["type"] = h.Base.internalType
	raw["name"] = h.Base.Name
	raw["description"] = h.Base.Description
	raw["pathName"] = h.Base.PathName
	raw["interval"] = h.Base.Interval
	raw["retryInterval"] = h.Base.RetryInterval
	raw["resendInterval"] = h.Base.ResendInterval
	raw["maxretries"] = h.Base.MaxRetries
	raw["upsideDown"] = h.Base.UpsideDown
	raw["active"] = h.Base.IsActive

	// Update notification IDs
	ids := map[string]bool{}
	for _, id := range h.Base.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}
	raw["notificationIDList"] = ids

	// Update HTTP-specific fields
	raw["url"] = h.HTTPDetails.URL
	raw["timeout"] = h.HTTPDetails.Timeout
	raw["expiryNotification"] = h.HTTPDetails.ExpiryNotification
	raw["ignoreTls"] = h.HTTPDetails.IgnoreTLS
	raw["maxredirects"] = h.HTTPDetails.MaxRedirects
	raw["accepted_statuscodes"] = h.HTTPDetails.AcceptedStatusCodes
	raw["proxyID"] = h.HTTPDetails.ProxyID
	raw["method"] = h.HTTPDetails.Method
	raw["httpBodyEncoding"] = h.HTTPDetails.HTTPBodyEncoding
	raw["body"] = h.HTTPDetails.Body
	raw["headers"] = h.HTTPDetails.Headers
	raw["authMethod"] = h.HTTPDetails.AuthMethod
	raw["basic_auth_user"] = h.HTTPDetails.BasicAuthUser
	raw["basic_auth_pass"] = h.HTTPDetails.BasicAuthPass
	raw["authDomain"] = h.HTTPDetails.AuthDomain
	raw["authWorkstation"] = h.HTTPDetails.AuthWorkstation
	raw["tlsCert"] = h.HTTPDetails.TLSCert
	raw["tlsKey"] = h.HTTPDetails.TLSKey
	raw["tlsCa"] = h.HTTPDetails.TLSCa
	raw["oauth_auth_method"] = h.HTTPDetails.OAuthAuthMethod
	raw["oauth_token_url"] = h.HTTPDetails.OAuthTokenURL
	raw["oauth_client_id"] = h.HTTPDetails.OAuthClientID
	raw["oauth_client_secret"] = h.HTTPDetails.OAuthClientSecret
	raw["oauth_scopes"] = h.HTTPDetails.OAuthScopes

	return json.Marshal(raw)
}

type HTTPDetails struct {
	URL                 string     `json:"url"`
	Timeout             int64      `json:"timeout"`
	ExpiryNotification  bool       `json:"expiryNotification"`
	IgnoreTLS           bool       `json:"ignoreTls"`
	MaxRedirects        int        `json:"maxredirects"`
	AcceptedStatusCodes []string   `json:"accepted_statuscodes"`
	ProxyID             int64      `json:"proxyID"`
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
}

func (h HTTPDetails) Type() string {
	return "http"
}

// AuthMethod represents the authentication method for monitors
type AuthMethod string

// Auth methods
const (
	AuthMethodNone     AuthMethod = ""
	AuthMethodBasic    AuthMethod = "basic"
	AuthMethodNTLM     AuthMethod = "ntlm"
	AuthMethodMTLS     AuthMethod = "mtls"
	AuthMethodOAtuh2CC AuthMethod = "oauth2-cc"
)
