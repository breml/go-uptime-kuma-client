package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestBase_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want monitor.Base
	}{
		{
			name: "success with parent",
			data: []byte(`{"id":2,"name":"foobar.com","description":null,"pathName":"group / foobar.com","parent":1,"childrenIDs":[],"url":"https://www.foobar.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"http","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.Base{
				ID:              2,
				Name:            "foobar.com",
				Description:     nil,
				PathName:        "group / foobar.com",
				Parent:          &parent1,
				Interval:        60,
				RetryInterval:   60,
				ResendInterval:  0,
				MaxRetries:      2,
				UpsideDown:      false,
				NotificationIDs: []int64{1},
				IsActive:        true,
			},
		},
		{
			name: "parent null",
			data: []byte(`{"id":1,"name":"top-level","description":null,"pathName":"top-level","parent":null,"proxyId":null,"childrenIDs":[],"maxretries":0,"weight":2000,"active":true,"type":"group","interval":60,"retryInterval":60,"resendInterval":0,"upsideDown":false,"notificationIDList":{},"tags":[],"maintenance":false}`),

			want: monitor.Base{
				ID:              1,
				Name:            "top-level",
				Description:     nil,
				PathName:        "top-level",
				Parent:          nil,
				ProxyID:         nil,
				Interval:        60,
				RetryInterval:   60,
				ResendInterval:  0,
				MaxRetries:      0,
				UpsideDown:      false,
				NotificationIDs: nil,
				IsActive:        true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			someMonitor := monitor.Base{}

			err := json.Unmarshal(tc.data, &someMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, someMonitor)

			data, err := json.Marshal(someMonitor)
			require.NoError(t, err)

			require.JSONEq(t, string(tc.data), string(data))
		})
	}
}
