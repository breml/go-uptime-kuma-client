package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorGroup_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     monitor.Group
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"group","description":null,"pathName":"group","parent":null,"childrenIDs":[2],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":0,"weight":2000,"active":false,"forceInactive":false,"type":"group","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.Group{
				Base: monitor.Base{
					ID:             1,
					Name:           "group",
					Description:    nil,
					PathName:       "group",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       false,
				},
				GroupDetails: monitor.GroupDetails{},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":false,"description":null,"id":1,"interval":60,"maxretries":0,"name":"group","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"group","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			groupMonitor := monitor.Group{}

			err := json.Unmarshal(tc.data, &groupMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, groupMonitor)

			data, err := json.Marshal(groupMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
