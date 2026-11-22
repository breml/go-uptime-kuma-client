package kuma

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maldikhan/go.socket.io/socket.io/v5/client/emit"
)

// The two-factor management commands and the code generation are internal:
// they administer an account rather than authenticate a client, and there is
// no request for them as public API yet. The end-to-end tests need them to put
// a real server into the state they exercise, and those live in the external
// test package, so they are exported for the tests alone.
//
//nolint:gochecknoglobals // the standard way to reach internals from an external test package.
var (
	Enable2FA   = (*Client).enable2FA
	Confirm2FA  = (*Client).confirm2FA
	Disable2FA  = (*Client).disable2FA
	TwoFAStatus = (*Client).twoFAStatus

	NormalizeTOTPSecret = normalizeTOTPSecret
	TOTPCodeAt          = totpCodeAt
)

// SetSetupEventGrace widens the window setupPending gives a setup event that
// has not arrived yet, and returns the function that puts it back.
//
// The window settles a race whose timing a fake server cannot hit reliably, so
// a test that wants the late-setup path takes the timing out of the assertion
// rather than tuning a delay against it.
func SetSetupEventGrace(d time.Duration) func() {
	previous := setupEventGrace
	setupEventGrace = d

	return func() { setupEventGrace = previous }
}

// SetDefaultOperationTimeout shortens the budget New applies when the caller
// configures none, and returns the function that puts it back.
//
// It exists so a test can exercise the default itself rather than a value it
// passed in: without it, removing the default from New would leave every
// command unbounded and the suite green.
func SetDefaultOperationTimeout(d time.Duration) func() {
	previous := defaultOperationTimeout
	defaultOperationTimeout = d

	return func() { defaultOperationTimeout = previous }
}

// AdaptMonitorPayload runs the payload rewrite of CreateMonitor and
// UpdateMonitor for a client that authenticated with scheme.
func AdaptMonitorPayload(scheme AuthScheme, monitorData map[string]any) {
	(&Client{scheme: scheme}).adaptMonitorPayload(monitorData)
}

// ModuleVersion is the version lookup the default User-Agent makes, for a test
// that hands it build information of its own.
//
//nolint:gochecknoglobals // the standard way to reach internals from an external test package.
var (
	ModuleVersion = moduleVersion
	AddAPIKey     = (*Client).addAPIKey
	SetAPIEnabled = (*Client).setAPIEnabled
)

// addAPIKey, exported as AddAPIKey, creates an API key for the logged in account and returns it, the
// way the settings page of an Uptime Kuma 3.x server does. It is for the
// end-to-end tests of WithAPIKey, which need a key the client cannot create
// for itself: the server wants a login session for the command.
func (c *Client) addAPIKey(ctx context.Context, name string) (string, error) {
	release, err := c.lockConn(ctx)
	if err != nil {
		return "", err
	}

	defer release()

	res := make(chan map[string]any, 1)

	err = c.currentConn().socket.Emit(
		"addAPIKey",
		map[string]any{"name": name},
		emit.WithAck(func(response map[string]any) { res <- response }),
	)
	if err != nil {
		return "", fmt.Errorf("addAPIKey: %w", err)
	}

	select {
	case response := <-res:
		key, hasKey := response["key"].(string)
		ok, hasOK := response["ok"].(bool)

		if !hasKey || !hasOK || !ok || key == "" {
			return "", fmt.Errorf("addAPIKey: %v", response)
		}

		return key, nil

	case <-ctx.Done():
		return "", fmt.Errorf("addAPIKey: %w", ctx.Err())
	}
}

// setAPIEnabled, exported as SetAPIEnabled, turns the "API and MCP" setting (apiEnabled) of an Uptime Kuma
// 3.x server on or off, keeping every other setting as it is. The Settings
// type has no field for it, since it is a 3.x setting in alpha upstream.
func (c *Client) setAPIEnabled(ctx context.Context, enabled bool, password string) error {
	response, err := c.syncEmit(ctx, "getSettings")
	if err != nil {
		return err
	}

	if response.Data == nil {
		return errors.New("getSettings: no data")
	}

	response.Data["apiEnabled"] = enabled

	_, err = c.syncEmit(ctx, "setSettings", response.Data, password)

	return err
}
