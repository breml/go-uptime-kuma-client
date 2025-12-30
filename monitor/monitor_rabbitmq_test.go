package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorRabbitMQ_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	timeout48 := int64(48)
	username := "guest"
	password := "guest"

	tests := []struct {
		name string
		data []byte

		want     monitor.RabbitMQ
		wantJSON string
	}{
		{
			name: "success with single node",
			data: []byte(
				`{"id":1,"name":"rabbitmq-monitor","description":"Test RabbitMQ monitor","pathName":"group / rabbitmq-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"rabbitmq","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"rabbitmqNodes":"[\"http://rabbitmq.example.com:15672/\"]","rabbitmqUsername":"guest","rabbitmqPassword":"guest","includeSensitiveData":true}`,
			),

			want: monitor.RabbitMQ{
				Base: monitor.Base{
					ID:              1,
					Name:            "rabbitmq-monitor",
					Description:     ptr.To("Test RabbitMQ monitor"),
					PathName:        "group / rabbitmq-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				RabbitMQDetails: monitor.RabbitMQDetails{
					Nodes:    "[\"http://rabbitmq.example.com:15672/\"]",
					Username: &username,
					Password: &password,
					Timeout:  &timeout48,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test RabbitMQ monitor","id":1,"interval":60,"maxretries":2,"name":"rabbitmq-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"rabbitmqNodes":"[\"http://rabbitmq.example.com:15672/\"]","rabbitmqPassword":"guest","rabbitmqUsername":"guest","resendInterval":0,"retryInterval":60,"timeout":48,"type":"rabbitmq","upsideDown":false}`,
		},
		{
			name: "success with multiple nodes and no authentication",
			data: []byte(
				`{"id":2,"name":"rabbitmq-cluster","description":"Test RabbitMQ cluster monitor","pathName":"group / rabbitmq-cluster","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"rabbitmq","timeout":30,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"rabbitmqNodes":"[\"http://rabbitmq1.example.com:15672/\",\"http://rabbitmq2.example.com:15672/\",\"http://rabbitmq3.example.com:15672/\"]","rabbitmqUsername":null,"rabbitmqPassword":null,"includeSensitiveData":true}`,
			),

			want: monitor.RabbitMQ{
				Base: monitor.Base{
					ID:              2,
					Name:            "rabbitmq-cluster",
					Description:     ptr.To("Test RabbitMQ cluster monitor"),
					PathName:        "group / rabbitmq-cluster",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				RabbitMQDetails: monitor.RabbitMQDetails{
					Nodes:    "[\"http://rabbitmq1.example.com:15672/\",\"http://rabbitmq2.example.com:15672/\",\"http://rabbitmq3.example.com:15672/\"]",
					Username: nil,
					Password: nil,
					Timeout:  ptr.To(int64(30)),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test RabbitMQ cluster monitor","id":2,"interval":120,"maxretries":3,"name":"rabbitmq-cluster","notificationIDList":{},"parent":1,"rabbitmqNodes":"[\"http://rabbitmq1.example.com:15672/\",\"http://rabbitmq2.example.com:15672/\",\"http://rabbitmq3.example.com:15672/\"]","rabbitmqPassword":null,"rabbitmqUsername":null,"resendInterval":0,"retryInterval":60,"timeout":30,"type":"rabbitmq","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rmqMonitor := monitor.RabbitMQ{}

			err := json.Unmarshal(tc.data, &rmqMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, rmqMonitor)

			data, err := json.Marshal(rmqMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
