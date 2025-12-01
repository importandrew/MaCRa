package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/engine/game"
	"github.com/importandrew/MaCRa/models"
)

// type DBGameLive struct {
// 	ID         int       `db:"id"`
// 	Round      int       `db:"round"`
// 	EastScore  int       `db:"east_score"`
// 	SouthScore int       `db:"south_score"`
// 	WestScore  int       `db:"west_score"`
// 	NorthScore int       `db:"north_score"`
// 	CreatedAt  time.Time `db:"created_at"`
// }

// type DBGameFinal struct {
// 	ID        int `db:"id"`
// 	EastID    int `db:"east_id"`
// 	EastTP    int `db:"east_tp"`
// 	EastGP    int `db:"east_gp"`
// 	SouthID   int `db:"south_id"`
// 	SouthTP   int `db:"south_tp"`
// 	SouthGP   int `db:"south_gp"`
// 	WestID    int `db:"west_id"`
// 	WestTP    int `db:"west_tp"`
// 	WestGP    int `db:"west_gp"`
// 	NorthID   int `db:"north_id"`
// 	NorthTP   int `db:"north_tp"`
// 	NorthGP   int `db:"north_gp"`
// 	CreatedAt int `db:"created_at"`
// }

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

	g, err := game.NewGame(gameCfg)
	if err != nil {
		return err
	}

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
