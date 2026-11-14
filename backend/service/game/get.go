package game

import (
	"context"
	"fmt"

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

	for _, log := range roundLogs {
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

	return g, roundLogs, nil
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
