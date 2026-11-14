package game

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// Penalty enters a penalty or a false mahjong into the player's live game. A
// penalty must be one of models.PenaltyAmounts (models.ErrValidatePenalty);
// the engine checks the rest, also the false mahjong amounts.
func (s *Service) Penalty(ctx context.Context, penaltyType models.Outcome, playerID, penalizedID, value int, forceNextRound bool, eventID int, seenLogSize int) (models.GameNow, error) {
	if penaltyType == models.Penalty && !slices.Contains(models.PenaltyAmounts, value) {
		return models.GameNow{}, fmt.Errorf("%w: penalty should be one of %v, got %d",
			models.ErrValidatePenalty, models.PenaltyAmounts, value)
	}

	game, _, err := s.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	logSize := game.LogSize()
	if err := checkLogSize(logSize, seenLogSize); err != nil {
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
		LogSize:         logSize,
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
