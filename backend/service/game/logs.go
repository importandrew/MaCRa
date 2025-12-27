package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Logs(ctx context.Context, id int) (models.GameLog, error) {
	// maybe player header should be obtained by separate api call
	players, err := s.Players(ctx, id)
	if err != nil {
		return models.GameLog{}, err
	}

	roundLogs, err := s.gameRepo.GameLog(ctx, id)
	if err != nil {
		return models.GameLog{}, err
	}

	playerHeaders := make([]models.PlayerHeader, 0, len(players.Players))
	for _, player := range players.Players {
		playerHeaders = append(playerHeaders, models.PlayerHeader(player))
	}

	gameLog := models.GameLog{
		ID:        id,
		Players:   playerHeaders,
		RoundLogs: roundLogs,
	}

	return gameLog, nil
}
