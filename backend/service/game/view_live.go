package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) ViewLive(ctx context.Context, gameID int) (models.GameView, error) {
	game, log, err := s.LiveGame(ctx, gameID)
	if err != nil {
		return models.GameView{}, err
	}

	score := game.Score()
	round := game.Round()

	// Seats keep their players' IDs; the names are of whoever sits there.
	names, err := s.seatNames(ctx, game.ID())
	if err != nil {
		return models.GameView{}, err
	}
	playerScores := make([]models.PlayerWindScore, 0, len(models.Winds))
	for _, wind := range models.Winds {
		playerID := game.PlayerBySeatInRound(round, wind)
		playerScores = append(playerScores, models.PlayerWindScore{
			ID:    playerID,
			Name:  names[playerID],
			Wind:  game.CurrentWind(playerID).Kanji(),
			Score: score[playerID].GP,
		})
	}

	undone, err := s.gameRepo.UndoneEntries(ctx, game.ID())
	if err != nil {
		return models.GameView{}, err
	}

	limit, err := s.gameRepo.TimeLimit(ctx, game.ID())
	if err != nil {
		return models.GameView{}, err
	}

	view := models.GameView{
		ID: game.ID(),
		// LiveGame loaded a live game: one stopped at its tour's time has
		// ended in the replay but stays live until it is finished.
		Status:  models.GameLive,
		LogSize: game.LogSize(),
		Rules:   game.Rules(),
		Round: models.Round{
			Num:   round,
			Kanji: models.RoundKanji(round),
		},
		PlayerScores: playerScores,
		Log:          gameLog(game, log, undone, names),
		Timer:        gameTimer(game, limit),
	}

	return view, nil
}
