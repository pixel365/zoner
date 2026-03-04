package zone

import (
	"context"

	"github.com/pixel365/zoner/internal/repository"
)

var _ ZoneService = (*Service)(nil)

type ZoneService interface {
	Resolve(context.Context) error
	Policy(context.Context) error
}

type Service struct {
	repo repository.ZonesRepository
}

func MustService(repo repository.ZonesRepository) *Service {
	if repo == nil {
		panic("zones repository is nil")
	}
	return &Service{repo}
}
