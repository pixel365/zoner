package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func rootCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
	}

	cmd.AddCommand(
		newUpCommand(ctx),
		newDownCommand(ctx),
	)

	return cmd
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCommand(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
