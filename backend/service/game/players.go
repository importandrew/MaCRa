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

	east, err := s.playerRepo.Get(ctx, players.East)
	if err != nil {
		return models.GamePlayers{}, err
	}

	south, err := s.playerRepo.Get(ctx, players.South)
	if err != nil {
		return models.GamePlayers{}, err
	}

	west, err := s.playerRepo.Get(ctx, players.West)
	if err != nil {
		return models.GamePlayers{}, err
	}

	north, err := s.playerRepo.Get(ctx, players.North)
	if err != nil {
		return models.GamePlayers{}, err
	}

	gamePlayers := models.GamePlayers{
		ID:      gameID,
		Players: []models.Player{east, south, west, north},
	}

	return gamePlayers, nil
}
