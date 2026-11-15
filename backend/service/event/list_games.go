package event

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
)

const (
	defaultGamesListLimit = 50
	maxGamesListLimit     = 100
)

func (s *Service) ListGamesByStatus(
	ctx context.Context,
	eventID, callerID int,
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

	hiddenTour, err := s.hiddenTour(ctx, eventID, callerID, status)
	if err != nil {
		return models.EventGames{}, err
	}
	// Fetch one extra row to detect whether more pages exist.
	rows, err := s.repo.ListGamesByStatus(ctx, eventID, status, hiddenTour.ID, callerID, limit+1, offset)
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
	eventGames.HiddenTour = hiddenTour.Number

	return eventGames, nil
}

// hiddenTour returns the live tour whose tables the caller may not see in
// the list of live games (tourTableHidden: not staff, and their own game in
// it hasn't ended); the zero LiveTour when none is hidden.
func (s *Service) hiddenTour(ctx context.Context, eventID, callerID int, status models.GameStatus) (models.LiveTour, error) {
	if status != models.GameLive {
		return models.LiveTour{}, nil
	}
	event, err := s.repo.Get(ctx, eventID)
	if err != nil || event.Type != models.EventTournament {
		return models.LiveTour{}, err
	}
	authCtx, err := s.buildAuthContext(ctx, event, callerID, 0)
	if err != nil || authCtx.isStaff() {
		return models.LiveTour{}, err
	}
	tour, err := s.repo.LiveTour(ctx, eventID, callerID)
	switch {
	case errors.Is(err, models.ErrTourNotFound) || err == nil && tour.Finished:
		return models.LiveTour{}, nil
	case err != nil:
		return models.LiveTour{}, err
	}
	return tour, nil
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
			Hands:   view.Rules.Hands,
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
