package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorMQTT_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	port1883 := int64(1883)
	port8883 := int64(8883)
	port443 := int64(443)
	emptyString := ""
	successMsg := "OK"

	tests := []struct {
		name string
		data []byte

		want     monitor.MQTT
		wantJSON string
	}{
		{
			name: "success with keyword check type",
			data: []byte(
				`{"id":5,"name":"mqtt-monitor","description":"Test MQTT monitor","pathName":"group / mqtt-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mqtt.example.com","port":1883,"mqttTopic":"home/temperature","mqttUsername":"","mqttPassword":"","mqttWebsocketPath":null,"mqttCheckType":"keyword","mqttSuccessMessage":"OK","jsonPath":null,"expectedValue":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"mqtt","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"httpBodyEncoding":"json","kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.MQTT{
				Base: monitor.Base{
					ID:              5,
					Name:            "mqtt-monitor",
					Description:     ptr.To("Test MQTT monitor"),
					PathName:        "group / mqtt-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				MQTTDetails: monitor.MQTTDetails{
					Hostname:           "mqtt.example.com",
					Port:               &port1883,
					MQTTTopic:          "home/temperature",
					MQTTUsername:       &emptyString,
					MQTTPassword:       &emptyString,
					MQTTWebsocketPath:  nil,
					MQTTCheckType:      monitor.MQTTCheckTypeKeyword,
					MQTTSuccessMessage: &successMsg,
					JSONPath:           nil,
					ExpectedValue:      nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test MQTT monitor","expectedValue":null,"hostname":"mqtt.example.com","id":5,"interval":60,"jsonPath":null,"maxretries":2,"mqttCheckType":"keyword","mqttPassword":"","mqttSuccessMessage":"OK","mqttTopic":"home/temperature","mqttUsername":"","mqttWebsocketPath":null,"name":"mqtt-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":1883,"resendInterval":0,"retryInterval":60,"type":"mqtt","upsideDown":false}`,
		},
		{
			name: "success with json-query check type",
			data: []byte(
				`{"id":6,"name":"mqtt-json-monitor","description":null,"pathName":"mqtt-json-monitor","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"mqtt.local","port":8883,"mqttTopic":"home/data","mqttUsername":null,"mqttPassword":null,"mqttWebsocketPath":null,"mqttCheckType":"json-query","mqttSuccessMessage":null,"jsonPath":"temperature","expectedValue":"25","maxretries":3,"weight":2000,"active":false,"forceInactive":false,"type":"mqtt","timeout":null,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"httpBodyEncoding":"json","kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.MQTT{
				Base: monitor.Base{
					ID:              6,
					Name:            "mqtt-json-monitor",
					Description:     nil,
					PathName:        "mqtt-json-monitor",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        false,
				},
				MQTTDetails: monitor.MQTTDetails{
					Hostname:           "mqtt.local",
					Port:               &port8883,
					MQTTTopic:          "home/data",
					MQTTUsername:       nil,
					MQTTPassword:       nil,
					MQTTWebsocketPath:  nil,
					MQTTCheckType:      monitor.MQTTCheckTypeJSONQuery,
					MQTTSuccessMessage: nil,
					JSONPath:           ptr.To("temperature"),
					ExpectedValue:      ptr.To("25"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":false,"conditions":[],"description":null,"expectedValue":"25","hostname":"mqtt.local","id":6,"interval":60,"jsonPath":"temperature","maxretries":3,"mqttCheckType":"json-query","mqttPassword":null,"mqttSuccessMessage":null,"mqttTopic":"home/data","mqttUsername":null,"mqttWebsocketPath":null,"name":"mqtt-json-monitor","notificationIDList":{},"parent":null,"port":8883,"resendInterval":0,"retryInterval":60,"type":"mqtt","upsideDown":false}`,
		},
		{
			name: "success with websocket path",
			data: []byte(
				`{"id":7,"name":"mqtt-ws-monitor","description":null,"pathName":"mqtt-ws-monitor","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"wss://mqtt.cloud.com","port":443,"mqttTopic":"sensor/humidity","mqttUsername":"clouduser","mqttPassword":"cloudpass","mqttWebsocketPath":"/mqtt","mqttCheckType":"keyword","mqttSuccessMessage":"valid","jsonPath":null,"expectedValue":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"mqtt","timeout":30,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"httpBodyEncoding":"json","kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.MQTT{
				Base: monitor.Base{
					ID:              7,
					Name:            "mqtt-ws-monitor",
					Description:     nil,
					PathName:        "mqtt-ws-monitor",
					Parent:          nil,
					Interval:        120,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				MQTTDetails: monitor.MQTTDetails{
					Hostname:           "wss://mqtt.cloud.com",
					Port:               &port443,
					MQTTTopic:          "sensor/humidity",
					MQTTUsername:       ptr.To("clouduser"),
					MQTTPassword:       ptr.To("cloudpass"),
					MQTTWebsocketPath:  ptr.To("/mqtt"),
					MQTTCheckType:      monitor.MQTTCheckTypeKeyword,
					MQTTSuccessMessage: ptr.To("valid"),
					JSONPath:           nil,
					ExpectedValue:      nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"expectedValue":null,"hostname":"wss://mqtt.cloud.com","id":7,"interval":120,"jsonPath":null,"maxretries":1,"mqttCheckType":"keyword","mqttPassword":"cloudpass","mqttSuccessMessage":"valid","mqttTopic":"sensor/humidity","mqttUsername":"clouduser","mqttWebsocketPath":"/mqtt","name":"mqtt-ws-monitor","notificationIDList":{},"parent":null,"port":443,"resendInterval":0,"retryInterval":120,"type":"mqtt","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mqttMonitor := monitor.MQTT{}

			err := json.Unmarshal(tc.data, &mqttMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, mqttMonitor)

			data, err := json.Marshal(mqttMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
