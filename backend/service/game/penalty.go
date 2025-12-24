package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Penalty(ctx context.Context, playerID, penalizedID, value int, forceNextRound bool) (models.GameNow, error) {
	game, err := s.CurrentGame(playerID)
	if err != nil {
		return models.GameNow{}, err
	}

	err = game.Penalty(penalizedID, value, forceNextRound)
	if err != nil {
		return models.GameNow{}, err
	}

	players := game.Players()
	score := game.Score()

	roundOutcome := models.RoundData{
		Round:           game.Round(),
		Outcome:         models.Penalty,
		LoserID:         penalizedID,
		Value:           value,
		EastStartScore:  score[players[models.East]].GP,
		SouthStartScore: score[players[models.South]].GP,
		WestStartScore:  score[players[models.West]].GP,
		NorthStartScore: score[players[models.North]].GP,
		ForceNextRound:  forceNextRound,
	}

	err = s.gameRepo.SaveRoundOutcome(ctx, game.ID(), roundOutcome)
	if err != nil {
		return models.GameNow{}, err
	}

	if game.Status() == models.GameEnded {
		err := s.finalize(ctx, game.ID(), score, players)
		if err != nil {
			return models.GameNow{}, err
		}

		return models.GameNow{
			ID:     game.ID(),
			Status: models.GameEnded,
		}, nil
	}

	now, err := s.Now(ctx, playerID)
	if err != nil {
		return models.GameNow{}, err
	}

	return now, nil
}
