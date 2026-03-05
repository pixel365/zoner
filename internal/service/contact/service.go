package contact

import (
	"context"

	"github.com/pixel365/goepp/command/create"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/internal/logger"

	"github.com/pixel365/zoner/internal/repository"
)

var _ ContactService = (*Service)(nil)

type ContactService interface {
	Info(context.Context) error
	Check(context.Context) error
	Create(context.Context, create.Contact, int64) response.Marshaller
}

type Service struct {
	repo repository.ContactRepository
	log  logger.Logger
}

func MustService(repo repository.ContactRepository, log logger.Logger) *Service {
	if repo == nil {
		panic("contact repository is nil")
	}
	return &Service{
		repo: repo,
		log:  log,
	}
}
