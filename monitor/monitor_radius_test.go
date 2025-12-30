package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorRadius_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	port1812 := int64(1812)
	port1813 := int64(1813)
	calledStation := "555-1234"
	callingStation := "555-9999"

	tests := []struct {
		name string
		data []byte

		want     monitor.Radius
		wantJSON string
	}{
		{
			name: "success with default port",
			data: []byte(
				`{"id":1,"name":"radius-monitor","description":"Test Radius monitor","pathName":"group / radius-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"radius.example.com","port":1812,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"radius","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":"testuser","radiusPassword":"testpass","radiusSecret":"sharedsecret","radiusCalledStationId":null,"radiusCallingStationId":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Radius{
				Base: monitor.Base{
					ID:              1,
					Name:            "radius-monitor",
					Description:     ptr.To("Test Radius monitor"),
					PathName:        "group / radius-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				RadiusDetails: monitor.RadiusDetails{
					Hostname:         "radius.example.com",
					Port:             &port1812,
					Username:         "testuser",
					Password:         "testpass",
					Secret:           "sharedsecret",
					CalledStationID:  nil,
					CallingStationID: nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Radius monitor","hostname":"radius.example.com","id":1,"interval":60,"maxretries":2,"name":"radius-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":1812,"radiusCalledStationId":null,"radiusCallingStationId":null,"radiusPassword":"testpass","radiusSecret":"sharedsecret","radiusUsername":"testuser","resendInterval":0,"retryInterval":60,"type":"radius","upsideDown":false}`,
		},
		{
			name: "success with custom port and station IDs",
			data: []byte(
				`{"id":2,"name":"radius-advanced","description":"Test Radius with station IDs","pathName":"group / radius-advanced","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"auth.example.com","port":1813,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"radius","timeout":60,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":"555-1234","radiusCallingStationId":"555-9999","game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":"admin","radiusPassword":"adminpass","radiusSecret":"secret123","radiusCalledStationId":"555-1234","radiusCallingStationId":"555-9999","mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Radius{
				Base: monitor.Base{
					ID:              2,
					Name:            "radius-advanced",
					Description:     ptr.To("Test Radius with station IDs"),
					PathName:        "group / radius-advanced",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				RadiusDetails: monitor.RadiusDetails{
					Hostname:         "auth.example.com",
					Port:             &port1813,
					Username:         "admin",
					Password:         "adminpass",
					Secret:           "secret123",
					CalledStationID:  &calledStation,
					CallingStationID: &callingStation,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test Radius with station IDs","hostname":"auth.example.com","id":2,"interval":120,"maxretries":3,"name":"radius-advanced","notificationIDList":{},"parent":1,"port":1813,"radiusCalledStationId":"555-1234","radiusCallingStationId":"555-9999","radiusPassword":"adminpass","radiusSecret":"secret123","radiusUsername":"admin","resendInterval":0,"retryInterval":60,"type":"radius","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			radiusMonitor := monitor.Radius{}

			err := json.Unmarshal(tc.data, &radiusMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, radiusMonitor)

			data, err := json.Marshal(radiusMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
