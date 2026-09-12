package event

import (
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Create(ctx context.Context, params models.EventCreateParams) (int, error) {
	if !slices.Contains(models.ValidEventTypes, params.Type) {
		return 0, models.ErrInvalidEventType
	}

	if params.Name == "" {
		return 0, models.ErrEmptyName
	}

	visibility := models.EventVisibilityPublic
	if params.Visibility != "" {
		visibility = params.Visibility
	}

	event := models.Event{
		Type:         params.Type,
		Name:         params.Name,
		Description:  params.Description,
		OwnerUserID:  params.OwnerID,
		Visibility:   visibility,
		RatingsReach: models.DefaultContentReach(visibility),
		GamesReach:   models.DefaultContentReach(visibility),
	}

	if params.RatingsReach != nil {
		if !slices.Contains(models.ValidContentReach, *params.RatingsReach) {
			return 0, models.ErrInvalidContentReach
		}
		event.RatingsReach = *params.RatingsReach
	}
	if params.GamesReach != nil {
		if !slices.Contains(models.ValidContentReach, *params.GamesReach) {
			return 0, models.ErrInvalidContentReach
		}
		event.GamesReach = *params.GamesReach
	}

	return s.repo.Create(ctx, event)
}
