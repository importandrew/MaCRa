package game

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// Penalty enters a penalty or a false mahjong into the player's live game in
// the event, at their table, and returns the table. A penalty must be one of
// models.PenaltyAmounts (models.ErrValidatePenalty); the engine checks the
// rest, also the false mahjong amounts.
func (s *Service) Penalty(ctx context.Context, penaltyType models.Outcome, playerID, penalizedID, value int, forceNextRound bool, eventID int, seenLogSize int) (models.GameNow, error) {
	if err := validatePenaltyAmount(penaltyType, value); err != nil {
		return models.GameNow{}, err
	}

	game, _, err := s.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	if err := s.applyPenalty(ctx, game, penaltyType, penalizedID, value, forceNextRound, seenLogSize); err != nil {
		return models.GameNow{}, err
	}

	if game.Status() == models.GameEnded {
		return models.GameNow{ID: game.ID(), Status: models.GameEnded}, nil
	}

	return s.Now(ctx, playerID, eventID)
}

// PenaltyInGame enters a penalty or a false mahjong into the live game, as a
// referee does from the game's view, and returns that view. The amounts are
// checked as in Penalty; a game that isn't live is models.ErrGameNotLive.
func (s *Service) PenaltyInGame(ctx context.Context, penaltyType models.Outcome, gameID, penalizedID, value int, forceNextRound bool, seenLogSize int) (models.GameView, error) {
	if err := validatePenaltyAmount(penaltyType, value); err != nil {
		return models.GameView{}, err
	}

	game, _, err := s.LiveGame(ctx, gameID)
	if err != nil {
		return models.GameView{}, err
	}

	if err := s.applyPenalty(ctx, game, penaltyType, penalizedID, value, forceNextRound, seenLogSize); err != nil {
		return models.GameView{}, err
	}

	if game.Status() == models.GameEnded {
		return models.GameView{ID: game.ID(), Status: models.GameEnded}, nil
	}

	return s.ViewLive(ctx, gameID)
}

// validatePenaltyAmount checks a new penalty against models.PenaltyAmounts.
// Only new entries are checked: a game is replayed with the amounts it has.
func validatePenaltyAmount(penaltyType models.Outcome, value int) error {
	if penaltyType == models.Penalty && !slices.Contains(models.PenaltyAmounts, value) {
		return fmt.Errorf("%w: penalty should be one of %v, got %d",
			models.ErrValidatePenalty, models.PenaltyAmounts, value)
	}
	return nil
}

// applyPenalty enters the penalty into the loaded game, saves it as the next
// log entry (refused with models.ErrLogChanged unless the game still has
// seenLogSize entries) and finalizes the game if it has ended.
func (s *Service) applyPenalty(ctx context.Context, game GameEngine, penaltyType models.Outcome, penalizedID, value int, forceNextRound bool, seenLogSize int) error {
	logSize := game.LogSize()
	if err := checkLogSize(logSize, seenLogSize); err != nil {
		return err
	}

	var err error
	switch penaltyType {
	case models.Penalty:
		err = game.Penalty(penalizedID, value, forceNextRound)
	case models.FalseMahjong:
		err = game.FalseMahjong(penalizedID, value, forceNextRound)
	default:
		return fmt.Errorf("unexpected penalty type: %v", penaltyType)
	}
	if err != nil {
		return err
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

	if err := s.gameRepo.SaveRoundOutcome(ctx, game.ID(), roundOutcome); err != nil {
		return err
	}

	if game.Status() == models.GameEnded {
		return s.finalize(ctx, game.ID(), score, players)
	}

	return nil
}
