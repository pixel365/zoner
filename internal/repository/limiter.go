package repository

import (
	"context"
	"time"
)

type LimiterRepository interface {
	Reserve(context.Context, string, int64, time.Duration) (bool, error)
	Release(context.Context, string) error
}
