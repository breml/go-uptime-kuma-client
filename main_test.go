package kuma_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	kuma "github.com/breml/go-uptime-kuma-client"
)

//nolint:gochecknoglobals // client is used across multiple tests.
var client *kuma.Client

func TestMain(m *testing.M) {
	code, err := testMainSetup(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Test setup failed: %v\n", err)
		//revive:disable:redundant-test-main-exit
		os.Exit(code)
		//revive:enable:redundant-test-main-exit
	}
}

func testMainSetup(m *testing.M) (int, error) {
	dockerTimeout := uint(60)

	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if e2eTest {
		dockerTimeout = 600
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return 1, fmt.Errorf("could not construct pool: %w", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return 1, fmt.Errorf("could not connect to Docker: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       fmt.Sprintf("uptime-kuma-%s", randomString(8)),
		Repository: "louislam/uptime-kuma",
		Tag:        "2",
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return 1, fmt.Errorf("could not start resource: %w", err)
	}

	err = resource.Expire(dockerTimeout)
	if err != nil {
		return 1, fmt.Errorf("could not set resource expiry: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup client disconnect defer before attempting connection
	defer func() {
		if client != nil {
			err := client.Disconnect()
			if err != nil {
				log.Printf("Failed to disconnect from uptime kuma: %v", err)
			}
		}
	}()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	retryErr := pool.Retry(func() error {
		var err error
		client, err = kuma.New(
			ctx,
			fmt.Sprintf("http://localhost:%s", resource.GetPort("3001/tcp")),
			"admin", "admin1",
			kuma.WithAutosetup(),
			kuma.WithLogLevel(kuma.LogLevel(os.Getenv("SOCKETIO_LOG_LEVEL"))),
			kuma.WithConnectTimeout(10*time.Second),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if retryErr != nil {
		return 1, fmt.Errorf("could not connect to uptime kuma: %w", retryErr)
	}

	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
	defer func() {
		err := pool.Purge(resource)
		if err != nil {
			log.Printf("Could not purge resource: %v", err)
		}
	}()

	return m.Run(), nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	result := make([]byte, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
