package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) EnterRound(
	ctx context.Context,
	playerID int,
	outcome models.Outcome,
	value, winnerID, loserID int,
) error {
	game, err := s.CurrentGame(playerID)
	if err != nil {
		return err
	}

	err = game.EnterRoundResult(value, loserID, winnerID)
	if err != nil {
		return err
	}

	players := game.Players()

	east := players[models.East]
	south := players[models.South]
	west := players[models.West]
	north := players[models.North]

	score := game.Score()

	round := game.Round()

	roundOutcome := models.RoundData{
		Round:           round,
		Outcome:         outcome,
		WinnerID:        winnerID,
		LoserID:         loserID,
		Value:           value,
		EastStartScore:  score[east].GP,
		SouthStartScore: score[south].GP,
		WestStartScore:  score[west].GP,
		NorthStartScore: score[north].GP,
	}

	err = s.gameRepo.SaveRoundOutcome(ctx, game.ID(), roundOutcome)
	if err != nil {
		return err
	}

	if round == 16 {
		finalData := models.FinalGameData{
			GameID: game.ID(),
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

		err = s.gameRepo.Finalize(ctx, finalData)
		if err != nil {
			return err
		}

		delete(s.liveGames, game.ID())
	}

	return nil
}
