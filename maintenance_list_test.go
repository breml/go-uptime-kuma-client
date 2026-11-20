package kuma_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/maintenance"
)

// TestMaintenanceListKeepsDecodableEntries pins that a maintenance window the
// client cannot decode costs only that window, not the whole maintenanceList:
// the others reach the cache and GetMaintenances names the one left out.
func TestMaintenanceListKeepsDecodableEntries(t *testing.T) {
	t.Parallel()

	fake := &fakeAckServer{messages: make(chan []byte, 32)}
	kumaClient := newFakeAckClient(t, fake)

	// Uptime Kuma accepts any date JavaScript's new Date(...) reads, such as an
	// offset without a colon, and stores it as sent.
	fake.send([]byte(`42["maintenanceList",{` +
		`"1":{"id":1,"title":"good","strategy":"manual","active":true,"dateRange":[]},` +
		`"2":{"id":2,"title":"bad","strategy":"single","active":true,` +
		`"dateRange":["2026-06-29T13:39:00+0200","2026-06-29T14:39:00+0200"]}}]`))

	var (
		maintenances []maintenance.Maintenance
		err          error
	)

	require.Eventually(t, func() bool {
		maintenances, err = kumaClient.GetMaintenances(t.Context())

		return len(maintenances) == 1
	}, 5*time.Second, 10*time.Millisecond, "the decodable window must reach the cache")

	require.ErrorIs(t, err, kuma.ErrMaintenanceDecode)
	require.ErrorContains(t, err, "maintenance 2")
	require.Equal(t, int64(1), maintenances[0].ID)

	fake.send([]byte(`42["maintenanceList",{` +
		`"1":{"id":1,"title":"good","strategy":"manual","active":true,"dateRange":[]},` +
		`"3":{"id":3,"title":"also good","strategy":"manual","active":true,"dateRange":[]}}]`))

	require.Eventually(t, func() bool {
		maintenances, err = kumaClient.GetMaintenances(t.Context())

		return len(maintenances) == 2
	}, 5*time.Second, 10*time.Millisecond, "the next list must replace the cache")

	require.NoError(t, err, "a list that decodes completely clears the earlier error")
}
