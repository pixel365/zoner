package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/pixel365/zoner/internal/db/postgres"
)

var pool *pgxpool.Pool

func init() {
	_ = godotenv.Load()
}

func rootCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "superuser",
		Short: "Superuser commands",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			pgConfig := postgres.NewConfigFromEnv()
			pgPool, err := postgres.NewPool(ctx, pgConfig)
			if err != nil {
				return err
			}

			pool = pgPool

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if pool != nil {
				pool.Close()
			}
		},
	}

	cmd.AddCommand(newAddCommand(ctx))

	return cmd
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCommand(ctx).Execute(); err != nil {
		stop()
		os.Exit(1)
	}
}
