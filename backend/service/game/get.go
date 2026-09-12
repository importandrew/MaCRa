package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/engine/game"
	"github.com/importandrew/MaCRa/models"
)

func (s *Service) load(ctx context.Context, gameID int) (GameEngine, error) {
	players, err := s.gameRepo.Players(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get players: %w", err)
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
		Rules:   game.Rules{},
	}

	g := game.NewGame(gameCfg)

	roundLogs, err := s.gameRepo.GameLog(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game log: %w", err)
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
			return nil, fmt.Errorf("log record %d game %d: %w", log.ID, gameID, err)
		}
	}

	return g, nil
}

func (s *Service) LiveGame(ctx context.Context, gameID int) (GameEngine, error) {
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if status != models.GameLive {
		return nil, models.ErrGameNotLive
	}

	game, err := s.load(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return game, nil
}

func (s *Service) CurrentGame(ctx context.Context, playerID, eventID int) (GameEngine, error) {
	gameID, err := s.gameRepo.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return nil, fmt.Errorf("get current game: %w", err)
	}

	game, err := s.load(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return game, nil
}

func (s *Service) CurrentGameID(ctx context.Context, playerID, eventID int) (int, error) {
	return s.gameRepo.CurrentGame(ctx, playerID, eventID)
}
