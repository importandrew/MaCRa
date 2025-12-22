package game

import (
	"github.com/importandrew/MaCRa/models"
)

func (s *Service) GetLiveGame(gameID int) (GameEngine, error) {
	g, ok := s.liveGames[gameID]
	if !ok {
		return nil, models.ErrGameNotLive
	}

	return g, nil
}

func (s *Service) CurrentGame(playerID int) (GameEngine, error) {
	for _, g := range s.liveGames {
		for _, id := range g.Players() {
			if playerID == id {
				return g, nil
			}
		}
	}

	return nil, models.ErrGameNotFound
}
