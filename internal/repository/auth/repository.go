package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pixel365/zoner/internal/repository"
)

var _ repository.AuthRepository = (*Auth)(nil)

type Auth struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Auth {
	return &Auth{db}
}
