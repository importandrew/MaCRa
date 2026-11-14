package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) List(ctx context.Context) ([]models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, type, name, description, owner_user_id,
			visibility, ratings_reach, games_reach, rules, settings,
			created_at, finished_at
		FROM events e
		ORDER BY e.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		return nil, err
	}

	return events, nil
}
