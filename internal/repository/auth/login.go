package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/pixel365/zoner/internal/db/postgres"

	"github.com/pixel365/zoner/internal/stringutils/password"
)

func (a *Auth) Login(ctx context.Context, username, psw, newPassword string) (int64, int64, error) {
	var userId int64
	var maxActiveSessions int64
	var isActive bool

	username = strings.ToLower(username)
	username = strings.TrimSpace(username)
	if username == "" {
		return userId, maxActiveSessions, ErrInvalidCredentials
	}

	psw = strings.TrimSpace(psw)
	if psw == "" {
		return userId, maxActiveSessions, ErrInvalidCredentials
	}

	var passwordHash string
	err := a.db.QueryRow(ctx,
		`SELECT id, password_hash, is_active, max_active_sessions 
FROM registrars 
WHERE username = $1`,
		username).Scan(&userId, &passwordHash, &isActive, &maxActiveSessions)
	if errors.Is(err, pgx.ErrNoRows) {
		return userId, maxActiveSessions, ErrInvalidCredentials
	}

	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrInternalError, err)
	}

	if !isActive {
		return 0, 0, ErrRegistrarIsBlocked
	}

	ok, err := password.Verify(psw, passwordHash)
	if err != nil {
		return userId, maxActiveSessions, fmt.Errorf("%w: %w", ErrInternalError, err)
	}

	if !ok {
		return userId, maxActiveSessions, ErrInvalidCredentials
	}

	if newPassword != "" {
		if err = a.changePassword(ctx, userId, newPassword); err != nil {
			return userId, maxActiveSessions, err
		}
	}

	return userId, maxActiveSessions, nil
}

func (a *Auth) changePassword(ctx context.Context, userId int64, newPassword string) error {
	sql := `UPDATE registrars SET password_hash = $2 WHERE id = $1`
	psw, _ := password.Hash(newPassword, password.DefaultParams)

	err := postgres.Tx(ctx, a.db, pgx.ReadCommitted, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx, sql, userId, psw)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCannotChangePassword, err)
		}

		if res.RowsAffected() != 1 {
			return ErrCannotChangePassword
		}

		return nil
	})

	return err
}
