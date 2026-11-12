package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// Update changes the given settings of the event. After the event is
// finished only its visibility can change: a change of anything else is
// refused with models.ErrEventAlreadyFinished (sending a value unchanged is
// fine, as a form sends all its fields).
func (s *Service) Update(ctx context.Context, params models.UpdateEventParams) error {
	event, err := s.repo.Get(ctx, params.ID)
	if err != nil {
		return err
	}
	before := event

	if params.Name != nil {
		if event.Name, err = eventName(*params.Name); err != nil {
			return err
		}
	}
	if params.Description != nil {
		if event.Description, err = eventDescription(*params.Description); err != nil {
			return err
		}
	}
	if params.Visibility != nil {
		if err := validateVisibility(*params.Visibility); err != nil {
			return err
		}
		event.Visibility = *params.Visibility
	}
	if params.RatingsReach != nil {
		if err := validateReach(*params.RatingsReach); err != nil {
			return err
		}
		event.RatingsReach = *params.RatingsReach
	}
	if params.GamesReach != nil {
		if err := validateReach(*params.GamesReach); err != nil {
			return err
		}
		event.GamesReach = *params.GamesReach
	}

	if before.FinishedAt != nil {
		beforeOtherThanVisibility := before
		beforeOtherThanVisibility.Visibility = event.Visibility
		if beforeOtherThanVisibility != event {
			return models.ErrEventAlreadyFinished
		}
	}

	return s.repo.Update(ctx, event)
}
