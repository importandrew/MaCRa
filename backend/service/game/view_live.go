package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) ViewLive(ctx context.Context, gameID int) (models.GameView, error) {
	game, err := s.GetLiveGame(gameID)
	if err != nil {
		return models.GameView{}, err
	}

	score := game.Score()
	round := game.Round()

	playerScores := make([]models.PlayerWindScore, 0, len(models.Winds))
	for _, wind := range models.Winds {
		playerID := game.PlayerBySeatInRound(round, wind)

		player, err := s.playerRepo.Get(ctx, playerID)
		if err != nil {
			return models.GameView{}, err
		}

		playerScores = append(playerScores, models.PlayerWindScore{
			ID:    playerID,
			Name:  player.Name,
			Wind:  game.CurrentWind(playerID).Kanji(),
			Score: score[playerID].GP,
		})
	}

	view := models.GameView{
		ID:     game.ID(),
		Status: game.Status(),
		Round: models.Round{
			Num:   round,
			Kanji: models.RoundKanji(round),
		},
		PlayerScores: playerScores,
	}

	return view, nil
}
