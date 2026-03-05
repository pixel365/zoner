package auth

import (
	"context"

	"github.com/pixel365/goepp/command/login"
)

func (s *Service) Login(ctx context.Context, creds *login.Login) (int64, int64, error) {
	return s.repo.Login(ctx, creds.ClientID, creds.Password, creds.NewPassword)
}
