package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorNTP_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.NTP
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(
				`{"id":4,"name":"ntp-monitor","description":"Test NTP monitor","pathName":"ntp-monitor","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"pool.ntp.org","port":123,"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"ntp","timeout":10,"interval":300,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"ntpStratumThreshold":null,"ntpTimeOffsetThreshold":null,"ntpRootDispersionThreshold":null}`,
			),

			want: monitor.NTP{
				Base: monitor.Base{
					ID:             4,
					Name:           "ntp-monitor",
					Description:    ptr.To("Test NTP monitor"),
					PathName:       "ntp-monitor",
					Interval:       300,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     2,
					UpsideDown:     false,
					IsActive:       true,
				},
				NTPDetails: monitor.NTPDetails{
					Hostname: "pool.ntp.org",
					Port:     ptr.To(int64(123)),
					Timeout:  ptr.To(int64(10)),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Test NTP monitor","hostname":"pool.ntp.org","id":4,"interval":300,"maxretries":2,"name":"ntp-monitor","notificationIDList":{},"ntpRootDispersionThreshold":null,"ntpStratumThreshold":null,"ntpTimeOffsetThreshold":null,"parent":null,"port":123,"resendInterval":0,"retryInterval":60,"timeout":10,"type":"ntp","upsideDown":false}`,
		},
		{
			name: "with thresholds",
			data: []byte(
				`{"id":5,"name":"ntp-thresholds-monitor","description":"NTP monitor with thresholds","pathName":"group / ntp-thresholds-monitor","parent":1,"hostname":"time.cloudflare.com","port":1123,"maxretries":3,"active":true,"type":"ntp","timeout":20,"interval":600,"retryInterval":120,"resendInterval":0,"upsideDown":false,"notificationIDList":{"1":true,"2":true},"accepted_statuscodes":["200-299"],"ntpStratumThreshold":3,"ntpTimeOffsetThreshold":250,"ntpRootDispersionThreshold":100}`,
			),

			want: monitor.NTP{
				Base: monitor.Base{
					ID:              5,
					Name:            "ntp-thresholds-monitor",
					Description:     ptr.To("NTP monitor with thresholds"),
					PathName:        "group / ntp-thresholds-monitor",
					Parent:          &parent1,
					Interval:        600,
					RetryInterval:   120,
					ResendInterval:  0,
					MaxRetries:      3,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				NTPDetails: monitor.NTPDetails{
					Hostname:                   "time.cloudflare.com",
					Port:                       ptr.To(int64(1123)),
					Timeout:                    ptr.To(int64(20)),
					NTPStratumThreshold:        ptr.To(int64(3)),
					NTPTimeOffsetThreshold:     ptr.To(int64(250)),
					NTPRootDispersionThreshold: ptr.To(int64(100)),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"NTP monitor with thresholds","hostname":"time.cloudflare.com","id":5,"interval":600,"maxretries":3,"name":"ntp-thresholds-monitor","notificationIDList":{"1":true,"2":true},"ntpRootDispersionThreshold":100,"ntpStratumThreshold":3,"ntpTimeOffsetThreshold":250,"parent":1,"port":1123,"resendInterval":0,"retryInterval":120,"timeout":20,"type":"ntp","upsideDown":false}`,
		},
		{
			name: "success with unset optional fields",
			data: []byte(
				`{"id":6,"name":"ntp-minimal","description":null,"pathName":"ntp-minimal","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":"time.example.com","port":null,"maxretries":0,"weight":2000,"active":true,"forceInactive":false,"type":"ntp","timeout":null,"interval":300,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"ntpStratumThreshold":null,"ntpTimeOffsetThreshold":null,"ntpRootDispersionThreshold":null}`,
			),

			want: monitor.NTP{
				Base: monitor.Base{
					ID:             6,
					Name:           "ntp-minimal",
					PathName:       "ntp-minimal",
					Interval:       300,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
				NTPDetails: monitor.NTPDetails{
					Hostname: "time.example.com",
				},
			},
			// Port and the thresholds stay null, but timeout is substituted
			// because the server column is NOT NULL.
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"hostname":"time.example.com","id":6,"interval":300,"maxretries":0,"name":"ntp-minimal","notificationIDList":{},"ntpRootDispersionThreshold":null,"ntpStratumThreshold":null,"ntpTimeOffsetThreshold":null,"parent":null,"port":null,"resendInterval":0,"retryInterval":60,"timeout":10,"type":"ntp","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ntpMonitor := monitor.NTP{}

			err := json.Unmarshal(tc.data, &ntpMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, ntpMonitor)

			data, err := json.Marshal(ntpMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}

func TestMonitorNTP_MarshalRequiresHostname(t *testing.T) {
	ntpMonitor := monitor.NTP{
		Base: monitor.Base{Name: "ntp-without-hostname"},
	}

	_, err := json.Marshal(ntpMonitor)
	require.Error(t, err)
	require.ErrorContains(t, err, "hostname is required")
}
