package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateTournament sets a tournament's settings. It refuses a regular event
// (models.ErrNotTournament) and a finished one
// (models.ErrEventAlreadyFinished).
func (r *Repo) UpdateTournament(ctx context.Context, eventID int, settings models.TournamentSettings) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, err := tx.lockOpenEvent(ctx, eventID); err != nil {
			return err
		}
		res, err := tx.db.Exec(ctx, `
			UPDATE events SET tournament = $1 WHERE id = $2 AND type = $3`,
			settings, eventID, models.EventTournament)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return models.ErrNotTournament
		}
		return nil
	})
}
