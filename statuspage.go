package kuma

import (
	"context"
	"fmt"

	"github.com/breml/go-uptime-kuma-client/statuspage"
)

// GetStatusPages retrieves all status pages from the client cache.
func (c *Client) GetStatusPages(ctx context.Context) (map[int64]statuspage.StatusPage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state.statusPages == nil {
		return map[int64]statuspage.StatusPage{}, nil
	}

	statusPages := make(map[int64]statuspage.StatusPage, len(c.state.statusPages))
	for k, v := range c.state.statusPages {
		statusPages[k] = v
	}

	return statusPages, nil
}

// GetStatusPage retrieves a specific status page by slug.
// Note: The server does not return PublicGroupList in this endpoint.
// PublicGroupList must be maintained separately when calling SaveStatusPage.
func (c *Client) GetStatusPage(ctx context.Context, slug string) (*statuspage.StatusPage, error) {
	response, err := c.syncEmit(ctx, "getStatusPage", slug)
	if err != nil {
		return nil, fmt.Errorf("get status page %s: %v", slug, err)
	}

	if response.Config == nil {
		return nil, fmt.Errorf("get status page %s: config not found in response", slug)
	}

	var sp statuspage.StatusPage
	err = convertToStruct(response.Config, &sp)
	if err != nil {
		return nil, fmt.Errorf("get status page %s: %v", slug, err)
	}

	return &sp, nil
}

// AddStatusPage creates a new status page with the given title and slug.
func (c *Client) AddStatusPage(ctx context.Context, title, slug string) error {
	_, err := c.syncEmit(ctx, "addStatusPage", title, slug)
	if err != nil {
		return fmt.Errorf("add status page: %v", err)
	}

	return nil
}

// SaveStatusPage updates an existing status page configuration and returns the updated public group list with IDs.
func (c *Client) SaveStatusPage(ctx context.Context, sp *statuspage.StatusPage) ([]statuspage.PublicGroup, error) {
	config := map[string]any{
		"slug":                  sp.Slug,
		"title":                 sp.Title,
		"description":           sp.Description,
		"theme":                 sp.Theme,
		"published":             sp.Published,
		"showTags":              sp.ShowTags,
		"domainNameList":        sp.DomainNameList,
		"googleAnalyticsId":     sp.GoogleAnalyticsID,
		"customCSS":             sp.CustomCSS,
		"footerText":            sp.FooterText,
		"showPoweredBy":         sp.ShowPoweredBy,
		"showCertificateExpiry": sp.ShowCertificateExpiry,
	}

	imgDataURL := sp.Icon

	publicGroupList := make([]map[string]any, len(sp.PublicGroupList))
	for i, group := range sp.PublicGroupList {
		monitorList := make([]map[string]any, len(group.MonitorList))
		for j, monitor := range group.MonitorList {
			monitorData := map[string]any{
				"id": monitor.ID,
			}
			if monitor.SendURL != nil {
				monitorData["sendUrl"] = *monitor.SendURL
			}

			monitorList[j] = monitorData
		}

		publicGroupList[i] = map[string]any{
			"id":          group.ID,
			"name":        group.Name,
			"weight":      group.Weight,
			"monitorList": monitorList,
		}
	}

	response, err := c.syncEmit(ctx, "saveStatusPage", sp.Slug, config, imgDataURL, publicGroupList)
	if err != nil {
		return nil, fmt.Errorf("save status page: %v", err)
	}

	// Parse the returned public group list with IDs
	var groups []statuspage.PublicGroup
	if response.PublicGroupList != nil {
		err = convertToStruct(response.PublicGroupList, &groups)
		if err != nil {
			return nil, fmt.Errorf("save status page: failed to parse response public group list: %v", err)
		}
	}

	return groups, nil
}

// DeleteStatusPage deletes a status page by slug.
func (c *Client) DeleteStatusPage(ctx context.Context, slug string) error {
	_, err := c.syncEmit(ctx, "deleteStatusPage", slug)
	if err != nil {
		return fmt.Errorf("delete status page: %v", err)
	}

	return nil
}

// PostIncident posts or updates an incident on a status page.
func (c *Client) PostIncident(ctx context.Context, slug string, incident *statuspage.Incident) error {
	incidentData, err := structToMap(incident)
	if err != nil {
		return fmt.Errorf("post incident: %v", err)
	}

	_, err = c.syncEmit(ctx, "postIncident", slug, incidentData)
	if err != nil {
		return fmt.Errorf("post incident: %v", err)
	}

	return nil
}

// UnpinIncident unpins the currently pinned incident on a status page.
func (c *Client) UnpinIncident(ctx context.Context, slug string) error {
	_, err := c.syncEmit(ctx, "unpinIncident", slug)
	if err != nil {
		return fmt.Errorf("unpin incident: %v", err)
	}

	return nil
}
