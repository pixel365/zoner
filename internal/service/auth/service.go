package auth

import (
	"context"

	"github.com/pixel365/goepp/command/login"

	"github.com/pixel365/zoner/internal/repository"
)

var _ AuthService = (*Service)(nil)

type AuthService interface {
	Login(context.Context, *login.Login) (int64, int64, error)
	Logout() error
}

type Service struct {
	repo repository.AuthRepository
}

func MustService(repo repository.AuthRepository) *Service {
	if repo == nil {
		panic("auth repository is nil")
	}
	return &Service{repo}
}
