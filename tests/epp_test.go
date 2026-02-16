package tests

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StdoutLogConsumer struct{}

func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}

func eppUp(ctx context.Context) {
	name := "epp"

	_, err := os.ReadFile("../server.crt")
	if err != nil {
		panic(err)
	}

	_, err = os.ReadFile("../server.key")
	if err != nil {
		panic(err)
	}

	req := testcontainers.ContainerRequest{
		Name: name,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../",
			Dockerfile: "testing.dockerfile",
			Tag:        "0.0.1",
			BuildArgs: map[string]*string{
				"SERVICE_NAME": &name,
			},
			BuildLogWriter: os.Stdout,
		},
		WaitingFor:      wait.ForListeningPort("7000/tcp"),
		AlwaysPullImage: false,
		Networks:        []string{testingNetwork.Name},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{
				testcontainers.WithLogProductionTimeout(10 * time.Second),
			},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{}},
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../server.crt",
				ContainerFilePath: "/app/server.crt",
				FileMode:          0644,
			},
			{
				HostFilePath:      "../server.key",
				ContainerFilePath: "/app/server.key",
				FileMode:          0644,
			},
			{
				HostFilePath:      "../config.dev.yaml",
				ContainerFilePath: "/app/config.dev.yaml",
				FileMode:          0644,
			},
		},
		Env: map[string]string{
			"CONFIG_PATH": "/app/config.dev.yaml",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	eppServerContainer = container

	time.Sleep(5 * time.Second)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "7000/tcp")

	eppHost = host
	eppPort = port.Port()

	slog.Info("epp server started", "host", host, "port", port.Port())
}
