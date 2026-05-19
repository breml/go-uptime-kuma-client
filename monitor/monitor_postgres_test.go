package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorPostgres_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.Postgres
		wantJSON string
	}{
		{
			name: "success with connection string only",
			data: []byte(
				`{"id":5,"name":"postgres-default","description":"PostgreSQL with default query","pathName":"group / postgres-default","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"postgres","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"postgres://user:pass@localhost:5432/app","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Postgres{
				Base: monitor.Base{
					ID:              5,
					Name:            "postgres-default",
					Description:     ptr.To("PostgreSQL with default query"),
					PathName:        "group / postgres-default",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				PostgresDetails: monitor.PostgresDetails{
					DatabaseConnectionString: "postgres://user:pass@localhost:5432/app",
					DatabaseQuery:            nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"databaseConnectionString":"postgres://user:pass@localhost:5432/app","databaseQuery":null,"description":"PostgreSQL with default query","id":5,"interval":60,"maxretries":1,"name":"postgres-default","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"postgres","upsideDown":false}`,
		},
		{
			name: "success",
			data: []byte(
				`{"id":6,"name":"postgres-monitor","description":"Test PostgreSQL monitor","pathName":"group / postgres-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"postgres","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":"SELECT 1","authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"postgres://username:password@host:port/database","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.Postgres{
				Base: monitor.Base{
					ID:              6,
					Name:            "postgres-monitor",
					Description:     ptr.To("Test PostgreSQL monitor"),
					PathName:        "group / postgres-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				PostgresDetails: monitor.PostgresDetails{
					DatabaseConnectionString: "postgres://username:password@host:port/database",
					DatabaseQuery:            ptr.To("SELECT 1"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"databaseConnectionString":"postgres://username:password@host:port/database","databaseQuery":"SELECT 1","description":"Test PostgreSQL monitor","id":6,"interval":60,"maxretries":2,"name":"postgres-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"postgres","upsideDown":false}`,
		},
		{
			name: "success with conditions",
			data: []byte(
				`{"id":7,"name":"postgres-conditions","description":null,"pathName":"postgres-conditions","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"postgres","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"databaseQuery":"SELECT version()","databaseConnectionString":"postgres://user:pass@localhost:5432/app","conditions":[{"type":"expression","variable":"result","operator":"contains","value":"PostgreSQL","andOr":"and"}]}`,
			),

			want: monitor.Postgres{
				Base: monitor.Base{
					ID:              7,
					Name:            "postgres-conditions",
					Description:     nil,
					PathName:        "postgres-conditions",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				PostgresDetails: monitor.PostgresDetails{
					DatabaseConnectionString: "postgres://user:pass@localhost:5432/app",
					DatabaseQuery:            ptr.To("SELECT version()"),
					Conditions: []monitor.Condition{
						{Variable: "result", Operator: "contains", Value: "PostgreSQL", AndOr: monitor.ConditionAnd},
					},
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[{"type":"expression","variable":"result","operator":"contains","value":"PostgreSQL","andOr":"and"}],"databaseConnectionString":"postgres://user:pass@localhost:5432/app","databaseQuery":"SELECT version()","description":null,"id":7,"interval":60,"maxretries":1,"name":"postgres-conditions","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"postgres","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			postgresMonitor := monitor.Postgres{}

			err := json.Unmarshal(tc.data, &postgresMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, postgresMonitor)

			data, err := json.Marshal(postgresMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
