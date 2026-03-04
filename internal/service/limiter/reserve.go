package limiter

import (
	"context"
	"time"
)

func (s *Service) Reserve(
	ctx context.Context,
	key string,
	limit int64,
	ttl time.Duration,
) (bool, error) {
	return s.repo.Reserve(ctx, key, limit, ttl)
}
