package proxy

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProxy_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": 1,
		"userId": 100,
		"protocol": "http",
		"host": "proxy.example.com",
		"port": 8080,
		"auth": 1,
		"username": "user",
		"password": "pass",
		"active": 1,
		"default": 0,
		"createdDate": "2024-01-01T00:00:00.000Z"
	}`

	var proxy Proxy
	err := json.Unmarshal([]byte(jsonData), &proxy)
	require.NoError(t, err)

	require.Equal(t, int64(1), proxy.ID)
	require.Equal(t, int64(100), proxy.UserID)
	require.Equal(t, "http", proxy.Protocol)
	require.Equal(t, "proxy.example.com", proxy.Host)
	require.Equal(t, 8080, proxy.Port)
	require.True(t, proxy.Auth)
	require.Equal(t, "user", proxy.Username)
	require.Equal(t, "pass", proxy.Password)
	require.True(t, proxy.Active)
	require.False(t, proxy.Default)
}

func TestProxy_MarshalJSON(t *testing.T) {
	proxy := Proxy{
		ID:          1,
		UserID:      100,
		Protocol:    "socks5",
		Host:        "socks.example.com",
		Port:        1080,
		Auth:        true,
		Username:    "admin",
		Password:    "secret",
		Active:      true,
		Default:     true,
		CreatedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(proxy)
	require.NoError(t, err)

	var unmarshaled Proxy
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	require.Equal(t, proxy.ID, unmarshaled.ID)
	require.Equal(t, proxy.Protocol, unmarshaled.Protocol)
	require.Equal(t, proxy.Host, unmarshaled.Host)
	require.Equal(t, proxy.Default, unmarshaled.Default)
}

func TestConfig_MarshalJSON(t *testing.T) {
	config := Config{
		Protocol:      "https",
		Host:          "proxy.test.com",
		Port:          3128,
		Auth:          false,
		Active:        true,
		Default:       false,
		ApplyExisting: true,
	}

	data, err := json.Marshal(config)
	require.NoError(t, err)

	var unmarshaled Config
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	require.Equal(t, config.Protocol, unmarshaled.Protocol)
	require.Equal(t, config.ApplyExisting, unmarshaled.ApplyExisting)
}

func TestConfig_WithAuth(t *testing.T) {
	config := Config{
		Protocol: "socks5",
		Host:     "localhost",
		Port:     1080,
		Auth:     true,
		Username: "testuser",
		Password: "testpass",
		Active:   true,
		Default:  false,
	}

	data, err := json.Marshal(config)
	require.NoError(t, err)

	var unmarshaled Config
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	require.True(t, unmarshaled.Auth)
	require.Equal(t, config.Username, unmarshaled.Username)
	require.Equal(t, config.Password, unmarshaled.Password)
}

func TestProxy_GetID(t *testing.T) {
	proxy := Proxy{ID: 42}
	require.Equal(t, int64(42), proxy.GetID())
}
