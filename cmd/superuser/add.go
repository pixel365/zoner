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
INSERT INTO staffs (name, last_name, email, password_hash, roles, is_active) 
VALUES ($1, $2, $3, $4, $5::staff_role[], true)
`

			name, _ := cmd.Flags().GetString("name")
			lastName, _ := cmd.Flags().GetString("lastname")
			email, _ := cmd.Flags().GetString("email")

			psw := rand.Text()
			passwordHash, err := password.Hash(psw, password.DefaultParams)
			if err != nil {
				return err
			}

			_, err = pool.Exec(ctx, query, name, lastName, email, passwordHash, []string{"root"})
			if err != nil {
				return err
			}

			fmt.Printf("User %s added with password: %s\n", email, psw)
			fmt.Println("Press Enter to clear...")
			_, _ = fmt.Scanln()
			fmt.Print("\033[2J\033[3J\033[H")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name")
	cmd.Flags().StringP("lastname", "l", "", "Lastname")
	cmd.Flags().StringP("email", "e", "", "Email")

	_ = cmd.MarkFlagRequired("email")

	return cmd
}
