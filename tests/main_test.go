package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

const minFrameLength = 4
const maxFrameSize uint32 = 4 * 1024 * 1024

var (
	eppHost string
	eppPort string

	testingNetwork     *testcontainers.DockerNetwork
	eppServerContainer testcontainers.Container
	redisContainer     testcontainers.Container
	postgresContainer  testcontainers.Container

	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword = "password"

	PostgresDb       = "test"
	PostgresHost     string
	PostgresPort     string
	PostgresUser     = "test"
	PostgresPassword = "test"

	testingRegistrarUsername = "test"
	testingRegistrarPassword = "test"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := mustSetupNetwork(ctx); err != nil {
		slog.Error("failed to setup network", slog.Any("error", err))
	}

	defer func() {
		if eppServerContainer != nil {
			_ = eppServerContainer.Terminate(ctx)
		}
		if redisContainer != nil {
			_ = redisContainer.Terminate(ctx)
		}
		if postgresContainer != nil {
			_ = postgresContainer.Terminate(ctx)
		}

		if testingNetwork != nil {
			_ = testingNetwork.Remove(ctx)
		}
	}()

	redisUp(ctx)
	postgresUp(ctx)
	applyMigrations(ctx)
	eppUp(ctx)

	code := m.Run()

	os.Exit(code)
}

func mustSetupNetwork(ctx context.Context) error {
	var err error
	testingNetwork, err = network.New(ctx)
	return err
}
