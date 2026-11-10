package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

// TestSystemServiceNameFromConfig covers the read side of the name a pm2 or
// system-service monitor carries: a 2.x server sends system_service_name, a
// 3.x server config.systemServiceName.
func TestSystemServiceNameFromConfig(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "2.x top-level field",
			data: `{"system_service_name":"nginx.service"}`,
			want: "nginx.service",
		},
		{
			name: "3.x config",
			data: `{"config":{"systemServiceName":"nginx.service"}}`,
			want: "nginx.service",
		},
		{
			name: "config wins over the top-level field",
			data: `{"system_service_name":"old.service","config":{"systemServiceName":"nginx.service"}}`,
			want: "nginx.service",
		},
		{
			name: "config without the key falls back to the top-level field",
			data: `{"system_service_name":"nginx.service","config":{}}`,
			want: "nginx.service",
		},
		{
			name: "neither",
			data: `{"config":{}}`,
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var service monitor.SystemService

			err := json.Unmarshal([]byte(tc.data), &service)
			require.NoError(t, err)
			require.Equal(t, tc.want, service.SystemServiceName)

			var pm2 monitor.PM2

			err = json.Unmarshal([]byte(tc.data), &pm2)
			require.NoError(t, err)
			require.Equal(t, tc.want, pm2.ProcessName)
		})
	}
}

// TestSystemServiceNameFromConfigViaBase covers the path GetMonitorAs takes: a
// 3.x payload decoded into monitor.Base and converted with As.
func TestSystemServiceNameFromConfigViaBase(t *testing.T) {
	data := []byte(`{"id":7,"name":"pm2","type":"pm2","interval":60,"active":true,` +
		`"notificationIDList":{},"config":{"systemServiceName":"worker 1"}}`)

	var base monitor.Base

	err := json.Unmarshal(data, &base)
	require.NoError(t, err)

	var pm2 monitor.PM2

	err = base.As(&pm2)
	require.NoError(t, err)
	require.Equal(t, int64(7), pm2.ID)
	require.Equal(t, "worker 1", pm2.ProcessName)
}
