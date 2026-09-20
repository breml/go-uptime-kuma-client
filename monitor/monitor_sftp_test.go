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

		want           monitor.SFTP
		wantJSON       string
		wantMarshalErr string
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
			// the server representation. No socket path of the server sends
			// such a payload to a client -- includeSensitiveData defaults to
			// true and only notification templating passes false -- but were
			// one ever unmarshalled, the credentials come back as their zero
			// values.
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
			// Marshalling such a monitor back is refused: the write would
			// replace every credential column with the zero values above.
			wantMarshalErr: "ssh username is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sftpMonitor := monitor.SFTP{}

			err := json.Unmarshal(tc.data, &sftpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, sftpMonitor)

			data, err := json.Marshal(sftpMonitor)

			if tc.wantMarshalErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantMarshalErr)

				return
			}

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

func TestMonitorSFTP_MarshalRequiresHostname(t *testing.T) {
	sftpMonitor := monitor.SFTP{
		Base:        monitor.Base{Name: "sftp-without-hostname"},
		SFTPDetails: monitor.SFTPDetails{SSHUsername: "sftpuser"},
	}

	_, err := json.Marshal(sftpMonitor)
	require.Error(t, err)
	require.ErrorContains(t, err, "hostname is required")
}

func TestMonitorSFTP_MarshalRequiresSSHUsername(t *testing.T) {
	sftpMonitor := monitor.SFTP{
		Base:        monitor.Base{Name: "sftp-without-ssh-username"},
		SFTPDetails: monitor.SFTPDetails{Hostname: "sftp.example.com"},
	}

	_, err := json.Marshal(sftpMonitor)
	require.Error(t, err)
	require.ErrorContains(t, err, "ssh username is required")
}

func TestMonitorSFTP_MarshalRejectsUnknownAuthMethod(t *testing.T) {
	tests := []struct {
		name       string
		authMethod monitor.SFTPAuthMethod
	}{
		{
			name:       "misspelled private key",
			authMethod: "privatekey",
		},
		{
			name:       "unknown method",
			authMethod: "keyboard-interactive",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sftpMonitor := monitor.SFTP{
				Base: monitor.Base{Name: "sftp-with-unknown-auth-method"},
				SFTPDetails: monitor.SFTPDetails{
					Hostname:      "sftp.example.com",
					SSHUsername:   "sftpuser",
					SSHAuthMethod: tc.authMethod,
				},
			}

			_, err := json.Marshal(sftpMonitor)
			require.Error(t, err)
			require.ErrorContains(t, err, "invalid ssh auth method")
		})
	}
}

func TestMonitorSFTP_String(t *testing.T) {
	tests := []struct {
		name string

		details monitor.SFTPDetails

		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "all optional fields set",
			details: monitor.SFTPDetails{
				Hostname:      "sftp.example.com",
				Port:          ptr.To(int64(22)),
				Timeout:       ptr.To(float64(10)),
				SSHUsername:   "sftpuser",
				SSHAuthMethod: monitor.SFTPAuthMethodPrivateKey,
				SSHPassword:   ptr.To("sftppass"),
				SSHPrivateKey: ptr.To("-----BEGIN OPENSSH PRIVATE KEY-----\nabc"),
				SSHPassphrase: ptr.To("secret"),
				SFTPPath:      ptr.To("/upload"),
			},
			wantContains: []string{
				`hostname: "sftp.example.com"`,
				"port: 22",
				"timeout: 10",
				`sshUsername: "sftpuser"`,
				`sshAuthMethod: "privateKey"`,
				`sshPassword: "***"`,
				`sshPrivateKey: "***"`,
				`sshPassphrase: "***"`,
				`sftpPath: "/upload"`,
			},
			wantNotContains: []string{
				"0x",
				"sftppass",
				"BEGIN OPENSSH PRIVATE KEY",
				"secret",
			},
		},
		{
			name: "all optional fields nil",
			details: monitor.SFTPDetails{
				Hostname:      "sftp.example.com",
				SSHUsername:   "sftpuser",
				SSHAuthMethod: monitor.SFTPAuthMethodPassword,
			},
			wantContains: []string{
				`hostname: "sftp.example.com"`,
				"port: <nil>",
				"timeout: <nil>",
				"sshPassword: <nil>",
				"sshPrivateKey: <nil>",
				"sshPassphrase: <nil>",
				"sftpPath: <nil>",
			},
			wantNotContains: []string{"0x"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := monitor.SFTP{SFTPDetails: tc.details}.String()

			for _, want := range tc.wantContains {
				require.Contains(t, got, want)
			}

			for _, notWant := range tc.wantNotContains {
				require.NotContains(t, got, notWant)
			}
		})
	}
}
