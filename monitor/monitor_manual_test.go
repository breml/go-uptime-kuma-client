package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestMonitorManual_Unmarshal(t *testing.T) {
	parent1 := int64(1)

	tests := []struct {
		name string
		data []byte

		want     monitor.Manual
		wantJSON string
	}{
		{
			name: "minimal",
			data: []byte(
				`{"id":7,"name":"manual-monitor","description":null,"pathName":"manual-monitor","parent":null,"childrenIDs":[],"url":null,"method":"GET","hostname":null,"port":null,"maxretries":0,"weight":2000,"active":true,"forceInactive":false,"type":"manual","timeout":48,"interval":60,"retryInterval":60,"resendInterval":0,"keyword":null,"invertKeyword":false,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"includeSensitiveData":true}`,
			),

			want: monitor.Manual{
				Base: monitor.Base{
					ID:             7,
					Name:           "manual-monitor",
					PathName:       "manual-monitor",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"id":7,"interval":60,"manual_status":null,"maxretries":0,"name":"manual-monitor","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"manual","upsideDown":false}`,
		},
		{
			// The server never emits manual_status, these cases only cover the
			// client side round trip of the field.
			name: "with status down",
			data: []byte(
				`{"id":9,"name":"manual-down","description":null,"pathName":"manual-down","parent":null,"childrenIDs":[],"maxretries":0,"weight":2000,"active":true,"forceInactive":false,"type":"manual","interval":60,"retryInterval":60,"resendInterval":0,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{},"tags":[],"maintenance":false,"conditions":[],"manual_status":0}`,
			),

			want: monitor.Manual{
				Base: monitor.Base{
					ID:             9,
					Name:           "manual-down",
					PathName:       "manual-down",
					Interval:       60,
					RetryInterval:  60,
					ResendInterval: 0,
					MaxRetries:     0,
					UpsideDown:     false,
					IsActive:       true,
				},
				ManualDetails: monitor.ManualDetails{
					ManualStatus: new(monitor.ManualStatusDown),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":null,"id":9,"interval":60,"manual_status":0,"maxretries":0,"name":"manual-down","notificationIDList":{},"parent":null,"resendInterval":0,"retryInterval":60,"type":"manual","upsideDown":false}`,
		},
		{
			// The server never emits manual_status, this case only covers the
			// client side round trip of the field.
			name: "with status",
			data: []byte(
				`{"id":8,"name":"manual-up","description":"Manually controlled monitor","pathName":"group / manual-up","parent":1,"childrenIDs":[],"maxretries":2,"weight":2000,"active":true,"forceInactive":false,"type":"manual","interval":120,"retryInterval":30,"resendInterval":5,"upsideDown":false,"accepted_statuscodes":["200-299"],"notificationIDList":{"1":true,"2":true},"tags":[],"maintenance":false,"conditions":[],"manual_status":1}`,
			),

			want: monitor.Manual{
				Base: monitor.Base{
					ID:              8,
					Name:            "manual-up",
					Description:     new("Manually controlled monitor"),
					PathName:        "group / manual-up",
					Parent:          &parent1,
					Interval:        120,
					RetryInterval:   30,
					ResendInterval:  5,
					MaxRetries:      2,
					UpsideDown:      false,
					NotificationIDs: []int64{1, 2},
					IsActive:        true,
				},
				ManualDetails: monitor.ManualDetails{
					ManualStatus: new(monitor.ManualStatusUp),
				},
			},
			wantJSON: `{"accepted_statuscodes":[],"active":true,"conditions":[],"description":"Manually controlled monitor","id":8,"interval":120,"manual_status":1,"maxretries":2,"name":"manual-up","notificationIDList":{"1":true,"2":true},"parent":1,"resendInterval":5,"retryInterval":30,"type":"manual","upsideDown":false}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manualMonitor := monitor.Manual{}

			err := json.Unmarshal(tc.data, &manualMonitor)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, manualMonitor)

			data, err := json.Marshal(manualMonitor)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
