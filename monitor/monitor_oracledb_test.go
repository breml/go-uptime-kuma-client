package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorOracleDB_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.OracleDB
		wantJSON string
	}{
		{
			name: "success with connection string only",
			data: []byte(
				`{"id":1,"name":"oracledb-monitor","description":"Test OracleDB monitor","pathName":"group / oracledb-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"oracledb","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":"oracle","basic_auth_pass":"secret","oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"localhost:1521/XEPDB1","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.OracleDB{
				Base: monitor.Base{
					ID:              1,
					Name:            "oracledb-monitor",
					Description:     ptr.To("Test OracleDB monitor"),
					PathName:        "group / oracledb-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				OracleDBDetails: monitor.OracleDBDetails{
					DatabaseConnectionString: "localhost:1521/XEPDB1",
					DatabaseQuery:            nil,
					Username:                 "oracle",
					Password:                 "secret",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"basic_auth_pass":"secret","basic_auth_user":"oracle","conditions":[],"databaseConnectionString":"localhost:1521/XEPDB1","databaseQuery":null,"description":"Test OracleDB monitor","id":1,"interval":60,"maxretries":2,"name":"oracledb-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":0,"retryInterval":60,"type":"oracledb","upsideDown":false}`,
		},
		{
			name: "success with connection string and query",
			data: []byte(
				`{"id":2,"name":"oracledb-query","description":"Test OracleDB with query","pathName":"group / oracledb-query","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"oracledb","timeout":60,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":"SELECT COUNT(*) FROM user_tables","authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":"admin","basic_auth_pass":"adminpass","oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":"oracle.example.com:1521/PROD","radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.OracleDB{
				Base: monitor.Base{
					ID:              2,
					Name:            "oracledb-query",
					Description:     ptr.To("Test OracleDB with query"),
					PathName:        "group / oracledb-query",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				OracleDBDetails: monitor.OracleDBDetails{
					DatabaseConnectionString: "oracle.example.com:1521/PROD",
					DatabaseQuery:            ptr.To("SELECT COUNT(*) FROM user_tables"),
					Username:                 "admin",
					Password:                 "adminpass",
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"basic_auth_pass":"adminpass","basic_auth_user":"admin","conditions":[],"databaseConnectionString":"oracle.example.com:1521/PROD","databaseQuery":"SELECT COUNT(*) FROM user_tables","description":"Test OracleDB with query","id":2,"interval":120,"maxretries":3,"name":"oracledb-query","notificationIDList":{},"parent":1,"resendInterval":0,"retryInterval":60,"type":"oracledb","upsideDown":false}`,
		},
		{
			name: "success with conditions",
			data: []byte(
				`{"id":3,"name":"oracledb-conditions","description":null,"pathName":"oracledb-conditions","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"oracledb","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":null,"dns_resolve_server":null,"dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"databaseQuery":"SELECT banner FROM v$version WHERE banner LIKE 'Oracle%'","databaseConnectionString":"db.example.com:1521/ORCL","basic_auth_user":"system","basic_auth_pass":"oracle","conditions":[{"type":"expression","variable":"result","operator":"contains","value":"Oracle","andOr":"and"}]}`,
			),

			want: monitor.OracleDB{
				Base: monitor.Base{
					ID:              3,
					Name:            "oracledb-conditions",
					Description:     nil,
					PathName:        "oracledb-conditions",
					Parent:          nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				OracleDBDetails: monitor.OracleDBDetails{
					DatabaseConnectionString: "db.example.com:1521/ORCL",
					DatabaseQuery:            ptr.To("SELECT banner FROM v$version WHERE banner LIKE 'Oracle%'"),
					Username:                 "system",
					Password:                 "oracle",
					Conditions: []monitor.Condition{
						{Variable: "result", Operator: "contains", Value: "Oracle", AndOr: monitor.ConditionAnd},
					},
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"basic_auth_pass":"oracle","basic_auth_user":"system","conditions":[{"type":"expression","variable":"result","operator":"contains","value":"Oracle","andOr":"and"}],"databaseConnectionString":"db.example.com:1521/ORCL","databaseQuery":"SELECT banner FROM v$version WHERE banner LIKE 'Oracle%'","description":null,"id":3,"interval":60,"maxretries":1,"name":"oracledb-conditions","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"oracledb","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oracleMonitor := monitor.OracleDB{}

			err := json.Unmarshal(tc.data, &oracleMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, oracleMonitor)

			data, err := json.Marshal(oracleMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestMonitorOracleDB_String(t *testing.T) {
	tests := []struct {
		name string

		details monitor.OracleDBDetails

		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "query set",
			details: monitor.OracleDBDetails{
				DatabaseConnectionString: "localhost:1521/XEPDB1",
				DatabaseQuery:            ptr.To("SELECT 1 FROM DUAL"),
				Username:                 "oracle",
				Password:                 "secret",
			},
			wantContains: []string{
				`databaseConnectionString: "localhost:1521/XEPDB1"`,
				`databaseQuery: "SELECT 1 FROM DUAL"`,
			},
			wantNotContains: []string{"0x"},
		},
		{
			name: "query nil",
			details: monitor.OracleDBDetails{
				DatabaseConnectionString: "localhost:1521/XEPDB1",
				DatabaseQuery:            nil,
				Username:                 "oracle",
				Password:                 "secret",
			},
			wantContains: []string{
				`databaseConnectionString: "localhost:1521/XEPDB1"`,
				"databaseQuery: <nil>",
			},
			wantNotContains: []string{"0x"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := monitor.OracleDB{OracleDBDetails: tc.details}.String()

			for _, want := range tc.wantContains {
				require.Contains(t, got, want)
			}

			for _, notWant := range tc.wantNotContains {
				require.NotContains(t, got, notWant)
			}
		})
	}
}
