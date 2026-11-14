package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateSettings sets the event's settings, or returns
// models.ErrEventNotFound.
func (r *Repo) UpdateSettings(ctx context.Context, eventID int, settings models.EventSettings) error {
	res, err := r.db.Exec(ctx, `UPDATE events SET settings = $1 WHERE id = $2`, settings, eventID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrEventNotFound
	}

	return nil
}
