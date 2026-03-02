package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type QueryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
