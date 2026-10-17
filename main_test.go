package kuma_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	kuma "github.com/breml/go-uptime-kuma-client"
)

//nolint:gochecknoglobals // client is used across multiple tests.
var client *kuma.Client

// The account the tests log in as. The password satisfies the rules an Uptime
// Kuma 3.x server applies to the first user (at least eight characters) and
// the leak check is turned off on the container, see kumaContainerEnv.
const (
	testUsername = "admin"
	testPassword = "Kuma-Test-8f4Q2xR9p"
)

// defaultKumaImage is the Uptime Kuma the tests run against unless
// UPTIME_KUMA_IMAGE names another one, such as louislam/uptime-kuma:nightly3
// for the 3.x line.
const defaultKumaImage = "louislam/uptime-kuma:2.5.6"

// kumaImage returns the repository and tag of the Uptime Kuma image to test
// against.
func kumaImage() (repository string, tag string) {
	image := os.Getenv("UPTIME_KUMA_IMAGE")
	if image == "" {
		image = defaultKumaImage
	}

	repository, tag, ok := strings.Cut(image, ":")
	if !ok {
		return image, "latest"
	}

	return repository, tag
}

// kumaContainerEnv is the environment every test container gets: a 3.x server
// checks the first user's password against the Have I Been Pwned database,
// which the test password has no business with.
func kumaContainerEnv() []string {
	return []string{"UPTIME_KUMA_DISABLE_LEAK_PWD_CHECK=1"}
}

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
	// Long enough for the whole package to run against the container: the
	// integration tests alone take about a minute, and an expiry that lands in
	// the middle of them fails the run with a closed connection rather than a
	// real defect.
	dockerTimeout := uint(300)

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

	repository, tag := kumaImage()

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       fmt.Sprintf("uptime-kuma-%s", randomString(8)),
		Repository: repository,
		Tag:        tag,
		Env:        kumaContainerEnv(),
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
			testUsername, testPassword,
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
