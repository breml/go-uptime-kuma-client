package monitor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestIPFamily_String(t *testing.T) {
	tests := []struct {
		family monitor.IPFamily
		want   string
	}{
		{monitor.IPFamilyAuto, ""},
		{monitor.IPFamilyIPv4, "ipv4"},
		{monitor.IPFamilyIPv6, "ipv6"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			require.Equal(t, tc.want, tc.family.String())
		})
	}
}
