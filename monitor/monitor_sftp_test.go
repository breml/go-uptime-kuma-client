package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorSFTP_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.SFTP
		wantJSON string
	}{
		{
			name: "password authentication",
			data: []byte(
				`{"id":4,"name":"sftp-monitor","description":"Test SFTP monitor","pathName":"sftp-monitor","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"sftp.example.com","port":22,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"sftp","timeout":10,"interval":60,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"sftpPath":"/upload","sshAuthMethod":"password","sshUsername":"sftpuser","sshPassword":"sftppass","sshPrivateKey":null,"sshPassphrase":null}`,
			),

			want: monitor.SFTP{
				Base: monitor.Base{
					ID:             4,
					Name:           "sftp-monitor",
					Description:    ptr.To("Test SFTP monitor"),
					PathName:       "sftp-monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				SFTPDetails: monitor.SFTPDetails{
					Hostname:      "sftp.example.com",
					Port:          ptr.To(int64(22)),
					Timeout:       ptr.To(float64(10)),
					SSHUsername:   "sftpuser",
					SSHAuthMethod: monitor.SFTPAuthMethodPassword,
					SSHPassword:   ptr.To("sftppass"),
					SFTPPath:      ptr.To("/upload"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test SFTP monitor","hostname":"sftp.example.com","id":4,"interval":60,"maxretries":2,"name":"sftp-monitor","notificationIDList":{},"parent":null,"port":22,"resendInterval":0,"retryInterval":60,"sftpPath":"/upload","sshAuthMethod":"password","sshPassphrase":null,"sshPassword":"sftppass","sshPrivateKey":null,"sshUsername":"sftpuser","timeout":10,"type":"sftp","upsideDown":false}`,
		},
		{
			name: "private key authentication",
			data: []byte(
				`{"id":5,"name":"sftp-key-monitor","description":"SFTP monitor with key auth","pathName":"group / sftp-key-monitor","parent":1,"hostname":"files.example.com","port":2222,"maxretries":3,"active":true,"type":"sftp","timeout":2.5,"interval":120,"retryInterval":120,"resendInterval":0,"upsideDown":false,"notificationIDList":{"1":true,"2":true},"accepted_statuscodes":["200-299"],"sftpPath":"/srv/incoming","sshAuthMethod":"privateKey","sshUsername":"deploy","sshPassword":null,"sshPrivateKey":"-----BEGIN OPENSSH PRIVATE KEY-----\nabc\n-----END OPENSSH PRIVATE KEY-----","sshPassphrase":"secret"}`,
			),

			want: monitor.SFTP{
				Base: monitor.Base{
					ID:              5,
					Name:            "sftp-key-monitor",
					Description:     ptr.To("SFTP monitor with key auth"),
					PathName:        "group / sftp-key-monitor",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				SFTPDetails: monitor.SFTPDetails{
					Hostname:      "files.example.com",
					Port:          ptr.To(int64(2222)),
					Timeout:       ptr.To(2.5),
					SSHUsername:   "deploy",
					SSHAuthMethod: monitor.SFTPAuthMethodPrivateKey,
					SSHPrivateKey: ptr.To(
						"-----BEGIN OPENSSH PRIVATE KEY-----\nabc\n-----END OPENSSH PRIVATE KEY-----",
					),
					SSHPassphrase: ptr.To("secret"),
					SFTPPath:      ptr.To("/srv/incoming"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"SFTP monitor with key auth","hostname":"files.example.com","id":5,"interval":120,"maxretries":3,"name":"sftp-key-monitor","notificationIDList":{"1":true,"2":true},"parent":1,"port":2222,"resendInterval":0,"retryInterval":120,"sftpPath":"/srv/incoming","sshAuthMethod":"privateKey","sshPassphrase":"secret","sshPassword":null,"sshPrivateKey":"-----BEGIN OPENSSH PRIVATE KEY-----\nabc\n-----END OPENSSH PRIVATE KEY-----","sshUsername":"deploy","timeout":2.5,"type":"sftp","upsideDown":false}`,
		},
		{
			name: "success with unset optional fields",
			data: []byte(
				`{"id":6,"name":"sftp-minimal","description":null,"pathName":"sftp-minimal","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"sftp.example.com","port":null,"maxretries":0,"weight":2000,"active":true,"forceInactive":false,"type":"sftp","timeout":null,"interval":60,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"sftpPath":null,"sshAuthMethod":"password","sshUsername":"sftpuser","sshPassword":null,"sshPrivateKey":null,"sshPassphrase":null}`,
			),

			want: monitor.SFTP{
				Base: monitor.Base{
					ID:             6,
					Name:           "sftp-minimal",
					PathName:       "sftp-minimal",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
				SFTPDetails: monitor.SFTPDetails{
					Hostname:      "sftp.example.com",
					SSHUsername:   "sftpuser",
					SSHAuthMethod: monitor.SFTPAuthMethodPassword,
				},
			},
			// Port and the optional strings stay null, but timeout is
			// substituted because the server column is NOT NULL.
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"hostname":"sftp.example.com","id":6,"interval":60,"maxretries":0,"name":"sftp-minimal","notificationIDList":{},"parent":null,"port":null,"resendInterval":0,"retryInterval":60,"sftpPath":null,"sshAuthMethod":"password","sshPassphrase":null,"sshPassword":null,"sshPrivateKey":null,"sshUsername":"sftpuser","timeout":10,"type":"sftp","upsideDown":false}`,
		},
		{
			name: "payload without sensitive data",
			data: []byte(
				`{"id":7,"name":"sftp-foreign","description":null,"pathName":"sftp-foreign","parent":null,"hostname":"sftp.example.com","port":22,"maxretries":0,"active":true,"type":"sftp","timeout":10,"interval":60,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"conditions":[],"sftpPath":"/upload","sshAuthMethod":"privateKey","includeSensitiveData":false}`,
			),

			// The SSH credentials are only part of the sensitive section of
			// the server representation, so they unmarshal to their zero
			// values for anyone but the owner of the monitor.
			want: monitor.SFTP{
				Base: monitor.Base{
					ID:             7,
					Name:           "sftp-foreign",
					PathName:       "sftp-foreign",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
				SFTPDetails: monitor.SFTPDetails{
					Hostname:      "sftp.example.com",
					Port:          ptr.To(int64(22)),
					Timeout:       ptr.To(float64(10)),
					SSHAuthMethod: monitor.SFTPAuthMethodPrivateKey,
					SFTPPath:      ptr.To("/upload"),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"hostname":"sftp.example.com","id":7,"interval":60,"maxretries":0,"name":"sftp-foreign","notificationIDList":{},"parent":null,"port":22,"resendInterval":0,"retryInterval":60,"sftpPath":"/upload","sshAuthMethod":"privateKey","sshPassphrase":null,"sshPassword":null,"sshPrivateKey":null,"sshUsername":"","timeout":10,"type":"sftp","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sftpMonitor := monitor.SFTP{}

			err := json.Unmarshal(tc.data, &sftpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sftpMonitor)

			data, err := json.Marshal(sftpMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestMonitorSFTP_MarshalDefaultsAuthMethod(t *testing.T) {
	sftpMonitor := monitor.SFTP{
		Base:        monitor.Base{Name: "sftp-without-auth-method"},
		SFTPDetails: monitor.SFTPDetails{Hostname: "sftp.example.com", SSHUsername: "sftpuser"},
	}

	data, err := json.Marshal(sftpMonitor)
	require.NoError(t, err)

	var raw map[string]any

	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	require.Equal(t, string(monitor.SFTPAuthMethodPassword), raw["sshAuthMethod"])
}
