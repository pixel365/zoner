package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

var (
	eppHost string
	eppPort string

	testingNetwork     *testcontainers.DockerNetwork
	eppServerContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := mustSetupNetwork(ctx); err != nil {
		slog.Error("failed to setup network", slog.Any("error", err))
	}

	eppUp(ctx)

	code := m.Run()

	if err := eppServerContainer.Terminate(ctx); err != nil {
		slog.Error("failed to terminate container", slog.Any("error", err))
	}

	if err := testingNetwork.Remove(ctx); err != nil {
		slog.Error("failed to remove network", slog.Any("error", err))
	}

	os.Exit(code)
}

func mustSetupNetwork(ctx context.Context) error {
	var err error
	testingNetwork, err = network.New(ctx)
	return err
}
