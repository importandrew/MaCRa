package game

import "context"

// Finish ends a live game, finalizes its result and removes it from liveGames.
func (s *Service) Finish(ctx context.Context, gameID int) error {
	game, err := s.GetLiveGame(gameID)
	if err != nil {
		return err
	}

	if err := game.End(); err != nil {
		return err
	}

	score := game.Score()
	players := game.Players()

	if err := s.finalize(ctx, gameID, score, players); err != nil {
		return err
	}

	return nil
}
