package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSIPOptions_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.SIPOptions
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":4,"name":"sip-options-monitor","description":"Test SIP Options monitor","pathName":"group / sip-options-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"sip.example.com","port":5060,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"sip-options","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.SIPOptions{
				Base: monitor.Base{
					ID:              4,
					Name:            "sip-options-monitor",
					Description:     ptr.To("Test SIP Options monitor"),
					PathName:        "group / sip-options-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SIPOptionsDetails: monitor.SIPOptionsDetails{
					Hostname: "sip.example.com",
					Port:     5060,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SIP Options monitor","hostname":"sip.example.com","id":4,"interval":60,"maxretries":2,"name":"sip-options-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":5060,"resendInterval":0,"retryInterval":60,"type":"sip-options","upsideDown":false}`,
		},
		{
			name: "success without parent",
			data: []byte(
				`{"id":5,"name":"sip-options-monitor-no-parent","description":null,"pathName":"sip-options-monitor-no-parent","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"192.168.1.10","port":5061,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"sip-options","timeout":null,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.SIPOptions{
				Base: monitor.Base{
					ID:             5,
					Name:           "sip-options-monitor-no-parent",
					Description:    nil,
					PathName:       "sip-options-monitor-no-parent",
					Parent:         nil,
					Interval:       120,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				SIPOptionsDetails: monitor.SIPOptionsDetails{
					Hostname: "192.168.1.10",
					Port:     5061,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"hostname":"192.168.1.10","id":5,"interval":120,"maxretries":3,"name":"sip-options-monitor-no-parent","notificationIDList":{},"parent":null,"port":5061,"resendInterval":0,"retryInterval":60,"type":"sip-options","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sipOptionsMonitor := monitor.SIPOptions{}

			err := json.Unmarshal(tc.data, &sipOptionsMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sipOptionsMonitor)

			data, err := json.Marshal(sipOptionsMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
