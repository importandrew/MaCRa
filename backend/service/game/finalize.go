package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) finalize(ctx context.Context, gameID int, score models.GameScore, players [4]int) error {
	east := players[models.East]
	south := players[models.South]
	west := players[models.West]
	north := players[models.North]

	finalData := models.FinalGameData{
		GameID: gameID,
		East: models.FinalPlayerData{
			ID: east,
			GP: score[east].GP,
			TP: score[east].TP,
		},
		South: models.FinalPlayerData{
			ID: south,
			GP: score[south].GP,
			TP: score[south].TP,
		},
		West: models.FinalPlayerData{
			ID: west,
			GP: score[west].GP,
			TP: score[west].TP,
		},
		North: models.FinalPlayerData{
			ID: north,
			GP: score[north].GP,
			TP: score[north].TP,
		},
	}

	err := s.gameRepo.Finalize(ctx, finalData)
	if err != nil {
		return err
	}

	return nil
}
