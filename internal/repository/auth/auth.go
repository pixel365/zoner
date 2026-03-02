package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	e "github.com/pixel365/zoner/internal/errors"

	"github.com/pixel365/zoner/internal/repository"
	"github.com/pixel365/zoner/internal/stringutils/password"
)

var _ repository.Authenticator = (*Auth)(nil)

type Auth struct {
	db repository.QueryRower
}

func (a *Auth) Login(ctx context.Context, username, psw string) error {
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)
	if username == "" {
		return e.ErrInvalidCredentials
	}

	psw = strings.TrimSpace(psw)
	if psw == "" {
		return e.ErrInvalidCredentials
	}

	var passwordHash string
	err := a.db.QueryRow(ctx, "SELECT password_hash FROM users WHERE username = $1", username).
		Scan(&passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return e.ErrInvalidCredentials
	}

	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrInternalError, err)
	}

	ok, err := password.Verify(psw, passwordHash)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrInternalError, err)
	}

	if !ok {
		return e.ErrInvalidCredentials
	}

	return nil
}

func (a *Auth) Logout() error {
	return nil
}

func NewAuth(db repository.QueryRower) *Auth {
	return &Auth{db}
}
