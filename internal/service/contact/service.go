package contact

import (
	"context"

	"github.com/pixel365/zoner/internal/repository"
)

var _ ContactService = (*Service)(nil)

type ContactService interface {
	Info(context.Context) error
	Check(context.Context) error
}

type Service struct {
	repo repository.ContactRepository
}

func MustService(repo repository.ContactRepository) *Service {
	if repo == nil {
		panic("contact repository is nil")
	}
	return &Service{repo}
}
