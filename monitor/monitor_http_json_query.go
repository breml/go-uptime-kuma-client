package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type HTTPJSONQuery struct {
	Base
	HTTPDetails
	HTTPJSONQueryDetails
}

func (h HTTPJSONQuery) Type() string {
	return h.HTTPJSONQueryDetails.Type()
}

func (h HTTPJSONQuery) String() string {
	return fmt.Sprintf(
		"%s, %s, %s",
		formatMonitor(h.Base, false),
		formatMonitor(h.HTTPDetails, true),
		formatMonitor(h.HTTPJSONQueryDetails, true),
	)
}

func (h *HTTPJSONQuery) UnmarshalJSON(data []byte) error {
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

	jsonQueryDetails := HTTPJSONQueryDetails{}
	err = json.Unmarshal(data, &jsonQueryDetails)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*h = HTTPJSONQuery{
		Base:                 base,
		HTTPDetails:          httpDetails,
		HTTPJSONQueryDetails: jsonQueryDetails,
	}

	return nil
}

func (h HTTPJSONQuery) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = h.ID
	raw["type"] = "json-query"
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

	// Always override with current JSON Query-specific field values.
	raw["jsonPath"] = h.JSONPath
	raw["expectedValue"] = h.ExpectedValue
	raw["jsonPathOperator"] = h.JSONPathOperator

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return data, nil
}

type HTTPJSONQueryDetails struct {
	JSONPath         string `json:"jsonPath"`
	ExpectedValue    string `json:"expectedValue"`
	JSONPathOperator string `json:"jsonPathOperator"`
}

func (h HTTPJSONQueryDetails) Type() string {
	return "json-query"
}
