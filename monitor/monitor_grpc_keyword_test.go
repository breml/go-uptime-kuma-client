package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorGrpcKeyword_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.GrpcKeyword
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":2,"name":"grpc-test","description":null,"pathName":"group / grpc-test","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"grpc-keyword","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":"SERVING","invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":"localhost:50051","grpcProtobuf":"syntax = \"proto3\";\n\npackage grpc.health.v1;\n\nservice Health {\n  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);\n}\n\nmessage HealthCheckRequest {\n  string service = 1;\n}\n\nmessage HealthCheckResponse {\n  enum ServingStatus {\n    UNKNOWN = 0;\n    SERVING = 1;\n    NOT_SERVING = 2;\n  }\n  ServingStatus status = 1;\n}\n","grpcMethod":"Check","grpcServiceName":"Health","grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":"{\"service\":\"\"}","grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GrpcKeyword{
				Base: monitor.Base{
					ID:              2,
					Name:            "grpc-test",
					Description:     nil,
					PathName:        "group / grpc-test",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1},
					IsActive:        true,
				},
				GrpcKeywordDetails: monitor.GrpcKeywordDetails{
					GrpcURL:         "localhost:50051",
					GrpcProtobuf:    "syntax = \"proto3\";\n\npackage grpc.health.v1;\n\nservice Health {\n  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);\n}\n\nmessage HealthCheckRequest {\n  string service = 1;\n}\n\nmessage HealthCheckResponse {\n  enum ServingStatus {\n    UNKNOWN = 0;\n    SERVING = 1;\n    NOT_SERVING = 2;\n  }\n  ServingStatus status = 1;\n}\n",
					GrpcServiceName: "Health",
					GrpcMethod:      "Check",
					GrpcEnableTLS:   false,
					GrpcBody:        "{\"service\":\"\"}",
					Keyword:         "SERVING",
					InvertKeyword:   false,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"grpcBody":"{\"service\":\"\"}","grpcEnableTls":false,"grpcMethod":"Check","grpcProtobuf":"syntax = \"proto3\";\n\npackage grpc.health.v1;\n\nservice Health {\n  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);\n}\n\nmessage HealthCheckRequest {\n  string service = 1;\n}\n\nmessage HealthCheckResponse {\n  enum ServingStatus {\n    UNKNOWN = 0;\n    SERVING = 1;\n    NOT_SERVING = 2;\n  }\n  ServingStatus status = 1;\n}\n","grpcServiceName":"Health","grpcUrl":"localhost:50051","id":2,"interval":60,"invertKeyword":false,"keyword":"SERVING","maxretries":2,"name":"grpc-test","notificationIDList":{"1":true},"parent":1,"description":null,"resendInterval":0,"retryInterval":60,"type":"grpc-keyword","upsideDown":false}`,
		},
		{
			name: "with invert keyword",
			data: []byte(`{"id":3,"name":"grpc-inverted","description":"Test inverted keyword","pathName":"grpc-inverted","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"grpc-keyword","timeout":48,"interval":120,"retryInterval":60,"resendInterval":0,"keyword":"NOT_SERVING","invertKeyword":true,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":5,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":"example.com:443","grpcProtobuf":"syntax = \"proto3\";","grpcMethod":"Status","grpcServiceName":"MyService","grpcEnableTls":true,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":"{}","grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.GrpcKeyword{
				Base: monitor.Base{
					ID:              3,
					Name:            "grpc-inverted",
					Description:     strPtr("Test inverted keyword"),
					PathName:        "grpc-inverted",
					Parent:          nil,
					Interval:        120,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      1,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				GrpcKeywordDetails: monitor.GrpcKeywordDetails{
					GrpcURL:         "example.com:443",
					GrpcProtobuf:    "syntax = \"proto3\";",
					GrpcServiceName: "MyService",
					GrpcMethod:      "Status",
					GrpcEnableTLS:   true,
					GrpcBody:        "{}",
					Keyword:         "NOT_SERVING",
					InvertKeyword:   true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test inverted keyword","grpcBody":"{}","grpcEnableTls":true,"grpcMethod":"Status","grpcProtobuf":"syntax = \"proto3\";","grpcServiceName":"MyService","grpcUrl":"example.com:443","id":3,"interval":120,"invertKeyword":true,"keyword":"NOT_SERVING","maxretries":1,"name":"grpc-inverted","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"grpc-keyword","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grpcKeywordMonitor := monitor.GrpcKeyword{}

			err := json.Unmarshal(tc.data, &grpcKeywordMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, grpcKeywordMonitor)

			data, err := json.Marshal(grpcKeywordMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func strPtr(s string) *string {
	return &s
}
