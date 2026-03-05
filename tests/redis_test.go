package tests

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func redisUp(ctx context.Context) {
	req := testcontainers.ContainerRequest{
		Name:         "testing-bitnami-redis",
		Image:        "bitnami/redis:8.0.2",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
		Env: map[string]string{
			"REDIS_PASSWORD": RedisPassword,
		},
		Networks: []string{testingNetwork.Name},
		NetworkAliases: map[string][]string{
			testingNetwork.Name: {"testing-redis"},
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Println("Failed to start container:", err)
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	redisContainer = container

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")
	RedisHost = host
	RedisPort = port.Port()
}
