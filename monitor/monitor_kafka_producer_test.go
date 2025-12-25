package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorKafkaProducer_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	saslOptionsBasic := map[string]any{"mechanism": "plain", "username": "user", "password": "pass"}
	saslOptionsNone := map[string]any{"mechanism": "None"}

	tests := []struct {
		name string
		data []byte

		want     monitor.KafkaProducer
		wantJSON string
	}{
		{
			name: "success with single broker and SSL",
			data: []byte(`{"id":1,"name":"kafka-monitor","description":"Test Kafka monitor","pathName":"group / kafka-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"kafka-producer","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":"test-topic","kafkaProducerBrokers":["localhost:9092"],"kafkaProducerSsl":true,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":"test message","screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.KafkaProducer{
				Base: monitor.Base{
					ID:              1,
					Name:            "kafka-monitor",
					Description:     ptr.To("Test Kafka monitor"),
					PathName:        "group / kafka-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				KafkaProducerDetails: monitor.KafkaProducerDetails{
					Brokers:                []string{"localhost:9092"},
					Topic:                  "test-topic",
					Message:                "test message",
					SSL:                    true,
					AllowAutoTopicCreation: false,
					SASLOptions:            &saslOptionsNone,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Kafka monitor","id":1,"interval":60,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerBrokers":["localhost:9092"],"kafkaProducerMessage":"test message","kafkaProducerSaslOptions":{"mechanism":"None"},"kafkaProducerSsl":true,"kafkaProducerTopic":"test-topic","maxretries":2,"name":"kafka-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"kafka-producer","upsideDown":false}`,
		},
		{
			name: "success with multiple brokers and SASL",
			data: []byte(`{"id":2,"name":"kafka-cluster","description":"Test Kafka cluster","pathName":"group / kafka-cluster","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"kafka-producer","timeout":30,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":"events","kafkaProducerBrokers":["kafka1:9092","kafka2:9092","kafka3:9092"],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":true,"kafkaProducerMessage":"health check","screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"plain","username":"user","password":"pass"},"includeSensitiveData":true}`),

			want: monitor.KafkaProducer{
				Base: monitor.Base{
					ID:              2,
					Name:            "kafka-cluster",
					Description:     ptr.To("Test Kafka cluster"),
					PathName:        "group / kafka-cluster",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				KafkaProducerDetails: monitor.KafkaProducerDetails{
					Brokers:                []string{"kafka1:9092", "kafka2:9092", "kafka3:9092"},
					Topic:                  "events",
					Message:                "health check",
					SSL:                    false,
					AllowAutoTopicCreation: true,
					SASLOptions:            &saslOptionsBasic,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Kafka cluster","id":2,"interval":120,"kafkaProducerAllowAutoTopicCreation":true,"kafkaProducerBrokers":["kafka1:9092","kafka2:9092","kafka3:9092"],"kafkaProducerMessage":"health check","kafkaProducerSaslOptions":{"mechanism":"plain","password":"pass","username":"user"},"kafkaProducerSsl":false,"kafkaProducerTopic":"events","maxretries":3,"name":"kafka-cluster","notificationIDList":{},"parent":1,"resendInterval":0,"retryInterval":60,"type":"kafka-producer","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kafkaMonitor := monitor.KafkaProducer{}

			err := json.Unmarshal(tc.data, &kafkaMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, kafkaMonitor)

			data, err := json.Marshal(kafkaMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
