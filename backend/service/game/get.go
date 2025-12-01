package game

import (
	"errors"
	"fmt"
)

var errGameNotFound = errors.New("game not found")

func (s *Service) GetLiveGame(gameID int) (GameEngine, error) {
	g, ok := s.liveGames[gameID]
	if !ok {
		return nil, errGameNotFound
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

	return nil, fmt.Errorf("no current game for player ID %d", playerID)
}
