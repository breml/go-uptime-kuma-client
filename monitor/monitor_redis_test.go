package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorRedis_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.Redis
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":6,"name":"redis-monitor","description":"Test Redis monitor","pathName":"group / redis-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"redis","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"redis://user:password@localhost:6379","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Redis{
				Base: monitor.Base{
					ID:              6,
					Name:            "redis-monitor",
					Description:     ptr.To("Test Redis monitor"),
					PathName:        "group / redis-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				RedisDetails: monitor.RedisDetails{
					ConnectionString: "redis://user:password@localhost:6379",
					IgnoreTLS:        false,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"databaseConnectionString":"redis://user:password@localhost:6379","description":"Test Redis monitor","id":6,"ignoreTls":false,"interval":60,"maxretries":2,"name":"redis-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"redis","upsideDown":false}`,
		},
		{
			name: "success with conditions",
			data: []byte(
				`{"id":7,"name":"redis-conditions","description":null,"pathName":"redis-conditions","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"redis","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":true,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"databaseQuery":null,"databaseConnectionString":"rediss://localhost:6380","conditions":[{"type":"expression","variable":"result","operator":"==","value":"PONG","andOr":"and"}]}`,
			),

			want: monitor.Redis{
				Base: monitor.Base{
					ID:              7,
					Name:            "redis-conditions",
					Description:     nil,
					PathName:        "redis-conditions",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				RedisDetails: monitor.RedisDetails{
					ConnectionString: "rediss://localhost:6380",
					IgnoreTLS:        true,
					Conditions: []monitor.Condition{
						{Variable: "result", Operator: "==", Value: "PONG", AndOr: monitor.ConditionAnd},
					},
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[{"type":"expression","variable":"result","operator":"==","value":"PONG","andOr":"and"}],"databaseConnectionString":"rediss://localhost:6380","description":null,"id":7,"ignoreTls":true,"interval":60,"maxretries":1,"name":"redis-conditions","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"redis","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			redisMonitor := monitor.Redis{}

			err := json.Unmarshal(tc.data, &redisMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, redisMonitor)

			data, err := json.Marshal(redisMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
