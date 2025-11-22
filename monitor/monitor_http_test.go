package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorHTTP_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.HTTP
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":2,"name":"foobar.com","description":null,"pathName":"group / foobar.com","parent":1,"childrenIDs":[],"url":"https://www.foobar.com","method":"GET","hostname":null,"port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"http","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true}`),

			want: monitor.HTTP{
				Base: monitor.Base{
					ID:              2,
					Name:            "foobar.com",
					Description:     nil,
					PathName:        "group / foobar.com",
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
					URL:                 "https://www.foobar.com",
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
			},
			wantJSON: `{"accepted_statuscodes":["200-299"],"active":true,"authDomain":"","authMethod":"","authWorkstation":"","basic_auth_pass":"","basic_auth_user":"","body":"","conditions":[],"description":null,"expiryNotification":false,"headers":"","httpBodyEncoding":"json","id":2,"ignoreTls":false,"interval":60,"maxredirects":10,"maxretries":2,"method":"GET","name":"foobar.com","notificationIDList":{"1":true},"oauth_auth_method":"client_secret_basic","oauth_client_id":"","oauth_client_secret":"","oauth_scopes":"","oauth_token_url":"","parent":1,"proxyId":null,"resendInterval":0,"retryInterval":60,"timeout":48,"tlsCa":"","tlsCert":"","tlsKey":"","type":"http","upsideDown":false,"url":"https://www.foobar.com"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpMonitor := monitor.HTTP{}

			err := json.Unmarshal(tc.data, &httpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, httpMonitor)

			data, err := json.Marshal(httpMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
