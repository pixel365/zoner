package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pixel365/zoner/internal/stringutils/password"
)

func newAddCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new superuser",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if pool == nil {
				return fmt.Errorf("pool is nil")
			}

			query := `
INSERT INTO users (username, password_hash, email, is_superuser, max_active_sessions) 
VALUES ($1, $2, $3, true, 10)
`

			username, _ := cmd.Flags().GetString("username")
			email, _ := cmd.Flags().GetString("email")

			psw := rand.Text()
			passwordHash, err := password.Hash(psw, password.DefaultParams)
			if err != nil {
				return err
			}

			_, err = pool.Exec(ctx, query, username, passwordHash, email)
			if err != nil {
				return err
			}

			fmt.Printf("User %s added with password: %s\n", username, psw)
			fmt.Println("Press Enter to clear...")
			_, _ = fmt.Scanln()
			fmt.Print("\033[2J\033[3J\033[H")

			return nil
		},
	}

	cmd.Flags().StringP("username", "u", "", "Username")
	cmd.Flags().StringP("email", "e", "", "Email")

	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}
