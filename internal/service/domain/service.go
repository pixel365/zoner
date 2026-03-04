package domain

import (
	"context"

	"github.com/pixel365/zoner/internal/repository"
)

var _ DomainService = (*Service)(nil)

type DomainService interface {
	Create(context.Context) error
	Check(context.Context) error
	Info(context.Context) error
	Update(context.Context) error
	Delete(context.Context) error
	Renew(context.Context) error
	Transfer(context.Context) error
	Poll(context.Context) error
}

type Service struct {
	repo repository.DomainRepository
}

func MustService(repo repository.DomainRepository) *Service {
	if repo == nil {
		panic("domain repository is nil")
	}
	return &Service{repo}
}
