package kuma_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	kuma "github.com/breml/go-uptime-kuma-client"
	"github.com/breml/go-uptime-kuma-client/dockerhost"
)

func TestDockerHostCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var err error

	t.Run("initial_state", func(t *testing.T) {
		hosts := client.GetDockerHostList(ctx)
		t.Logf("Initial Docker hosts count: %d", len(hosts))
	})

	var socketHostID int64
	t.Run("create_socket_docker_host", func(t *testing.T) {
		initialHosts := client.GetDockerHostList(ctx)
		initialCount := len(initialHosts)

		config := dockerhost.Config{
			Name:         "Local Docker",
			DockerDaemon: "unix:///var/run/docker.sock",
			DockerType:   "socket",
		}

		socketHostID, err = client.CreateDockerHost(ctx, config)
		require.NoError(t, err)
		require.Positive(t, socketHostID)

		hosts := client.GetDockerHostList(ctx)
		require.Len(t, hosts, initialCount+1)

		createdHost, err := client.GetDockerHost(ctx, socketHostID)
		require.NoError(t, err)
		require.Equal(t, socketHostID, createdHost.GetID())
		require.Equal(t, "Local Docker", createdHost.Name)
		require.Equal(t, "unix:///var/run/docker.sock", createdHost.DockerDaemon)
		require.Equal(t, "socket", createdHost.DockerType)
	})

	var tcpHostID int64
	t.Run("create_tcp_docker_host", func(t *testing.T) {
		initialHosts := client.GetDockerHostList(ctx)
		initialCount := len(initialHosts)

		config := dockerhost.Config{
			Name:         "Remote Docker",
			DockerDaemon: "tcp://192.168.1.100:2375",
			DockerType:   "tcp",
		}

		tcpHostID, err = client.CreateDockerHost(ctx, config)
		require.NoError(t, err)
		require.Positive(t, tcpHostID)

		hosts := client.GetDockerHostList(ctx)
		require.Len(t, hosts, initialCount+1)

		createdHost, err := client.GetDockerHost(ctx, tcpHostID)
		require.NoError(t, err)
		require.Equal(t, tcpHostID, createdHost.GetID())
		require.Equal(t, "Remote Docker", createdHost.Name)
		require.Equal(t, "tcp://192.168.1.100:2375", createdHost.DockerDaemon)
		require.Equal(t, "tcp", createdHost.DockerType)
	})

	t.Run("update_docker_host", func(t *testing.T) {
		config := dockerhost.Config{
			ID:           socketHostID,
			Name:         "Local Docker Updated",
			DockerDaemon: "unix:///var/run/docker.sock",
			DockerType:   "socket",
		}

		err = client.UpdateDockerHost(ctx, config)
		require.NoError(t, err)

		updatedHost, err := client.GetDockerHost(ctx, socketHostID)
		require.NoError(t, err)
		require.Equal(t, "Local Docker Updated", updatedHost.Name)
		require.Equal(t, "unix:///var/run/docker.sock", updatedHost.DockerDaemon)
		require.Equal(t, "socket", updatedHost.DockerType)
	})

	t.Run("update_docker_host_change_connection", func(t *testing.T) {
		config := dockerhost.Config{
			ID:           tcpHostID,
			Name:         "Remote Docker Updated",
			DockerDaemon: "tcp://192.168.1.200:2376",
			DockerType:   "tcp",
		}

		err = client.UpdateDockerHost(ctx, config)
		require.NoError(t, err)

		updatedHost, err := client.GetDockerHost(ctx, tcpHostID)
		require.NoError(t, err)
		require.Equal(t, "Remote Docker Updated", updatedHost.Name)
		require.Equal(t, "tcp://192.168.1.200:2376", updatedHost.DockerDaemon)
		require.Equal(t, "tcp", updatedHost.DockerType)
	})

	t.Run("delete_docker_host", func(t *testing.T) {
		preDeleteHosts := client.GetDockerHostList(ctx)
		preDeleteCount := len(preDeleteHosts)

		err := client.DeleteDockerHost(ctx, socketHostID)
		require.NoError(t, err)

		hosts := client.GetDockerHostList(ctx)
		require.Len(t, hosts, preDeleteCount-1)

		_, err = client.GetDockerHost(ctx, socketHostID)
		require.Error(t, err)
		require.ErrorIs(t, err, kuma.ErrNotFound)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteDockerHost(ctx, tcpHostID)
		require.NoError(t, err)
	})
}

func TestDockerHostTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	t.Run("test_unreachable_tcp_host", func(t *testing.T) {
		config := dockerhost.Config{
			Name:         "Unreachable Docker",
			DockerDaemon: "tcp://192.168.255.255:2375",
			DockerType:   "tcp",
		}

		result, err := client.TestDockerHost(ctx, config)
		// The server returns an error for unreachable hosts
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("test_invalid_socket_path", func(t *testing.T) {
		config := dockerhost.Config{
			Name:         "Invalid Socket",
			DockerDaemon: "unix:///nonexistent/docker.sock",
			DockerType:   "socket",
		}

		result, err := client.TestDockerHost(ctx, config)
		// The server returns an error for invalid socket paths
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestDockerHostString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var hostID int64
	var err error
	t.Run("create_and_verify_string", func(t *testing.T) {
		config := dockerhost.Config{
			Name:         "Test String Docker",
			DockerDaemon: "unix:///var/run/docker.sock",
			DockerType:   "socket",
		}

		hostID, err = client.CreateDockerHost(ctx, config)
		require.NoError(t, err)
		require.Positive(t, hostID)

		host, err := client.GetDockerHost(ctx, hostID)
		require.NoError(t, err)

		str := host.String()
		require.NotEmpty(t, str)
		t.Logf("Docker host string: %s", str)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := client.DeleteDockerHost(ctx, hostID)
		require.NoError(t, err)
	})
}

func TestDockerHostGetNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	t.Run("get_nonexistent_docker_host", func(t *testing.T) {
		_, err := client.GetDockerHost(ctx, 999999)
		require.Error(t, err)
		require.ErrorIs(t, err, kuma.ErrNotFound)
	})
}
