package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet2FAStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := client.Get2FAStatus(ctx)
	require.NoError(t, err)
	assert.False(t, status, "2FA should be disabled initially")
}

func TestPrepare2FA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := client.Prepare2FA(ctx)
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.NotEmpty(t, config.Secret, "Secret should not be empty")
	assert.NotEmpty(t, config.URI, "URI should not be empty")
	assert.Contains(t, config.URI, "otpauth://", "URI should be an otpauth URI")
}

func Test2FAFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 1: Verify 2FA is disabled
	status, err := client.Get2FAStatus(ctx)
	require.NoError(t, err)
	assert.False(t, status, "2FA should be disabled initially")

	// Step 2: Prepare 2FA
	config, err := client.Prepare2FA(ctx)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Step 3: Generate TOTP token from secret
	token, err := totp.GenerateCode(config.Secret, time.Now())
	require.NoError(t, err)

	// Step 4: Save/enable 2FA
	err = client.Save2FA(ctx, config.Secret, token)
	require.NoError(t, err)

	// Step 5: Verify 2FA is now enabled
	status, err = client.Get2FAStatus(ctx)
	require.NoError(t, err)
	assert.True(t, status, "2FA should be enabled after Save2FA")

	// Step 6: Verify token
	token, err = totp.GenerateCode(config.Secret, time.Now())
	require.NoError(t, err)
	valid, err := client.Verify2FAToken(ctx, token)
	require.NoError(t, err)
	assert.True(t, valid, "Generated token should be valid")

	// Step 7: Disable 2FA (cleanup)
	err = client.Disable2FA(ctx, "admin1")
	require.NoError(t, err)

	// Step 8: Verify 2FA is disabled again
	status, err = client.Get2FAStatus(ctx)
	require.NoError(t, err)
	assert.False(t, status, "2FA should be disabled after Disable2FA")
}

func TestChangePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Change password
	err := client.ChangePassword(ctx, "admin1", "newpassword123")
	require.NoError(t, err)

	// Change it back
	err = client.ChangePassword(ctx, "newpassword123", "admin1")
	require.NoError(t, err)
}

func TestLogout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logout should succeed
	err := client.Logout(ctx)
	require.NoError(t, err)

	// After logout, operations should fail
	// Note: We can't test this easily because we're using a shared client
	// In a real scenario, subsequent operations would fail after logout
}
