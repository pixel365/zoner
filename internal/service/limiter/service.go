package limiter

import (
	"context"
	"time"

	"github.com/pixel365/zoner/internal/repository"
)

var _ LimiterService = (*Service)(nil)

type LimiterService interface {
	Reserve(context.Context, string, int64, time.Duration) (bool, error)
	Release(context.Context, string) error
}

type Service struct {
	repo repository.LimiterRepository
}

func MustService(repo repository.LimiterRepository) *Service {
	if repo == nil {
		panic("limiter repository is nil")
	}
	return &Service{repo}
}
