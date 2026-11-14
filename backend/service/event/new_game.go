package event

import (
	"context"
)

// NewGame creates a game in the event with the seating. The game repo checks,
// under its locks, that the event is open, that all four have the player
// role (models.ErrMissingPlayerRole) and that none of them is in another
// live game; the game service checks the seating itself.
func (s *Service) NewGame(ctx context.Context, eventID int, playerIDs []int) (int, error) {
	return s.gameService.NewGame(ctx, playerIDs, eventID)
}
