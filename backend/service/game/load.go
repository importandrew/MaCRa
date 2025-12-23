package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/engine/game"
	"github.com/importandrew/MaCRa/models"
)

func (s *Service) LoadLiveGames(ctx context.Context) error {
	ids, err := s.gameRepo.LiveGames(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = s.loadGame(ctx, id)
		if err != nil {
			return fmt.Errorf("load game %d: %w", id, err)
		}
	}

	return nil
}

func (s *Service) loadGame(ctx context.Context, gameID int) error {
	players, err := s.gameRepo.Players(ctx, gameID)
	if err != nil {
		return err
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
		return err
	}

	for _, log := range roundLogs {
		switch log.Outcome {
		case models.FalseMahjong:
			err = g.FalseMahjong(log.Loser, log.Value, log.ForceNextRound)
		case models.Penalty:
			err = g.Penalty(log.Loser, log.Value, log.ForceNextRound)
		default:
			err = g.EnterRoundResult(log.Value, log.Loser, log.Winner)
		}
		if err != nil {
			return fmt.Errorf("log record %d game %d: %w", log.ID, gameID, err)
		}
	}

	s.liveGames[gameID] = g

	return nil
}
