package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorHTTPKeyword_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.HTTPKeyword
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":3,"name":"keyword-test","description":null,"pathName":"group / keyword-test","parent":1,"childrenIDs":[],"url":"https://www.example.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"keyword","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":"success","invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.HTTPKeyword{
				Base: monitor.Base{
					ID:              3,
					Name:            "keyword-test",
					Description:     nil,
					PathName:        "group / keyword-test",
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
					URL:                 "https://www.example.com",
					Timeout:             48,
					ExpiryNotification:  false,
					IgnoreTLS:           false,
					MaxRedirects:        10,
					AcceptedStatusCodes: []string{"200-299"},
					Method:              "GET",
					HTTPBodyEncoding:    "json",
					Body:                "",
					Headers:             "",
					AuthMethod:          monitor.AuthMethodNone,
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
				},
				HTTPKeywordDetails: monitor.HTTPKeywordDetails{
					Keyword:       "success",
					InvertKeyword: false,
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","cacheBust":false,"conditions":[],"description":null,"expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":3,"ignoreTls":false,"interval":60,"invertKeyword":false,"keyword":"success","maxredirects":10,"maxretries":2,"method":"GET","name":"keyword-test","notificationIDList":{"1":true},"oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":1,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"keyword","upsideDown":false,"url":"https://www.example.com"}`,
		},
		{
			name: "success_inverted",
			data: []byte(
				`{"id":4,"name":"keyword-inverted","description":"Testing inverted keyword","pathName":"keyword-inverted","parent":null,"childrenIDs":[],"url":"https://www.example.com","method":"POST","hostname":null,"port":null,"maxretries":3,"weight":2000,"active":true,"forceInactive":false,"type":"keyword","timeout":60,"interval":120,"retryInterval":120,"resendInterval":0,"keyword":"error","invertKeyword":true,"expiryNotification":true,"ignoreTls":true,"upsideDown":false,"packetSize":56,"maxredirects":5,"accepted_statuscodes":["200-299","300-399"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":"basic","grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"xml","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":"X-Custom: test","body":"test data","grpcBody":null,"grpcMetadata":null,"basic_auth_user":"testuser","basic_auth_pass":"testpass","oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`,
			),

			want: monitor.HTTPKeyword{
				Base: monitor.Base{
					ID:              4,
					Name:            "keyword-inverted",
					Description:     ptr.To("Testing inverted keyword"),
					PathName:        "keyword-inverted",
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
					URL:                 "https://www.example.com",
					Timeout:             60,
					ExpiryNotification:  true,
					IgnoreTLS:           true,
					MaxRedirects:        5,
					AcceptedStatusCodes: []string{"200-299", "300-399"},
					Method:              "POST",
					HTTPBodyEncoding:    "xml",
					Body:                "test data",
					Headers:             "X-Custom: test",
					AuthMethod:          monitor.AuthMethodBasic,
					BasicAuthUser:       "testuser",
					BasicAuthPass:       "testpass",
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
				},
				HTTPKeywordDetails: monitor.HTTPKeywordDetails{
					Keyword:       "error",
					InvertKeyword: true,
				},
			},
			wantJSON: `{"accepted_statuscodes":["200-299","300-399"],"active":true,"authDomain":"","authMethod":"basic","authWorkstation":"","basic_auth_pass":"testpass","basic_auth_user":"testuser","body":"test data","cacheBust":false,"conditions":[],"description":"Testing inverted keyword","expiryNotification":true,"headers":"X-Custom: test","httpBodyEncoding":"xml","id":4,"ignoreTls":true,"interval":120,"invertKeyword":true,"keyword":"error","maxredirects":5,"maxretries":3,"method":"POST","name":"keyword-inverted","notificationIDList":{},"oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":null,"proxyId":null,"resendInterval":0,"retryInterval":120,"timeout":60,"tlsCa":"","tlsCert":"","tlsKey":"","type":"keyword","upsideDown":false,"url":"https://www.example.com"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpKeywordMonitor := monitor.HTTPKeyword{}

			err := json.Unmarshal(tc.data, &httpKeywordMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, httpKeywordMonitor)

			data, err := json.Marshal(httpKeywordMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
