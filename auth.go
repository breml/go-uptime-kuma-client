package kuma

import (
	"context"
	"fmt"
	"sync"
	"time"

	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
	"github.com/maldikhan/go.socket.io/utils"
	"github.com/maniartech/signals"
)

// TwoFactorConfig contains 2FA setup information returned when preparing 2FA.
type TwoFactorConfig struct {
	Secret string `json:"secret"` // TOTP secret for 2FA
	URI    string `json:"uri"`    // otpauth:// URI for QR code generation
}

// TwoFactorStatus contains the current 2FA status for a user.
type TwoFactorStatus struct {
	Status bool `json:"status"` // Whether 2FA is enabled
}

// TokenInfo contains authentication token details.
type TokenInfo struct {
	Token     string    `json:"token"`     // Session token
	ExpiresAt time.Time `json:"expiresAt"` // Token expiration time
}

// Get2FAStatus returns whether 2FA is enabled for the current user.
func (c *Client) Get2FAStatus(ctx context.Context) (bool, error) {
	response, err := c.syncEmit(ctx, "twoFAStatus")
	if err != nil {
		return false, err
	}

	var status TwoFactorStatus
	if statusVal, ok := response.Data["status"].(bool); ok {
		status.Status = statusVal
	}

	c.mu.Lock()
	c.has2FA = status.Status
	c.mu.Unlock()

	return status.Status, nil
}

// Prepare2FA prepares 2FA setup and returns the secret and QR code URI.
// This should be called before Save2FA to get the TOTP secret.
func (c *Client) Prepare2FA(ctx context.Context) (*TwoFactorConfig, error) {
	response, err := c.syncEmit(ctx, "prepare2FA")
	if err != nil {
		return nil, err
	}

	config := &TwoFactorConfig{}
	if secret, ok := response.Data["secret"].(string); ok {
		config.Secret = secret
	}
	if uri, ok := response.Data["uri"].(string); ok {
		config.URI = uri
	}

	if config.Secret == "" {
		return nil, fmt.Errorf("prepare2FA: no secret returned")
	}

	return config, nil
}

// Save2FA enables 2FA using the secret from Prepare2FA and a TOTP token.
// The token parameter should be a valid TOTP token generated from the secret.
func (c *Client) Save2FA(ctx context.Context, secret string, token string) error {
	_, err := c.syncEmit(ctx, "save2FA", secret, token)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.has2FA = true
	c.mu.Unlock()

	return nil
}

// Disable2FA disables 2FA for the current user.
// Requires the user's password for verification.
func (c *Client) Disable2FA(ctx context.Context, password string) error {
	_, err := c.syncEmit(ctx, "disable2FA", password)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.has2FA = false
	c.mu.Unlock()

	return nil
}

// Verify2FAToken verifies that a TOTP token is valid for the current user.
// Returns true if the token is valid, false otherwise.
func (c *Client) Verify2FAToken(ctx context.Context, token string) (bool, error) {
	response, err := c.syncEmit(ctx, "verifyToken", token)
	if err != nil {
		return false, err
	}

	if valid, ok := response.Data["valid"].(bool); ok {
		return valid, nil
	}

	return false, nil
}

// ChangePassword changes the user's password.
// Requires the current (old) password and the new password.
func (c *Client) ChangePassword(ctx context.Context, oldPassword string, newPassword string) error {
	_, err := c.syncEmit(ctx, "changePassword", map[string]any{
		"currentPassword": oldPassword,
		"newPassword":     newPassword,
	})
	return err
}

// LoginByToken authenticates using a session token instead of username/password.
// This is used by the NewWithToken constructor.
func (c *Client) LoginByToken(ctx context.Context, token string) error {
	response, err := c.syncEmit(ctx, "loginByToken", token)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.authToken = token
	c.mu.Unlock()

	// Check if the response indicates the token was accepted
	if response.OK {
		return nil
	}

	return fmt.Errorf("loginByToken: authentication failed")
}

// GetCurrentToken retrieves the current session token.
// This token can be used with NewWithToken for subsequent connections.
func (c *Client) GetCurrentToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	token := c.authToken
	c.mu.Unlock()

	if token != "" {
		return token, nil
	}

	// If we don't have a cached token, try to get it from the server
	response, err := c.syncEmit(ctx, "getToken")
	if err != nil {
		return "", err
	}

	if tokenStr, ok := response.Data["token"].(string); ok {
		c.mu.Lock()
		c.authToken = tokenStr
		c.mu.Unlock()
		return tokenStr, nil
	}

	return "", fmt.Errorf("getToken: no token returned")
}

// Logout terminates the current session.
func (c *Client) Logout(ctx context.Context) error {
	_, err := c.syncEmit(ctx, "logout")
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.authToken = ""
	c.has2FA = false
	c.mu.Unlock()

	return nil
}

// NewWithToken creates a client authenticated with a session token.
// This allows reconnecting to an existing session without username/password.
func NewWithToken(ctx context.Context, baseURL string, token string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},
		mu:             &sync.Mutex{},
		updates:        signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
	}

	ctxWithConnectTimeout := ctx
	if c.socketioClientConnectTimeout != 0 {
		var cancel func()
		ctxWithConnectTimeout, cancel = context.WithTimeout(ctx, c.socketioClientConnectTimeout)
		defer cancel()
	}

	// Handle database setup for Uptime Kuma v2 if autosetup is enabled
	if c.autosetup {
		if err := setupDatabase(ctxWithConnectTimeout, baseURL); err != nil {
			return nil, fmt.Errorf("database setup: %w", err)
		}
	}

	client, err := socketio.NewClient(
		socketio.WithRawURL(baseURL),
		socketio.WithLogger(c.socketioLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %v", err)
	}

	c.socketioClient = client

	if err := c.setupEventHandlers(ctx); err != nil {
		return nil, err
	}

	connect := c.setupConnectHandler()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %v", err)
	}

	if err := c.waitForConnect(ctx, connect); err != nil {
		return nil, err
	}

	// Authenticate with token
	if err := c.LoginByToken(ctxWithConnectTimeout, token); err != nil {
		return nil, fmt.Errorf("login by token: %v", err)
	}

	if err := c.waitForReady(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

// NewWith2FA creates a client authenticated with username/password and a 2FA TOTP token.
// Use this constructor when the account has 2FA enabled.
func NewWith2FA(ctx context.Context, baseURL string, username string, password string, totpToken string, opts ...Option) (*Client, error) {
	c := &Client{
		socketioLogger: &utils.DefaultLogger{Level: utils.NONE},
		mu:             &sync.Mutex{},
		updates:        signals.New[string](),
	}

	for _, opt := range opts {
		opt(c)
	}

	ctxWithConnectTimeout := ctx
	if c.socketioClientConnectTimeout != 0 {
		var cancel func()
		ctxWithConnectTimeout, cancel = context.WithTimeout(ctx, c.socketioClientConnectTimeout)
		defer cancel()
	}

	// Handle database setup for Uptime Kuma v2 if autosetup is enabled
	if c.autosetup {
		if err := setupDatabase(ctxWithConnectTimeout, baseURL); err != nil {
			return nil, fmt.Errorf("database setup: %w", err)
		}
	}

	client, err := socketio.NewClient(
		socketio.WithRawURL(baseURL),
		socketio.WithLogger(c.socketioLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("create socketio client: %v", err)
	}

	c.socketioClient = client

	if err := c.setupEventHandlers(ctx); err != nil {
		return nil, err
	}

	connect := c.setupConnectHandler()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %v", err)
	}

	if err := c.waitForConnect(ctx, connect); err != nil {
		return nil, err
	}

	// Authenticate with username/password and 2FA token
	_, err = c.syncEmit(ctxWithConnectTimeout, "login", map[string]any{
		"username": username,
		"password": password,
		"token":    totpToken,
	})
	if err != nil {
		return nil, fmt.Errorf("login with 2FA: %v", err)
	}

	c.mu.Lock()
	c.has2FA = true
	c.mu.Unlock()

	if err := c.waitForReady(ctx); err != nil {
		return nil, err
	}

	return c, nil
}
