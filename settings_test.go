package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSettingsGetAndSet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	t.Run("get_default_settings", func(t *testing.T) {
		s, err := client.GetSettings(ctx)
		require.NoError(t, err)
		require.NotNil(t, s)

		require.NotEmpty(t, s.ServerTimezone)
	})

	t.Run("update_timezone", func(t *testing.T) {
		original, err := client.GetSettings(ctx)
		require.NoError(t, err)

		modified := *original
		modified.ServerTimezone = "Europe/Berlin"

		err = client.SetSettings(ctx, modified, "")
		require.NoError(t, err)

		updated, err := client.GetSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, "Europe/Berlin", updated.ServerTimezone)

		// Restore original
		err = client.SetSettings(ctx, *original, "")
		require.NoError(t, err)
	})

	t.Run("update_keep_data_period_days", func(t *testing.T) {
		original, err := client.GetSettings(ctx)
		require.NoError(t, err)

		modified := *original
		modified.KeepDataPeriodDays = 7

		err = client.SetSettings(ctx, modified, "")
		require.NoError(t, err)

		updated, err := client.GetSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, 7, updated.KeepDataPeriodDays)

		// Restore original
		err = client.SetSettings(ctx, *original, "")
		require.NoError(t, err)
	})

	t.Run("update_tls_expiry_notify_days", func(t *testing.T) {
		original, err := client.GetSettings(ctx)
		require.NoError(t, err)

		modified := *original
		modified.TLSExpiryNotifyDays = []int{3, 7, 30}

		err = client.SetSettings(ctx, modified, "")
		require.NoError(t, err)

		updated, err := client.GetSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, []int{3, 7, 30}, updated.TLSExpiryNotifyDays)

		// Restore original
		err = client.SetSettings(ctx, *original, "")
		require.NoError(t, err)
	})

	t.Run("update_multiple_settings", func(t *testing.T) {
		original, err := client.GetSettings(ctx)
		require.NoError(t, err)

		modified := *original
		modified.SearchEngineIndex = !original.SearchEngineIndex
		modified.TrustProxy = !original.TrustProxy
		modified.EntryPage = "dashboard"

		err = client.SetSettings(ctx, modified, "")
		require.NoError(t, err)

		updated, err := client.GetSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, modified.SearchEngineIndex, updated.SearchEngineIndex)
		require.Equal(t, modified.TrustProxy, updated.TrustProxy)
		require.Equal(t, "dashboard", updated.EntryPage)

		// Restore original
		err = client.SetSettings(ctx, *original, "")
		require.NoError(t, err)
	})
}
