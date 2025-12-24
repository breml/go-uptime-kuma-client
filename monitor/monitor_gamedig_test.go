package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorGameDig_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.GameDig
		wantJSON string
	}{
		{
			name: "success with parent and given port only",
			data: []byte(`{"id":5,"name":"gamedig-monitor","description":"Test GameDig monitor","pathName":"group / gamedig-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"192.168.1.100","port":25565,"game":"minecraft","gamedigGivenPortOnly":true,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"gamedig","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              5,
					Name:            "gamedig-monitor",
					Description:     func() *string { s := "Test GameDig monitor"; return &s }(),
					PathName:        "group / gamedig-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				GameDigDetails: monitor.GameDigDetails{
					Hostname:             "192.168.1.100",
					Port:                 25565,
					Game:                 "minecraft",
					GameDigGivenPortOnly: true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test GameDig monitor","game":"minecraft","gamedigGivenPortOnly":true,"hostname":"192.168.1.100","id":5,"interval":60,"maxretries":2,"name":"gamedig-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":25565,"resendInterval":0,"retryInterval":60,"type":"gamedig","upsideDown":false}`,
		},
		{
			name: "success with null description and given port only false",
			data: []byte(`{"id":6,"name":"gamedig-monitor-csgo","description":null,"pathName":"gamedig-monitor-csgo","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"10.0.0.1","port":27015,"game":"csgo","gamedigGivenPortOnly":false,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"gamedig","timeout":null,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              6,
					Name:            "gamedig-monitor-csgo",
					Description:     nil,
					PathName:        "gamedig-monitor-csgo",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				GameDigDetails: monitor.GameDigDetails{
					Hostname:             "10.0.0.1",
					Port:                 27015,
					Game:                 "csgo",
					GameDigGivenPortOnly: false,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"game":"csgo","gamedigGivenPortOnly":false,"hostname":"10.0.0.1","id":6,"interval":60,"maxretries":3,"name":"gamedig-monitor-csgo","notificationIDList":{},"parent":null,"port":27015,"resendInterval":0,"retryInterval":60,"type":"gamedig","upsideDown":false}`,
		},
		{
			name: "success with different game type and inactive",
			data: []byte(`{"id":7,"name":"gamedig-monitor-rust","description":null,"pathName":"gamedig-monitor-rust","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":28015,"game":"rust","gamedigGivenPortOnly":true,"maxretries":1,"weight":2000,"active":false,"forceInactive":false,"type":"gamedig","timeout":30,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              7,
					Name:            "gamedig-monitor-rust",
					Description:     nil,
					PathName:        "gamedig-monitor-rust",
					Parent:          nil,
					Interval:        120,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        false,
				},
				GameDigDetails: monitor.GameDigDetails{
					Hostname:             "example.com",
					Port:                 28015,
					Game:                 "rust",
					GameDigGivenPortOnly: true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":false,"conditions":[],"description":null,"game":"rust","gamedigGivenPortOnly":true,"hostname":"example.com","id":7,"interval":120,"maxretries":1,"name":"gamedig-monitor-rust","notificationIDList":{},"parent":null,"port":28015,"resendInterval":0,"retryInterval":120,"type":"gamedig","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gamedigMonitor := monitor.GameDig{}

			err := json.Unmarshal(tc.data, &gamedigMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, gamedigMonitor)

			data, err := json.Marshal(gamedigMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
