package event

import (
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// Create creates an event. Without a visibility it is private; reaches not
// given follow the visibility (models.DefaultContentReach). Without rules it
// plays by the MCR standard ones; the owner may change them until the first
// game.
func (s *Service) Create(ctx context.Context, params models.EventCreateParams) (int, error) {
	if !slices.Contains(models.ValidEventTypes, params.Type) {
		return 0, models.ErrInvalidEventType
	}

	name, err := eventName(params.Name)
	if err != nil {
		return 0, err
	}
	description, err := eventDescription(params.Description)
	if err != nil {
		return 0, err
	}

	visibility := models.EventVisibilityPrivate
	if params.Visibility != "" {
		if err := validateVisibility(params.Visibility); err != nil {
			return 0, err
		}
		visibility = params.Visibility
	}

	event := models.Event{
		Type:         params.Type,
		Name:         name,
		Description:  description,
		OwnerUserID:  params.OwnerID,
		Visibility:   visibility,
		RatingsReach: models.DefaultContentReach(visibility),
		GamesReach:   models.DefaultContentReach(visibility),
		Rules:        models.StandardGameRules(),
	}

	if params.RatingsReach != nil {
		if err := validateReach(*params.RatingsReach); err != nil {
			return 0, err
		}
		event.RatingsReach = *params.RatingsReach
	}
	if params.GamesReach != nil {
		if err := validateReach(*params.GamesReach); err != nil {
			return 0, err
		}
		event.GamesReach = *params.GamesReach
	}

	if params.Rules != nil {
		if err := params.Rules.Validate(); err != nil {
			return 0, err
		}
		event.Rules = *params.Rules
	}

	return s.repo.Create(ctx, event)
}
