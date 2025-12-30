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

var client *kuma.Client

func TestMain(m *testing.M) {
	dockerTimeout := uint(60)

	e2eTest, _ := strconv.ParseBool(os.Getenv("E2E_TEST"))
	if e2eTest {
		dockerTimeout = 600
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %v", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
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
		log.Fatalf("Could not start resource: %v", err)
	}

	err = resource.Expire(dockerTimeout)
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
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
	}); err != nil {
		log.Fatalf("Could not connect to uptime kuma: %v", err)
	}

	defer func() {
		err = client.Disconnect()
		if err != nil {
			log.Fatalf("Failed to connect to uptime kuma: %v", err)
		}
	}()

	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %v", err)
		}
	}()

	m.Run()
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
