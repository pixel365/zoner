package main

import (
	"context"

	"github.com/spf13/cobra"
)

func newUpCommand(_ context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate up",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
