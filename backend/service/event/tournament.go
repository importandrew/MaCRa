package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateTournament sets a tournament's settings (Manage → Settings). Until
// tours exist, all of them may change; the repo refuses a regular or a
// finished event.
func (s *Service) UpdateTournament(ctx context.Context, eventID int, settings models.TournamentSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateTournament(ctx, eventID, settings)
}
