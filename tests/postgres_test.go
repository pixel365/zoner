package tests

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pixel365/zoner/internal/stringutils/password"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
)

func postgresUp(ctx context.Context) {
	req := testcontainers.ContainerRequest{
		Name:         "testing-postgres",
		Image:        "postgres:18.3",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     PostgresUser,
			"POSTGRES_PASSWORD": PostgresPassword,
			"POSTGRES_DB":       PostgresDb,
			"POSTGRES_SSL_MODE": "disable",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		Networks:   []string{testingNetwork.Name},
		NetworkAliases: map[string][]string{
			testingNetwork.Name: {"testing-postgres"},
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

	postgresContainer = container

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	PostgresHost = host
	PostgresPort = port.Port()

	time.Sleep(5 * time.Second)

	fmt.Printf("PostgreSQL is running at %s:%s\n", PostgresHost, PostgresPort)
}

func applyMigrations(ctx context.Context) {
	user := url.UserPassword(PostgresUser, PostgresPassword)
	dsn := fmt.Sprintf(
		"postgres://%s@%s:%s/%s?sslmode=%s",
		user,
		PostgresHost,
		PostgresPort,
		PostgresDb,
		"disable",
	)

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	if err = goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	err = goose.UpContext(ctx, db, "../migrations")
	if err != nil {
		panic(err)
	}

	sql := `
INSERT INTO registrars (username, password_hash, email, is_active, max_active_sessions) 
VALUES ($1, $2, $3, true, $4)
`

	passwordHash, _ := password.Hash(testingRegistrarPassword, password.DefaultParams)
	_, err = db.ExecContext(ctx, sql, testingRegistrarUsername, passwordHash, "test@test.com", 100)
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
}
