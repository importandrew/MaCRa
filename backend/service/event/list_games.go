package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

const (
	defaultGamesListLimit = 50
	maxGamesListLimit     = 100
)

func (s *Service) ListGamesByStatus(
	ctx context.Context,
	eventID int,
	status models.GameStatus,
	limit, offset int,
) (models.EventGames, error) {
	if status != models.GameLive && status != models.GameEnded {
		return models.EventGames{}, models.ErrInvalidGameStatus
	}

	if limit <= 0 {
		limit = defaultGamesListLimit
	}
	if limit > maxGamesListLimit {
		limit = maxGamesListLimit
	}
	if offset < 0 {
		offset = 0
	}

	// Fetch one extra row to detect whether more pages exist.
	rows, err := s.repo.ListGamesByStatus(ctx, eventID, status, limit+1, offset)
	if err != nil {
		return models.EventGames{}, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	gameIDs := make([]int, len(rows))
	for i, row := range rows {
		gameIDs[i] = row.ID
	}

	var games []models.EventGameItem
	switch status {
	case models.GameLive:
		games, err = s.listLiveGames(ctx, gameIDs)
	case models.GameEnded:
		games, err = s.listEndedGames(ctx, gameIDs)
	default:
		return models.EventGames{}, models.ErrInvalidGameStatus
	}
	if err != nil {
		return models.EventGames{}, err
	}

	eventGames := models.EventGames{
		Games:   games,
		HasMore: hasMore,
	}

	return eventGames, nil
}

func (s *Service) listLiveGames(
	ctx context.Context,
	gameIDs []int,
) ([]models.EventGameItem, error) {
	games := make([]models.EventGameItem, 0, len(gameIDs))

	for _, gameID := range gameIDs {
		view, err := s.gameService.ViewLive(ctx, gameID)
		if err != nil {
			return nil, err
		}

		playerScores := make([]models.EventGamePlayerScore, 0, len(view.PlayerScores))
		for _, ps := range view.PlayerScores {
			playerScores = append(playerScores, models.EventGamePlayerScore{
				ID:   ps.ID,
				Name: ps.Name,
				Wind: ps.Wind,
				GP:   ps.Score,
			})
		}

		round := view.Round
		games = append(games, models.EventGameItem{
			ID:      view.ID,
			Status:  view.Status,
			Round:   &round,
			Players: playerScores,
		})
	}

	return games, nil
}

func (s *Service) listEndedGames(
	ctx context.Context,
	gameIDs []int,
) ([]models.EventGameItem, error) {
	results, err := s.gameService.GameResults(ctx, gameIDs)
	if err != nil {
		return nil, err
	}

	games := make([]models.EventGameItem, 0, len(gameIDs))

	for _, gameID := range gameIDs {
		result := results[gameID]
		playerScores := make([]models.EventGamePlayerScore, 0, len(result.PlayerResults))
		for _, pr := range result.PlayerResults {
			playerScores = append(playerScores, models.EventGamePlayerScore{
				ID:   pr.ID,
				Name: pr.Name,
				GP:   pr.GP,
			})
		}

		games = append(games, models.EventGameItem{
			ID:      gameID,
			Status:  models.GameEnded,
			Players: playerScores,
		})
	}

	return games, nil
}
