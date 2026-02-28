package main

import (
	"context"

	"github.com/spf13/cobra"
)

func newDownCommand(_ context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Migrate down",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
