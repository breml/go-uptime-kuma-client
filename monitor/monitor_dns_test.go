package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

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
			data: []byte(`{"id":5,"name":"dns-monitor","description":"Test DNS monitor","pathName":"group / dns-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":53,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"dns","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.DNS{
				Base: monitor.Base{
					ID:              5,
					Name:            "dns-monitor",
					Description:     stringPtr("Test DNS monitor"),
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
