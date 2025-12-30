package dockerhost_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/dockerhost"
)

func TestDockerHost_GetID(t *testing.T) {
	d := dockerhost.DockerHost{
		ID: 42,
	}

	require.Equal(t, int64(42), d.GetID())
}

func TestDockerHost_String(t *testing.T) {
	d := dockerhost.DockerHost{
		ID:           1,
		UserID:       100,
		DockerDaemon: "unix:///var/run/docker.sock",
		DockerType:   "socket",
		Name:         "Local Docker",
	}

	got := d.String()
	require.Contains(t, got, "id: 1")
	require.Contains(t, got, `dockerDaemon: "unix:///var/run/docker.sock"`)
	require.Contains(t, got, `name: "Local Docker"`)
}

func TestDockerHost_MarshalJSON(t *testing.T) {
	d := dockerhost.DockerHost{
		ID:           1,
		UserID:       100,
		DockerDaemon: "unix:///var/run/docker.sock",
		DockerType:   "socket",
		Name:         "Local Docker",
	}

	data, err := json.Marshal(d)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.InEpsilon(t, float64(1), result["id"], 0)
	require.Equal(t, "unix:///var/run/docker.sock", result["dockerDaemon"])
}

func TestDockerHost_UnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"userId": 100,
		"dockerDaemon": "unix:///var/run/docker.sock",
		"dockerType": "socket",
		"name": "Local Docker"
	}`

	var d dockerhost.DockerHost
	err := json.Unmarshal([]byte(jsonStr), &d)
	require.NoError(t, err)

	require.Equal(t, int64(1), d.ID)
	require.Equal(t, int64(100), d.UserID)
	require.Equal(t, "unix:///var/run/docker.sock", d.DockerDaemon)
	require.Equal(t, "socket", d.DockerType)
	require.Equal(t, "Local Docker", d.Name)
}

func TestConfig_MarshalJSON(t *testing.T) {
	c := dockerhost.Config{
		ID:           1,
		Name:         "Test Docker",
		DockerDaemon: "tcp://192.168.1.100:2375",
		DockerType:   "tcp",
	}

	data, err := json.Marshal(c)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.Equal(t, "Test Docker", result["name"])
	require.Equal(t, "tcp", result["dockerType"])
}

func TestTestResult_UnmarshalJSON_StringVersion(t *testing.T) {
	jsonStr := `{
		"ok": true,
		"msg": "Connected Successfully.",
		"version": "20.10.17"
	}`

	var tr dockerhost.TestResult
	err := json.Unmarshal([]byte(jsonStr), &tr)
	require.NoError(t, err)

	require.True(t, tr.OK)
	require.Equal(t, "20.10.17", tr.Version)
}

func TestTestResult_UnmarshalJSON_ObjectVersion(t *testing.T) {
	jsonStr := `{
		"ok": true,
		"msg": "Connected Successfully.",
		"version": {
			"Version": "20.10.17",
			"ApiVersion": "1.41"
		}
	}`

	var tr dockerhost.TestResult
	err := json.Unmarshal([]byte(jsonStr), &tr)
	require.NoError(t, err)

	require.True(t, tr.OK)
	require.Equal(t, "20.10.17", tr.Version)
}

func TestTestResult_UnmarshalJSON_MissingVersion(t *testing.T) {
	jsonStr := `{
		"ok": false,
		"msg": "Connection failed"
	}`

	var tr dockerhost.TestResult
	err := json.Unmarshal([]byte(jsonStr), &tr)
	require.NoError(t, err)

	require.False(t, tr.OK)
	require.Empty(t, tr.Version)
	require.Equal(t, "Connection failed", tr.Msg)
}
