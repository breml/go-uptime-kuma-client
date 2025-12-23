package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

// monitorTestCase defines a single monitor type's CRUD test scenario
type monitorTestCase struct {
	name              string                                                   // Test name (e.g., "HTTP", "Ping")
	create            monitor.Monitor                                          // Monitor to create
	updateFunc        func(monitor.Monitor)                                    // Function to modify monitor for update test
	verifyCreatedFunc func(t *testing.T, actual monitor.Monitor, id int64)     // Function to verify created monitor
	createTypedFunc   func(t *testing.T, base monitor.Monitor) monitor.Monitor // Function to create typed monitor
	verifyUpdatedFunc func(t *testing.T, actual monitor.Monitor)               // Function to verify updated monitor
	testPauseResume   bool                                                     // Whether to test pause/resume functionality
}

func TestMonitorCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	testCases := []monitorTestCase{
		{
			name: "HTTP",
			create: monitor.HTTP{
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
			},
			updateFunc: func(m monitor.Monitor) {
				http := m.(*monitor.HTTP)
				http.Name = "Updated HTTP Monitor"
				http.URL = "https://httpbin.org/status/201"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var http monitor.HTTP
				err := actual.As(&http)
				require.NoError(t, err)
				require.Equal(t, id, http.ID)
				require.Equal(t, "Test HTTP Monitor", http.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var http monitor.HTTP
				err := base.As(&http)
				require.NoError(t, err)
				return &http
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var http monitor.HTTP
				err := actual.As(&http)
				require.NoError(t, err)
				require.Equal(t, "Updated HTTP Monitor", http.Name)
				require.Equal(t, "https://httpbin.org/status/201", http.URL)
			},
			testPauseResume: true,
		},
		{
			name: "Monitor Group",
			create: monitor.Group{
				Base: monitor.Base{
					Name:           "Test Monitor Group",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
			},
			updateFunc: func(m monitor.Monitor) {
				group := m.(*monitor.Group)
				group.Name = "Updated Monitor Group"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var group monitor.Group
				err := actual.As(&group)
				require.NoError(t, err)
				require.Equal(t, id, group.ID)
				require.Equal(t, "Test Monitor Group", group.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var group monitor.Group
				err := base.As(&group)
				require.NoError(t, err)
				return &group
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var group monitor.Group
				err := actual.As(&group)
				require.NoError(t, err)
				require.Equal(t, "Updated Monitor Group", group.Name)
			},
			testPauseResume: true,
		},
		{
			name: "Ping",
			create: monitor.Ping{
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
			},
			updateFunc: func(m monitor.Monitor) {
				ping := m.(*monitor.Ping)
				ping.Name = "Updated Ping Monitor"
				ping.Hostname = "1.1.1.1"
				ping.PacketSize = 64
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var ping monitor.Ping
				err := actual.As(&ping)
				require.NoError(t, err)
				require.Equal(t, id, ping.ID)
				require.Equal(t, "Test Ping Monitor", ping.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var ping monitor.Ping
				err := base.As(&ping)
				require.NoError(t, err)
				return &ping
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var ping monitor.Ping
				err := actual.As(&ping)
				require.NoError(t, err)
				require.Equal(t, "Updated Ping Monitor", ping.Name)
				require.Equal(t, "1.1.1.1", ping.Hostname)
				require.Equal(t, 64, ping.PacketSize)
			},
			testPauseResume: true,
		},
		{
			name: "Push",
			create: monitor.Push{
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
			},
			updateFunc: func(m monitor.Monitor) {
				push := m.(*monitor.Push)
				push.Name = "Updated Push Monitor"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var push monitor.Push
				err := actual.As(&push)
				require.NoError(t, err)
				require.Equal(t, id, push.ID)
				require.Equal(t, "Test Push Monitor", push.Name)
				require.NotEmpty(t, push.PushToken)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var push monitor.Push
				err := base.As(&push)
				require.NoError(t, err)
				return &push
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var push monitor.Push
				err := actual.As(&push)
				require.NoError(t, err)
				require.Equal(t, "Updated Push Monitor", push.Name)
			},
			testPauseResume: false,
		},
		{
			name: "TCP Port",
			create: monitor.TCPPort{
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
			},
			updateFunc: func(m monitor.Monitor) {
				tcp := m.(*monitor.TCPPort)
				tcp.Name = "Updated TCP Port Monitor"
				tcp.Hostname = "cloudflare.com"
				tcp.Port = 80
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var tcp monitor.TCPPort
				err := actual.As(&tcp)
				require.NoError(t, err)
				require.Equal(t, id, tcp.ID)
				require.Equal(t, "Test TCP Port Monitor", tcp.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var tcp monitor.TCPPort
				err := base.As(&tcp)
				require.NoError(t, err)
				return &tcp
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var tcp monitor.TCPPort
				err := actual.As(&tcp)
				require.NoError(t, err)
				require.Equal(t, "Updated TCP Port Monitor", tcp.Name)
				require.Equal(t, "cloudflare.com", tcp.Hostname)
				require.Equal(t, 80, tcp.Port)
			},
			testPauseResume: true,
		},
		{
			name: "HTTP Keyword",
			create: monitor.HTTPKeyword{
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
			},
			updateFunc: func(m monitor.Monitor) {
				http := m.(*monitor.HTTPKeyword)
				http.Name = "Updated HTTP Keyword Monitor"
				http.Keyword = "Moby Dick"
				http.InvertKeyword = true
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var http monitor.HTTPKeyword
				err := actual.As(&http)
				require.NoError(t, err)
				require.Equal(t, id, http.ID)
				require.Equal(t, "Test HTTP Keyword Monitor", http.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var http monitor.HTTPKeyword
				err := base.As(&http)
				require.NoError(t, err)
				return &http
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var http monitor.HTTPKeyword
				err := actual.As(&http)
				require.NoError(t, err)
				require.Equal(t, "Updated HTTP Keyword Monitor", http.Name)
				require.Equal(t, "Moby Dick", http.Keyword)
				require.Equal(t, true, http.InvertKeyword)
			},
			testPauseResume: true,
		},
		{
			name: "DNS",
			create: monitor.DNS{
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
			},
			updateFunc: func(m monitor.Monitor) {
				dns := m.(*monitor.DNS)
				dns.Name = "Updated DNS Monitor"
				dns.Hostname = "google.com"
				dns.ResolverServer = "8.8.8.8"
				dns.ResolveType = monitor.DNSResolveTypeAAAA
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var dns monitor.DNS
				err := actual.As(&dns)
				require.NoError(t, err)
				require.Equal(t, id, dns.ID)
				require.Equal(t, "Test DNS Monitor", dns.Name)
				require.Equal(t, monitor.DNSResolveTypeA, dns.ResolveType)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var dns monitor.DNS
				err := base.As(&dns)
				require.NoError(t, err)
				return &dns
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var dns monitor.DNS
				err := actual.As(&dns)
				require.NoError(t, err)
				require.Equal(t, "Updated DNS Monitor", dns.Name)
				require.Equal(t, "google.com", dns.Hostname)
				require.Equal(t, "8.8.8.8", dns.ResolverServer)
				require.Equal(t, monitor.DNSResolveTypeAAAA, dns.ResolveType)
			},
			testPauseResume: true,
		},
		{
			name: "HTTP JSON Query",
			create: monitor.HTTPJSONQuery{
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
			},
			updateFunc: func(m monitor.Monitor) {
				json := m.(*monitor.HTTPJSONQuery)
				json.Name = "Updated JSON Query Monitor"
				json.JSONPath = "slideshow.author"
				json.ExpectedValue = "Yours Truly"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var json monitor.HTTPJSONQuery
				err := actual.As(&json)
				require.NoError(t, err)
				require.Equal(t, id, json.ID)
				require.Equal(t, "Test JSON Query Monitor", json.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var json monitor.HTTPJSONQuery
				err := base.As(&json)
				require.NoError(t, err)
				return &json
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var json monitor.HTTPJSONQuery
				err := actual.As(&json)
				require.NoError(t, err)
				require.Equal(t, "Updated JSON Query Monitor", json.Name)
				require.Equal(t, "slideshow.author", json.JSONPath)
				require.Equal(t, "Yours Truly", json.ExpectedValue)
			},
			testPauseResume: true,
		},
		{
			name: "Postgres",
			create: monitor.Postgres{
				Base: monitor.Base{
					Name:           "Test Postgres Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       false,
				},
				PostgresDetails: monitor.PostgresDetails{
					DatabaseConnectionString: "postgres://testuser:testpass@localhost:5432/testdb",
					DatabaseQuery:            "SELECT 1",
				},
			},
			updateFunc: func(m monitor.Monitor) {
				postgres := m.(*monitor.Postgres)
				postgres.Name = "Updated Postgres Monitor"
				postgres.DatabaseConnectionString = "postgres://newuser:newpass@localhost:5432/newdb"
				postgres.DatabaseQuery = "SELECT version()"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var postgres monitor.Postgres
				err := actual.As(&postgres)
				require.NoError(t, err)
				require.Equal(t, id, postgres.ID)
				require.Equal(t, "Test Postgres Monitor", postgres.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var postgres monitor.Postgres
				err := base.As(&postgres)
				require.NoError(t, err)
				return &postgres
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var postgres monitor.Postgres
				err := actual.As(&postgres)
				require.NoError(t, err)
				require.Equal(t, "Updated Postgres Monitor", postgres.Name)
				require.Equal(t, "postgres://newuser:newpass@localhost:5432/newdb", postgres.DatabaseConnectionString)
				require.Equal(t, "SELECT version()", postgres.DatabaseQuery)
			},
			testPauseResume: false,
		},
		{
			name: "Real Browser",
			create: monitor.RealBrowser{
				Base: monitor.Base{
					Name:           "Test RealBrowser Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				RealBrowserDetails: monitor.RealBrowserDetails{
					URL:                 "https://httpbin.org/status/200",
					Timeout:             48,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
				},
			},
			updateFunc: func(m monitor.Monitor) {
				browser := m.(*monitor.RealBrowser)
				browser.Name = "Updated RealBrowser Monitor"
				browser.URL = "https://httpbin.org/status/201"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var browser monitor.RealBrowser
				err := actual.As(&browser)
				require.NoError(t, err)
				require.Equal(t, id, browser.ID)
				require.Equal(t, "Test RealBrowser Monitor", browser.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var browser monitor.RealBrowser
				err := base.As(&browser)
				require.NoError(t, err)
				return &browser
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var browser monitor.RealBrowser
				err := actual.As(&browser)
				require.NoError(t, err)
				require.Equal(t, "Updated RealBrowser Monitor", browser.Name)
				require.Equal(t, "https://httpbin.org/status/201", browser.URL)
			},
			testPauseResume: true,
		},
		{
			name: "Redis",
			create: monitor.Redis{
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
			},
			updateFunc: func(m monitor.Monitor) {
				redis := m.(*monitor.Redis)
				redis.Name = "Updated Redis Monitor"
				redis.ConnectionString = "redis://user:password@localhost:6380"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var redis monitor.Redis
				err := actual.As(&redis)
				require.NoError(t, err)
				require.Equal(t, id, redis.ID)
				require.Equal(t, "Test Redis Monitor", redis.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var redis monitor.Redis
				err := base.As(&redis)
				require.NoError(t, err)
				return &redis
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var redis monitor.Redis
				err := actual.As(&redis)
				require.NoError(t, err)
				require.Equal(t, "Updated Redis Monitor", redis.Name)
				require.Equal(t, "redis://user:password@localhost:6380", redis.ConnectionString)
			},
			testPauseResume: true,
		},
		{
			name: "SMTP",
			create: monitor.SMTP{
				Base: monitor.Base{
					Name:           "Test SMTP Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       false,
				},
				SMTPDetails: monitor.SMTPDetails{
					Hostname:     "mail.example.com",
					Port:         nil,
					SMTPSecurity: nil,
				},
			},
			updateFunc: func(m monitor.Monitor) {
				smtp := m.(*monitor.SMTP)
				smtp.Name = "Updated SMTP Monitor"
				port465 := int64(465)
				securitySecure := "secure"
				smtp.Hostname = "smtp.newserver.com"
				smtp.Port = &port465
				smtp.SMTPSecurity = &securitySecure
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var smtp monitor.SMTP
				err := actual.As(&smtp)
				require.NoError(t, err)
				require.Equal(t, id, smtp.ID)
				require.Equal(t, "Test SMTP Monitor", smtp.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var smtp monitor.SMTP
				err := base.As(&smtp)
				require.NoError(t, err)
				return &smtp
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var smtp monitor.SMTP
				err := actual.As(&smtp)
				require.NoError(t, err)
				require.Equal(t, "Updated SMTP Monitor", smtp.Name)
				require.Equal(t, "smtp.newserver.com", smtp.Hostname)
				require.Equal(t, int64(465), *smtp.Port)
				require.Equal(t, "secure", *smtp.SMTPSecurity)
			},
			testPauseResume: true,
		},
		{
			name: "gRPC Keyword",
			create: monitor.GrpcKeyword{
				Base: monitor.Base{
					Name:           "Test gRPC Keyword Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       false,
				},
				GrpcKeywordDetails: monitor.GrpcKeywordDetails{
					GrpcURL:         "localhost:50051",
					GrpcProtobuf:    "syntax = \"proto3\";\n\npackage grpc.health.v1;\n\nservice Health {\n  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);\n}\n\nmessage HealthCheckRequest {\n  string service = 1;\n}\n\nmessage HealthCheckResponse {\n  enum ServingStatus {\n    UNKNOWN = 0;\n    SERVING = 1;\n    NOT_SERVING = 2;\n  }\n  ServingStatus status = 1;\n}\n",
					GrpcServiceName: "Health",
					GrpcMethod:      "Check",
					GrpcEnableTLS:   false,
					GrpcBody:        "{\"service\":\"\"}",
					Keyword:         "SERVING",
					InvertKeyword:   false,
				},
			},
			updateFunc: func(m monitor.Monitor) {
				grpc := m.(*monitor.GrpcKeyword)
				grpc.Name = "Updated gRPC Keyword Monitor"
				grpc.GrpcURL = "example.com:443"
				grpc.Keyword = "NOT_SERVING"
				grpc.InvertKeyword = true
				grpc.GrpcEnableTLS = true
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var grpc monitor.GrpcKeyword
				err := actual.As(&grpc)
				require.NoError(t, err)
				require.Equal(t, id, grpc.ID)
				require.Equal(t, "Test gRPC Keyword Monitor", grpc.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var grpc monitor.GrpcKeyword
				err := base.As(&grpc)
				require.NoError(t, err)
				return &grpc
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var grpc monitor.GrpcKeyword
				err := actual.As(&grpc)
				require.NoError(t, err)
				require.Equal(t, "Updated gRPC Keyword Monitor", grpc.Name)
				require.Equal(t, "example.com:443", grpc.GrpcURL)
				require.Equal(t, "NOT_SERVING", grpc.Keyword)
				require.Equal(t, true, grpc.InvertKeyword)
				require.Equal(t, true, grpc.GrpcEnableTLS)
			},
			testPauseResume: false,
		},
		{
			name: "SNMP",
			create: monitor.SNMP{
				Base: monitor.Base{
					Name:           "Test SNMP Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				SNMPDetails: monitor.SNMPDetails{
					Hostname:      "192.168.1.1",
					Port:          ptr.To(int64(161)),
					SNMPVersion:   "2c",
					SNMPOID:       "1.3.6.1.2.1.1.3.0",
					SNMPCommunity: "public",
				},
			},
			updateFunc: func(m monitor.Monitor) {
				snmp := m.(*monitor.SNMP)
				snmp.Name = "Updated SNMP Monitor"
				snmp.Hostname = "10.0.0.1"
				snmp.SNMPOID = "1.3.6.1.2.1.2.2.1.5.1"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var snmp monitor.SNMP
				err := actual.As(&snmp)
				require.NoError(t, err)
				require.Equal(t, id, snmp.ID)
				require.Equal(t, "Test SNMP Monitor", snmp.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var snmp monitor.SNMP
				err := base.As(&snmp)
				require.NoError(t, err)
				return &snmp
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var snmp monitor.SNMP
				err := actual.As(&snmp)
				require.NoError(t, err)
				require.Equal(t, "Updated SNMP Monitor", snmp.Name)
				require.Equal(t, "10.0.0.1", snmp.Hostname)
				require.Equal(t, "1.3.6.1.2.1.2.2.1.5.1", snmp.SNMPOID)
			},
			testPauseResume: true,
		},
		{
			name: "Docker",
			create: monitor.Docker{
				Base: monitor.Base{
					Name:           "Test Docker Monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				DockerDetails: monitor.DockerDetails{
					DockerHost:      1,
					DockerContainer: "my-container",
				},
			},
			updateFunc: func(m monitor.Monitor) {
				docker := m.(*monitor.Docker)
				docker.Name = "Updated Docker Monitor"
				docker.DockerContainer = "updated-container"
			},
			verifyCreatedFunc: func(t *testing.T, actual monitor.Monitor, id int64) {
				t.Helper()
				var docker monitor.Docker
				err := actual.As(&docker)
				require.NoError(t, err)
				require.Equal(t, id, docker.ID)
				require.Equal(t, "Test Docker Monitor", docker.Name)
			},
			createTypedFunc: func(t *testing.T, base monitor.Monitor) monitor.Monitor {
				t.Helper()
				var docker monitor.Docker
				err := base.As(&docker)
				require.NoError(t, err)
				return &docker
			},
			verifyUpdatedFunc: func(t *testing.T, actual monitor.Monitor) {
				t.Helper()
				var docker monitor.Docker
				err := actual.As(&docker)
				require.NoError(t, err)
				require.Equal(t, "Updated Docker Monitor", docker.Name)
				require.Equal(t, "updated-container", docker.DockerContainer)
			},
			testPauseResume: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			var initialCount int
			monitors, err := client.GetMonitors(ctx)
			require.NoError(t, err)
			initialCount = len(monitors)

			// Create
			created := tc.create
			monitorID, err := client.CreateMonitor(ctx, created)
			require.NoError(t, err)
			require.Greater(t, monitorID, int64(0))

			// Verify count increased
			monitors, err = client.GetMonitors(ctx)
			require.NoError(t, err)
			require.Equal(t, initialCount+1, len(monitors))

			// Retrieve
			retrievedMonitor, err := client.GetMonitor(ctx, monitorID)
			require.NoError(t, err)
			require.Equal(t, monitorID, retrievedMonitor.ID)

			// Verify retrieved monitor
			retrieved := tc.createTypedFunc(t, retrievedMonitor)
			tc.verifyCreatedFunc(t, retrieved, monitorID)

			// Update
			tc.updateFunc(retrieved)
			err = client.UpdateMonitor(ctx, retrieved)
			require.NoError(t, err)

			// Verify update
			updated, err := client.GetMonitor(ctx, monitorID)
			require.NoError(t, err)
			updatedTyped := tc.createTypedFunc(t, updated)
			tc.verifyUpdatedFunc(t, updatedTyped)

			// Pause/Resume if supported
			if tc.testPauseResume {
				err = client.PauseMonitor(ctx, monitorID)
				require.NoError(t, err)

				err = client.ResumeMonitor(ctx, monitorID)
				require.NoError(t, err)
			}

			// Delete
			err = client.DeleteMonitor(ctx, monitorID)
			require.NoError(t, err)

			// Verify count restored
			monitors, err = client.GetMonitors(ctx)
			require.NoError(t, err)
			require.Equal(t, initialCount, len(monitors))
		})
	}
}
