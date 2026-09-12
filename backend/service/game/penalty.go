package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Penalty(ctx context.Context, penaltyType models.Outcome, playerID, penalizedID, value int, forceNextRound bool, eventID int) (models.GameNow, error) {
	game, err := s.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	switch penaltyType {
	case models.Penalty:
		err = game.Penalty(penalizedID, value, forceNextRound)
	case models.FalseMahjong:
		err = game.FalseMahjong(penalizedID, value, forceNextRound)
	default:
		return models.GameNow{}, fmt.Errorf("unexpected penalty type: %v", penaltyType)
	}
	if err != nil {
		return models.GameNow{}, err
	}

	players := game.Players()
	score := game.Score()

	roundOutcome := models.RoundData{
		Round:           game.Round(),
		Outcome:         penaltyType,
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

	now, err := s.Now(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	return now, nil
}
