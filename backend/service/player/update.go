package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

var ()

func (s *Service) Update(ctx context.Context, id int, name string) error {
	if name == "" {
		return models.ErrEmptyName
	}

	if len(name) > 50 {
		return models.ErrLongName
	}

	if err := s.repo.Update(ctx, id, name); err != nil {
		return err
	}

	return nil
}
