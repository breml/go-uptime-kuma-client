package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSQLServer_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.SQLServer
		wantJSON string
	}{
		{
			name: "success with connection string only",
			data: []byte(`{"id":1,"name":"sqlserver-monitor","description":"Test SQL Server monitor","pathName":"group / sqlserver-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"sqlserver","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"Server=localhost;User Id=sa;Password=password123;","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.SQLServer{
				Base: monitor.Base{
					ID:              1,
					Name:            "sqlserver-monitor",
					Description:     ptr.To("Test SQL Server monitor"),
					PathName:        "group / sqlserver-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SQLServerDetails: monitor.SQLServerDetails{
					DatabaseConnectionString: "Server=localhost;User Id=sa;Password=password123;",
					DatabaseQuery:            nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"databaseConnectionString":"Server=localhost;User Id=sa;Password=password123;","databaseQuery":null,"description":"Test SQL Server monitor","id":1,"interval":60,"maxretries":2,"name":"sqlserver-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"sqlserver","upsideDown":false}`,
		},
		{
			name: "success with connection string and query",
			data: []byte(`{"id":2,"name":"sqlserver-query","description":"Test SQL Server with query","pathName":"group / sqlserver-query","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"sqlserver","timeout":60,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":"SELECT COUNT(*) FROM sys.tables;","authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"Server=sqlserver.example.com,1433;Database=testdb;User Id=user;Password=pass;","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.SQLServer{
				Base: monitor.Base{
					ID:              2,
					Name:            "sqlserver-query",
					Description:     ptr.To("Test SQL Server with query"),
					PathName:        "group / sqlserver-query",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				SQLServerDetails: monitor.SQLServerDetails{
					DatabaseConnectionString: "Server=sqlserver.example.com,1433;Database=testdb;User Id=user;Password=pass;",
					DatabaseQuery:            ptr.To("SELECT COUNT(*) FROM sys.tables;"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"databaseConnectionString":"Server=sqlserver.example.com,1433;Database=testdb;User Id=user;Password=pass;","databaseQuery":"SELECT COUNT(*) FROM sys.tables;","description":"Test SQL Server with query","id":2,"interval":120,"maxretries":3,"name":"sqlserver-query","notificationIDList":{},"parent":1,"resendInterval":0,"retryInterval":60,"type":"sqlserver","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sqlserverMonitor := monitor.SQLServer{}

			err := json.Unmarshal(tc.data, &sqlserverMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sqlserverMonitor)

			data, err := json.Marshal(sqlserverMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
