package game

import (
	"context"
	"fmt"
	"time"

	"github.com/importandrew/MaCRa/engine/game"
	"github.com/importandrew/MaCRa/models"
)

// load rebuilds the game, by the rules it started with, from its round log and returns it with that log, so
// a response can show the table and the log of one state.
func (s *Service) load(ctx context.Context, gameID int) (GameEngine, []models.RoundLog, error) {
	players, err := s.gameRepo.Players(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("get players: %w", err)
	}

	rules, err := s.gameRepo.Rules(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("get rules: %w", err)
	}
	// Stored rules were valid when saved; invalid ones mean broken data, a
	// server error (%v: not the client's invalid_game_rules), not a game
	// that never ends.
	if err := rules.Validate(); err != nil {
		return nil, nil, fmt.Errorf("game %d: stored rules: %v", gameID, err) //nolint:errorlint // see above
	}

	seating := game.Seating([4]int{
		players.East,
		players.South,
		players.West,
		players.North,
	})

	gameCfg := game.Config{
		ID:      gameID,
		Players: seating,
		Rules:   rules,
	}

	g := game.NewGame(gameCfg)

	roundLogs, err := s.gameRepo.GameLog(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("get game log: %w", err)
	}

	limit, err := s.gameRepo.TimeLimit(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("get time limit: %w", err)
	}
	timeUp := timeUpRule(g, limit)

	for _, log := range roundLogs {
		// The first entry after the time is up: the rule applies from the
		// hand that was in progress then.
		if limit != nil && log.CreatedAt.After(limit.EndsAt) {
			if err := timeUp(); err != nil {
				return nil, nil, err
			}
		}
		switch log.Outcome {
		case models.FalseMahjong:
			err = g.FalseMahjong(log.Loser, log.Value, log.ForceNextRound)
		case models.Penalty:
			err = g.Penalty(log.Loser, log.Value, log.ForceNextRound)
		default:
			err = g.EnterRoundResult(log.Outcome, log.Value, log.Loser, log.Winner)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("log record %d game %d: %w", log.ID, gameID, err)
		}
	}
	if limit != nil && time.Now().After(limit.EndsAt) {
		if err := timeUp(); err != nil {
			return nil, nil, err
		}
	}

	return g, roundLogs, nil
}

// timeUpRule returns what happens to g once its tour's time is up, applied
// once however often it is called: the hand in progress is the last one
// (models.TimeUpFinishHand), or a few more after it (models.TimeUpExtraHands),
// or the game stops at once (models.TimeUpStopNow). Without a limit it does
// nothing.
func timeUpRule(g GameEngine, limit *models.GameTimeLimit) func() error {
	applied := false
	return func() error {
		if limit == nil || applied {
			return nil
		}
		applied = true
		switch limit.TimeUp {
		case models.TimeUpExtraHands:
			g.LimitHands(g.Round() + limit.ExtraHands)
		case models.TimeUpStopNow:
			return g.StopAtTime()
		default:
			g.LimitHands(g.Round())
		}
		return nil
	}
}

// gameTimer is the timer a table shows for a started tournament game; nil
// without a time limit.
func gameTimer(g GameEngine, limit *models.GameTimeLimit) *models.GameTimer {
	if limit == nil {
		return nil
	}
	timer := &models.GameTimer{EndsAt: limit.EndsAt, TimeUp: limit.TimeUp, Stopped: g.StoppedAtTime()}
	if time.Now().After(limit.EndsAt) && !timer.Stopped {
		timer.LastHand = g.LastHand()
	}
	return timer
}

// LiveGame loads a live game, with its round log (see load).
func (s *Service) LiveGame(ctx context.Context, gameID int) (GameEngine, []models.RoundLog, error) {
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}

	if status != models.GameLive {
		return nil, nil, models.ErrGameNotLive
	}

	game, log, err := s.load(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return game, log, nil
}

// CurrentGame loads the player's live game in the event, with its round log
// (see load).
func (s *Service) CurrentGame(ctx context.Context, playerID, eventID int) (GameEngine, []models.RoundLog, error) {
	gameID, err := s.gameRepo.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return nil, nil, fmt.Errorf("get current game: %w", err)
	}

	game, log, err := s.load(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return game, log, nil
}

func (s *Service) CurrentGameID(ctx context.Context, playerID, eventID int) (int, error) {
	return s.gameRepo.CurrentGame(ctx, playerID, eventID)
}

// SeatedGameID returns the player's game that isn't over: live, or a
// tournament game whose tour is published but not started.
func (s *Service) SeatedGameID(ctx context.Context, playerID, eventID int) (int, error) {
	return s.gameRepo.SeatedGame(ctx, playerID, eventID)
}
