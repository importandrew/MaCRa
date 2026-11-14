package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateSettings sets the event's settings (who may finish or abort a game).
// They can change at any time, but not in a finished event
// (models.ErrEventAlreadyFinished): it has no games left to act on.
func (s *Service) UpdateSettings(ctx context.Context, eventID int, settings models.EventSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}

	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}

	return s.repo.UpdateSettings(ctx, eventID, settings)
}
