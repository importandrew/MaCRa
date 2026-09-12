package event

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) Get(ctx context.Context, eventID int) (models.Event, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.Event{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, type, name, description, owner_user_id,
			visibility, ratings_reach, games_reach,
			created_at, finished_at
		FROM events WHERE id = $1`,
		eventID,
	)
	if err != nil {
		return models.Event{}, err
	}

	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Event{}, models.ErrEventNotFound
		}
		return models.Event{}, err
	}

	return event, nil
}
