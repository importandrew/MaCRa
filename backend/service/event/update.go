package event

import (
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Update(ctx context.Context, params models.UpdateEventParams) error {
	event, err := s.repo.Get(ctx, params.ID)
	if err != nil {
		return err
	}

	if params.Name != nil {
		event.Name = *params.Name
	}
	if params.Description != nil {
		event.Description = *params.Description
	}
	if params.Visibility != nil {
		event.Visibility = *params.Visibility
	}
	if params.RatingsReach != nil {
		if !slices.Contains(models.ValidContentReach, *params.RatingsReach) {
			return models.ErrInvalidContentReach
		}
		event.RatingsReach = *params.RatingsReach
	}
	if params.GamesReach != nil {
		if !slices.Contains(models.ValidContentReach, *params.GamesReach) {
			return models.ErrInvalidContentReach
		}
		event.GamesReach = *params.GamesReach
	}

	return s.repo.Update(ctx, event)
}
