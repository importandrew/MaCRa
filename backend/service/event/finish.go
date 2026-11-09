package event

import (
	"context"
)

// Finish marks the event finished. The repo checks, under a lock, that it
// isn't finished yet and has no live games.
func (s *Service) Finish(ctx context.Context, eventID int) error {
	return s.repo.Finish(ctx, eventID)
}
