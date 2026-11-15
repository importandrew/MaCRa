package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// Participants returns a tournament's players with their numbers, for its
// desk.
func (s *Service) Participants(ctx context.Context, eventID int) (models.TournamentParticipants, error) {
	return s.repo.ListParticipants(ctx, eventID)
}

// Draw gives the tournament's players new random numbers.
func (s *Service) Draw(ctx context.Context, eventID int) error {
	return s.repo.Draw(ctx, eventID)
}

// Withdraw marks a participant of a tournament in progress as withdrawn.
func (s *Service) Withdraw(ctx context.Context, eventID, userID int) error {
	return s.repo.Withdraw(ctx, eventID, userID)
}
