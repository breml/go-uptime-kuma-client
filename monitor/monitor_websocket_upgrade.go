package monitor

import (
	"encoding/json"
	"fmt"
)

// WebsocketUpgrade represents a websocket-upgrade monitor.
type WebsocketUpgrade struct {
	Base
	HTTPDetails
	WebsocketUpgradeDetails
}

// Type returns the monitor type.
func (w WebsocketUpgrade) Type() string {
	return w.WebsocketUpgradeDetails.Type()
}

// String returns a string representation of the monitor.
func (w WebsocketUpgrade) String() string {
	return fmt.Sprintf(
		"%s, %s, %s",
		formatMonitor(w.Base, false),
		formatMonitor(w.HTTPDetails, true),
		formatMonitor(w.WebsocketUpgradeDetails, true),
	)
}

// UnmarshalJSON unmarshals a JSON byte slice into a monitor.
func (w *WebsocketUpgrade) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal websocket-upgrade base: %w", err)
	}

	httpDetails := HTTPDetails{}
	err = json.Unmarshal(data, &httpDetails)
	if err != nil {
		return fmt.Errorf("unmarshal websocket-upgrade http details: %w", err)
	}

	wsDetails := WebsocketUpgradeDetails{}
	err = json.Unmarshal(data, &wsDetails)
	if err != nil {
		return fmt.Errorf("unmarshal websocket-upgrade details: %w", err)
	}

	*w = WebsocketUpgrade{
		Base:                    base,
		HTTPDetails:             httpDetails,
		WebsocketUpgradeDetails: wsDetails,
	}

	return nil
}

// MarshalJSON marshals a monitor into a JSON byte slice.
func (w WebsocketUpgrade) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = w.ID
	raw["type"] = "websocket-upgrade"
	raw["name"] = w.Name
	raw["description"] = w.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = w.PathName
	raw["parent"] = w.Parent
	raw["interval"] = w.Interval
	raw["retryInterval"] = w.RetryInterval
	raw["resendInterval"] = w.ResendInterval
	raw["maxretries"] = w.MaxRetries
	raw["upsideDown"] = w.UpsideDown
	raw["active"] = w.IsActive

	raw["notificationIDList"] = notificationIDMap(w.NotificationIDs)

	// Always override with current HTTP-specific field values.
	raw["url"] = w.URL
	raw["timeout"] = w.Timeout
	raw["expiryNotification"] = w.ExpiryNotification
	raw["ignoreTls"] = w.IgnoreTLS
	raw["maxredirects"] = w.MaxRedirects
	raw["accepted_statuscodes"] = w.AcceptedStatusCodes
	raw["proxyId"] = w.ProxyID
	raw["method"] = w.Method
	raw["httpBodyEncoding"] = w.HTTPBodyEncoding
	raw["body"] = w.Body
	raw["headers"] = w.Headers
	raw["authMethod"] = w.AuthMethod
	raw["bearer_token"] = w.BearerToken
	raw["basic_auth_user"] = w.BasicAuthUser
	raw["basic_auth_pass"] = w.BasicAuthPass
	raw["authDomain"] = w.AuthDomain
	raw["authWorkstation"] = w.AuthWorkstation
	raw["tlsCert"] = w.TLSCert
	raw["tlsKey"] = w.TLSKey
	raw["tlsCa"] = w.TLSCa
	raw["oauth_auth_method"] = w.OAuthAuthMethod
	raw["oauth_token_url"] = w.OAuthTokenURL
	raw["oauth_client_id"] = w.OAuthClientID
	raw["oauth_client_secret"] = w.OAuthClientSecret
	raw["oauth_scopes"] = w.OAuthScopes
	raw["oauth_audience"] = w.OAuthAudience
	raw["cacheBust"] = w.CacheBust

	// Always override with current WebSocket-specific field values.
	raw["wsIgnoreSecWebsocketAcceptHeader"] = w.IgnoreSecWebsocketAcceptHeader
	raw["wsSubprotocol"] = w.Subprotocol

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal websocket-upgrade monitor: %w", err)
	}

	return data, nil
}

// WebsocketUpgradeDetails contains websocket-upgrade-specific monitor configuration.
type WebsocketUpgradeDetails struct {
	IgnoreSecWebsocketAcceptHeader bool   `json:"wsIgnoreSecWebsocketAcceptHeader"`
	Subprotocol                    string `json:"wsSubprotocol"`
}

// Type returns the monitor type.
func (WebsocketUpgradeDetails) Type() string {
	return "websocket-upgrade"
}
