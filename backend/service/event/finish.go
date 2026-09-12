package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Finish(ctx context.Context, eventID int) error {
	event, err := s.Get(ctx, eventID)
	if err != nil {
		return err
	}

	if event.FinishedAt != nil {
		return models.ErrEventAlreadyFinished
	}

	games, err := s.repo.ListGames(ctx, eventID)
	if err != nil {
		return err
	}

	for _, game := range games {
		if game.Status == models.GameLive {
			return models.ErrEventHasOngoingGames
		}
	}

	return s.repo.Finish(ctx, eventID)
}
