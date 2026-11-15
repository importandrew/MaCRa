package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Players(ctx context.Context, gameID int) (models.GamePlayers, error) {
	players, err := s.gameRepo.Players(ctx, gameID)
	if err != nil {
		return models.GamePlayers{}, err
	}

	// Seats keep their players' IDs; the names are of whoever sits there.
	names, err := s.seatNames(ctx, gameID)
	if err != nil {
		return models.GamePlayers{}, err
	}
	seat := func(id int) models.Player { return models.Player{ID: id, Name: names[id]} }
	east, south, west, north := seat(players.East), seat(players.South), seat(players.West), seat(players.North)

	gamePlayers := models.GamePlayers{
		ID:      gameID,
		Players: []models.Player{east, south, west, north},
	}

	return gamePlayers, nil
}
