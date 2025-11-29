package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type RealBrowser struct {
	Base
	RealBrowserDetails
}

func (r RealBrowser) Type() string {
	return r.RealBrowserDetails.Type()
}

func (r RealBrowser) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(r.Base, false), formatMonitor(r.RealBrowserDetails, true))
}

func (r *RealBrowser) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := RealBrowserDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*r = RealBrowser{
		Base:               base,
		RealBrowserDetails: details,
	}

	return nil
}

func (r RealBrowser) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = r.ID
	raw["type"] = "real-browser"
	raw["name"] = r.Name
	raw["description"] = r.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = r.PathName
	raw["parent"] = r.Parent
	raw["interval"] = r.Interval
	raw["retryInterval"] = r.RetryInterval
	raw["resendInterval"] = r.ResendInterval
	raw["maxretries"] = r.MaxRetries
	raw["upsideDown"] = r.UpsideDown
	raw["active"] = r.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range r.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}
	raw["notificationIDList"] = ids

	// Always override with current RealBrowser-specific field values.
	raw["url"] = r.URL
	raw["timeout"] = r.Timeout
	raw["ignoreTls"] = r.IgnoreTLS
	raw["maxredirects"] = r.MaxRedirects
	raw["accepted_statuscodes"] = r.AcceptedStatusCodes
	raw["proxyId"] = r.ProxyID
	raw["remote_browser"] = r.RemoteBrowser

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	return json.Marshal(raw)
}

type RealBrowserDetails struct {
	URL                 string   `json:"url"`
	Timeout             int64    `json:"timeout"`
	IgnoreTLS           bool     `json:"ignoreTls"`
	MaxRedirects        int      `json:"maxredirects"`
	AcceptedStatusCodes []string `json:"accepted_statuscodes"`
	RemoteBrowser       *int64   `json:"remote_browser,omitempty"`
}

func (r RealBrowserDetails) Type() string {
	return "real-browser"
}
