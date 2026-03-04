package limiter

import "context"

func (s *Service) Release(ctx context.Context, key string) error {
	return s.repo.Release(ctx, key)
}
