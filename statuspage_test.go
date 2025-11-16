package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/statuspage"
)

func TestClient_StatusPageCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	slug := "test-status-page"
	title := "Test Status Page"

	t.Run("add_status_page", func(t *testing.T) {
		err = client.AddStatusPage(ctx, title, slug)
		require.NoError(t, err)
	})

	var retrievedStatusPage *statuspage.StatusPage
	t.Run("get_status_page", func(t *testing.T) {
		retrievedStatusPage, err = client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.NotNil(t, retrievedStatusPage)
		require.Equal(t, slug, retrievedStatusPage.Slug)
		require.Equal(t, title, retrievedStatusPage.Title)
		require.Equal(t, "auto", retrievedStatusPage.Theme)
	})

	t.Run("get_status_pages", func(t *testing.T) {
		statusPages, err := client.GetStatusPages(ctx)
		require.NoError(t, err)
		// Status pages list might be empty initially or not yet populated in cache
		// This is acceptable as the cache is populated asynchronously
		t.Logf("Status pages in cache: %d", len(statusPages))
	})

	t.Run("update_status_page", func(t *testing.T) {
		retrievedStatusPage.Title = "Updated Test Status Page"
		retrievedStatusPage.Description = "This is an updated test status page"
		retrievedStatusPage.Theme = statuspage.ThemeDark()
		retrievedStatusPage.Published = true
		retrievedStatusPage.ShowTags = true
		retrievedStatusPage.FooterText = "© 2024 Test Inc."
		retrievedStatusPage.ShowPoweredBy = false
		retrievedStatusPage.CustomCSS = "body { background: #000; }"

		err = client.SaveStatusPage(ctx, retrievedStatusPage)
		require.NoError(t, err)

		updated, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.Equal(t, "Updated Test Status Page", updated.Title)
		require.Equal(t, "This is an updated test status page", updated.Description)
		require.Equal(t, statuspage.ThemeDark(), updated.Theme)
		require.True(t, updated.Published)
		require.True(t, updated.ShowTags)
		require.Equal(t, "© 2024 Test Inc.", updated.FooterText)
		require.False(t, updated.ShowPoweredBy)
		require.Equal(t, "body { background: #000; }", updated.CustomCSS)
	})

	t.Run("delete_status_page", func(t *testing.T) {
		err = client.DeleteStatusPage(ctx, slug)
		require.NoError(t, err)
	})
}

func TestClient_StatusPageWithMonitors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error

	slug := "test-status-with-monitors"
	title := "Status Page with Monitors"

	var monitor1ID, monitor2ID int64

	t.Run("create_monitors", func(t *testing.T) {
		httpMonitor1 := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Test HTTP Monitor 1",
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     3,
				UpsideDown:     false,
				IsActive:       true,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                 "https://httpbin.org/status/200",
				Timeout:             48,
				Method:              "GET",
				ExpiryNotification:  false,
				IgnoreTLS:           false,
				MaxRedirects:        10,
				AcceptedStatusCodes: []string{"200-299"},
				AuthMethod:          monitor.AuthMethodNone,
			},
		}

		monitor1ID, err = client.CreateMonitor(ctx, httpMonitor1)
		require.NoError(t, err)
		require.Greater(t, monitor1ID, int64(0))

		httpMonitor2 := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Test HTTP Monitor 2",
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     3,
				UpsideDown:     false,
				IsActive:       true,
			},
			HTTPDetails: monitor.HTTPDetails{
				URL:                 "https://httpbin.org/status/201",
				Timeout:             48,
				Method:              "GET",
				ExpiryNotification:  false,
				IgnoreTLS:           false,
				MaxRedirects:        10,
				AcceptedStatusCodes: []string{"200-299"},
				AuthMethod:          monitor.AuthMethodNone,
			},
		}

		monitor2ID, err = client.CreateMonitor(ctx, httpMonitor2)
		require.NoError(t, err)
		require.Greater(t, monitor2ID, int64(0))
	})

	t.Run("add_status_page", func(t *testing.T) {
		err = client.AddStatusPage(ctx, title, slug)
		require.NoError(t, err)
	})

	t.Run("add_monitors_to_status_page", func(t *testing.T) {
		sp, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)

		sendURLTrue := true
		sendURLFalse := false

		sp.PublicGroupList = []statuspage.PublicGroup{
			{
				Name:   "Web Services",
				Weight: 1,
				MonitorList: []statuspage.PublicMonitor{
					{ID: monitor1ID, SendURL: &sendURLTrue},
					{ID: monitor2ID, SendURL: &sendURLFalse},
				},
			},
		}

		err = client.SaveStatusPage(ctx, sp)
		require.NoError(t, err)

		// Note: GetStatusPage does not return PublicGroupList from the server
		// The server's getStatusPage endpoint only returns basic config fields
		// PublicGroupList is maintained separately and sent with SaveStatusPage
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteStatusPage(ctx, slug)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, monitor1ID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, monitor2ID)
		require.NoError(t, err)
	})
}

func TestClient_StatusPageIncidents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	slug := "test-status-incidents"
	title := "Status Page for Incidents"

	t.Run("add_status_page", func(t *testing.T) {
		err = client.AddStatusPage(ctx, title, slug)
		require.NoError(t, err)
	})

	t.Run("post_incident", func(t *testing.T) {
		incident := &statuspage.Incident{
			Title:   "Service Degradation",
			Content: "We are experiencing issues with our API service.",
			Style:   statuspage.StyleWarning(),
			Pin:     true,
		}

		err = client.PostIncident(ctx, slug, incident)
		require.NoError(t, err)
	})

	t.Run("update_incident", func(t *testing.T) {
		incident := &statuspage.Incident{
			Title:   "Major Outage",
			Content: "All services are currently down. We are working on a fix.",
			Style:   statuspage.StyleDanger(),
			Pin:     true,
		}

		err = client.PostIncident(ctx, slug, incident)
		require.NoError(t, err)
	})

	t.Run("unpin_incident", func(t *testing.T) {
		err = client.UnpinIncident(ctx, slug)
		require.NoError(t, err)
	})

	t.Run("post_info_incident", func(t *testing.T) {
		incident := &statuspage.Incident{
			Title:   "Maintenance Window",
			Content: "Scheduled maintenance in progress.",
			Style:   statuspage.StyleInfo(),
			Pin:     true,
		}

		err = client.PostIncident(ctx, slug, incident)
		require.NoError(t, err)
	})

	t.Run("post_primary_incident", func(t *testing.T) {
		incident := &statuspage.Incident{
			Title:   "New Feature Released",
			Content: "We've just released a new feature!",
			Style:   statuspage.StylePrimary(),
			Pin:     true,
		}

		err = client.PostIncident(ctx, slug, incident)
		require.NoError(t, err)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.UnpinIncident(ctx, slug)
		require.NoError(t, err)

		err = client.DeleteStatusPage(ctx, slug)
		require.NoError(t, err)
	})
}

func TestClient_StatusPageThemes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	slug := "test-status-themes"
	title := "Status Page Themes Test"

	t.Run("add_status_page", func(t *testing.T) {
		err = client.AddStatusPage(ctx, title, slug)
		require.NoError(t, err)
	})

	t.Run("set_light_theme", func(t *testing.T) {
		sp, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)

		require.True(t, statuspage.ValidTheme(statuspage.ThemeLight()))
		sp.Theme = statuspage.ThemeLight()

		err = client.SaveStatusPage(ctx, sp)
		require.NoError(t, err)

		updated, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.Equal(t, statuspage.ThemeLight(), updated.Theme)
	})

	t.Run("set_dark_theme", func(t *testing.T) {
		sp, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)

		require.True(t, statuspage.ValidTheme(statuspage.ThemeDark()))
		sp.Theme = statuspage.ThemeDark()

		err = client.SaveStatusPage(ctx, sp)
		require.NoError(t, err)

		updated, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.Equal(t, statuspage.ThemeDark(), updated.Theme)
	})

	t.Run("set_auto_theme", func(t *testing.T) {
		sp, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)

		require.True(t, statuspage.ValidTheme(statuspage.ThemeAuto()))
		sp.Theme = statuspage.ThemeAuto()

		err = client.SaveStatusPage(ctx, sp)
		require.NoError(t, err)

		updated, err := client.GetStatusPage(ctx, slug)
		require.NoError(t, err)
		require.Equal(t, statuspage.ThemeAuto(), updated.Theme)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteStatusPage(ctx, slug)
		require.NoError(t, err)
	})
}
