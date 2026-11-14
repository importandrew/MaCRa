package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

// TimeUpCheckInterval is how often the server finishes the tournament games
// stopped at their tour's time.
const TimeUpCheckInterval = 15 * time.Second

// FinishTimedOut finishes the tournament games whose time is up with the
// rule "the game ends at once" (they have ended in the replay), the same way
// the table's Finish does, and returns how many it finished. A game finished
// meanwhile (by its table or a referee) is skipped.
func (s *Service) FinishTimedOut(ctx context.Context) (int, error) {
	gameIDs, err := s.gameRepo.TimedOutGames(ctx, time.Now().UTC())
	if err != nil {
		return 0, fmt.Errorf("timed-out games: %w", err)
	}

	finished := 0
	var errs []error
	for _, gameID := range gameIDs {
		err := s.Finish(ctx, gameID)
		switch {
		case err == nil:
			finished++
		case errors.Is(err, models.ErrGameNotLive):
		default:
			errs = append(errs, fmt.Errorf("finish game %d: %w", gameID, err))
		}
	}
	return finished, errors.Join(errs...)
}

// RunTimeUpChecks calls FinishTimedOut every interval until ctx is done, so
// a stopped game ends without anyone tapping Finish. It also catches up on
// what was missed while the server was down: each check finds every such
// game.
func (s *Service) RunTimeUpChecks(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		finished, err := s.FinishTimedOut(ctx)
		if err != nil {
			logger.Errorf("time-up check: %v", err)
		}
		if finished > 0 {
			logger.Infof("time-up check: %d games finished", finished)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
