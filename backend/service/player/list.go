package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error) {
	players, err := s.repo.ListByNameLimit(ctx, name, limit)
	if err != nil {
		return nil, err
	}

	return players, nil
}
