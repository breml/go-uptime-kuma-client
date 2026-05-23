package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorDNS_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.DNS
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":5,"name":"dns-monitor","description":"Test DNS monitor","pathName":"group / dns-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":53,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"dns","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.DNS{
				Base: monitor.Base{
					ID:              5,
					Name:            "dns-monitor",
					Description:     ptr.To("Test DNS monitor"),
					PathName:        "group / dns-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				DNSDetails: monitor.DNSDetails{
					Hostname:       "example.com",
					ResolverServer: "1.1.1.1",
					ResolveType:    monitor.DNSResolveTypeA,
					Port:           53,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test DNS monitor","dns_resolve_server":"1.1.1.1","dns_resolve_type":"A","hostname":"example.com","id":5,"interval":60,"maxretries":2,"name":"dns-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":53,"resendInterval":0,"retryInterval":60,"type":"dns","upsideDown":false}`,
		},
		{
			name: "with conditions",
			data: []byte(
				`{"id":7,"name":"dns-cname","description":"Test CNAME DNS monitor","pathName":"dns-cname","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":53,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"dns","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"CNAME","dns_resolve_server":"9.9.9.9","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[{"type":"expression","andOr":"and","variable":"record","operator":"equals","value":"target.example.com."}],"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.DNS{
				Base: monitor.Base{
					ID:             7,
					Name:           "dns-cname",
					Description:    ptr.To("Test CNAME DNS monitor"),
					PathName:       "dns-cname",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				DNSDetails: monitor.DNSDetails{
					Hostname:       "example.com",
					ResolverServer: "9.9.9.9",
					ResolveType:    monitor.DNSResolveTypeCNAME,
					Port:           53,
					Conditions: []monitor.Condition{
						{
							Variable: "record",
							Operator: "equals",
							Value:    "target.example.com.",
							AndOr:    monitor.ConditionAnd,
						},
					},
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[{"type":"expression","variable":"record","operator":"equals","value":"target.example.com.","andOr":"and"}],"description":"Test CNAME DNS monitor","dns_resolve_server":"9.9.9.9","dns_resolve_type":"CNAME","hostname":"example.com","id":7,"interval":60,"maxretries":2,"name":"dns-cname","notificationIDList":{},"parent":null,"port":53,"resendInterval":0,"retryInterval":60,"type":"dns","upsideDown":false}`,
		},
		{
			name: "multi resolver",
			data: []byte(
				`{"id":6,"name":"dns-multi-resolver","description":"Test multi-resolver DNS monitor","pathName":"group / dns-multi-resolver","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":53,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"dns","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1,8.8.8.8","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.DNS{
				Base: monitor.Base{
					ID:              6,
					Name:            "dns-multi-resolver",
					Description:     ptr.To("Test multi-resolver DNS monitor"),
					PathName:        "group / dns-multi-resolver",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				DNSDetails: monitor.DNSDetails{
					Hostname:       "example.com",
					ResolverServer: "1.1.1.1,8.8.8.8",
					ResolveType:    monitor.DNSResolveTypeA,
					Port:           53,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test multi-resolver DNS monitor","dns_resolve_server":"1.1.1.1,8.8.8.8","dns_resolve_type":"A","hostname":"example.com","id":6,"interval":60,"maxretries":2,"name":"dns-multi-resolver","notificationIDList":{"1":true,"2":true},"parent":1,"port":53,"resendInterval":0,"retryInterval":60,"type":"dns","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dnsMonitor := monitor.DNS{}

			err := json.Unmarshal(tc.data, &dnsMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, dnsMonitor)

			data, err := json.Marshal(dnsMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestDNSDetails_ResolverServers(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  []string
	}{
		{
			name:  "empty",
			field: "",
			want:  nil,
		},
		{
			name:  "single",
			field: "1.1.1.1",
			want:  []string{"1.1.1.1"},
		},
		{
			name:  "multiple",
			field: "1.1.1.1,8.8.8.8",
			want:  []string{"1.1.1.1", "8.8.8.8"},
		},
		{
			name:  "trim and skip empty",
			field: " 1.1.1.1 , , 8.8.8.8 ,",
			want:  []string{"1.1.1.1", "8.8.8.8"},
		},
		{
			name:  "only separators",
			field: " , , ",
			want:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := monitor.DNSDetails{ResolverServer: tc.field}
			require.Equal(t, tc.want, d.ResolverServers())
		})
	}
}

func TestDNSDetails_SetResolverServers(t *testing.T) {
	tests := []struct {
		name    string
		servers []string
		want    string
	}{
		{
			name:    "nil",
			servers: nil,
			want:    "",
		},
		{
			name:    "single",
			servers: []string{"1.1.1.1"},
			want:    "1.1.1.1",
		},
		{
			name:    "multiple",
			servers: []string{"1.1.1.1", "8.8.8.8"},
			want:    "1.1.1.1,8.8.8.8",
		},
		{
			name:    "trim and skip empty",
			servers: []string{" 1.1.1.1 ", "", " ", " 8.8.8.8"},
			want:    "1.1.1.1,8.8.8.8",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := monitor.DNSDetails{}
			d.SetResolverServers(tc.servers)
			require.Equal(t, tc.want, d.ResolverServer)
		})
	}
}
