package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSNMP_Unmarshal(t *testing.T) {
	port161 := int64(161)
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.SNMP
		wantJSON string
	}{
		{
			name: "success with all fields",
			data: []byte(`{"id":10,"name":"snmp-monitor","description":"Test SNMP monitor","pathName":"group / snmp-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"192.168.1.1","port":161,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"snmp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":"ifSpeed","jsonPathOperator":"==","expectedValue":"1000000000","kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":"public","radiusSecret":null,"snmpVersion":"2c","snmpOid":"1.3.6.1.2.1.2.2.1.5.1","mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.SNMP{
				Base: monitor.Base{
					ID:              10,
					Name:            "snmp-monitor",
					Description:     ptr.To("Test SNMP monitor"),
					PathName:        "group / snmp-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SNMPDetails: monitor.SNMPDetails{
					Hostname:         "192.168.1.1",
					Port:             &port161,
					SNMPVersion:      "2c",
					SNMPOID:          "1.3.6.1.2.1.2.2.1.5.1",
					SNMPCommunity:   "public",
					JSONPath:         ptr.To("ifSpeed"),
					JSONPathOperator: ptr.To("=="),
					ExpectedValue:    ptr.To("1000000000"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SNMP monitor","expectedValue":"1000000000","hostname":"192.168.1.1","id":10,"interval":60,"jsonPath":"ifSpeed","jsonPathOperator":"==","maxretries":2,"name":"snmp-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":161,"radiusPassword":"public","resendInterval":0,"retryInterval":60,"snmpOid":"1.3.6.1.2.1.2.2.1.5.1","snmpVersion":"2c","type":"snmp","upsideDown":false}`,
		},
		{
			name: "success with minimal fields",
			data: []byte(`{"id":11,"name":"snmp-minimal","description":null,"pathName":"snmp-minimal","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"10.0.0.1","port":161,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"snmp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"jsonPathOperator":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":"public","radiusSecret":null,"snmpVersion":"1","snmpOid":"1.3.6.1.2.1.1.3.0","mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.SNMP{
				Base: monitor.Base{
					ID:              11,
					Name:            "snmp-minimal",
					Description:     nil,
					PathName:        "snmp-minimal",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				SNMPDetails: monitor.SNMPDetails{
					Hostname:         "10.0.0.1",
					Port:             &port161,
					SNMPVersion:      "1",
					SNMPOID:          "1.3.6.1.2.1.1.3.0",
					SNMPCommunity:   "public",
					JSONPath:         nil,
					JSONPathOperator: nil,
					ExpectedValue:    nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"expectedValue":null,"hostname":"10.0.0.1","id":11,"interval":60,"jsonPath":null,"jsonPathOperator":null,"maxretries":3,"name":"snmp-minimal","notificationIDList":{},"parent":null,"port":161,"radiusPassword":"public","resendInterval":0,"retryInterval":60,"snmpOid":"1.3.6.1.2.1.1.3.0","snmpVersion":"1","type":"snmp","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			snmpMonitor := monitor.SNMP{}

			err := json.Unmarshal(tc.data, &snmpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, snmpMonitor)

			data, err := json.Marshal(snmpMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
