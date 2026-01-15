package game

import (
	"context"
	"sort"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) GameResult(ctx context.Context, gameID int) (models.GameResult, error) {
	status, err := s.Status(ctx, gameID)
	if err != nil {
		return models.GameResult{}, err
	}

	if status != models.GameEnded {
		return models.GameResult{}, models.ErrGameNotEnded
	}

	result, err := s.gameRepo.GameResult(ctx, gameID)
	if err != nil {
		return models.GameResult{}, err
	}

	// sort from first place to last
	sort.Slice(result.PlayerResults, func(i, j int) bool {
		return result.PlayerResults[i].TP > result.PlayerResults[j].TP
	})

	return result, nil
}
