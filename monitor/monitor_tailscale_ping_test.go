package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorTailscalePing_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.TailscalePing
		wantJSON string
	}{
		{
			name: "success with tailscale IP",
			data: []byte(`{"id":1,"name":"tailscale-monitor","description":"Test Tailscale Ping monitor","pathName":"group / tailscale-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"100.100.100.100","port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"tailscale-ping","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.TailscalePing{
				Base: monitor.Base{
					ID:              1,
					Name:            "tailscale-monitor",
					Description:     ptr.To("Test Tailscale Ping monitor"),
					PathName:        "group / tailscale-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				TailscalePingDetails: monitor.TailscalePingDetails{
					Hostname: "100.100.100.100",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Tailscale Ping monitor","hostname":"100.100.100.100","id":1,"interval":60,"maxretries":2,"name":"tailscale-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"tailscale-ping","upsideDown":false}`,
		},
		{
			name: "success with tailscale hostname",
			data: []byte(`{"id":2,"name":"tailscale-hostname","description":"Test Tailscale with hostname","pathName":"group / tailscale-hostname","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mydevice.mydomain.ts.net","port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"tailscale-ping","timeout":60,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.TailscalePing{
				Base: monitor.Base{
					ID:              2,
					Name:            "tailscale-hostname",
					Description:     ptr.To("Test Tailscale with hostname"),
					PathName:        "group / tailscale-hostname",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				TailscalePingDetails: monitor.TailscalePingDetails{
					Hostname: "mydevice.mydomain.ts.net",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Tailscale with hostname","hostname":"mydevice.mydomain.ts.net","id":2,"interval":120,"maxretries":3,"name":"tailscale-hostname","notificationIDList":{},"parent":1,"resendInterval":0,"retryInterval":60,"type":"tailscale-ping","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tailscaleMonitor := monitor.TailscalePing{}

			err := json.Unmarshal(tc.data, &tailscaleMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, tailscaleMonitor)

			data, err := json.Marshal(tailscaleMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
