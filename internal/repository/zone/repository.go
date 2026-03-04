package zone

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pixel365/zoner/internal/repository"
)

var _ repository.ZonesRepository = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}
