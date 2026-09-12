package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) Update(ctx context.Context, event models.Event) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, `
		UPDATE events
		SET name = $2,
			description = $3,
			visibility = $4,
			ratings_reach = $5,
			games_reach = $6
		WHERE id = $1`,
		event.ID, event.Name, event.Description,
		event.Visibility, event.RatingsReach, event.GamesReach,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrEventNotFound
	}

	return nil
}
