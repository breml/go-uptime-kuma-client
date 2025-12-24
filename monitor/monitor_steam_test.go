package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSteam_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	timeout48 := int64(48)
	timeout30 := int64(30)

	tests := []struct {
		name string
		data []byte

		want     monitor.Steam
		wantJSON string
	}{
		{
			name: "success with parent and timeout",
			data: []byte(`{"id":5,"name":"steam-monitor","description":"Test Steam monitor","pathName":"group / steam-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"192.168.1.100","port":27015,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"steam","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.Steam{
				Base: monitor.Base{
					ID:              5,
					Name:            "steam-monitor",
					Description:     ptr.To("Test Steam monitor"),
					PathName:        "group / steam-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SteamDetails: monitor.SteamDetails{
					Hostname: "192.168.1.100",
					Port:     27015,
					Timeout:  &timeout48,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Steam monitor","hostname":"192.168.1.100","id":5,"interval":60,"maxretries":2,"name":"steam-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":27015,"resendInterval":0,"retryInterval":60,"timeout":48,"type":"steam","upsideDown":false}`,
		},
		{
			name: "success with null timeout",
			data: []byte(`{"id":6,"name":"steam-monitor-no-timeout","description":null,"pathName":"steam-monitor-no-timeout","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"10.0.0.1","port":27016,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"steam","timeout":null,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.Steam{
				Base: monitor.Base{
					ID:              6,
					Name:            "steam-monitor-no-timeout",
					Description:     nil,
					PathName:        "steam-monitor-no-timeout",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				SteamDetails: monitor.SteamDetails{
					Hostname: "10.0.0.1",
					Port:     27016,
					Timeout:  nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"hostname":"10.0.0.1","id":6,"interval":60,"maxretries":3,"name":"steam-monitor-no-timeout","notificationIDList":{},"parent":null,"port":27016,"resendInterval":0,"retryInterval":60,"timeout":null,"type":"steam","upsideDown":false}`,
		},
		{
			name: "success with different timeout",
			data: []byte(`{"id":7,"name":"steam-monitor-custom-timeout","description":null,"pathName":"steam-monitor-custom-timeout","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":27020,"maxretries":1,"weight":2000,"active":false,"forceInactive":false,"type":"steam","timeout":30,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.Steam{
				Base: monitor.Base{
					ID:              7,
					Name:            "steam-monitor-custom-timeout",
					Description:     nil,
					PathName:        "steam-monitor-custom-timeout",
					Parent:          nil,
					Interval:        120,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        false,
				},
				SteamDetails: monitor.SteamDetails{
					Hostname: "example.com",
					Port:     27020,
					Timeout:  &timeout30,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":false,"conditions":[],"description":null,"hostname":"example.com","id":7,"interval":120,"maxretries":1,"name":"steam-monitor-custom-timeout","notificationIDList":{},"parent":null,"port":27020,"resendInterval":0,"retryInterval":120,"timeout":30,"type":"steam","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			steamMonitor := monitor.Steam{}

			err := json.Unmarshal(tc.data, &steamMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, steamMonitor)

			data, err := json.Marshal(steamMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
