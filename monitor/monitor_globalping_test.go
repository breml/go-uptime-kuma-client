package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorGlobalping_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.Globalping
		wantJSON string
	}{
		{
			name: "ping",
			data: []byte(
				`{"id":10,"name":"globalping-ping","description":"Globalping ping","pathName":"globalping-ping","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"globalping","subtype":"ping","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"location":"world","ipFamily":null,"protocol":"ICMP","ping_count":3,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":"","oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Globalping{
				Base: monitor.Base{
					ID:              10,
					Name:            "globalping-ping",
					Description:     ptr.To("Globalping ping"),
					PathName:        "globalping-ping",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1},
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					Timeout:             48,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					AuthMethod:          monitor.AuthMethodNone,
					OAuthAuthMethod:     "client_secret_basic",
				},
				GlobalpingDetails: monitor.GlobalpingDetails{
					Subtype:          monitor.GlobalpingSubtypePing,
					Location:         "world",
					IPFamily:         monitor.GlobalpingIPFamilyAuto,
					Protocol:         "ICMP",
					PingCount:        3,
					Hostname:         "example.com",
					DNSResolveType:   monitor.DNSResolveTypeA,
					DNSResolveServer: "1.1.1.1",
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","cacheBust":false,"conditions":[],"description":"Globalping ping","dns_resolve_server":"1.1.1.1","dns_resolve_type":"A","expectedValue":"","expiryNotification":false,"headers":"","hostname":"example.com","httpBodyEncoding":"json","id":10,"ignoreTls":false,"interval":60,"invertKeyword":false,"ipFamily":"","jsonPath":"","jsonPathOperator":"","keyword":"","location":"world","maxredirects":10,"maxretries":2,"method":"GET","name":"globalping-ping","notificationIDList":{"1":true},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":1,"ping_count":3,"port":0,"protocol":"ICMP","proxyId":null,"resendInterval":0,"retryInterval":60,"subtype":"ping","timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"globalping","upsideDown":false,"url":""}`,
		},
		{
			name: "ping_tcp_with_port",
			data: []byte(
				`{"id":11,"name":"globalping-ping-tcp","description":null,"pathName":"globalping-ping-tcp","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":443,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"globalping","subtype":"ping","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"location":"europe","ipFamily":"ipv4","protocol":"TCP","ping_count":5,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":"","oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Globalping{
				Base: monitor.Base{
					ID:             11,
					Name:           "globalping-ping-tcp",
					Description:    nil,
					PathName:       "globalping-ping-tcp",
					Parent:         nil,
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				HTTPDetails: monitor.HTTPDetails{
					Timeout:             48,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					AuthMethod:          monitor.AuthMethodNone,
					OAuthAuthMethod:     "client_secret_basic",
				},
				GlobalpingDetails: monitor.GlobalpingDetails{
					Subtype:          monitor.GlobalpingSubtypePing,
					Location:         "europe",
					IPFamily:         monitor.GlobalpingIPFamilyIPv4,
					Protocol:         "TCP",
					PingCount:        5,
					Hostname:         "example.com",
					Port:             443,
					DNSResolveType:   monitor.DNSResolveTypeA,
					DNSResolveServer: "1.1.1.1",
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","cacheBust":false,"conditions":[],"description":null,"dns_resolve_server":"1.1.1.1","dns_resolve_type":"A","expectedValue":"","expiryNotification":false,"headers":"","hostname":"example.com","httpBodyEncoding":"json","id":11,"ignoreTls":false,"interval":60,"invertKeyword":false,"ipFamily":"ipv4","jsonPath":"","jsonPathOperator":"","keyword":"","location":"europe","maxredirects":10,"maxretries":3,"method":"GET","name":"globalping-ping-tcp","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"ping_count":5,"port":443,"protocol":"TCP","proxyId":null,"resendInterval":0,"retryInterval":60,"subtype":"ping","timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"globalping","upsideDown":false,"url":""}`,
		},
		{
			name: "http",
			data: []byte(
				`{"id":12,"name":"globalping-http","description":"Globalping HTTP","pathName":"globalping-http","parent":null,"childrenIDs":[],"url":"https://www.example.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"globalping","subtype":"http","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":"success","invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"location":"world","ipFamily":"ipv6","protocol":"HTTP2","ping_count":0,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":"basic","grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":"X-Custom: test","body":"","grpcBody":null,"grpcMetadata":null,"basic_auth_user":"u","basic_auth_pass":"p","oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":"","oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Globalping{
				Base: monitor.Base{
					ID:             12,
					Name:           "globalping-http",
					Description:    ptr.To("Globalping HTTP"),
					PathName:       "globalping-http",
					Parent:         nil,
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "https://www.example.com",
					Timeout:             48,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Headers:             "X-Custom: test",
					AuthMethod:          monitor.AuthMethodBasic,
					BasicAuthUser:       "u",
					BasicAuthPass:       "p",
					OAuthAuthMethod:     "client_secret_basic",
				},
				GlobalpingDetails: monitor.GlobalpingDetails{
					Subtype:          monitor.GlobalpingSubtypeHTTP,
					Location:         "world",
					IPFamily:         monitor.GlobalpingIPFamilyIPv6,
					Protocol:         "HTTP2",
					Keyword:          "success",
					DNSResolveType:   monitor.DNSResolveTypeA,
					DNSResolveServer: "1.1.1.1",
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"basic","authWorkstation":"","basic_auth_pass":"p","basic_auth_user":"u","body":"","cacheBust":false,"conditions":[],"description":"Globalping HTTP","dns_resolve_server":"1.1.1.1","dns_resolve_type":"A","expectedValue":"","expiryNotification":false,"headers":"X-Custom: test","hostname":"","httpBodyEncoding":"json","id":12,"ignoreTls":false,"interval":60,"invertKeyword":false,"ipFamily":"ipv6","jsonPath":"","jsonPathOperator":"","keyword":"success","location":"world","maxredirects":10,"maxretries":2,"method":"GET","name":"globalping-http","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"ping_count":0,"port":0,"protocol":"HTTP2","proxyId":null,"resendInterval":0,"retryInterval":60,"subtype":"http","timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"globalping","upsideDown":false,"url":"https://www.example.com"}`,
		},
		{
			name: "dns",
			data: []byte(
				`{"id":13,"name":"globalping-dns","description":null,"pathName":"globalping-dns","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":53,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"globalping","subtype":"dns","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":"127\\.0\\.0\\.1","invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"AAAA","dns_resolve_server":"8.8.8.8","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"location":"asia","ipFamily":"ipv4","protocol":"UDP","ping_count":0,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":"","oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Globalping{
				Base: monitor.Base{
					ID:             13,
					Name:           "globalping-dns",
					Description:    nil,
					PathName:       "globalping-dns",
					Parent:         nil,
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				HTTPDetails: monitor.HTTPDetails{
					Timeout:             48,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					AuthMethod:          monitor.AuthMethodNone,
					OAuthAuthMethod:     "client_secret_basic",
				},
				GlobalpingDetails: monitor.GlobalpingDetails{
					Subtype:          monitor.GlobalpingSubtypeDNS,
					Location:         "asia",
					IPFamily:         monitor.GlobalpingIPFamilyIPv4,
					Protocol:         "UDP",
					Hostname:         "example.com",
					Port:             53,
					DNSResolveType:   monitor.DNSResolveTypeAAAA,
					DNSResolveServer: "8.8.8.8",
					Keyword:          `127\.0\.0\.1`,
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","cacheBust":false,"conditions":[],"description":null,"dns_resolve_server":"8.8.8.8","dns_resolve_type":"AAAA","expectedValue":"","expiryNotification":false,"headers":"","hostname":"example.com","httpBodyEncoding":"json","id":13,"ignoreTls":false,"interval":60,"invertKeyword":false,"ipFamily":"ipv4","jsonPath":"","jsonPathOperator":"","keyword":"127\\.0\\.0\\.1","location":"asia","maxredirects":10,"maxretries":2,"method":"GET","name":"globalping-dns","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"ping_count":0,"port":53,"protocol":"UDP","proxyId":null,"resendInterval":0,"retryInterval":60,"subtype":"dns","timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"globalping","upsideDown":false,"url":""}`,
		},
		{
			name: "traceroute",
			data: []byte(
				`{"id":14,"name":"globalping-traceroute","description":null,"pathName":"globalping-traceroute","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"globalping","subtype":"traceroute","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"location":"north-america","ipFamily":null,"protocol":"ICMP","ping_count":0,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":"","oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Globalping{
				Base: monitor.Base{
					ID:             14,
					Name:           "globalping-traceroute",
					Description:    nil,
					PathName:       "globalping-traceroute",
					Parent:         nil,
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				HTTPDetails: monitor.HTTPDetails{
					Timeout:             48,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					AuthMethod:          monitor.AuthMethodNone,
					OAuthAuthMethod:     "client_secret_basic",
				},
				GlobalpingDetails: monitor.GlobalpingDetails{
					Subtype:          monitor.GlobalpingSubtypeTraceroute,
					Location:         "north-america",
					IPFamily:         monitor.GlobalpingIPFamilyAuto,
					Protocol:         "ICMP",
					Hostname:         "example.com",
					DNSResolveType:   monitor.DNSResolveTypeA,
					DNSResolveServer: "1.1.1.1",
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","cacheBust":false,"conditions":[],"description":null,"dns_resolve_server":"1.1.1.1","dns_resolve_type":"A","expectedValue":"","expiryNotification":false,"headers":"","hostname":"example.com","httpBodyEncoding":"json","id":14,"ignoreTls":false,"interval":60,"invertKeyword":false,"ipFamily":"","jsonPath":"","jsonPathOperator":"","keyword":"","location":"north-america","maxredirects":10,"maxretries":2,"method":"GET","name":"globalping-traceroute","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"ping_count":0,"port":0,"protocol":"ICMP","proxyId":null,"resendInterval":0,"retryInterval":60,"subtype":"traceroute","timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"globalping","upsideDown":false,"url":""}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			globalpingMonitor := monitor.Globalping{}

			err := json.Unmarshal(tc.data, &globalpingMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, globalpingMonitor)

			data, err := json.Marshal(globalpingMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestGlobalpingSubtype_String(t *testing.T) {
	tests := []struct {
		subtype monitor.GlobalpingSubtype
		want    string
	}{
		{monitor.GlobalpingSubtypePing, "ping"},
		{monitor.GlobalpingSubtypeTraceroute, "traceroute"},
		{monitor.GlobalpingSubtypeDNS, "dns"},
		{monitor.GlobalpingSubtypeHTTP, "http"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			require.Equal(t, tc.want, tc.subtype.String())
		})
	}
}

func TestGlobalpingIPFamily_String(t *testing.T) {
	tests := []struct {
		family monitor.GlobalpingIPFamily
		want   string
	}{
		{monitor.GlobalpingIPFamilyAuto, ""},
		{monitor.GlobalpingIPFamilyIPv4, "ipv4"},
		{monitor.GlobalpingIPFamilyIPv6, "ipv6"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			require.Equal(t, tc.want, tc.family.String())
		})
	}
}
