package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorWebsocketUpgrade_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.WebsocketUpgrade
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":20,"name":"ws-test","description":null,"pathName":"group / ws-test","parent":1,"childrenIDs":[],"url":"wss://echo.example.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"websocket-upgrade","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["1000"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":"","kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.WebsocketUpgrade{
				Base: monitor.Base{
					ID:              20,
					Name:            "ws-test",
					Description:     nil,
					PathName:        "group / ws-test",
					Parent:          &parent1,
					ProxyID:         nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1},
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "wss://echo.example.com",
					Timeout:             48,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"1000"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             "",
					AuthMethod:          monitor.AuthMethodNone,
					BearerToken:         "",
					BasicAuthUser:       "",
					BasicAuthPass:       "",
					AuthDomain:          "",
					AuthWorkstation:     "",
					TLSCert:             "",
					TLSKey:              "",
					TLSCa:               "",
					OAuthAuthMethod:     "client_secret_basic",
					OAuthTokenURL:       "",
					OAuthClientID:       "",
					OAuthClientSecret:   "",
					OAuthScopes:         "",
					OAuthAudience:       "",
				},
				WebsocketUpgradeDetails: monitor.WebsocketUpgradeDetails{
					IgnoreSecWebsocketAcceptHeader: false,
					Subprotocol:                    "",
				},
			},
			wantJSON: `{"accepted_statuscodes":["1000"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","bearer_token":"","body":"","cacheBust":false,"conditions":[],"description":null,"expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":20,"ignoreTls":false,"interval":60,"maxredirects":10,"maxretries":2,"method":"GET","name":"ws-test","notificationIDList":{"1":true},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":1,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"websocket-upgrade","upsideDown":false,"url":"wss://echo.example.com","wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":""}`,
		},
		{
			name: "success_basic_auth",
			data: []byte(
				`{"id":21,"name":"ws-basic","description":"WebSocket with basic auth","pathName":"ws-basic","parent":null,"childrenIDs":[],"url":"wss://secure.example.com/ws","method":"GET","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"websocket-upgrade","timeout":30,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":true,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["1000","1001"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":"basic","bearer_token":"ws-bearer-token","grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":"{\"X-Custom\":\"value\"}","body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":"wsuser","basic_auth_pass":"wspass","oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":"","kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.WebsocketUpgrade{
				Base: monitor.Base{
					ID:              21,
					Name:            "ws-basic",
					Description:     ptr.To("WebSocket with basic auth"),
					PathName:        "ws-basic",
					Parent:          nil,
					ProxyID:         nil,
					Interval:        120,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "wss://secure.example.com/ws",
					Timeout:             30,
					ExpiryNotification:  false,
					IgnoreTLS:           true,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"1000", "1001"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             `{"X-Custom":"value"}`,
					AuthMethod:          monitor.AuthMethodBasic,
					BearerToken:         "ws-bearer-token",
					BasicAuthUser:       "wsuser",
					BasicAuthPass:       "wspass",
					AuthDomain:          "",
					AuthWorkstation:     "",
					TLSCert:             "",
					TLSKey:              "",
					TLSCa:               "",
					OAuthAuthMethod:     "client_secret_basic",
					OAuthTokenURL:       "",
					OAuthClientID:       "",
					OAuthClientSecret:   "",
					OAuthScopes:         "",
					OAuthAudience:       "",
				},
				WebsocketUpgradeDetails: monitor.WebsocketUpgradeDetails{
					IgnoreSecWebsocketAcceptHeader: false,
					Subprotocol:                    "",
				},
			},
			wantJSON: `{"accepted_statuscodes":["1000","1001"],"active":true,"authDomain":"","authMethod":"basic","authWorkstation":"","basic_auth_pass":"wspass","basic_auth_user":"wsuser","bearer_token":"ws-bearer-token","body":"","cacheBust":false,"conditions":[],"description":"WebSocket with basic auth","expiryNotification":false,"headers":"{\"X-Custom\":\"value\"}","httpBodyEncoding":"json","id":21,"ignoreTls":true,"interval":120,"maxredirects":10,"maxretries":3,"method":"GET","name":"ws-basic","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"proxyId":null,"resendInterval":0,"retryInterval":120,"timeout":30,"tlsCa":"","tlsCert":"","tlsKey":"","type":"websocket-upgrade","upsideDown":false,"url":"wss://secure.example.com/ws","wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":""}`,
		},
		{
			name: "success_oauth2_cc",
			data: []byte(
				`{"id":22,"name":"ws-oauth","description":null,"pathName":"ws-oauth","parent":null,"childrenIDs":[],"url":"wss://api.example.com/ws","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"websocket-upgrade","timeout":20,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["1000"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":"oauth2-cc","grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":"my-client-id","oauth_client_secret":"my-client-secret","oauth_token_url":"https://auth.example.com/token","oauth_scopes":"read write","oauth_audience":"https://api.example.com","oauth_auth_method":"client_secret_post","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":"","kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.WebsocketUpgrade{
				Base: monitor.Base{
					ID:              22,
					Name:            "ws-oauth",
					Description:     nil,
					PathName:        "ws-oauth",
					Parent:          nil,
					ProxyID:         nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "wss://api.example.com/ws",
					Timeout:             20,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"1000"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             "",
					AuthMethod:          monitor.AuthMethodOAtuh2CC,
					BearerToken:         "",
					BasicAuthUser:       "",
					BasicAuthPass:       "",
					AuthDomain:          "",
					AuthWorkstation:     "",
					TLSCert:             "",
					TLSKey:              "",
					TLSCa:               "",
					OAuthAuthMethod:     "client_secret_post",
					OAuthTokenURL:       "https://auth.example.com/token",
					OAuthClientID:       "my-client-id",
					OAuthClientSecret:   "my-client-secret",
					OAuthScopes:         "read write",
					OAuthAudience:       "https://api.example.com",
				},
				WebsocketUpgradeDetails: monitor.WebsocketUpgradeDetails{
					IgnoreSecWebsocketAcceptHeader: false,
					Subprotocol:                    "",
				},
			},
			wantJSON: `{"accepted_statuscodes":["1000"],"active":true,"authDomain":"","authMethod":"oauth2-cc","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","bearer_token":"","body":"","cacheBust":false,"conditions":[],"description":null,"expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":22,"ignoreTls":false,"interval":60,"maxredirects":10,"maxretries":2,"method":"GET","name":"ws-oauth","notificationIDList":{},"oauth_audience":"https://api.example.com","oauth_auth_method":"client_secret_post","oauth_client_id":"my-client-id","oauth_client_secret":"my-client-secret","oauth_scopes":"read write","oauth_token_url":"https://auth.example.com/token","parent":null,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":20,"tlsCa":"","tlsCert":"","tlsKey":"","type":"websocket-upgrade","upsideDown":false,"url":"wss://api.example.com/ws","wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":""}`,
		},
		{
			name: "success_mtls",
			data: []byte(
				`{"id":23,"name":"ws-mtls","description":"mTLS WebSocket","pathName":"ws-mtls","parent":null,"childrenIDs":[],"url":"wss://secure.example.com/ws","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"websocket-upgrade","timeout":20,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["1000"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":"mtls","grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":"-----BEGIN CERTIFICATE-----\nca-cert\n-----END CERTIFICATE-----","tlsCert":"-----BEGIN CERTIFICATE-----\nclient-cert\n-----END CERTIFICATE-----","tlsKey":"-----BEGIN PRIVATE KEY-----\nclient-key\n-----END PRIVATE KEY-----","wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":"","kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.WebsocketUpgrade{
				Base: monitor.Base{
					ID:              23,
					Name:            "ws-mtls",
					Description:     ptr.To("mTLS WebSocket"),
					PathName:        "ws-mtls",
					Parent:          nil,
					ProxyID:         nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "wss://secure.example.com/ws",
					Timeout:             20,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"1000"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             "",
					AuthMethod:          monitor.AuthMethodMTLS,
					BearerToken:         "",
					BasicAuthUser:       "",
					BasicAuthPass:       "",
					AuthDomain:          "",
					AuthWorkstation:     "",
					TLSCert:             "-----BEGIN CERTIFICATE-----\nclient-cert\n-----END CERTIFICATE-----",
					TLSKey:              "-----BEGIN PRIVATE KEY-----\nclient-key\n-----END PRIVATE KEY-----",
					TLSCa:               "-----BEGIN CERTIFICATE-----\nca-cert\n-----END CERTIFICATE-----",
					OAuthAuthMethod:     "client_secret_basic",
					OAuthTokenURL:       "",
					OAuthClientID:       "",
					OAuthClientSecret:   "",
					OAuthScopes:         "",
					OAuthAudience:       "",
				},
				WebsocketUpgradeDetails: monitor.WebsocketUpgradeDetails{
					IgnoreSecWebsocketAcceptHeader: false,
					Subprotocol:                    "",
				},
			},
			wantJSON: `{"accepted_statuscodes":["1000"],"active":true,"authDomain":"","authMethod":"mtls","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","bearer_token":"","body":"","cacheBust":false,"conditions":[],"description":"mTLS WebSocket","expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":23,"ignoreTls":false,"interval":60,"maxredirects":10,"maxretries":2,"method":"GET","name":"ws-mtls","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":20,"tlsCa":"-----BEGIN CERTIFICATE-----\nca-cert\n-----END CERTIFICATE-----","tlsCert":"-----BEGIN CERTIFICATE-----\nclient-cert\n-----END CERTIFICATE-----","tlsKey":"-----BEGIN PRIVATE KEY-----\nclient-key\n-----END PRIVATE KEY-----","type":"websocket-upgrade","upsideDown":false,"url":"wss://secure.example.com/ws","wsIgnoreSecWebsocketAcceptHeader":false,"wsSubprotocol":""}`,
		},
		{
			name: "success_ws_specific",
			data: []byte(
				`{"id":24,"name":"ws-specific","description":null,"pathName":"ws-specific","parent":null,"childrenIDs":[],"url":"wss://ws.example.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"websocket-upgrade","timeout":20,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["1000"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_audience":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"wsIgnoreSecWebsocketAcceptHeader":true,"wsSubprotocol":"chat,superchat","kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.WebsocketUpgrade{
				Base: monitor.Base{
					ID:              24,
					Name:            "ws-specific",
					Description:     nil,
					PathName:        "ws-specific",
					Parent:          nil,
					ProxyID:         nil,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				HTTPDetails: monitor.HTTPDetails{
					URL:                 "wss://ws.example.com",
					Timeout:             20,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"1000"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             "",
					AuthMethod:          monitor.AuthMethodNone,
					BearerToken:         "",
					BasicAuthUser:       "",
					BasicAuthPass:       "",
					AuthDomain:          "",
					AuthWorkstation:     "",
					TLSCert:             "",
					TLSKey:              "",
					TLSCa:               "",
					OAuthAuthMethod:     "client_secret_basic",
					OAuthTokenURL:       "",
					OAuthClientID:       "",
					OAuthClientSecret:   "",
					OAuthScopes:         "",
					OAuthAudience:       "",
				},
				WebsocketUpgradeDetails: monitor.WebsocketUpgradeDetails{
					IgnoreSecWebsocketAcceptHeader: true,
					Subprotocol:                    "chat,superchat",
				},
			},
			wantJSON: `{"accepted_statuscodes":["1000"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","bearer_token":"","body":"","cacheBust":false,"conditions":[],"description":null,"expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":24,"ignoreTls":false,"interval":60,"maxredirects":10,"maxretries":2,"method":"GET","name":"ws-specific","notificationIDList":{},"oauth_audience":"","oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":20,"tlsCa":"","tlsCert":"","tlsKey":"","type":"websocket-upgrade","upsideDown":false,"url":"wss://ws.example.com","wsIgnoreSecWebsocketAcceptHeader":true,"wsSubprotocol":"chat,superchat"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wsMonitor := monitor.WebsocketUpgrade{}

			err := json.Unmarshal(tc.data, &wsMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, wsMonitor)

			data, err := json.Marshal(wsMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
