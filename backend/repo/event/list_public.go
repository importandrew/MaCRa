package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListPublic(ctx context.Context) ([]models.Event, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, type, name, description, owner_user_id,
			visibility, ratings_reach, games_reach,
			created_at, finished_at
		FROM events
		WHERE visibility = 'public'
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		return nil, err
	}

	return events, nil
}
