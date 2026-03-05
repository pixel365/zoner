package contact

import (
	"context"
	"errors"

	"github.com/pixel365/zoner/internal/model"
)

func (s *Service) Create(ctx context.Context, data model.ContactCreateInput) error {
	roid, err := s.repo.Create(ctx, data)

	if err != nil {
		//TODO: error response
		return err
	}

	if roid == "" {
		//TODO: error response
		return errors.New("empty roid")
	}

	return nil
}
