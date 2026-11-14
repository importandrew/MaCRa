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

	undone, err := s.gameRepo.UndoneEntries(ctx, id)
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
		Undone:    undone,
	}

	return gameLog, nil
}

// gameLog is the round log of game as the client gets it: the players in
// starting seat order (east, south, west, north) with their names, the
// entries and the undone ones. names maps each player's ID to the name.
func gameLog(game GameEngine, log []models.RoundLog, undone []models.UndoneEntry, names map[int]string) models.GameLog {
	seating := game.Players()
	players := make([]models.PlayerHeader, 0, len(seating))
	for _, id := range seating {
		players = append(players, models.PlayerHeader{ID: id, Name: names[id]})
	}

	return models.GameLog{
		ID:        game.ID(),
		Players:   players,
		RoundLogs: log,
		Undone:    undone,
	}
}
