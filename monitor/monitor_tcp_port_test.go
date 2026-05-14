package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorTCPPort_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.TCPPort
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":4,"name":"tcp-port-monitor","description":"Test TCP Port monitor","pathName":"group / tcp-port-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"example.com","port":8080,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"port","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"includeSensitiveData":true,"smtpSecurity":null,"expectedTlsAlert":null}`,
			),

			want: monitor.TCPPort{
				Base: monitor.Base{
					ID:              4,
					Name:            "tcp-port-monitor",
					Description:     ptr.To("Test TCP Port monitor"),
					PathName:        "group / tcp-port-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				TCPPortDetails: monitor.TCPPortDetails{
					Hostname: "example.com",
					Port:     8080,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test TCP Port monitor","expectedTlsAlert":null,"expiryNotification":false,"hostname":"example.com","id":4,"interval":60,"maxretries":2,"name":"tcp-port-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":8080,"resendInterval":0,"retryInterval":60,"smtpSecurity":null,"type":"port","upsideDown":false}`,
		},
		{
			name: "with TLS fields",
			data: []byte(
				`{"id":5,"name":"tcp-tls-monitor","description":"TCP monitor with TLS","pathName":"tcp-tls-monitor","parent":null,"hostname":"smtp.example.com","port":465,"maxretries":3,"active":true,"type":"port","interval":120,"retryInterval":60,"resendInterval":0,"upsideDown":false,"smtpSecurity":"secure","expiryNotification":true,"expectedTlsAlert":"certificate_required","notificationIDList":{},"accepted_statuscodes":["200-299"]}`,
			),

			want: monitor.TCPPort{
				Base: monitor.Base{
					ID:             5,
					Name:           "tcp-tls-monitor",
					Description:    ptr.To("TCP monitor with TLS"),
					PathName:       "tcp-tls-monitor",
					Interval:       120,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     3,
					UpsideDown:     false,
					IsActive:       true,
				},
				TCPPortDetails: monitor.TCPPortDetails{
					Hostname:           "smtp.example.com",
					Port:               465,
					SMTPSecurity:       ptr.To("secure"),
					ExpiryNotification: true,
					ExpectedTLSAlert:   ptr.To("certificate_required"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"TCP monitor with TLS","expectedTlsAlert":"certificate_required","expiryNotification":true,"hostname":"smtp.example.com","id":5,"interval":120,"maxretries":3,"name":"tcp-tls-monitor","notificationIDList":{},"parent":null,"port":465,"resendInterval":0,"retryInterval":60,"smtpSecurity":"secure","type":"port","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcpPortMonitor := monitor.TCPPort{}

			err := json.Unmarshal(tc.data, &tcpPortMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, tcpPortMonitor)

			data, err := json.Marshal(tcpPortMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
