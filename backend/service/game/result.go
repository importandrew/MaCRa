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

func (s *Service) GameResults(ctx context.Context, gameIDs []int) (map[int]models.GameResult, error) {
	if len(gameIDs) == 0 {
		return map[int]models.GameResult{}, nil
	}

	results, err := s.gameRepo.GameResults(ctx, gameIDs)
	if err != nil {
		return nil, err
	}

	for id, result := range results {
		sort.Slice(result.PlayerResults, func(i, j int) bool {
			return result.PlayerResults[i].TP > result.PlayerResults[j].TP
		})
		results[id] = result
	}

	return results, nil
}
