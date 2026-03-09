package contact

import (
	"context"

	"github.com/pixel365/goepp/command/check"
	"github.com/pixel365/goepp/command/create"
	"github.com/pixel365/goepp/command/info"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/internal/logger"

	"github.com/pixel365/zoner/internal/repository"
)

type DiscloseField string

func (d DiscloseField) String() string { return string(d) }

const (
	DiscloseAddr    DiscloseField = "addr"
	DiscloseAddrInt DiscloseField = "addr:int"
	DiscloseAddrLoc DiscloseField = "addr:loc"
	DiscloseName    DiscloseField = "name"
	DiscloseOrg     DiscloseField = "org"
	DiscloseVoice   DiscloseField = "voice"
	DiscloseFax     DiscloseField = "fax"
	DiscloseEmail   DiscloseField = "email"
)

var _ ContactService = (*Service)(nil)

type ContactService interface {
	Info(context.Context, info.ContactInfo, int64, string) response.Marshaller
	Check(context.Context, check.ContactCheck, int64) response.Marshaller
	Create(context.Context, create.Contact, int64, string) response.Marshaller
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
