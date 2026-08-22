package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorPM2_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.PM2
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":4,"name":"pm2-monitor","description":"Test PM2 monitor","pathName":"group / pm2-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"pm2","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"system_service_name":"api-server","includeSensitiveData":true}`,
			),

			want: monitor.PM2{
				Base: monitor.Base{
					ID:              4,
					Name:            "pm2-monitor",
					Description:     ptr.To("Test PM2 monitor"),
					PathName:        "group / pm2-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				PM2Details: monitor.PM2Details{
					ProcessName: "api-server",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test PM2 monitor","id":4,"interval":60,"maxretries":2,"name":"pm2-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"system_service_name":"api-server","type":"pm2","upsideDown":false}`,
		},
		{
			// The "worker 1" fixture contains a space: the system-service
			// monitor rejects it, pm2 forbids only ASCII control characters.
			name: "success without parent",
			data: []byte(
				`{"id":5,"name":"pm2-monitor-no-parent","description":null,"pathName":"pm2-monitor-no-parent","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"pm2","timeout":null,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"system_service_name":"worker 1","includeSensitiveData":true}`,
			),

			want: monitor.PM2{
				Base: monitor.Base{
					ID:             5,
					Name:           "pm2-monitor-no-parent",
					Description:    nil,
					PathName:       "pm2-monitor-no-parent",
					Parent:         nil,
					Interval:       120,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				PM2Details: monitor.PM2Details{
					ProcessName: "worker 1",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"id":5,"interval":120,"maxretries":3,"name":"pm2-monitor-no-parent","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"system_service_name":"worker 1","type":"pm2","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pm2Monitor := monitor.PM2{}

			err := json.Unmarshal(tc.data, &pm2Monitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, pm2Monitor)

			data, err := json.Marshal(pm2Monitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestMonitorPM2_MarshalRequiresProcessName(t *testing.T) {
	tests := []struct {
		name        string
		processName string
	}{
		{
			name:        "empty",
			processName: "",
		},
		{
			name:        "whitespace only",
			processName: "   ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pm2Monitor := monitor.PM2{
				Base:       monitor.Base{Name: "pm2-without-process-name"},
				PM2Details: monitor.PM2Details{ProcessName: tc.processName},
			}

			_, err := json.Marshal(pm2Monitor)
			require.Error(t, err)
			require.ErrorContains(t, err, "process name is required")
		})
	}
}
