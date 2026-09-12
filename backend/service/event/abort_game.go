package event

import (
	"context"
)

func (s *Service) AbortGame(ctx context.Context, gameID int) error {
	return s.gameService.Abort(ctx, gameID)
}

func (s *Service) UndoLastRound(ctx context.Context, gameID int) error {
	return s.gameService.UndoLastRound(ctx, gameID)
}
