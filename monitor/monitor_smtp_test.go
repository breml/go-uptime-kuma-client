package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSMTP_Unmarshal(t *testing.T) {
	parent1 := int64(1)
	port25 := int64(25)
	portSecure := int64(465)
	portSTARTTLS := int64(587)

	tests := []struct {
		name string
		data []byte

		want     monitor.SMTP
		wantJSON string
	}{
		{
			name: "success with default port and no security",
			data: []byte(
				`{"id":7,"name":"smtp-monitor","description":"Test SMTP monitor","pathName":"group / smtp-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mail.example.com","port":null,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"smtp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"smtpSecurity":null,"includeSensitiveData":true}`,
			),

			want: monitor.SMTP{
				Base: monitor.Base{
					ID:              7,
					Name:            "smtp-monitor",
					Description:     ptr.To("Test SMTP monitor"),
					PathName:        "group / smtp-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SMTPDetails: monitor.SMTPDetails{
					Hostname:     "mail.example.com",
					Port:         nil,
					SMTPSecurity: nil,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SMTP monitor","hostname":"mail.example.com","id":7,"interval":60,"maxretries":2,"name":"smtp-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":null,"resendInterval":0,"retryInterval":60,"smtpSecurity":null,"type":"smtp","upsideDown":false}`,
		},
		{
			name: "success with secure port and secure security",
			data: []byte(
				`{"id":8,"name":"smtp-secure-monitor","description":"Test SMTP secure monitor","pathName":"group / smtp-secure-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mail.example.com","port":465,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"smtp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"smtpSecurity":"secure","includeSensitiveData":true}`,
			),

			want: monitor.SMTP{
				Base: monitor.Base{
					ID:              8,
					Name:            "smtp-secure-monitor",
					Description:     ptr.To("Test SMTP secure monitor"),
					PathName:        "group / smtp-secure-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1},
					IsActive:        true,
				},
				SMTPDetails: monitor.SMTPDetails{
					Hostname:     "mail.example.com",
					Port:         &portSecure,
					SMTPSecurity: ptr.To("secure"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SMTP secure monitor","hostname":"mail.example.com","id":8,"interval":60,"maxretries":2,"name":"smtp-secure-monitor","notificationIDList":{"1":true},"parent":1,"port":465,"resendInterval":0,"retryInterval":60,"smtpSecurity":"secure","type":"smtp","upsideDown":false}`,
		},
		{
			name: "success with starttls security",
			data: []byte(
				`{"id":9,"name":"smtp-starttls-monitor","description":"Test SMTP STARTTLS monitor","pathName":"group / smtp-starttls-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mail.example.com","port":587,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"smtp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"smtpSecurity":"starttls","includeSensitiveData":true}`,
			),

			want: monitor.SMTP{
				Base: monitor.Base{
					ID:              9,
					Name:            "smtp-starttls-monitor",
					Description:     ptr.To("Test SMTP STARTTLS monitor"),
					PathName:        "group / smtp-starttls-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: nil,
					IsActive:        true,
				},
				SMTPDetails: monitor.SMTPDetails{
					Hostname:     "mail.example.com",
					Port:         &portSTARTTLS,
					SMTPSecurity: ptr.To("starttls"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SMTP STARTTLS monitor","hostname":"mail.example.com","id":9,"interval":60,"maxretries":2,"name":"smtp-starttls-monitor","notificationIDList":{},"parent":1,"port":587,"resendInterval":0,"retryInterval":60,"smtpSecurity":"starttls","type":"smtp","upsideDown":false}`,
		},
		{
			name: "success with nostarttls security",
			data: []byte(
				`{"id":10,"name":"smtp-nostarttls-monitor","description":"Test SMTP no STARTTLS monitor","pathName":"group / smtp-nostarttls-monitor","parent":1,"childrenIDs":[],"url":null,"method":"GET","hostname":"mail.example.com","port":25,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"smtp","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"gamedigGivenPortOnly":true,"httpBodyEncoding":"json","jsonPath":null,"expectedValue":null,"kafkaProducerTopic":null,"kafkaProducerBrokers":[],"kafkaProducerSsl":false,"kafkaProducerAllowAutoTopicCreation":false,"kafkaProducerMessage":null,"screenshot":null,"headers":null,"body":null,"grpcBody":null,"grpcMetadata":null,"basic_auth_user":null,"basic_auth_pass":null,"oauth_client_id":null,"oauth_client_secret":null,"oauth_token_url":null,"oauth_scopes":null,"oauth_auth_method":"client_secret_basic","pushToken":null,"databaseConnectionString":null,"radiusUsername":null,"radiusPassword":null,"radiusSecret":null,"mqttUsername":"","mqttPassword":"","authWorkstation":null,"authDomain":null,"tlsCa":null,"tlsCert":null,"tlsKey":null,"kafkaProducerSaslOptions":{"mechanism":"None"},"smtpSecurity":"nostarttls","includeSensitiveData":true}`,
			),

			want: monitor.SMTP{
				Base: monitor.Base{
					ID:              10,
					Name:            "smtp-nostarttls-monitor",
					Description:     ptr.To("Test SMTP no STARTTLS monitor"),
					PathName:        "group / smtp-nostarttls-monitor",
					Parent:          &parent1,
					Interval:        60,
					RetryInterval:   60,
					ResendInterval:  0,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SMTPDetails: monitor.SMTPDetails{
					Hostname:     "mail.example.com",
					Port:         &port25,
					SMTPSecurity: ptr.To("nostarttls"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SMTP no STARTTLS monitor","hostname":"mail.example.com","id":10,"interval":60,"maxretries":2,"name":"smtp-nostarttls-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":25,"resendInterval":0,"retryInterval":60,"smtpSecurity":"nostarttls","type":"smtp","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smtpMonitor := monitor.SMTP{}

			err := json.Unmarshal(tc.data, &smtpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, smtpMonitor)

			data, err := json.Marshal(smtpMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
