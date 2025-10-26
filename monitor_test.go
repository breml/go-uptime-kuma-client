package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestClient_MonitorCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test HTTP monitor
	httpMonitor := monitor.HTTP{
		Base: monitor.Base{
			Name:           "Test HTTP Monitor",
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

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var httpMonitorRetrieved monitor.HTTP
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, httpMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test HTTP Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs

		err = client.GetMonitorAs(ctx, monitorID, &httpMonitorRetrieved)
		require.NoError(t, err)
		httpMonitor.ID = monitorID
		httpMonitor.PathName = httpMonitor.Name
		require.EqualExportedValues(t, httpMonitor, httpMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		httpMonitorRetrieved.Name = "Updated HTTP Monitor"
		httpMonitorRetrieved.URL = "https://httpbin.org/status/201"
		err := client.UpdateMonitor(ctx, httpMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated HTTP Monitor", updatedMonitor.Name)

		var updatedHTTP monitor.HTTP
		err = client.GetMonitorAs(ctx, monitorID, &updatedHTTP)
		require.NoError(t, err)
		require.Equal(t, "https://httpbin.org/status/201", updatedHTTP.URL)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorGroupCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test Group monitor
	groupMonitor := monitor.Group{
		Base: monitor.Base{
			Name:           "Test Group Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     0,
			UpsideDown:     false,
			IsActive:       false,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var groupMonitorRetrieved monitor.Group
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, groupMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test Group Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &groupMonitorRetrieved)
		require.NoError(t, err)
		groupMonitor.ID = monitorID
		groupMonitor.PathName = groupMonitor.Name
		require.EqualExportedValues(t, groupMonitor, groupMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		groupMonitorRetrieved.Name = "Updated Group Monitor"
		err := client.UpdateMonitor(ctx, groupMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated Group Monitor", updatedMonitor.Name)

		var updatedGroup monitor.Group
		err = client.GetMonitorAs(ctx, monitorID, &updatedGroup)
		require.NoError(t, err)
		require.Equal(t, "Updated Group Monitor", updatedGroup.Name)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorPingCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test Ping monitor
	pingMonitor := monitor.Ping{
		Base: monitor.Base{
			Name:           "Test Ping Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		PingDetails: monitor.PingDetails{
			Hostname:   "8.8.8.8",
			PacketSize: 56,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var pingMonitorRetrieved monitor.Ping
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, pingMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test Ping Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &pingMonitorRetrieved)
		require.NoError(t, err)
		pingMonitor.ID = monitorID
		pingMonitor.PathName = pingMonitor.Name
		require.EqualExportedValues(t, pingMonitor, pingMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		pingMonitorRetrieved.Name = "Updated Ping Monitor"
		pingMonitorRetrieved.Hostname = "1.1.1.1"
		pingMonitorRetrieved.PacketSize = 64
		err := client.UpdateMonitor(ctx, pingMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated Ping Monitor", updatedMonitor.Name)

		var updatedPing monitor.Ping
		err = client.GetMonitorAs(ctx, monitorID, &updatedPing)
		require.NoError(t, err)
		require.Equal(t, "1.1.1.1", updatedPing.Hostname)
		require.Equal(t, 64, updatedPing.PacketSize)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorPushCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test Push monitor
	pushMonitor := monitor.Push{
		Base: monitor.Base{
			Name:           "Test Push Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     0,
			UpsideDown:     false,
			IsActive:       true,
		},
		PushDetails: monitor.PushDetails{
			PushToken: "testtoken123",
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var pushMonitorRetrieved monitor.Push
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, pushMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test Push Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &pushMonitorRetrieved)
		require.NoError(t, err)
		pushMonitor.ID = monitorID
		pushMonitor.PathName = pushMonitor.Name
		require.EqualExportedValues(t, pushMonitor, pushMonitorRetrieved)
		require.NotEmpty(t, pushMonitorRetrieved.PushToken)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		pushMonitorRetrieved.Name = "Updated Push Monitor"
		err := client.UpdateMonitor(ctx, pushMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated Push Monitor", updatedMonitor.Name)

		var updatedPush monitor.Push
		err = client.GetMonitorAs(ctx, monitorID, &updatedPush)
		require.NoError(t, err)
		require.Equal(t, "Updated Push Monitor", updatedPush.Name)
		require.Equal(t, pushMonitorRetrieved.PushToken, updatedPush.PushToken)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorTCPPortCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test TCP Port monitor
	tcpPortMonitor := monitor.TCPPort{
		Base: monitor.Base{
			Name:           "Test TCP Port Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		TCPPortDetails: monitor.TCPPortDetails{
			Hostname: "example.com",
			Port:     443,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var tcpPortMonitorRetrieved monitor.TCPPort
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, tcpPortMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test TCP Port Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &tcpPortMonitorRetrieved)
		require.NoError(t, err)
		tcpPortMonitor.ID = monitorID
		tcpPortMonitor.PathName = tcpPortMonitor.Name
		require.EqualExportedValues(t, tcpPortMonitor, tcpPortMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		tcpPortMonitorRetrieved.Name = "Updated TCP Port Monitor"
		tcpPortMonitorRetrieved.Hostname = "cloudflare.com"
		tcpPortMonitorRetrieved.Port = 80
		err := client.UpdateMonitor(ctx, tcpPortMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated TCP Port Monitor", updatedMonitor.Name)

		var updatedTCPPort monitor.TCPPort
		err = client.GetMonitorAs(ctx, monitorID, &updatedTCPPort)
		require.NoError(t, err)
		require.Equal(t, "cloudflare.com", updatedTCPPort.Hostname)
		require.Equal(t, 80, updatedTCPPort.Port)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorParent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	groupMonitor := monitor.Group{
		Base: monitor.Base{
			Name:           "Parent Group",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     0,
			UpsideDown:     false,
			IsActive:       true,
		},
	}

	var parentID int64
	t.Run("create_parent_group", func(t *testing.T) {
		parentID, err = client.CreateMonitor(ctx, groupMonitor)
		require.NoError(t, err)
		require.Greater(t, parentID, int64(0))
	})

	var childID int64
	t.Run("create_monitor_with_parent", func(t *testing.T) {
		httpMonitor := monitor.HTTP{
			Base: monitor.Base{
				Name:           "Child Monitor",
				Parent:         &parentID,
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

		childID, err = client.CreateMonitor(ctx, httpMonitor)
		require.NoError(t, err)
		require.Greater(t, childID, int64(0))

		var childMonitorRetrieved monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &childMonitorRetrieved)
		require.NoError(t, err)
		require.NotNil(t, childMonitorRetrieved.Parent)
		require.Equal(t, parentID, *childMonitorRetrieved.Parent)
	})

	var newParentID int64
	t.Run("update_parent", func(t *testing.T) {
		newGroupMonitor := monitor.Group{
			Base: monitor.Base{
				Name:           "New Parent Group",
				Interval:       60,
				RetryInterval:  60,
				ResendInterval: 0,
				MaxRetries:     0,
				UpsideDown:     false,
				IsActive:       true,
			},
		}

		newParentID, err = client.CreateMonitor(ctx, newGroupMonitor)
		require.NoError(t, err)
		require.Greater(t, newParentID, int64(0))

		var httpMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &httpMonitor)
		require.NoError(t, err)

		httpMonitor.Parent = &newParentID

		err = client.UpdateMonitor(ctx, httpMonitor)
		require.NoError(t, err)

		var updatedMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &updatedMonitor)
		require.NoError(t, err)
		require.NotNil(t, updatedMonitor.Parent)
		require.Equal(t, newParentID, *updatedMonitor.Parent)
	})

	t.Run("remove_parent", func(t *testing.T) {
		var httpMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &httpMonitor)
		require.NoError(t, err)

		httpMonitor.Parent = nil

		err = client.UpdateMonitor(ctx, httpMonitor)
		require.NoError(t, err)

		var updatedMonitor monitor.HTTP
		err = client.GetMonitorAs(ctx, childID, &updatedMonitor)
		require.NoError(t, err)
		require.Nil(t, updatedMonitor.Parent)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteMonitor(ctx, childID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, newParentID)
		require.NoError(t, err)

		err = client.DeleteMonitor(ctx, parentID)
		require.NoError(t, err)
	})
}

func TestClient_MonitorHTTPKeywordCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	httpKeywordMonitor := monitor.HTTPKeyword{
		Base: monitor.Base{
			Name:           "Test HTTP Keyword Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 "https://httpbin.org/html",
			Timeout:             48,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
			AuthMethod:          monitor.AuthMethodNone,
		},
		HTTPKeywordDetails: monitor.HTTPKeywordDetails{
			Keyword:       "Herman Melville",
			InvertKeyword: false,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var httpKeywordMonitorRetrieved monitor.HTTPKeyword
	t.Run("create", func(t *testing.T) {
		monitorID, err = client.CreateMonitor(ctx, httpKeywordMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test HTTP Keyword Monitor", retrievedMonitor.Name)

		err = client.GetMonitorAs(ctx, monitorID, &httpKeywordMonitorRetrieved)
		require.NoError(t, err)
		httpKeywordMonitor.ID = monitorID
		httpKeywordMonitor.PathName = httpKeywordMonitor.Name
		require.EqualExportedValues(t, httpKeywordMonitor, httpKeywordMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		httpKeywordMonitorRetrieved.Name = "Updated HTTP Keyword Monitor"
		httpKeywordMonitorRetrieved.Keyword = "Moby Dick"
		httpKeywordMonitorRetrieved.InvertKeyword = true
		err := client.UpdateMonitor(ctx, httpKeywordMonitorRetrieved)
		require.NoError(t, err)

		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated HTTP Keyword Monitor", updatedMonitor.Name)

		var updatedHTTPKeyword monitor.HTTPKeyword
		err = client.GetMonitorAs(ctx, monitorID, &updatedHTTPKeyword)
		require.NoError(t, err)
		require.Equal(t, "Moby Dick", updatedHTTPKeyword.Keyword)
		require.Equal(t, true, updatedHTTPKeyword.InvertKeyword)
	})

	t.Run("pause", func(t *testing.T) {
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorDNSCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test DNS monitor
	dnsMonitor := monitor.DNS{
		Base: monitor.Base{
			Name:           "Test DNS Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		DNSDetails: monitor.DNSDetails{
			Hostname:       "example.com",
			ResolverServer: "1.1.1.1",
			ResolveType:    monitor.DNSResolveTypeA,
			Port:           53,
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var dnsMonitorRetrieved monitor.DNS
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, dnsMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test DNS Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &dnsMonitorRetrieved)
		require.NoError(t, err)
		dnsMonitor.ID = monitorID
		dnsMonitor.PathName = dnsMonitor.Name
		require.EqualExportedValues(t, dnsMonitor, dnsMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		dnsMonitorRetrieved.Name = "Updated DNS Monitor"
		dnsMonitorRetrieved.Hostname = "cloudflare.com"
		dnsMonitorRetrieved.ResolveType = monitor.DNSResolveTypeAAAA
		err := client.UpdateMonitor(ctx, dnsMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated DNS Monitor", updatedMonitor.Name)

		var updatedDNS monitor.DNS
		err = client.GetMonitorAs(ctx, monitorID, &updatedDNS)
		require.NoError(t, err)
		require.Equal(t, "cloudflare.com", updatedDNS.Hostname)
		require.Equal(t, monitor.DNSResolveTypeAAAA, updatedDNS.ResolveType)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorHTTPJSONQueryCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test HTTP JSON Query monitor
	jsonQueryMonitor := monitor.HTTPJSONQuery{
		Base: monitor.Base{
			Name:           "Test JSON Query Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		HTTPDetails: monitor.HTTPDetails{
			URL:                 "https://httpbin.org/json",
			Timeout:             48,
			Method:              "GET",
			ExpiryNotification:  false,
			IgnoreTLS:           false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
			AuthMethod:          monitor.AuthMethodNone,
		},
		HTTPJSONQueryDetails: monitor.HTTPJSONQueryDetails{
			JSONPath:      "slideshow.title",
			ExpectedValue: "Sample Slide Show",
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var jsonQueryMonitorRetrieved monitor.HTTPJSONQuery
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, jsonQueryMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test JSON Query Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &jsonQueryMonitorRetrieved)
		require.NoError(t, err)
		jsonQueryMonitor.ID = monitorID
		jsonQueryMonitor.PathName = jsonQueryMonitor.Name
		require.EqualExportedValues(t, jsonQueryMonitor, jsonQueryMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		jsonQueryMonitorRetrieved.Name = "Updated JSON Query Monitor"
		jsonQueryMonitorRetrieved.JSONPath = "slideshow.author"
		jsonQueryMonitorRetrieved.ExpectedValue = "Yours Truly"
		err := client.UpdateMonitor(ctx, jsonQueryMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated JSON Query Monitor", updatedMonitor.Name)

		var updatedJSONQuery monitor.HTTPJSONQuery
		err = client.GetMonitorAs(ctx, monitorID, &updatedJSONQuery)
		require.NoError(t, err)
		require.Equal(t, "slideshow.author", updatedJSONQuery.JSONPath)
		require.Equal(t, "Yours Truly", updatedJSONQuery.ExpectedValue)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorRedisCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test Redis monitor
	redisMonitor := monitor.Redis{
		Base: monitor.Base{
			Name:           "Test Redis Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       true,
		},
		RedisDetails: monitor.RedisDetails{
			ConnectionString: "redis://localhost:6379",
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var redisMonitorRetrieved monitor.Redis
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, redisMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test Redis Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &redisMonitorRetrieved)
		require.NoError(t, err)
		redisMonitor.ID = monitorID
		redisMonitor.PathName = redisMonitor.Name
		require.EqualExportedValues(t, redisMonitor, redisMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		redisMonitorRetrieved.Name = "Updated Redis Monitor"
		redisMonitorRetrieved.ConnectionString = "redis://user:password@localhost:6380"
		err := client.UpdateMonitor(ctx, redisMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated Redis Monitor", updatedMonitor.Name)

		var updatedRedis monitor.Redis
		err = client.GetMonitorAs(ctx, monitorID, &updatedRedis)
		require.NoError(t, err)
		require.Equal(t, "redis://user:password@localhost:6380", updatedRedis.ConnectionString)
	})

	t.Run("pause", func(t *testing.T) {
		// Test PauseMonitor
		err := client.PauseMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("resume", func(t *testing.T) {
		// Test ResumeMonitor
		err := client.ResumeMonitor(ctx, monitorID)
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}

func TestClient_MonitorGrpcKeywordCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	// Create a test gRPC Keyword monitor
	grpcKeywordMonitor := monitor.GrpcKeyword{
		Base: monitor.Base{
			Name:           "Test gRPC Keyword Monitor",
			Interval:       60,
			RetryInterval:  60,
			ResendInterval: 0,
			MaxRetries:     3,
			UpsideDown:     false,
			IsActive:       false, // Set to false to avoid actual gRPC calls during testing
		},
		GrpcKeywordDetails: monitor.GrpcKeywordDetails{
			GrpcURL:             "localhost:50051",
			GrpcProtobuf:        "syntax = \"proto3\";\n\npackage grpc.health.v1;\n\nservice Health {\n  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);\n}\n\nmessage HealthCheckRequest {\n  string service = 1;\n}\n\nmessage HealthCheckResponse {\n  enum ServingStatus {\n    UNKNOWN = 0;\n    SERVING = 1;\n    NOT_SERVING = 2;\n  }\n  ServingStatus status = 1;\n}\n",
			GrpcServiceName:     "Health",
			GrpcMethod:          "Check",
			GrpcEnableTLS:       false,
			GrpcBody:            "{\"service\":\"\"}",
			Keyword:             "SERVING",
			InvertKeyword:       false,
			MaxRedirects:        10,
			AcceptedStatusCodes: []string{"200-299"},
		},
	}

	var initialCount int
	t.Run("initial_state", func(t *testing.T) {
		// Test GetMonitors (should work even if empty)
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		initialCount = len(monitors)
	})

	var monitorID int64
	var grpcKeywordMonitorRetrieved monitor.GrpcKeyword
	t.Run("create", func(t *testing.T) {
		// Test CreateMonitor
		monitorID, err = client.CreateMonitor(ctx, grpcKeywordMonitor)
		require.NoError(t, err)
		require.Greater(t, monitorID, int64(0))

		// Test GetMonitors after creation
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, len(monitors))

		// Test GetMonitor
		retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, monitorID, retrievedMonitor.ID)
		require.Equal(t, "Test gRPC Keyword Monitor", retrievedMonitor.Name)

		// Test GetMonitorAs
		err = client.GetMonitorAs(ctx, monitorID, &grpcKeywordMonitorRetrieved)
		require.NoError(t, err)
		grpcKeywordMonitor.ID = monitorID
		grpcKeywordMonitor.PathName = grpcKeywordMonitor.Name
		require.EqualExportedValues(t, grpcKeywordMonitor, grpcKeywordMonitorRetrieved)
	})

	t.Run("update", func(t *testing.T) {
		// Test UpdateMonitor
		grpcKeywordMonitorRetrieved.Name = "Updated gRPC Keyword Monitor"
		grpcKeywordMonitorRetrieved.GrpcURL = "example.com:443"
		grpcKeywordMonitorRetrieved.Keyword = "NOT_SERVING"
		grpcKeywordMonitorRetrieved.InvertKeyword = true
		grpcKeywordMonitorRetrieved.GrpcEnableTLS = true
		err := client.UpdateMonitor(ctx, grpcKeywordMonitorRetrieved)
		require.NoError(t, err)

		// Verify update
		updatedMonitor, err := client.GetMonitor(ctx, monitorID)
		require.NoError(t, err)
		require.Equal(t, "Updated gRPC Keyword Monitor", updatedMonitor.Name)

		var updatedGrpcKeyword monitor.GrpcKeyword
		err = client.GetMonitorAs(ctx, monitorID, &updatedGrpcKeyword)
		require.NoError(t, err)
		require.Equal(t, "example.com:443", updatedGrpcKeyword.GrpcURL)
		require.Equal(t, "NOT_SERVING", updatedGrpcKeyword.Keyword)
		require.Equal(t, true, updatedGrpcKeyword.InvertKeyword)
		require.Equal(t, true, updatedGrpcKeyword.GrpcEnableTLS)
	})

	t.Run("delete", func(t *testing.T) {
		// Test DeleteMonitor
		err := client.DeleteMonitor(ctx, monitorID)
		require.NoError(t, err)

		// Verify deletion
		monitors, err := client.GetMonitors(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount, len(monitors))
	})
}
