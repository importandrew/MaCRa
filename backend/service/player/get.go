package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Get(ctx context.Context, id int) (models.Player, error) {
	player, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}
