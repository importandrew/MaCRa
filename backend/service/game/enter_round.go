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
	eventID int,
) (models.GameNow, error) {
	game, err := s.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	round := game.Round()

	err = game.EnterRoundResult(outcome, value, loserID, winnerID)
	if err != nil {
		return models.GameNow{}, err
	}

	players := game.Players()
	score := game.Score()

	roundOutcome := models.RoundData{
		Round:           round,
		Outcome:         outcome,
		WinnerID:        winnerID,
		LoserID:         loserID,
		Value:           value,
		EastStartScore:  score[players[models.East]].GP,
		SouthStartScore: score[players[models.South]].GP,
		WestStartScore:  score[players[models.West]].GP,
		NorthStartScore: score[players[models.North]].GP,
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
