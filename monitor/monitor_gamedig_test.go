package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
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
			name: "success with minecraft",
			data: []byte(`{"id":5,"name":"minecraft-server","description":"Test Minecraft Server","pathName":"group / minecraft-server","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mc.example.com","port":25565,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"gamedig","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":"minecraft","gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              5,
					Name:            "minecraft-server",
					Description:     ptr.To("Test Minecraft Server"),
					PathName:        "group / minecraft-server",
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
					Hostname:             "mc.example.com",
					Port:                 25565,
					Game:                 "minecraft",
					GameDigGivenPortOnly: true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Minecraft Server","game":"minecraft","gamedigGivenPortOnly":true,"hostname":"mc.example.com","id":5,"interval":60,"maxretries":2,"name":"minecraft-server","notificationIDList":{"1":true,"2":true},"parent":1,"port":25565,"resendInterval":0,"retryInterval":60,"type":"gamedig","upsideDown":false}`,
		},
		{
			name: "success with csgo",
			data: []byte(`{"id":6,"name":"csgo-server","description":null,"pathName":"csgo-server","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"192.168.1.100","port":27015,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"gamedig","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":"csgo","gamedigGivenPortOnly":false,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              6,
					Name:            "csgo-server",
					Description:     nil,
					PathName:        "csgo-server",
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
					Hostname:             "192.168.1.100",
					Port:                 27015,
					Game:                 "csgo",
					GameDigGivenPortOnly: false,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"game":"csgo","gamedigGivenPortOnly":false,"hostname":"192.168.1.100","id":6,"interval":60,"maxretries":3,"name":"csgo-server","notificationIDList":{},"parent":null,"port":27015,"resendInterval":0,"retryInterval":60,"type":"gamedig","upsideDown":false}`,
		},
		{
			name: "success with ark",
			data: []byte(`{"id":7,"name":"ark-server","description":null,"pathName":"ark-server","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"10.0.0.1","port":27015,"maxretries":1,"weight":2000,"active":false,"forceInactive":false,"type":"gamedig","timeout":48,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":"ark","gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GameDig{
				Base: monitor.Base{
					ID:              7,
					Name:            "ark-server",
					Description:     nil,
					PathName:        "ark-server",
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
					Hostname:             "10.0.0.1",
					Port:                 27015,
					Game:                 "ark",
					GameDigGivenPortOnly: true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":false,"conditions":[],"description":null,"game":"ark","gamedigGivenPortOnly":true,"hostname":"10.0.0.1","id":7,"interval":120,"maxretries":1,"name":"ark-server","notificationIDList":{},"parent":null,"port":27015,"resendInterval":0,"retryInterval":120,"type":"gamedig","upsideDown":false}`,
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
