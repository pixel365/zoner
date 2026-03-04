package auth

import "context"

func (s *Service) Login(ctx context.Context, username, psw string) (int64, int64, error) {
	return s.repo.Login(ctx, username, psw)
}
